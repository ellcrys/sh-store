package session

import (
	"fmt"

	"github.com/ellcrys/util"
)

// ErrNotFound indicates that a session was not found
var ErrNotFound = fmt.Errorf("not found")

// RegItem defines a session registry item
type RegItem struct {
	Address string                 `json:"address"`
	Port    int                    `json:"port"`
	SID     string                 `json:"sid"`
	Meta    map[string]interface{} `json:"meta"`
}

// ToJSON returns JSON representation
func (i *RegItem) ToJSON() []byte {
	return util.MustStringify(i)
}

// SessionRegistry defines an interface for a session registory
type SessionRegistry interface {
	Add(item RegItem) error
	Get(sid string) (*RegItem, error)
	Del(sid string) error
}
