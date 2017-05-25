package servers

import (
	"github.com/ncodes/bifrost/servers/proto_rpc"
	"golang.org/x/net/context"
)

// CreateDBSession creates a new session and returns the session ID
func (s *RPC) CreateDBSession(ctx context.Context, req *proto_rpc.DBSession) (*proto_rpc.DBSession, error) {
	return nil, nil
}
