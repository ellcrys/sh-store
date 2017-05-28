package db

import (
	"github.com/ellcrys/consul/api"
	"github.com/ellcrys/util"
	"github.com/pkg/errors"
)

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
	kv := r.client.KV()
	_, err := kv.Put(&api.KVPair{
		Key:   makeKey(item.SID),
		Value: item.ToJSON(),
	}, nil)
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
	_, err = r.client.KV().Delete(makeKey(sid), nil)
	if err != nil {
		return errors.Wrap(err, "failed to delete registry item")
	}
	return
}
