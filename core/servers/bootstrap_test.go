package servers

import (
	"fmt"
	"net"
	"testing"

	"strconv"

	"os"

	"github.com/ellcrys/elldb/core/servers/common"
	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/ellcrys/elldb/core/servers/oauth"
	"github.com/ellcrys/elldb/core/session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/phayes/freeport"
)

var err error
var conStr = "postgresql://postgres@localhost:5432?sslmode=disable"
var rpcServerHost = "localhost"
var rpcServerPort = 4000

func setup(t *testing.T, f func(rpc, rpc2 *RPC)) {

	oauth.SigningSecret = "secret"
	rpcServerPort = freeport.GetPort()
	os.Setenv("ELLDB_RPC_ADDR", net.JoinHostPort(rpcServerHost, strconv.Itoa(rpcServerPort)))

	dbName, err := common.CreateRandomDB()
	if err != nil {
		panic(fmt.Errorf("failed to create test database: %s", err))
	}

	conStrWithDB := "postgresql://postgres@localhost:5432/" + dbName + "?sslmode=disable"
	dbCon, err := gorm.Open("postgres", conStrWithDB)
	if err != nil {
		t.Log(err)
		t.Fatal("failed to connect to test database")
	}

	db.ApplyCallbacks(dbCon)

	defer common.DropDB(dbName)
	defer dbCon.Close()

	err = dbCon.CreateTable(&db.Bucket{}, &db.Object{}, &db.Account{}, &db.Mapping{}, &db.Contract{}).Error
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
		t.Fatalf("failed to listen on rpc addr %s:%d: %s", rpcServerHost, rpcServerPort, err)
	}

	rpc2 := NewRPC()
	rpc2.db = dbCon
	rpc2.sessionReg = rpc.sessionReg
	rpc2.dbSession = session.NewSession(rpc.sessionReg)
	rpc2.dbSession.SetDB(rpc2.db)
	go rpc2.serve(lis)

	f(rpc, rpc2)
	rpc2.Stop()
	rpc.Stop()
}
