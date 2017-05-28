package servers

import (
	"fmt"

	"github.com/ncodes/safehold/db"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"golang.org/x/net/context"
)

// makeDBSessionID returns a the standard db session name
func makeDBSessionID(sid, identityID string) string {
	return fmt.Sprintf("%s.%s", identityID, sid)
}

// CreateDBSession creates a new session and returns the session ID
func (s *RPC) CreateDBSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {
	developerID := ctx.Value(CtxIdentity).(string)

	sessionID := makeDBSessionID(developerID, req.ID)
	sessionID, err := s.dbSession.CreateSession(sessionID, developerID)
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
		if err == db.ErrNotFound {
			return nil, common.NewSingleAPIErr(404, "", "", "session not found", nil)
		}
	}

	return req, nil
}
