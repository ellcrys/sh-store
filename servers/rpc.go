package servers

import (
	"fmt"
	"net"
	"time"

	"github.com/ncodes/bifrost/servers/proto_rpc"
	"github.com/ncodes/cocoon/core/config"
	"google.golang.org/grpc"
)

var log = config.MakeLogger("rpc")

// RPC defines structure for the RPC server
type RPC struct {
	server *grpc.Server
}

// NewRPC creates a new RPC server
func NewRPC() *RPC {
	return &RPC{}
}

// Start starts the server. It will bind to the address provided
func (s *RPC) Start(addr string, startedCB func(rpcServer *RPC)) error {

	lis, err := net.Listen("tcp", fmt.Sprintf("%s", addr))
	if err != nil {
		_, port, _ := net.SplitHostPort(addr)
		log.Infof("failed to listen on port=%s. Err: %s", port, err)
		return err
	}

	time.AfterFunc(1*time.Second, func() {
		log.Infof("Started RPC server @ %s", addr)
		startedCB(s)
	})

	s.server = grpc.NewServer()
	proto_rpc.RegisterAPIServer(s.server, s)
	return s.server.Serve(lis)
}

// Stop stops the RPC server and all other sub routine
func (s *RPC) Stop() error {
	log.Info("Stopping RPC server")
	if s.server != nil {
		s.server.Stop()
	}
	return nil
}
