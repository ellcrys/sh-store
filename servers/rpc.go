package servers

import (
	"net"
	"time"

	"github.com/ellcrys/cocoon/core/config"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/elldb/session"
	"github.com/ellcrys/util"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
)

var (
	logRPC = config.MakeLogger("rpc")

	// database name of the global patchain database
	databaseName = util.Env("EDB_DB_NAME", "elldb_dev")

	// sqlDBConStr connection string
	sqlDBConStr = util.Env("EDB_SQL_CONSTR", "postgresql://root@localhost:26257/"+databaseName+"?sslmode=disable")
)

// RPC defines structure for the RPC server
type RPC struct {
	server     *grpc.Server
	dbSession  *session.Session
	sessionReg session.SessionRegistry
	db         *gorm.DB
}

// NewRPC creates a new RPC server
func NewRPC() *RPC {
	return &RPC{}
}

// GetDB returns the db
func (s *RPC) GetDB() *gorm.DB {
	return s.db
}

// Start starts the server. It will bind to the address provided
func (s *RPC) Start(addr string, startedCB func(rpcServer *RPC)) error {

	var err error

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logRPC.Infof("failed to listen on port=%s. Err: %s", addr, err)
		return err
	}

	time.AfterFunc(1*time.Second, func() {

		logRPC.Infof("Started server @ %s", addr)

		// connect database
		s.db, err = db.Connect(sqlDBConStr)
		if err != nil {
			logRPC.Errorf("Database Connection Error: %s", err)
			s.Stop()
			return
		}

		// initialize and seed database
		if err = db.Initialize(s.db); err != nil {
			logRPC.Errorf("%v", err)
			return
		}

		// connect to session registry
		s.sessionReg, err = session.NewConsulRegistry()
		if err != nil {
			logRPC.Errorf("%v", err)
			s.Stop()
			return
		}

		// create db session manager
		s.dbSession = session.NewSession(s.sessionReg)
		s.dbSession.SetDB(s.db)

		startedCB(s)
	})

	return s.serve(lis)
}

func (s *RPC) serve(lis net.Listener) error {
	s.server = grpc.NewServer(grpc.UnaryInterceptor(s.Interceptors()))
	proto_rpc.RegisterAPIServer(s.server, s)
	return s.server.Serve(lis)
}

// Stop stops the RPC server and all other sub routines
func (s *RPC) Stop() error {
	logRPC.Info("Stopping server")
	if s.server != nil {
		s.server.Stop()
	}
	return nil
}
