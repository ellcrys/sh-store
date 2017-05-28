package servers

import (
	"github.com/ncodes/safehold/db"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"golang.org/x/net/context"
)

// CreateDBSession creates a new session and returns the session ID
func (s *RPC) CreateDBSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {
	developerID := ctx.Value(CtxIdentity).(string)

	sessionID, err := s.dbSession.CreateSession(req.ID, developerID)
	if err != nil {
		logRPC.Error("%+v", err)
		return nil, common.NewSingleAPIErr(500, "", "", "session not created", nil)
	}

	return &proto_rpc.DBSession{
		ID: sessionID,
	}, nil
}

// GetDBSession gets a database session.
func (s *RPC) GetDBSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {

	if req.ID == "" {
		return nil, common.NewSingleAPIErr(400, "", "", "session id is required", nil)
	}

	// check if session exists locally
	if s.dbSession.HasSession(req.ID) {
		return req, nil
	}

	// check session registry
	_, err := s.sessionReg.Get(req.ID)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
	}

	return req, nil
}
