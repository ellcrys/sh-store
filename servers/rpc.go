package servers

import (
	"fmt"
	"net"
	"time"

	"github.com/ellcrys/util"
	"github.com/ncodes/bifrost/servers/proto_rpc"
	"github.com/ncodes/bifrost/session/db"
	"github.com/ncodes/cocoon/core/config"
	"github.com/ncodes/patchain/cockroach"
	"google.golang.org/grpc"
)

var (
	log = config.MakeLogger("rpc")

	// database name of the global patchain database
	databaseName = util.Env("BF_PATCHAIN_DB_NAME", "safehold_dev")

	// partitionChainConStr connection string
	partitionChainConStr = util.Env("BF_PATCHAIN_CONSTR", "postgresql://root@localhost:26257/"+databaseName+"?sslmode=disable")
)

// RPC defines structure for the RPC server
type RPC struct {
	server    *grpc.Server
	dbSession *db.Session
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

		// create db session manager
		dbSession := db.NewSession(partitionChainConStr)
		cdb := cockroach.NewDB()
		cdb.ConnectionString = partitionChainConStr
		if err = dbSession.Connect(cdb); err != nil {
			log.Fatalf("failed to create database session manager. Err: %s", err)
		}

		startedCB(s)
	})

	s.server = grpc.NewServer()
	proto_rpc.RegisterAPIServer(s.server, s)
	return s.server.Serve(lis)
}

// Stop stops the RPC server and all other sub routines
func (s *RPC) Stop() error {
	log.Info("Stopping RPC server")
	if s.server != nil {
		s.server.Stop()
	}
	if s.dbSession != nil {
		s.dbSession.Stop()
	}
	return nil
}
