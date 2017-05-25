package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"net"

	"io/ioutil"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/ncodes/cocoon/core/config"
	db "github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/safehold/core/servers/common"
	"github.com/ncodes/safehold/core/servers/oauth"
	"github.com/ncodes/safehold/core/servers/rpc/proto_rpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var log = config.MakeLogger("http")

// Server defines a structure for an HTTP server
// that provides REST API services.
type Server struct {
	rpcServerAddr string
	oauth         *oauth.OAuth
	db            db.DB
}

// NewServer creates an new http server instance
func NewServer(db db.DB, rpcServerAddr string) *Server {
	return &Server{
		rpcServerAddr: rpcServerAddr,
		db:            db,
		oauth:         oauth.NewOAuth(db),
	}
}

// getRouter returns the router
func (s *Server) getRouter() *mux.Router {
	r := mux.NewRouter()

	// oauth endpoints
	r.HandleFunc("/token", common.EasyHandle(http.MethodPost, s.oauth.GetToken))

	// v1 endpoints
	g := r.PathPrefix("/v1").Subrouter()
	g.HandleFunc("/identities", common.EasyHandle(http.MethodPost, s.createIdentity))
	g.HandleFunc("/objects", common.EasyHandle(http.MethodPost, s.createObjects))
	g.HandleFunc("/objects/get", common.EasyHandle(http.MethodPost, s.query))
	return r
}

// dialRPC creates a client connection to the RPC API
// and passes it to the callback function. It will close the connection
// after the callback function has finished
func (s *Server) dialRPC(cb func(client proto_rpc.APIClient) error) error {
	conn, err := grpc.Dial(s.rpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to dial RPC API")
	}
	defer conn.Close()
	return cb(proto_rpc.NewAPIClient(conn))
}

// Start starts the http server. Passes true to the startedCh channel
// when started
func (s *Server) Start(addr string, startedCB func(*Server)) error {

	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Errorf("%+v", errors.Wrap(err, "failed to parse address"))
		return err
	}

	// attempt to connect to the RPC API
	if err := s.dialRPC(func(client proto_rpc.APIClient) error {
		return nil
	}); err != nil {
		log.Errorf("%+v", err)
		return err
	}

	time.AfterFunc(2*time.Second, func() {
		log.Infof("Started HTTP API server @ :%s", port)
		if startedCB != nil {
			startedCB(s)
		}
	})

	return http.ListenAndServe(addr, s.getRouter())
}

// createIdentity creates an identity
func (s *Server) createIdentity(w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var err error
	var body proto_rpc.CreateIdentityMsg
	var resp *proto_rpc.ObjectResponse

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		resp, err = client.CreateIdentity(context.Background(), &body)
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}

// createObjects creates an arbitrary object.
// Requires an app token included in the authorization header
func (s *Server) createObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var err error
	var body []*proto_rpc.Object
	var resp *proto_rpc.MultiObjectResponse

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md metadata.MD
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{Objects: body})
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}

// query performs query operations
func (s *Server) query(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	qm, err := s.db.NewQuery()
	if err != nil {
		return err, 500
	}
	qm.SetTable(tables.Object{}, true)
	b, _ := ioutil.ReadAll(r.Body)
	err = qm.Parse(string(b))
	if err != nil {
		return err, 400
	}
	var out []tables.Object
	if err = qm.Find(&out); err != nil {
		return fmt.Errorf("all: %s", err), 500
	}
	return out, 200
}
