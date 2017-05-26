package servers

import (
	"fmt"
	"net"
	"time"

	"github.com/ellcrys/util"
	"github.com/labstack/gommon/log"
	"github.com/ncodes/cocoon/core/config"
	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/patchain/object"
	"github.com/ncodes/safehold/db"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	logRPC = config.MakeLogger("rpc")

	// database name of the global patchain database
	databaseName = util.Env("SH_PATCHAIN_DB_NAME", "safehold_dev")

	// partitionChainConStr connection string
	partitionChainConStr = util.Env("SH_PATCHAIN_CONSTR", "postgresql://root@localhost:26257/"+databaseName+"?sslmode=disable")

	// systemEmail is the email used for the system identity
	systemEmail = util.Env("SH_SYS_EMAIL", "core@safehold.io")

	// system ledger name
	systemLedgerName = util.Env("SH_SYS_LEDGER", "core")

	// system default number of partitions
	defNumPartitions int64 = 2
)

// RPC defines structure for the RPC server
type RPC struct {
	server    *grpc.Server
	db        patchain.DB
	dbSession *db.Session
	object    *object.Object
}

// NewRPC creates a new RPC server
func NewRPC(db patchain.DB) *RPC {
	return &RPC{
		db:     db,
		object: object.NewObject(db),
	}
}

// Start starts the server. It will bind to the address provided
func (s *RPC) Start(addr string, startedCB func(rpcServer *RPC)) error {

	lis, err := net.Listen("tcp", fmt.Sprintf("%s", addr))
	if err != nil {
		_, port, _ := net.SplitHostPort(addr)
		logRPC.Infof("failed to listen on port=%s. Err: %s", port, err)
		return err
	}

	time.AfterFunc(1*time.Second, func() {
		logRPC.Infof("Started RPC server @ %s", addr)

		// connect database
		if err := s.db.Connect(1000, 10); err != nil {
			log.Errorf("%v", err)
			s.Stop()
			return
		}

		// create tables if necessary
		if err := s.db.CreateTables(); err != nil {
			log.Errorf("%v", err)
			s.Stop()
			return
		}

		// create system resources if necessary
		if err := s.createSystemResources(); err != nil {
			log.Errorf("%v", err)
			s.Stop()
			return
		}

		// create db session manager
		s.dbSession = db.NewSession(partitionChainConStr)
		s.dbSession.SetDB(s.db)

		startedCB(s)
	})

	s.server = grpc.NewServer(grpc.UnaryInterceptor(s.Interceptors()))
	proto_rpc.RegisterAPIServer(s.server, s)
	return s.server.Serve(lis)
}

// Stop stops the RPC server and all other sub routines
func (s *RPC) Stop() error {
	logRPC.Info("Stopping RPC server")
	if s.server != nil {
		s.server.Stop()
	}
	if s.dbSession != nil {
		s.dbSession.Stop()
	}
	return nil
}

// createSystemResources creates system identity, ledger and partitions
func (s *RPC) createSystemResources() error {

	identity := &tables.Object{
		Key:       object.MakeIdentityKey(systemEmail),
		Protected: true,
	}

	obj := object.NewObject(s.db)

	// Create system identity if it doesn't already exists. Set previous hash
	// to the objects hash since system identity is not associated to a partition
	identity.PrevHash = identity.ComputeHash().Hash
	if err := obj.CreateOnce(identity); err != nil {
		return errors.Wrap(err, "failed to create system identity")
	}

	// Create system partitions only if the number of system partitions is lower than the required number
	var numExistingPartitions int64
	if err := s.db.Count(&tables.Object{OwnerID: identity.ID, QueryParams: patchain.KeyStartsWith(object.PartitionPrefix)}, &numExistingPartitions); err != nil {
		return errors.Wrap(err, "failed to determine number of existing system partitions")
	}

	if numExistingPartitions < defNumPartitions {
		if _, err := obj.MustCreatePartitions(defNumPartitions-numExistingPartitions, identity.ID, identity.ID); err != nil {
			return errors.Wrap(err, "failed to create system ledger")
		}
	}

	return nil
}
