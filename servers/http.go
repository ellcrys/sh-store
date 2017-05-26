package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"net"

	"github.com/gorilla/mux"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/ncodes/cocoon/core/config"
	"github.com/ncodes/safehold/core/servers/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var logHTTP = config.MakeLogger("http")

// HTTP defines a structure for an HTTP server
// that provides REST API services.
type HTTP struct {
	rpcServerAddr string
}

// NewHTTP creates an new http server instance
func NewHTTP(rpcServerAddr string) *HTTP {
	return &HTTP{
		rpcServerAddr: rpcServerAddr,
	}
}

// getRouter returns the router
func (s *HTTP) getRouter() *mux.Router {
	r := mux.NewRouter()

	// v1 endpoints
	g := r.PathPrefix("/v1").Subrouter()
	g.HandleFunc("/sessions", common.EasyHandle(http.MethodPost, s.createSession))
	return r
}

// dialRPC creates a client connection to the RPC API
// and passes it to the callback function. It will close the connection
// after the callback function has finished
func (s *HTTP) dialRPC(cb func(client proto_rpc.APIClient) error) error {
	conn, err := grpc.Dial(s.rpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to dial RPC API")
	}
	defer conn.Close()
	return cb(proto_rpc.NewAPIClient(conn))
}

// Start starts the http server. Passes true to the startedCh channel
// when started
func (s *HTTP) Start(addr string, startedCB func(*HTTP)) error {

	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		logHTTP.Errorf("%+v", errors.Wrap(err, "failed to parse address"))
		return err
	}

	// attempt to connect to the RPC API
	if err := s.dialRPC(func(client proto_rpc.APIClient) error {
		return nil
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err
	}

	time.AfterFunc(2*time.Second, func() {
		logHTTP.Infof("Started HTTP API server @ :%s", port)
		if startedCB != nil {
			startedCB(s)
		}
	})

	return http.ListenAndServe(addr, s.getRouter())
}

// createSession creates a new session
func (s *HTTP) createSession(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var body proto_rpc.DBSession
	var resp *proto_rpc.DBSession

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		resp, err = client.CreateDBSession(context.Background(), &body)
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}
