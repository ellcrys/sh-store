package servers

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/imdario/mergo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/elldb/session"
	"github.com/ellcrys/util"
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
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
	autoFinish := false
	sessionID := util.FromIncomingMD(ctx, "session_id")
	developerID := ctx.Value(CtxAccount).(string)
	authorization := util.FromIncomingMD(ctx, "authorization")
	checkLocalOnly := util.FromIncomingMD(ctx, "check-local-only") == "true"

	if len(req.Bucket) == 0 {
		return nil, common.NewSingleAPIErr(400, "", "bucket", "bucket name is required", nil)
	}

	// check if bucket exists
	_, err = s.getBucket(req.Bucket)
	if err != nil {
		return nil, err
	}

	// Handle session
	if len(sessionID) > 0 {

		if !s.dbSession.HasSession(sessionID) {

			// abort further operations
			if checkLocalOnly {
				return nil, fmt.Errorf("session not found")
			}

			var resp *proto_rpc.GetObjectsResponse
			err = s.getRemoteConnection(ctx, developerID, authorization, sessionID, func(ctx context.Context, remote proto_rpc.APIClient) error {
				resp, err = remote.CreateObjects(ctx, req)
				if err != nil {
					if grpc.ErrorDesc(err) == "session not found" {
						return common.NewSingleAPIErr(404, "", "", "session not found", nil)
					}
					return common.ServerError
				}
				return nil
			})

			return resp, err
		}

		// check if session is owned by the developer, if not, return permission error
		if sessionAgent := s.dbSession.GetAgent(sessionID); sessionAgent != nil {
			if sessionAgent.OwnerID != developerID {
				return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
			}
		}

	} else { // session id not provided, create local, unregistered session
		sessionID = util.UUID4()
		s.dbSession.CreateUnregisteredSession(sessionID, developerID)
		autoFinish = true
		defer s.dbSession.CommitEnd(sessionID)
	}

	// coerce objects to be created into slice of maps
	var objs []map[string]interface{}
	if err := util.FromJSON2(req.Objects, &objs); err != nil {
		if autoFinish {
			s.dbSession.RollbackEnd(sessionID)
		}
		return nil, common.NewSingleAPIErr(400, "", "", "failed to parse objects", nil)
	}

	// add developer as the creator, owner if 'owner_id' is unset and set created_at
	for _, obj := range objs {
		obj["created_at"] = time.Now().UnixNano()
		obj["creator_id"] = developerID
		if obj["owner_id"] == nil {
			obj["owner_id"] = developerID
		}
	}

	var mapping map[string]string

	// if a mapping is provided, get the mapping and unmap the objects with it
	if len(req.Mapping) > 0 {

		var m db.Mapping
		if err = s.db.Where("name = ? AND account = ?", req.Mapping, developerID).First(&m).Error; err != nil {
			if autoFinish {
				s.dbSession.RollbackEnd(sessionID)
			}
			if err == gorm.ErrRecordNotFound {
				return nil, common.NewSingleAPIErr(400, "", "", "mapping not found", nil)
			}
			return nil, common.ServerError
		}

		util.FromJSON([]byte(m.Mapping), &mapping)

		// replace custom map fields in objects with the original/standard object column names
		if err := common.UnMapFields(mapping, objs); err != nil {
			if autoFinish {
				s.dbSession.RollbackEnd(sessionID)
			}
			logRPC.Errorf("%+v", errors.Wrap(err, "failed to unmap object"))
			return nil, common.ServerError
		}
	}

	// validate the objects
	vErrs, err := s.validateObjects(objs, mapping)
	if err != nil {
		if autoFinish {
			s.dbSession.RollbackEnd(sessionID)
		}
		logRPC.Errorf("%+v", errors.Wrap(err, "failed to validate objects"))
		return nil, common.ServerError
	} else if len(vErrs) > 0 {
		if autoFinish {
			s.dbSession.RollbackEnd(sessionID)
		}
		return nil, common.NewMultiAPIErr(400, "validation errors", vErrs)
	}

	// ensure caller has permission to PUT objects on behalf of the object owner
	// TODO: check oauth permission
	if developerID != objs[0]["owner_id"] {
		if autoFinish {
			s.dbSession.RollbackEnd(sessionID)
		}
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you are not authorized to create objects for the owner", nil)
	}

	// create db objects to be inserted
	var objectsToCreate []*db.Object
	for _, obj := range objs {
		var o db.Object
		util.CopyToStruct(&o, obj)
		objectsToCreate = append(objectsToCreate, &o)
	}

	// using the session, insert objects
	if err := session.SendPutOpWithSession(s.dbSession, sessionID, req.Bucket, objectsToCreate); err != nil {
		if autoFinish {
			s.dbSession.RollbackEnd(sessionID)
		}
		logRPC.Errorf("%+v", err)
		return nil, common.NewSingleAPIErr(500, common.CodePutError, "objects", err.Error(), nil)
	}

	// reapply custom mapped fields to inserted objects so the caller
	objs = nil
	for _, o := range objectsToCreate {
		obj := structs.New(o).Map()
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

// GetObjects fetches objects belonging to an account
func (s *RPC) GetObjects(ctx context.Context, req *proto_rpc.GetObjectMsg) (*proto_rpc.GetObjectsResponse, error) {

	var err error
	authorization := util.FromIncomingMD(ctx, "authorization")
	developerID := ctx.Value(CtxAccount).(string)
	sessionID := util.FromIncomingMD(ctx, "session_id")
	checkLocalOnly := util.FromIncomingMD(ctx, "check-local-only") == "true"

	// bucket is required
	if len(req.Bucket) == 0 {
		return nil, common.NewSingleAPIErr(400, "", "bucket", "bucket name is required", nil)
	}

	// check if bucket exists
	_, err = s.getBucket(req.Bucket)
	if err != nil {
		return nil, err
	}

	// Handle session if provided
	if len(sessionID) > 0 {

		// if session id does not exist locally, find it in the session registry and forward
		// request to the session host
		if !s.dbSession.HasSession(sessionID) {

			// abort further operations
			if checkLocalOnly {
				return nil, fmt.Errorf("session not found")
			}

			var err error
			var resp *proto_rpc.GetObjectsResponse
			err = s.getRemoteConnection(ctx, developerID, authorization, sessionID, func(ctx context.Context, remote proto_rpc.APIClient) error {
				resp, err = remote.GetObjects(ctx, req)
				if err != nil {
					if grpc.ErrorDesc(err) == "session not found" {
						return common.NewSingleAPIErr(404, "", "", "session not found", nil)
					}
					return common.ServerError
				}
				return nil
			})

			return resp, err
		}

		// check if session is owned by the developer, if not, return permission error
		if sessionAgent := s.dbSession.GetAgent(sessionID); sessionAgent != nil {
			if sessionAgent.OwnerID != developerID {
				return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
			}
		}

	} else { // session id not provided, create local, unregistered session
		sessionID = util.UUID4()
		s.dbSession.CreateUnregisteredSession(sessionID, developerID)
		defer s.dbSession.CommitEnd(sessionID)
	}

	// set owner as the developer if not set
	if len(req.Owner) == 0 {
		req.Owner = developerID
	} else {
		if _, err = s.getAccount(req.Owner); err != nil {
			if strings.Contains(err.Error(), "account not found") {
				return nil, common.NewSingleAPIErr(404, "", "owner", "owner not found", nil)
			}
		}
	}

	// developer is not the owner, this action requires permission
	// TODO: ensure auth token must be a user token from the owner
	// and the token authorizes access to the object created by the creator for this developer
	if developerID != req.Owner {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you are not authorized to access objects belonging to the owner", nil)
	}

	// include owner, creator and bucket filters
	var query map[string]interface{}
	util.FromJSON(req.Query, &query)
	mergo.Merge(&query, map[string]interface{}{
		"owner_id": req.Owner,
		"bucket":   req.Bucket,
	})

	var mapping map[string]string

	// if mapping is provided in request, fetch it
	if len(req.Mapping) > 0 {
		m, err := s.getMapping(req.Mapping, developerID)
		if err != nil {
			return nil, err
		}
		util.FromJSON([]byte(m.Mapping), &mapping)
	}

	// using the mapping, convert custom mapped fields to the original object column names
	common.UnMapFields(mapping, query)

	var fetchedObjs []*db.Object
	if err = session.SendQueryOpWithSession(s.dbSession, sessionID, req.Bucket, string(util.MustStringify(query)), nil, int(req.Limit), orderByToString(req.Order), &fetchedObjs); err != nil {
		if strings.Contains(err.Error(), "parser") {
			msg := strings.SplitN(err.Error(), ":", 2)[1]
			return nil, common.NewSingleAPIErr(400, common.CodeInvalidParam, "query", msg, nil)
		}
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

// Create a remote RPC connection to the host of a session
func (s *RPC) getRemoteConnection(ctx context.Context, developerID, authorization, sessionID string, cb func(ctx context.Context, remoteServer proto_rpc.APIClient) error) error {

	// find session in the registry
	sessionItem, err := s.sessionReg.Get(sessionID)
	if err != nil {
		if err == session.ErrNotFound {
			return common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
		return common.ServerError
	}

	// check if session is owned by the developer, if not, return permission error
	if sessionItem.Meta["account"] != developerID {
		return common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
	}

	sessionHostAddr := net.JoinHostPort(sessionItem.Address, strconv.Itoa(sessionItem.Port))
	client, err := grpc.Dial(sessionHostAddr, grpc.WithInsecure())
	if err != nil {
		return common.ServerError
	}
	defer client.Close()

	// make call to the session host server.
	// include the session_id, auth token of the current request
	// and set check-local-only to force the RPC method
	// to only perform local object creation
	server := proto_rpc.NewAPIClient(client)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("session_id", sessionID, "authorization", authorization, "check-local-only", "true"))
	return cb(ctx, server)
}

// CountObjects counts the number of objects that match a query
func (s *RPC) CountObjects(ctx context.Context, req *proto_rpc.GetObjectMsg) (*proto_rpc.ObjectCountResponse, error) {

	var err error
	authorization := util.FromIncomingMD(ctx, "authorization")
	developerID := ctx.Value(CtxAccount).(string)
	sessionID := util.FromIncomingMD(ctx, "session_id")
	checkLocalOnly := util.FromIncomingMD(ctx, "check-local-only") == "true"

	// bucket is required
	if len(req.Bucket) == 0 {
		return nil, common.NewSingleAPIErr(400, "", "bucket", "bucket name is required", nil)
	}

	// check if bucket exists
	_, err = s.getBucket(req.Bucket)
	if err != nil {
		return nil, err
	}

	// check if session exist in the in-memory session cache,
	// use it else check if it exist on the session registry. If it does,
	// forward the request to the associated host
	if sessionID != "" {
		if !s.dbSession.HasSession(sessionID) {

			if checkLocalOnly {
				return nil, fmt.Errorf("session not found")
			}

			var err error
			var resp *proto_rpc.ObjectCountResponse
			err = s.getRemoteConnection(ctx, developerID, authorization, sessionID, func(ctx context.Context, remote proto_rpc.APIClient) error {
				resp, err = remote.CountObjects(ctx, req)
				if err != nil {
					if grpc.ErrorDesc(err) == "session not found" {
						return common.NewSingleAPIErr(404, "", "", "session not found", nil)
					}
					return common.ServerError
				}
				return nil
			})

			return resp, err
		}

		// check if session is owned by the developer, if not, return permission error
		if sessionAgent := s.dbSession.GetAgent(sessionID); sessionAgent != nil {
			if sessionAgent.OwnerID != developerID {
				return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
			}
		}

	} else {
		sessionID = util.UUID4()
		s.dbSession.CreateUnregisteredSession(sessionID, developerID)
		defer s.dbSession.CommitEnd(sessionID)
	}

	if len(req.Owner) == 0 {
		req.Owner = developerID
	} else {
		if _, err = s.getAccount(req.Owner); err != nil {
			if strings.Contains(err.Error(), "account not found") {
				return nil, common.NewSingleAPIErr(404, "", "owner", "owner not found", nil)
			}
		}
	}

	// developer is not the owner, this action requires permission
	// TODO: ensure auth token must be a user token from the owner
	// and the token authorizes access to the object created by the creator for this developer
	if developerID != req.Owner {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you are not authorized to access objects belonging to the owner", nil)
	}

	// include owner and bucket filters
	var query map[string]interface{}
	util.FromJSON(req.Query, &query)
	mergo.Merge(&query, map[string]interface{}{
		"owner_id": req.Owner,
		"bucket":   req.Bucket,
	})

	var count int64
	if err := session.SendCountOpWithSession(s.dbSession, sessionID, req.Bucket, string(util.MustStringify(query)), nil, &count); err != nil {
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return &proto_rpc.ObjectCountResponse{
		Count: count,
	}, nil
}

// UpdateObjects updates objects of a mutable bucket
func (s *RPC) UpdateObjects(ctx context.Context, req *proto_rpc.UpdateObjectsMsg) (*proto_rpc.AffectedResponse, error) {

	var err error
	authorization := util.FromIncomingMD(ctx, "authorization")
	developerID := ctx.Value(CtxAccount).(string)
	sessionID := util.FromIncomingMD(ctx, "session_id")
	checkLocalOnly := util.FromIncomingMD(ctx, "check-local-only") == "true"

	// bucket is required
	if len(req.Bucket) == 0 {
		return nil, common.NewSingleAPIErr(400, "", "bucket", "bucket name is required", nil)
	}

	// check if bucket exists
	bucket, err := s.getBucket(req.Bucket)
	if err != nil {
		return nil, err
	}

	// only objects in mutable buckets can be updated
	if bucket.Immutable {
		return nil, common.NewSingleAPIErr(400, "", "bucket", "bucket is not mutable", nil)
	}

	// check if session exist in the in-memory session cache,
	// use it else check if it exist on the session registry. If it does,
	// forward the request to the associated host
	if sessionID != "" {
		if !s.dbSession.HasSession(sessionID) {

			if checkLocalOnly {
				return nil, fmt.Errorf("session not found")
			}

			var err error
			var resp *proto_rpc.AffectedResponse
			err = s.getRemoteConnection(ctx, developerID, authorization, sessionID, func(ctx context.Context, remote proto_rpc.APIClient) error {
				resp, err = remote.UpdateObjects(ctx, req)
				if err != nil {
					if grpc.ErrorDesc(err) == "session not found" {
						return common.NewSingleAPIErr(404, "", "", "session not found", nil)
					}
					return common.ServerError
				}
				return nil
			})

			return resp, err
		}

		// check if session is owned by the developer, if not, return permission error
		if sessionAgent := s.dbSession.GetAgent(sessionID); sessionAgent != nil {
			if sessionAgent.OwnerID != developerID {
				return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
			}
		}

	} else {
		sessionID = util.UUID4()
		s.dbSession.CreateUnregisteredSession(sessionID, developerID)
		defer s.dbSession.CommitEnd(sessionID)
	}

	// set owner as the developer if not set
	if len(req.Owner) == 0 {
		req.Owner = developerID
	} else {
		if _, err = s.getAccount(req.Owner); err != nil {
			if strings.Contains(err.Error(), "account not found") {
				return nil, common.NewSingleAPIErr(404, "", "owner", "owner not found", nil)
			}
		}
	}

	// developer is not the owner, this action requires permission
	// TODO: ensure auth token must be a user token from the owner
	// and the token authorizes access to the object created by the creator for this developer
	if developerID != req.Owner {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you are not authorized to access objects belonging to the owner", nil)
	}

	// include owner, creator and bucket filters
	var query map[string]interface{}
	util.FromJSON(req.Query, &query)
	mergo.Merge(&query, map[string]interface{}{
		"owner_id": req.Owner,
		"bucket":   req.Bucket,
	})

	var mapping map[string]string

	// if mapping is provided in request, fetch it
	if len(req.Mapping) > 0 {
		m, err := s.getMapping(req.Mapping, developerID)
		if err != nil {
			return nil, err
		}
		util.FromJSON([]byte(m.Mapping), &mapping)
	}

	var update map[string]interface{}
	util.FromJSON(req.Update, &update)

	// using the mapping, convert custom mapped fields of
	// query and update to their original object column names
	common.UnMapFields(mapping, query)
	common.UnMapFields(mapping, update)

	// remove fields that cannot be updated
	var blacklistFields = []string{"id", "sn", "bucket", "creator_id", "owner_id", "created_at", "updated_at"}
	for _, f := range blacklistFields {
		delete(update, f)
	}

	numAffected, err := session.SendUpdateOpWithSession(s.dbSession, sessionID, string(util.MustStringify(query)), nil, update)
	if err != nil {
		if strings.Contains(err.Error(), "parser") {
			msg := strings.SplitN(err.Error(), ":", 2)[1]
			return nil, common.NewSingleAPIErr(400, common.CodeInvalidParam, "query", msg, nil)
		}
		return nil, common.ServerError
	}

	return &proto_rpc.AffectedResponse{
		Affected: numAffected,
	}, nil
}

// DeleteObjects deletes objects from a mutable bucket
func (s *RPC) DeleteObjects(ctx context.Context, req *proto_rpc.DeleteObjectsMsg) (*proto_rpc.AffectedResponse, error) {

	var err error
	authorization := util.FromIncomingMD(ctx, "authorization")
	developerID := ctx.Value(CtxAccount).(string)
	sessionID := util.FromIncomingMD(ctx, "session_id")
	checkLocalOnly := util.FromIncomingMD(ctx, "check-local-only") == "true"

	// bucket is required
	if len(req.Bucket) == 0 {
		return nil, common.NewSingleAPIErr(400, "", "bucket", "bucket name is required", nil)
	}

	// check if bucket exists
	bucket, err := s.getBucket(req.Bucket)
	if err != nil {
		return nil, err
	}

	// only objects in mutable buckets can be deleted
	if bucket.Immutable {
		return nil, common.NewSingleAPIErr(400, "", "bucket", "bucket is not mutable", nil)
	}

	// check if session exist in the in-memory session cache,
	// use it else check if it exist on the session registry. If it does,
	// forward the request to the associated host
	if sessionID != "" {
		if !s.dbSession.HasSession(sessionID) {

			if checkLocalOnly {
				return nil, fmt.Errorf("session not found")
			}

			var err error
			var resp *proto_rpc.AffectedResponse
			err = s.getRemoteConnection(ctx, developerID, authorization, sessionID, func(ctx context.Context, remote proto_rpc.APIClient) error {
				resp, err = remote.DeleteObjects(ctx, req)
				if err != nil {
					if grpc.ErrorDesc(err) == "session not found" {
						return common.NewSingleAPIErr(404, "", "", "session not found", nil)
					}
					return common.ServerError
				}
				return nil
			})

			return resp, err
		}

		// check if session is owned by the developer, if not, return permission error
		if sessionAgent := s.dbSession.GetAgent(sessionID); sessionAgent != nil {
			if sessionAgent.OwnerID != developerID {
				return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you don't have permission to perform this operation", nil)
			}
		}

	} else {
		sessionID = util.UUID4()
		s.dbSession.CreateUnregisteredSession(sessionID, developerID)
		defer s.dbSession.CommitEnd(sessionID)
	}

	// set owner as the developer if not set
	if len(req.Owner) == 0 {
		req.Owner = developerID
	} else {
		if _, err = s.getAccount(req.Owner); err != nil {
			if strings.Contains(err.Error(), "account not found") {
				return nil, common.NewSingleAPIErr(404, "", "owner", "owner not found", nil)
			}
		}
	}

	// developer is not the owner, this action requires permission
	// TODO: ensure auth token must be a user token from the owner
	// and the token authorizes access to the object created by the creator for this developer
	if developerID != req.Owner {
		return nil, common.NewSingleAPIErr(401, "", "", "permission denied: you are not authorized to access objects belonging to the owner", nil)
	}

	// include owner, creator and bucket filters
	var query map[string]interface{}
	util.FromJSON(req.Query, &query)
	mergo.Merge(&query, map[string]interface{}{
		"owner_id": req.Owner,
		"bucket":   req.Bucket,
	})

	var mapping map[string]string

	// if mapping is provided in request, fetch it
	if len(req.Mapping) > 0 {
		m, err := s.getMapping(req.Mapping, developerID)
		if err != nil {
			return nil, err
		}
		util.FromJSON([]byte(m.Mapping), &mapping)
	}

	// using the mapping, convert custom mapped fields of query
	common.UnMapFields(mapping, query)

	numAffected, err := session.SendDeleteOpWithSession(s.dbSession, sessionID, string(util.MustStringify(query)), nil)
	if err != nil {
		if strings.Contains(err.Error(), "parser") {
			msg := strings.SplitN(err.Error(), ":", 2)[1]
			return nil, common.NewSingleAPIErr(400, common.CodeInvalidParam, "query", msg, nil)
		}
		return nil, common.ServerError
	}

	return &proto_rpc.AffectedResponse{
		Affected: numAffected,
	}, nil
}
