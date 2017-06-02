package servers

import (
	"fmt"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"net"

	"github.com/ellcrys/util"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/ncodes/safehold/session"
	"golang.org/x/net/context"
)

// makeDBSessionID returns a the standard db session name
func makeDBSessionID(identityID, sid string) string {
	return fmt.Sprintf("%s.%s", identityID, sid)
}

// CreateDBSession creates a new session and returns the session ID
func (s *RPC) CreateDBSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {

	// allocate id if not provided
	if req.ID == "" {
		req.ID = util.UUID4()
	}

	developerID := ctx.Value(CtxIdentity).(string)

	sessionID := makeDBSessionID(developerID, req.ID)
	sessionID, err := s.dbSession.CreateSession(sessionID, developerID)
	if err != nil {
		logRPC.Error("%+v", err)
		return nil, common.NewSingleAPIErr(500, "", "", "session not created", nil)
	}

	return &proto_rpc.DBSession{
		ID: req.ID,
	}, nil
}

// GetDBSession gets a database session.
func (s *RPC) GetDBSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {
	developerID := ctx.Value(CtxIdentity).(string)

	if req.ID == "" {
		return nil, common.NewSingleAPIErr(400, "", "", "session id is required", nil)
	}

	sessionID := makeDBSessionID(developerID, req.ID)

	// check if session exists locally
	if s.dbSession.HasSession(sessionID) {
		return req, nil
	}

	// check session registry
	_, err := s.sessionReg.Get(sessionID)
	if err != nil {
		if err == session.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
	}

	return req, nil
}

// DeleteDBSession deletes a existing database session
func (s *RPC) DeleteDBSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {

	developerID := ctx.Value(CtxIdentity).(string)
	authorization := util.FromIncomingMD(ctx, "authorization")
	localOnly := util.FromIncomingMD(ctx, "local-only") == "true"

	if req.ID == "" {
		return nil, common.NewSingleAPIErr(400, "", "", "session id is required", nil)
	}

	sessionID := makeDBSessionID(developerID, req.ID)

	// check if session exists locally, if so, delete immediately
	if s.dbSession.HasSession(sessionID) {
		s.dbSession.End(sessionID)
		return req, nil
	}

	// if localOnly is true, return error
	if localOnly {
		return nil, fmt.Errorf("session not found")
	}

	// get session from the session registry.
	// if we find it, we need to call the session's host server to delete.
	// if not found, return `not found` error
	item, err := s.sessionReg.Get(sessionID)
	if err != nil {
		if err == session.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
		return nil, common.ServerError
	}

	sessionHostAddr := net.JoinHostPort(item.Address, strconv.Itoa(item.Port))
	client, err := grpc.Dial(sessionHostAddr, grpc.WithInsecure())
	if err != nil {
		return nil, common.ServerError
	}
	defer client.Close()

	// make call to the session host server.
	// include the auth token of the current request
	// and set local-only to force the RPC method
	// to only check local/in-memory session cache
	server := proto_rpc.NewAPIClient(client)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", authorization, "local-only", "true"))
	resp, err := server.DeleteDBSession(ctx, req)
	if err != nil {
		if grpc.ErrorDesc(err) == "session not found" {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
		logRPC.Errorf("%+v", err)
		return nil, common.ServerError
	}

	return resp, nil
}
