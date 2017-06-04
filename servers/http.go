package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"net"

	"github.com/ellcrys/util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/labstack/gommon/log"
	"github.com/ncodes/cocoon/core/config"
	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/object"
	"github.com/ncodes/safehold/servers/common"
	"github.com/ncodes/safehold/servers/oauth"
	"github.com/ncodes/safehold/servers/proto_rpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var logHTTP = config.MakeLogger("http")

// HTTP defines a structure for an HTTP server
// that provides REST API services.
type HTTP struct {
	rpcServerAddr string
	oauth         *oauth.OAuth
	db            patchain.DB
}

// NewHTTP creates an new http server instance
func NewHTTP(db patchain.DB, rpcServerAddr string) *HTTP {
	return &HTTP{
		rpcServerAddr: rpcServerAddr,
		db:            db,
		oauth:         oauth.NewOAuth(db),
	}
}

// getRouter returns the router
func (s *HTTP) getRouter() *mux.Router {
	r := mux.NewRouter()

	// oauth endpoints
	r.HandleFunc("/token", common.EasyHandle(http.MethodPost, s.oauth.GetToken))
	g := r.PathPrefix("/v1").Subrouter()

	// v1 endpoints
	g.HandleFunc("/identities", common.EasyHandle(http.MethodPost, s.createIdentity))
	g.HandleFunc("/identities/{id}", common.EasyHandle(http.MethodGet, s.getIdentity))
	g.HandleFunc("/mappings", common.EasyHandle(http.MethodPost, s.createMapping)).Methods(http.MethodPost)
	g.HandleFunc("/mappings", common.EasyHandle(http.MethodGet, s.getAllMapping)).Methods(http.MethodGet)
	g.HandleFunc("/mappings/{name}", common.EasyHandle(http.MethodGet, s.getMapping))
	g.HandleFunc("/sessions", common.EasyHandle(http.MethodPost, s.createSession))
	g.HandleFunc("/sessions/{id}", common.EasyHandle(http.MethodGet, s.getSession)).Methods(http.MethodGet)
	g.HandleFunc("/sessions/{id}", common.EasyHandle(http.MethodDelete, s.deleteSession)).Methods(http.MethodDelete)
	g.HandleFunc("/objects", common.EasyHandle(http.MethodPost, s.createObjects)).Methods(http.MethodPost)
	g.HandleFunc("/objects/query", common.EasyHandle(http.MethodPost, s.getObjects)).Methods(http.MethodPost)
	g.HandleFunc("/objects/count", common.EasyHandle(http.MethodPost, s.countObjects)).Methods(http.MethodPost)

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
		var md metadata.MD
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CreateDBSession(ctx, &body)
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}

// getSession gets a new session
func (s *HTTP) getSession(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.DBSession

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md metadata.MD
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetDBSession(ctx, &proto_rpc.DBSession{
			ID: mux.Vars(r)["id"],
		})
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}

// deleteSession deletes a session
func (s *HTTP) deleteSession(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.DBSession

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md metadata.MD
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.DeleteDBSession(ctx, &proto_rpc.DBSession{
			ID: mux.Vars(r)["id"],
		})
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}

// createIdentity creates an identity
func (s *HTTP) createIdentity(w http.ResponseWriter, r *http.Request) (interface{}, int) {

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
// Requires an app token included in the authorization header.
// If mapping name is provided, pass the mapping name and body encoded in json to rpc.CreateObjects
// so that it can process/unmap the custom/mapped fields in the request body.
// Otherwise, copy the body into a slice of proto_rpc.Object and pass this.
func (s *HTTP) createObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var err error
	var bodyJSON []map[string]interface{}
	var resp *proto_rpc.MultiObjectResponse
	var mappingName = r.URL.Query().Get("mapping")
	var sessionID = r.URL.Query().Get("session")
	var md = metadata.Join(metadata.Pairs("session_id", sessionID))
	var authorization = r.Header.Get("Authorization")

	if err = json.NewDecoder(r.Body).Decode(&bodyJSON); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		if len(authorization) > 0 {
			md = metadata.Join(md, metadata.Pairs("authorization", authorization))
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{
			Objects: util.MustStringify(bodyJSON),
			Mapping: mappingName,
		})
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}

// createMapping creates a mapping for an identity
func (s *HTTP) createMapping(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var body struct {
		Name    string                 `json:"name"`
		Mapping map[string]interface{} `json:"mapping"`
	}
	var resp *proto_rpc.CreateMappingResponse

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		authorization := r.Header.Get("Authorization")
		var md metadata.MD
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{
			Name:    body.Name,
			Mapping: util.MustStringify(body.Mapping),
		})
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 201
}

// getMapping fetches a mapping belonging to the logged in developer
func (s *HTTP) getMapping(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var md metadata.MD
	var resp *proto_rpc.GetMappingResponse

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetMapping(ctx, &proto_rpc.GetMappingMsg{
			Name: mux.Vars(r)["name"],
		})
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	var respMap map[string]interface{}
	if err := util.FromJSON(resp.Mapping, &respMap); err != nil {
		logHTTP.Errorf("%+v", errors.Wrap(err, "failed to parse response to json"))
		return common.ServerError, 500
	}

	return respMap, 200
}

// getAllMapping fetches the most recent mappings belonging to the logged in developer
func (s *HTTP) getAllMapping(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var md metadata.MD
	var resp *proto_rpc.GetAllMappingResponse

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{
			Limit: 50,
			Name:  r.URL.Query().Get("name"),
		})
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	var respMaps []map[string]interface{}
	for _, mapping := range resp.Mappings {
		_, name, _ := object.SplitKey(mapping.Name)
		var _mapping = map[string]interface{}{
			"name": name,
		}
		var mappingMap map[string]interface{}
		if err := util.FromJSON(mapping.Mapping, &mappingMap); err != nil {
			logHTTP.Errorf("%+v", errors.Wrap(err, "failed to parse response to json"))
			return common.ServerError, 500
		}
		_mapping["mapping"] = mappingMap
		respMaps = append(respMaps, _mapping)
	}

	return respMaps, 200
}

// getIdentity gets an identity
func (s *HTTP) getIdentity(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.ObjectResponse

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md metadata.MD
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Pairs("authorization", authorization)
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetIdentity(ctx, &proto_rpc.GetIdentityMsg{
			ID: mux.Vars(r)["id"],
		})
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}

// getObjects performs query operations
func (s *HTTP) getObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.MultiObjectResponse
	var rpcBody proto_rpc.GetObjectMsg
	var body struct {
		Query   map[string]interface{} `json:"query"`
		Owner   string                 `json:"owner"`
		Creator string                 `json:"creator"`
		Limit   int32                  `json:"limit"`
		Order   []struct {
			Field string `json:"field"`
			Order int32  `json:"order"`
		} `json:"order"`
	}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	copier.Copy(&rpcBody, body)
	rpcBody.Query = util.MustStringify(body.Query)

	var sessionID = r.URL.Query().Get("session")

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md = metadata.Join(metadata.Pairs("session_id", sessionID))
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Join(md, metadata.Pairs("authorization", authorization))
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetObjects(ctx, &rpcBody)
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}

// countObjects counts objects matching a query
func (s *HTTP) countObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.ObjectCountResponse
	var rpcBody proto_rpc.GetObjectMsg
	var body struct {
		Query   map[string]interface{} `json:"query"`
		Owner   string                 `json:"owner"`
		Creator string                 `json:"creator"`
	}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	copier.Copy(&rpcBody, body)
	rpcBody.Query = util.MustStringify(body.Query)

	var sessionID = r.URL.Query().Get("session")

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		var md = metadata.Join(metadata.Pairs("session_id", sessionID))
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 {
			md = metadata.Join(md, metadata.Pairs("authorization", authorization))
		}
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CountObjects(ctx, &rpcBody)
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}
