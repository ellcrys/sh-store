package servers

import (
	"fmt"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"strings"

	"github.com/ellcrys/util"
	"github.com/fatih/structs"
	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/ncodes/safehold/session"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

var (
	// MaxObjectPerPut is the maximum number of objects allowed in a single PUT operation
	MaxObjectPerPut = 25
)

// CreateObjects creates one or more objects using a session.
// If session id is provide, it checks whether the session exists locally or
// forwards the request to the remote session host.
// Supports the use of mapping to unmap custom fields in the objects to be created.
func (s *RPC) CreateObjects(ctx context.Context, req *proto_rpc.CreateObjectsMsg) (*proto_rpc.GetObjectsResponse, error) {

	var err error
	sessionID := util.FromIncomingMD(ctx, "session_id")
	developerID := ctx.Value(CtxIdentity).(string)
	sid := makeDBSessionID(developerID, sessionID)
	authorization := util.FromIncomingMD(ctx, "authorization")
	localOnly := util.FromIncomingMD(ctx, "local-only") == "true"

	if len(sessionID) > 0 {

		if !s.dbSession.HasSession(sid) {

			// abort further operations
			if localOnly {
				return nil, fmt.Errorf("session not found")
			}

			// find session in the registry
			sessionItem, err := s.sessionReg.Get(sid)
			if err != nil {
				if err == session.ErrNotFound {
					return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
				}
				return nil, common.ServerError
			}

			sessionHostAddr := net.JoinHostPort(sessionItem.Address, strconv.Itoa(sessionItem.Port))
			client, err := grpc.Dial(sessionHostAddr, grpc.WithInsecure())
			if err != nil {
				return nil, common.ServerError
			}
			defer client.Close()

			// make call to the session host. Include the session_id, auth token of the current request
			// and set local-only to force the RPC method to only perform local object creation
			server := proto_rpc.NewAPIClient(client)
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", authorization, "local-only", "true"))
			resp, err := server.CreateObjects(ctx, req)
			if err != nil {
				if grpc.ErrorDesc(err) == "session not found" {
					return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
				}
				logRPC.Errorf("%+v", err)
				return nil, common.ServerError
			}

			return resp, nil
		}
	} else { // session id not provided, create local, unregistered session
		sid = makeDBSessionID(developerID, util.UUID4())
		s.dbSession.CreateUnregisteredSession(sid, developerID)
		defer s.dbSession.CommitEnd(sid)
	}

	var objs []map[string]interface{}
	if err := util.FromJSON(req.Objects, &objs); err != nil {
		s.dbSession.RollbackEnd(sid)
		return nil, common.NewSingleAPIErr(400, "", "", "failed to parse objects", nil)
	}

	// add creator_id and include developer as the owner_id if any object is missing it
	for _, obj := range objs {
		obj["creator_id"] = developerID
		if obj["owner_id"] == nil {
			obj["owner_id"] = developerID
		}
	}

	var mapping map[string]string

	// if mapping name is provided, get the mapping and unmap the objects
	if len(req.Mapping) > 0 {
		mapping, err = s.getMappingWithSession(sid, req.Mapping, developerID)
		if err != nil {
			s.dbSession.RollbackEnd(sid)
			return nil, err
		}

		if err := common.UnMapFields(mapping, objs); err != nil {
			s.dbSession.RollbackEnd(sid)
			logRPC.Errorf("%+v", errors.Wrap(err, "failed to unmap object"))
			return nil, common.ServerError
		}
	}

	vErrs, err := s.validateObjects(objs, mapping)
	if err != nil {
		s.dbSession.RollbackEnd(sid)
		logRPC.Errorf("%+v", errors.Wrap(err, "failed to validate objects"))
		return nil, common.ServerError
	}

	if len(vErrs) > 0 {
		s.dbSession.RollbackEnd(sid)
		return nil, common.NewMultiAPIErr(400, "validation errors", vErrs)
	}

	// ensure caller has permission to PUT objects on behalf of the object owner
	if developerID != objs[0]["owner_id"].(string) {
		s.dbSession.RollbackEnd(sid)
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you are not authorized to create objects for the owner", nil)
	}

	var patchainObjs []*tables.Object

	// create final patchain objects
	for _, obj := range objs {
		var o tables.Object
		util.CopyToStruct(&o, obj)
		patchainObjs = append(patchainObjs, &o)
	}

	if err := session.SendPutOpWithSession(s.dbSession, sid, patchainObjs); err != nil {
		logRPC.Errorf("%+v", err)
		s.dbSession.RollbackEnd(sid)
		return nil, common.NewSingleAPIErr(500, common.CodePutError, "objects", err.Error(), nil)
	}

	// get the updated object and re-apply mapping
	objs = nil
	for _, updatedObj := range patchainObjs {
		obj := structs.New(updatedObj).Map()
		common.MapFields(mapping, obj)
		objs = append(objs, obj)
	}

	return &proto_rpc.GetObjectsResponse{
		Objects: util.MustStringify(objs),
	}, nil
}

// orderByList takes a slice of *proto_rpc.OrderBy and returns
// a concatenated gorm.Order compatible string
func orderByToString(orderByList []*proto_rpc.OrderBy) string {
	var s []string
	for _, o := range orderByList {
		ord := "asc"
		if o.Order <= 0 {
			ord = "desc"
		}
		s = append(s, fmt.Sprintf("%s %s", o.Field, ord))
	}
	return strings.Join(s, ", ")
}

// GetObjects fetches objects belonging to an identity
func (s *RPC) GetObjects(ctx context.Context, req *proto_rpc.GetObjectMsg) (*proto_rpc.GetObjectsResponse, error) {

	var err error
	authorization := util.FromIncomingMD(ctx, "authorization")
	developerID := ctx.Value(CtxIdentity).(string)
	sessionID := util.FromIncomingMD(ctx, "session_id")
	sid := makeDBSessionID(developerID, sessionID)
	localOnly := util.FromIncomingMD(ctx, "local-only") == "true"

	// Handle session if provided
	if len(sessionID) > 0 {

		// if session id does not exist locally, find it in the session registry and forward
		// request to the session host
		if !s.dbSession.HasSession(sid) {

			// abort further operations
			if localOnly {
				return nil, fmt.Errorf("session not found")
			}

			// find session in the registry
			sessionItem, err := s.sessionReg.Get(sid)
			if err != nil {
				if err == session.ErrNotFound {
					return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
				}
				return nil, common.ServerError
			}

			sessionHostAddr := net.JoinHostPort(sessionItem.Address, strconv.Itoa(sessionItem.Port))
			client, err := grpc.Dial(sessionHostAddr, grpc.WithInsecure())
			if err != nil {
				return nil, common.ServerError
			}
			defer client.Close()

			// make call to the session host server.
			// include the session_id, auth token of the current request
			// and set local-only to force the RPC method
			// to only perform local object creation
			server := proto_rpc.NewAPIClient(client)
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", authorization, "local-only", "true"))
			resp, err := server.GetObjects(ctx, req)
			if err != nil {
				if grpc.ErrorDesc(err) == "session not found" {
					return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
				}
				logRPC.Errorf("%+v", err)
				return nil, common.ServerError
			}

			return resp, nil
		}

	} else { // session id not provided, create local, unregistered session
		sid = makeDBSessionID(developerID, util.UUID4())
		s.dbSession.CreateUnregisteredSession(sid, developerID)
		defer s.dbSession.CommitEnd(sid)
	}

	// if owner is set and not the same as the developer id, check if it exists
	if len(req.Owner) > 0 && req.Owner != developerID {
		if err = session.SendQueryOpWithSession(s.dbSession, sid, "", &tables.Object{ID: req.Owner}, 1, "", &tables.Object{}); err != nil {
			if err == patchain.ErrNotFound {
				return nil, common.NewSingleAPIErr(404, "", "", "owner not found", nil)
			}
			logRPC.Errorf("%+v", err)
			return nil, common.ServerError
		}
	}

	// if creator is set and not the same as the developer id, check if it exists
	if len(req.Creator) > 0 && req.Creator != developerID {
		if err = session.SendQueryOpWithSession(s.dbSession, sid, "", &tables.Object{ID: req.Creator}, 1, "", &tables.Object{}); err != nil {
			if err == patchain.ErrNotFound {
				return nil, common.NewSingleAPIErr(404, "", "", "creator not found", nil)
			}
			logRPC.Errorf("%+v", err)
			return nil, common.ServerError
		}
	}

	// set default owner as the developer if not set
	if len(req.Owner) == 0 {
		req.Owner = developerID
	}

	// set default creator as the developer if not set
	if len(req.Creator) == 0 {
		req.Creator = developerID
	}

	// developer is not the owner, this action requires permission
	// TODO: ensure auth token must be a user token from the owner
	// and the token authorizes access to the object created by the creator for this developer
	if developerID != req.Owner {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you are not authorized to access objects for the owner", nil)
	}

	// include owner and creator filters
	var query map[string]interface{}
	util.FromJSON(req.Query, &query)
	query["owner_id"] = req.Owner
	query["creator_id"] = req.Creator

	var mapping map[string]string

	// if mapping is provided in request, fetch it
	if len(req.Mapping) > 0 {
		mapping, err = s.getMappingWithSession(sid, req.Mapping, developerID)
		if err != nil {
			return nil, err
		}
	}

	// unmap query fields
	common.UnMapFields(mapping, query)

	var fetchedObjs []*tables.Object
	if err = session.SendQueryOpWithSession(s.dbSession, sid, string(util.MustStringify(query)), nil, int(req.Limit), orderByToString(req.Order), &fetchedObjs); err != nil {
		if strings.Contains(err.Error(), "unknown query field") {
			return nil, common.NewSingleAPIErr(400, common.CodeInvalidParam, "query", err.Error(), nil)
		}
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	// convert fetched objects to map and re-apply mapping
	var objs []map[string]interface{}
	for _, obj := range fetchedObjs {
		mObj := structs.New(obj).Map()
		common.MapFields(mapping, mObj)
		objs = append(objs, mObj)
	}

	return &proto_rpc.GetObjectsResponse{
		Objects: util.MustStringify(objs),
	}, nil
}

// CountObjects counts the number of objects that match a query
func (s *RPC) CountObjects(ctx context.Context, req *proto_rpc.GetObjectMsg) (*proto_rpc.ObjectCountResponse, error) {

	authorization := util.FromIncomingMD(ctx, "authorization")
	developerID := ctx.Value(CtxIdentity).(string)
	sessionID := util.FromIncomingMD(ctx, "session_id")
	sid := makeDBSessionID(developerID, sessionID)
	localOnly := util.FromIncomingMD(ctx, "local-only") == "true"

	// check if session exist in the in-memory session cache,
	// use it else check if it exist on the session registry. If it does,
	// forward the request to the associated host
	if sessionID != "" {
		if !s.dbSession.HasSession(sid) {

			if localOnly {
				return nil, fmt.Errorf("session not found")
			}

			// find session in the registry
			sessionItem, err := s.sessionReg.Get(sid)
			if err != nil {
				if err == session.ErrNotFound {
					return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
				}
				return nil, common.ServerError
			}

			sessionHostAddr := net.JoinHostPort(sessionItem.Address, strconv.Itoa(sessionItem.Port))
			client, err := grpc.Dial(sessionHostAddr, grpc.WithInsecure())
			if err != nil {
				return nil, common.ServerError
			}
			defer client.Close()

			// make call to the session host server.
			// include the session_id, auth token of the current request
			// and set local-only to force the RPC method
			// to only perform local object creation
			server := proto_rpc.NewAPIClient(client)
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", authorization, "local-only", "true"))
			resp, err := server.CountObjects(ctx, req)
			if err != nil {
				if grpc.ErrorDesc(err) == "session not found" {
					return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
				}
				logRPC.Errorf("%+v", err)
				return nil, common.ServerError
			}

			return resp, nil
		}
	} else {
		sid = makeDBSessionID(developerID, util.UUID4())
		s.dbSession.CreateUnregisteredSession(sid, developerID)
		defer s.dbSession.CommitEnd(sid)
	}

	// if owner is set and not the same as the developer id, check if it exists
	if len(req.Owner) > 0 && req.Owner != developerID {
		if err := session.SendQueryOpWithSession(s.dbSession, sid, "", &tables.Object{ID: req.Owner}, 1, "", &tables.Object{}); err != nil {
			if err == patchain.ErrNotFound {
				return nil, common.NewSingleAPIErr(404, "", "", "owner not found", nil)
			}
			logRPC.Errorf("%+v", err)
			return nil, common.ServerError
		}
	}

	// if creator is set and not the same as the developer id, check if it exists
	if len(req.Creator) > 0 && req.Creator != developerID {
		if err := session.SendQueryOpWithSession(s.dbSession, sid, "", &tables.Object{ID: req.Creator}, 1, "", &tables.Object{}); err != nil {
			if err == patchain.ErrNotFound {
				return nil, common.NewSingleAPIErr(404, "", "", "creator not found", nil)
			}
			logRPC.Errorf("%+v", err)
			return nil, common.ServerError
		}
	}

	// default owner and creator
	if len(req.Owner) == 0 {
		req.Owner = developerID
	}
	if len(req.Creator) == 0 {
		req.Creator = developerID
	}

	// developer is not the owner, this action requires permission
	// TODO: ensure auth token must be a user token from the owner
	// and the token authorizes access to the object created by the creator for this developer
	if developerID != req.Owner {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you are not authorized to access objects for the owner", nil)
	}

	// include owner and creator filters
	var query map[string]interface{}
	util.FromJSON(req.Query, &query)
	query["owner_id"] = req.Owner
	query["creator_id"] = req.Creator

	var count int64
	if err := session.SendCountOpWithSession(s.dbSession, sid, string(util.MustStringify(query)), nil, &count); err != nil {
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return &proto_rpc.ObjectCountResponse{
		Count: count,
	}, nil
}
