package servers

import (
	"fmt"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/util"
	"github.com/ncodes/mapvalidator"
	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/patchain/object"
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

// validateObjects validates a slice of objects returning all errors found.
// - empty object not allowed
// - max number of object must not be exceeded
// - owner id is required
// - owner id of all objects must be the same
// - owner id must exist
// - key is required
// Invalid field will be replaced with their mappings if mapping is provided
func (s *RPC) validateObjects(objs []map[string]interface{}, mapping map[string]string) ([]common.Error, error) {

	var errs []common.Error

	// get the mapped field for an object field or return the object field if it has no mapped field
	var getMappedField = func(objectField string) string {
		for mappedField, _objectField := range mapping {
			if objectField == _objectField {
				return mappedField
			}
		}
		return objectField
	}

	if len(objs) == 0 {
		errs = append(errs, common.Error{Message: "no object provided. At least one object is required"})
		return errs, nil
	}

	if len(objs) > MaxObjectPerPut {
		errs = append(errs, common.Error{Message: fmt.Sprintf("too many objects. Only a maximum of %d can be created at once", MaxObjectPerPut)})
		return errs, nil
	}

	var prevOwnerID string
	for i, o := range objs {
		rules := []mapvalidator.Rule{
			mapvalidator.RequiredWithMsg("owner_id", fmt.Sprintf("object %d: %s is required", i, getMappedField("owner_id"))),
			mapvalidator.TypeWithMsg("owner_id", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("owner_id"))),
			mapvalidator.RequiredWithMsg("key", fmt.Sprintf("object %d: %s is required", i, getMappedField("key"))),
			mapvalidator.TypeWithMsg("key", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("key"))),
			mapvalidator.TypeWithMsg("value", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("value"))),
			mapvalidator.TypeWithMsg("ref1", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref1"))),
			mapvalidator.TypeWithMsg("ref2", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref2"))),
			mapvalidator.TypeWithMsg("ref3", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref3"))),
			mapvalidator.TypeWithMsg("ref4", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref4"))),
			mapvalidator.TypeWithMsg("ref5", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref5"))),
			mapvalidator.TypeWithMsg("ref6", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref6"))),
			mapvalidator.TypeWithMsg("ref7", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref7"))),
			mapvalidator.TypeWithMsg("ref8", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref8"))),
			mapvalidator.TypeWithMsg("ref9", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref9"))),
			mapvalidator.TypeWithMsg("ref10", mapvalidator.TypeString, fmt.Sprintf("object %d: %s must be a string", i, getMappedField("ref10"))),
			mapvalidator.Custom("owner_id", fmt.Sprintf("object %d: all objects must share the same %s", i, getMappedField("owner_id")), func(fieldValue interface{}, fullMap map[string]interface{}) bool {
				if fv, ok := fieldValue.(string); ok {
					if prevOwnerID == "" {
						prevOwnerID = fv
					}
					return prevOwnerID == fv
				}
				return true
			}),
		}

		vErrs := mapvalidator.Validate(o, rules...)
		for _, err := range vErrs {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: err.Error()})
		}
	}

	if len(errs) > 0 {
		return errs, nil
	}

	_, err := s.object.GetLast(&tables.Object{ID: objs[0]["owner_id"].(string)})
	if err != nil {
		if err == patchain.ErrNotFound {
			errs = append(errs, common.Error{Code: common.CodeInvalidParam, Message: "owner of object(s) does not exist"})
		} else {
			return nil, err
		}
	}

	return errs, nil
}

// getMappingWithSession gets a mapping using an active session. It will
// find mapping belonging to the ownerID if the mapping name is not a UUIDv4
// string, otherwise it will find any mapping matching the mapping name.
func (s *RPC) getMappingWithSession(sid, mappingName, ownerID string) (map[string]string, error) {

	var mappingQuery = fmt.Sprintf(`{ "key":"%s", "owner_id":"%s" }`, object.MakeMappingKey(mappingName), ownerID)
	if govalidator.IsUUIDv4(mappingName) {
		mappingQuery = fmt.Sprintf(`{ "id": "` + mappingName + `" }`)
	}

	var mappingObj tables.Object
	if err := session.SendQueryOpWithSession(s.dbSession, sid, mappingQuery, 1, "", &mappingObj); err != nil {
		if err == patchain.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "mapping with name=`"+mappingName+"` does not exist", nil)
		}
		return nil, common.ServerError
	}

	var mapping map[string]string
	if err := util.FromJSON([]byte(mappingObj.Value), &mapping); err != nil {
		return nil, common.NewSingleAPIErr(500, "", "", "failed to parse mapping value", nil)
	}

	return mapping, nil
}

// CreateObjects creates one or more objects using a session.
// If session id is provide, it checks whether the session exists locally or
// forwards the request to the remote session host.
// Supports the use of mapping to unmap custom fields in the objects to be created.
func (s *RPC) CreateObjects(ctx context.Context, req *proto_rpc.CreateObjectsMsg) (*proto_rpc.MultiObjectResponse, error) {

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

	var mapping map[string]string

	// if mapping name is provided, get the mapping and unmap the objects
	if len(req.Mapping) > 0 {

		// Fetch mapping
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

	var finalObjs []*tables.Object

	// set developer as the creator id and create final patchain objects
	for _, obj := range objs {
		obj["creator_id"] = developerID
		var o tables.Object
		util.CopyToStruct(&o, obj)
		finalObjs = append(finalObjs, &o)
	}

	if err := session.SendPutOpWithSession(s.dbSession, sid, finalObjs); err != nil {
		logRPC.Errorf("%+v", err)
		s.dbSession.RollbackEnd(sid)
		return nil, common.NewSingleAPIErr(500, common.CodePutError, "objects", err.Error(), nil)
	}

	resp, _ := NewMultiObjectResponse("object", finalObjs)
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
	sid := makeDBSessionID(developerID, sessionID)
	localOnly := util.FromIncomingMD(ctx, "local-only") == "true"

	// use session if one is provided
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
		if err := session.SendQueryOpWithSession(s.dbSession, sid, `{ "id": "`+req.Owner+`" }`, 1, "", &tables.Object{}); err != nil {
			if err == patchain.ErrNotFound {
				return nil, common.NewSingleAPIErr(404, "", "", "owner not found", nil)
			}
			logRPC.Errorf("%+v", err)
			return nil, common.ServerError
		}
	}

	// if creator is set and not the same as the developer id, check if it exists
	if len(req.Creator) > 0 && req.Creator != developerID {
		if err := session.SendQueryOpWithSession(s.dbSession, sid, `{ "id": "`+req.Creator+`" }`, 1, "", &tables.Object{}); err != nil {
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
	if err := session.SendQueryOpWithSession(s.dbSession, sid, string(util.MustStringify(queryAsMap)), int(req.Limit), orderByToString(req.Order), &objs); err != nil {
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	resp, _ := NewMultiObjectResponse("object", objs)
	return resp, nil
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
		if err := session.SendQueryOpWithSession(s.dbSession, sid, `{ "id": "`+req.Owner+`" }`, 1, "", &tables.Object{}); err != nil {
			if err == patchain.ErrNotFound {
				return nil, common.NewSingleAPIErr(404, "", "", "owner not found", nil)
			}
			logRPC.Errorf("%+v", err)
			return nil, common.ServerError
		}
	}

	// if creator is set and not the same as the developer id, check if it exists
	if len(req.Creator) > 0 && req.Creator != developerID {
		if err := session.SendQueryOpWithSession(s.dbSession, sid, `{ "id": "`+req.Creator+`" }`, 1, "", &tables.Object{}); err != nil {
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

	var count int64
	if err := session.SendCountOpWithSession(s.dbSession, sid, string(util.MustStringify(queryAsMap)), &count); err != nil {
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return &proto_rpc.ObjectCountResponse{
		Count: count,
	}, nil
}
