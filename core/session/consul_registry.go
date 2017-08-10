package session

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/util"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

var SessionTTL = "900s"

// ConsulRegistry implements a session registry based on
// consul. It satisfies SessionRegistry interface.
type ConsulRegistry struct {
	client *api.Client
}

// NewConsulRegistry creates a new consul registry.
// Connection to consul is attempted, error is returned on failure
// to connect.
func NewConsulRegistry() (r *ConsulRegistry, err error) {
	r = new(ConsulRegistry)
	cfg := api.DefaultConfig()
	cfg.Address = util.Env("CONSUL_ADDR", "127.0.0.1:8500")
	r.client, err = api.NewClient(cfg)
	if err != nil {
		err = errors.Wrap(err, "failed to create client")
		return
	}
	if _, err = r.client.Status().Leader(); err != nil {
		err = errors.Wrap(err, "failed to connect to consul")
	}
	return
}

// makeKey creates a key to be used in consul KV store
func makeKey(sid string) string {
	return "db_session_" + sid
}

// Add registers a new session
func (r *ConsulRegistry) Add(item RegItem) error {

	if !govalidator.IsUUIDv4(item.SID) {
		return fmt.Errorf("session id must be a UUID")
	}

	session := r.client.Session()
	sessionID, _, err := session.Create(&api.SessionEntry{
		ID:       item.SID,
		TTL:      SessionTTL,
		Behavior: api.SessionBehaviorDelete,
	}, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create consul session")
	}

	acquired, _, err := r.client.KV().Acquire(&api.KVPair{
		Key:     makeKey(item.SID),
		Value:   item.ToJSON(),
		Session: sessionID,
	}, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create consul key for db session")
	}

	if !acquired {
		return fmt.Errorf("failed to acquire lock on session")
	}

	return err
}

// Get gets a session. Returns ErrNotFound if not found
func (r *ConsulRegistry) Get(sid string) (*RegItem, error) {

	kv, _, err := r.client.KV().Get(makeKey(sid), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session")
	}
	if kv == nil {
		return nil, ErrNotFound
	}
	var item RegItem
	if err = util.FromJSON(kv.Value, &item); err != nil {
		return nil, errors.Wrap(err, "failed to parse registry item data")
	}
	return &item, nil
}

// Del deletes a session. Returns nil if the session was
// deleted successfully or it doesn't exists
func (r *ConsulRegistry) Del(sid string) (err error) {

	if !govalidator.IsUUIDv4(sid) {
		return fmt.Errorf("session id must be a UUID")
	}

	_, err = r.client.Session().Destroy(sid, nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete registry item")
	}

	_, err = r.client.KV().Delete(makeKey(sid), nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete session key")
	}

	return
}
