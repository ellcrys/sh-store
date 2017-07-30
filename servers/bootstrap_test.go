package servers

import (
	"fmt"
	"net"
	"testing"

	"strconv"

	"os"

	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/ellcrys/elldb/servers/oauth"
	"github.com/ellcrys/elldb/session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var err error
var conStr = "postgresql://root@localhost:26257?sslmode=disable"
var conStrWithDB string
var rpcServerHost = "localhost"
var rpcServerPort = 4000

func init() {
	os.Setenv("ELLDB_RPC_ADDR", net.JoinHostPort(rpcServerHost, strconv.Itoa(rpcServerPort)))
	oauth.SigningSecret = "secret"
}

func setup(t *testing.T, f func(rpc, rpc2 *RPC)) {

	var dbName string
	testDB, err := common.OpenDB(conStr)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %s", err))
	}

	dbName, err = common.CreateRandomDB(testDB)
	if err != nil {
		panic(fmt.Errorf("failed to create test database: %s", err))
	}

	conStrWithDB = "postgresql://root@localhost:26257/" + dbName + "?sslmode=disable"

	dbCon, err := gorm.Open("postgres", conStrWithDB)
	if err != nil {
		t.Fatal("failed to connect to test database")
	}

	err = dbCon.CreateTable(&db.Bucket{}, &db.Object{}, &db.Account{}, &db.Mapping{}).Error
	if err != nil {
		t.Fatalf("failed to create database tables. %s", err)
	}

	rpc := NewRPC()
	rpc.db = dbCon

	rpc.sessionReg, err = session.NewConsulRegistry()
	if err != nil {
		t.Fatalf("failed to connect consul registry. %s", err)
	}

	rpc.dbSession = session.NewSession(rpc.sessionReg)
	rpc.dbSession.SetDB(rpc.db)

	lis, err := net.Listen("tcp", net.JoinHostPort(rpcServerHost, strconv.Itoa(rpcServerPort)))
	if err != nil {
		t.Fatalf("failed to listen on rpc addr %s:%d", rpcServerHost, rpcServerPort)
	}

	rpc2 := NewRPC()
	rpc2.db = dbCon
	rpc2.sessionReg = rpc.sessionReg
	rpc2.dbSession = session.NewSession(rpc.sessionReg)
	rpc2.dbSession.SetDB(rpc2.db)
	go rpc2.serve(lis)

	f(rpc, rpc2)
	rpc2.Stop()

	common.DropDB(testDB, dbName)
}