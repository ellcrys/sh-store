package servers

import (
	"fmt"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"strings"

	"github.com/ellcrys/util"
	"github.com/jinzhu/copier"
	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/ncodes/safehold/session"
	"golang.org/x/net/context"
)

var (
	// MaxObjectPerPut is the maximum number of objects allowed in a single PUT operation
	MaxObjectPerPut = 25
)

// validateObjects validates a slice of objects returning all errors found.
// - owner id is required
// - owner id of all objects must be the same
// - owner id must exist
// - key is required
func (s *RPC) validateObjects(obj []*proto_rpc.Object) []common.Error {
	var errs []common.Error
	var ownerID = obj[0].OwnerID

	for i, o := range obj {

		if len(strings.TrimSpace(o.OwnerID)) == 0 {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: fmt.Sprintf("object %d: owner id is required", i), Field: "objects"})
			continue
		}

		if ownerID != o.OwnerID {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: fmt.Sprintf("object %d: has a different owner. All objects must be owned by a single identity", i), Field: "objects"})
			continue
		}

		_, err := s.object.GetLast(&tables.Object{ID: o.OwnerID})
		if err != nil {
			if err == patchain.ErrNotFound {
				errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: fmt.Sprintf("object %d: owner id does not exist", i), Field: "objects"})
				continue
			}
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: fmt.Sprintf("object %d: server error", i), Field: "objects"})
			continue
		}

		if len(strings.TrimSpace(o.Key)) == 0 {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: fmt.Sprintf("object %d: key is required", i), Field: "objects"})
			continue
		}
	}
	return errs
}

// CreateObjects creates one or more objects
func (s *RPC) CreateObjects(ctx context.Context, req *proto_rpc.CreateObjectsMsg) (*proto_rpc.MultiObjectResponse, error) {

	sessionID := util.FromIncomingMD(ctx, "session_id")
	developerID := ctx.Value(CtxIdentity).(string)
	authorization := util.FromIncomingMD(ctx, "authorization")
	localOnly := util.FromIncomingMD(ctx, "local-only") == "true"

	if len(req.Objects) == 0 {
		return nil, common.NewSingleAPIErr(400, "", "", "no object provided. At least one object is required", nil)
	}

	if len(req.Objects) > MaxObjectPerPut {
		return nil, common.NewSingleAPIErr(400, "", "", fmt.Sprintf("too many objects. Only a maximum of %d can be created at once", MaxObjectPerPut), nil)
	}

	if errs := s.validateObjects(req.Objects); len(errs) > 0 {
		return nil, common.NewMultiAPIErr(400, "validation errors", errs)
	}

	// ensure caller has permission to PUT objects on behalf of the object owner
	if developerID != req.Objects[0].OwnerID {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you are not authorized to create objects for the owner", nil)
	}

	// set developer as the creator id
	for _, obj := range req.Objects {
		obj.CreatorID = developerID
	}

	fullSessionID := makeDBSessionID(developerID, sessionID)
	var objs []*tables.Object
	copier.Copy(&objs, req.Objects)
	sessionOp := &session.Op{
		OpType: session.OpPutObjects,
		Data:   objs,
		Done:   make(chan struct{}),
	}

	if sessionID != "" {

		// use session if it exists in this server's in-memory session cache
		if s.dbSession.HasSession(fullSessionID) {
			if err := s.dbSession.SendOp(fullSessionID, sessionOp); err != nil {
				logRPC.Errorf("%+v", err)
				return nil, common.NewSingleAPIErr(500, common.CodePutError, "objects", err.Error(), nil)
			}
			resp, _ := NewMultiObjectResponse("object", objs)
			return resp, nil
		}

		if localOnly {
			return nil, fmt.Errorf("session not found")
		}

		// find session in the registry
		sessionItem, err := s.sessionReg.Get(fullSessionID)
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

	fullSessionID = makeDBSessionID(developerID, util.UUID4())
	fullSessionID, err := s.dbSession.CreateSession(fullSessionID, developerID)
	if err != nil {
		logRPC.Debugf("%+v", err)
		return nil, common.NewSingleAPIErr(500, "", "", "session not created", nil)
	}
	defer s.dbSession.CommitEnd(fullSessionID)

	if err := s.dbSession.SendOp(fullSessionID, sessionOp); err != nil {
		logRPC.Errorf("%+v", err)
		return nil, common.NewSingleAPIErr(500, common.CodePutError, "objects", err.Error(), nil)
	}

	resp, _ := NewMultiObjectResponse("object", objs)
	return resp, nil
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
func (s *RPC) GetObjects(ctx context.Context, req *proto_rpc.GetObjectMsg) (*proto_rpc.MultiObjectResponse, error) {

	authorization := util.FromIncomingMD(ctx, "authorization")
	developerID := ctx.Value(CtxIdentity).(string)
	sessionID := util.FromIncomingMD(ctx, "session_id")
	fullSessionID := makeDBSessionID(developerID, sessionID)
	localOnly := util.FromIncomingMD(ctx, "local-only") == "true"

	// check if session exist in the in-memory session cache,
	// use it else check if it exist on the session registry. If it does,
	// forward the request to the associated host
	if sessionID != "" {

		if !s.dbSession.HasSession(fullSessionID) {

			if localOnly {
				return nil, fmt.Errorf("session not found")
			}

			// find session in the registry
			sessionItem, err := s.sessionReg.Get(fullSessionID)
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

	} else {
		fullSessionID = makeDBSessionID(developerID, util.UUID4())
		s.dbSession.CreateUnregisteredSession(fullSessionID, developerID)
		defer s.dbSession.CommitEnd(fullSessionID)
	}

	// if owner is set and not the same as the developer id, check if exists
	if len(req.Owner) > 0 && req.Owner != developerID {
		if err := session.SendQueryOpWithSession(s.dbSession, fullSessionID, `{ "id": "`+req.Owner+`" }`, 1, "", &tables.Object{}); err != nil {
			if err == patchain.ErrNotFound {
				return nil, common.NewSingleAPIErr(404, "", "", "owner not found", nil)
			}
			logRPC.Errorf("%+v", err)
			return nil, common.ServerError
		}
	}

	// if creator is set and not the same as the developer id, check if it exists
	if len(req.Creator) > 0 && req.Creator != developerID {
		if err := session.SendQueryOpWithSession(s.dbSession, fullSessionID, `{ "id": "`+req.Creator+`" }`, 1, "", &tables.Object{}); err != nil {
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
	var queryAsMap map[string]interface{}
	util.FromJSON(req.Query, &queryAsMap)
	queryAsMap["owner_id"] = req.Owner
	queryAsMap["creator_id"] = req.Creator

	var objs []*tables.Object
	if err := session.SendQueryOpWithSession(s.dbSession, fullSessionID, string(util.MustStringify(queryAsMap)), int(req.Limit), orderByToString(req.Order), &objs); err != nil {
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	resp, _ := NewMultiObjectResponse("object", objs)
	return resp, nil
}
