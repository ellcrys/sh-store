package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"net"

	"github.com/ellcrys/cocoon/core/config"
	"github.com/ellcrys/elldb/servers/common"
	"github.com/ellcrys/elldb/servers/oauth"
	"github.com/ellcrys/elldb/servers/proto_rpc"
	"github.com/ellcrys/util"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
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
}

// NewHTTP creates an new http server instance
func NewHTTP(rpcServerAddr string, db *gorm.DB) *HTTP {
	return &HTTP{
		rpcServerAddr: rpcServerAddr,
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
	g.HandleFunc("/buckets", common.EasyHandle(http.MethodPost, s.createBucket))
	g.HandleFunc("/accounts", common.EasyHandle(http.MethodPost, s.createAccount))
	g.HandleFunc("/buckets/{bucket}/objects", common.EasyHandle(http.MethodPost, s.createObjects)).Methods(http.MethodPost)
	g.HandleFunc("/buckets/{bucket}/objects/query", common.EasyHandle(http.MethodPost, s.getObjects)).Methods(http.MethodPost)
	g.HandleFunc("/buckets/{bucket}/objects/count", common.EasyHandle(http.MethodPost, s.countObjects)).Methods(http.MethodPost)
	g.HandleFunc("/buckets/{bucket}/objects/update", common.EasyHandle(http.MethodPost, s.updateObjects)).Methods(http.MethodPost)
	g.HandleFunc("/buckets/{bucket}/objects", common.EasyHandle(http.MethodDelete, s.deleteObjects)).Methods(http.MethodDelete)
	g.HandleFunc("/accounts/{id}", common.EasyHandle(http.MethodGet, s.getAccount))
	g.HandleFunc("/mappings", common.EasyHandle(http.MethodPost, s.createMapping)).Methods(http.MethodPost)
	g.HandleFunc("/mappings", common.EasyHandle(http.MethodGet, s.getAllMapping)).Methods(http.MethodGet)
	g.HandleFunc("/mappings/{name}", common.EasyHandle(http.MethodGet, s.getMapping))
	g.HandleFunc("/sessions", common.EasyHandle(http.MethodPost, s.createSession))
	g.HandleFunc("/sessions/{id}", common.EasyHandle(http.MethodGet, s.getSession)).Methods(http.MethodGet)
	g.HandleFunc("/sessions/{id}", common.EasyHandle(http.MethodDelete, s.deleteSession)).Methods(http.MethodDelete)
	g.HandleFunc("/sessions/{id}/commit", common.EasyHandle(http.MethodGet, s.commitSession)).Methods(http.MethodGet)
	g.HandleFunc("/sessions/{id}/rollback", common.EasyHandle(http.MethodGet, s.rollbackSession)).Methods(http.MethodGet)

	return r
}

// dialRPC creates a client connection to the RPC API
// and passes it to the callback function. It will close the connection
// after the callback function returns
func (s *HTTP) dialRPC(cb func(client proto_rpc.APIClient) error) error {
	conn, err := grpc.Dial(s.rpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "failed to dial RPC API")
	}
	defer conn.Close()
	if cb == nil {
		return nil
	}
	return cb(proto_rpc.NewAPIClient(conn))
}

// Start starts the http server
// Passes true to the startedCh channel when started
func (s *HTTP) Start(addr string, startedCB func(*HTTP)) error {

	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		logHTTP.Errorf("%+v", errors.Wrap(err, "failed to parse address"))
		return err
	}

	// attempt to connect to the RPC API
	if err := s.dialRPC(nil); err != nil {
		logHTTP.Errorf("%+v", err)
		return err
	}

	time.AfterFunc(2*time.Second, func() {
		logHTTP.Infof("Started server @ :%s", port)
		if startedCB != nil {
			startedCB(s)
		}
	})

	return http.ListenAndServe(addr, s.getRouter())
}

// createBucket creates a bucket
func (s *HTTP) createBucket(w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var err error
	var body proto_rpc.CreateBucketMsg
	var resp *proto_rpc.Bucket

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		md := metadata.Pairs("authorization", r.Header.Get("Authorization"))
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CreateBucket(ctx, &body)
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return common.SingleObjectResp("bucket", structs.New(resp).Map()), 201
}

// createSession creates a new session
func (s *HTTP) createSession(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var body proto_rpc.Session
	var resp *proto_rpc.Session

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		md := metadata.Pairs("authorization", r.Header.Get("Authorization"))
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CreateSession(ctx, &body)
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
	var resp *proto_rpc.Session

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		md := metadata.Pairs("authorization", r.Header.Get("Authorization"))
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetSession(ctx, &proto_rpc.Session{
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
	var resp *proto_rpc.Session

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		md := metadata.Pairs("authorization", r.Header.Get("Authorization"))
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.DeleteSession(ctx, &proto_rpc.Session{
			ID: mux.Vars(r)["id"],
		})
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}

// commitSession commits a session
func (s *HTTP) commitSession(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.Session

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		md := metadata.Pairs("authorization", r.Header.Get("Authorization"))
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CommitSession(ctx, &proto_rpc.Session{
			ID: mux.Vars(r)["id"],
		})
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}

// rollbackSession rolls back a session
func (s *HTTP) rollbackSession(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.Session

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		md := metadata.Pairs("authorization", r.Header.Get("Authorization"))
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.RollbackSession(ctx, &proto_rpc.Session{
			ID: mux.Vars(r)["id"],
		})
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}

// createAccount creates an account
func (s *HTTP) createAccount(w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var err error
	var body proto_rpc.CreateAccountMsg
	var resp *proto_rpc.GetObjectResponse

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		resp, err = client.CreateAccount(context.Background(), &body)
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	var account map[string]interface{}
	util.FromJSON2(resp.Object, &account)

	return common.SingleObjectResp("account", account), 201
}

// createMapping creates a mapping for an account
func (s *HTTP) createMapping(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.CreateMappingResponse
	var body struct {
		Bucket  string                 `json:"bucket"`
		Name    string                 `json:"name"`
		Mapping map[string]interface{} `json:"mapping"`
	}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		md := metadata.Pairs("authorization", r.Header.Get("Authorization"))
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CreateMapping(ctx, &proto_rpc.CreateMappingMsg{
			Bucket:  body.Bucket,
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
	var resp *proto_rpc.GetMappingResponse

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		md := metadata.Pairs("authorization", r.Header.Get("Authorization"))
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

	var mapping map[string]interface{}
	util.FromJSON([]byte(respMap["mapping"].(string)), &mapping)

	return common.SingleObjectResp("mapping", map[string]interface{}{
		"id":      respMap["id"],
		"name":    mux.Vars(r)["name"],
		"mapping": mapping,
	}), 200
}

// getAllMapping fetches the most recent mappings belonging to the logged in developer
func (s *HTTP) getAllMapping(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.GetAllMappingResponse

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		md := metadata.Pairs("authorization", r.Header.Get("Authorization"))
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetAllMapping(ctx, &proto_rpc.GetAllMappingMsg{
			Limit: 50,
		})
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	var respMaps []map[string]interface{}
	util.FromJSON(resp.Mappings, &respMaps)
	attrs := make([]map[string]interface{}, len(respMaps))

	for i, mapping := range respMaps {
		var mappingVal map[string]string
		util.FromJSON([]byte(mapping["mapping"].(string)), &mappingVal)
		attrs[i] = map[string]interface{}{
			"id":      mapping["id"],
			"name":    mapping["name"].(string),
			"mapping": mappingVal,
		}
	}

	return common.MultiObjectResp("mapping", attrs), 200
}

// getAccount gets an account
func (s *HTTP) getAccount(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.GetObjectResponse

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		md := metadata.Pairs("authorization", r.Header.Get("Authorization"))
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetAccount(ctx, &proto_rpc.GetAccountMsg{
			ID: mux.Vars(r)["id"],
		})
		return err
	}); err != nil {
		logHTTP.Errorf("%+v", err)
		return err, 0
	}

	var account map[string]interface{}
	util.FromJSON2(resp.Object, &account)

	return common.SingleObjectResp("account", account), 200
}

// createObjects creates an object and can optionally use an existing mapping
// to process custom object fields.
// Requires an app token.
// TODO: Requires user token for user-owned objects
func (s *HTTP) createObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {

	var err error
	var bodyJSON []map[string]interface{}
	var resp *proto_rpc.GetObjectsResponse
	var md = metadata.Join(
		metadata.Pairs("session_id", r.URL.Query().Get("session")),
		metadata.Pairs("authorization", r.Header.Get("Authorization")),
	)

	if err = util.DecodeJSON(r.Body, &bodyJSON); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CreateObjects(ctx, &proto_rpc.CreateObjectsMsg{
			Bucket:  mux.Vars(r)["bucket"],
			Objects: util.MustStringify(bodyJSON),
			Mapping: r.URL.Query().Get("mapping"),
		})
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	var objects []map[string]interface{}
	util.FromJSON2(resp.Objects, &objects)

	return common.MultiObjectResp("objects", objects), 201
}

// updateObjects updates objects of a mutable bucket
func (s *HTTP) updateObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.AffectedResponse
	var body struct {
		Query   map[string]interface{} `json:"query"`
		Owner   string                 `json:"owner"`
		Creator string                 `json:"creator"`
		Update  map[string]interface{} `json:"update"`
	}
	var md = metadata.Join(
		metadata.Pairs("session_id", r.URL.Query().Get("session")),
		metadata.Pairs("authorization", r.Header.Get("Authorization")),
	)

	if err = util.DecodeJSON(r.Body, &body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.UpdateObjects(ctx, &proto_rpc.UpdateObjectsMsg{
			Bucket:  mux.Vars(r)["bucket"],
			Query:   util.MustStringify(body.Query),
			Owner:   body.Owner,
			Creator: body.Creator,
			Update:  util.MustStringify(body.Update),
			Mapping: r.URL.Query().Get("mapping"),
		})
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return map[string]interface{}{
		"affected": resp.Affected,
	}, 200
}

// deleteObjects deletes an object in a mutable bucket
func (s *HTTP) deleteObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.AffectedResponse
	var body struct {
		Query   map[string]interface{} `json:"query"`
		Owner   string                 `json:"owner"`
		Creator string                 `json:"creator"`
	}
	var md = metadata.Join(
		metadata.Pairs("session_id", r.URL.Query().Get("session")),
		metadata.Pairs("authorization", r.Header.Get("Authorization")),
	)

	if err = util.DecodeJSON(r.Body, &body); err != nil {
		return common.BodyMalformedError, 400
	}

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.DeleteObjects(ctx, &proto_rpc.DeleteObjectsMsg{
			Bucket:  mux.Vars(r)["bucket"],
			Query:   util.MustStringify(body.Query),
			Owner:   body.Owner,
			Creator: body.Creator,
			Mapping: r.URL.Query().Get("mapping"),
		})
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return map[string]interface{}{
		"affected": resp.Affected,
	}, 200
}

// getObjects fetches objects
func (s *HTTP) getObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.GetObjectsResponse
	var rpcBody proto_rpc.GetObjectMsg
	var mappingName = r.URL.Query().Get("mapping")
	var sessionID = r.URL.Query().Get("session")
	var md = metadata.Join(metadata.Pairs("session_id", sessionID), metadata.Pairs("authorization", r.Header.Get("Authorization")))
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
	rpcBody.Bucket = mux.Vars(r)["bucket"]
	rpcBody.Mapping = mappingName

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.GetObjects(ctx, &rpcBody)
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	var objects []map[string]interface{}
	util.FromJSON(resp.Objects, &objects)

	return common.MultiObjectResp("object", objects), 200
}

// countObjects counts objects matching a query
func (s *HTTP) countObjects(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	var err error
	var resp *proto_rpc.ObjectCountResponse
	var rpcBody proto_rpc.GetObjectMsg
	var sessionID = r.URL.Query().Get("session")
	var md = metadata.Join(metadata.Pairs("session_id", sessionID), metadata.Pairs("authorization", r.Header.Get("Authorization")))
	var body struct {
		Query   map[string]interface{} `json:"query"`
		Owner   string                 `json:"owner"`
		Creator string                 `json:"creator"`
	}

	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		return common.BodyMalformedError, 400
	}

	copier.Copy(&rpcBody, body)
	rpcBody.Bucket = mux.Vars(r)["bucket"]
	rpcBody.Query = util.MustStringify(body.Query)

	if err = s.dialRPC(func(client proto_rpc.APIClient) error {
		ctx := metadata.NewContext(context.Background(), md)
		resp, err = client.CountObjects(ctx, &rpcBody)
		return err
	}); err != nil {
		log.Errorf("%+v", err)
		return err, 0
	}

	return resp, 200
}
