package db

import (
	"fmt"

	"github.com/ellcrys/util"
	"github.com/ncodes/patchain"
)

var (
	maxOpenConnection = 10
	maxIdleConnection = 5
)

// Session defines a structure for a partition chain session manager
type Session struct {
	ConnectionString string
	db               patchain.DB
	agents           map[string]*Agent
}

// NewSession creates a new partition chain session
func NewSession(conStr string) *Session {
	return &Session{
		ConnectionString: conStr,
		agents:           make(map[string]*Agent),
	}
}

// Connect creates a connection to the partition chain database
func (s *Session) Connect(db patchain.DB) (err error) {
	err = db.Connect(maxOpenConnection, maxIdleConnection)
	if err != nil {
		return err
	}
	s.db = db
	return
}

// SetDB sets the db connection to use directly
func (s *Session) SetDB(db patchain.DB) {
	s.db = db
}

// NumSessions returns the number of active sessions
func (s *Session) NumSessions() int {
	return len(s.agents)
}

// HasSession checks whether a session exists
func (s *Session) HasSession(id string) bool {
	if _, ok := s.agents[id]; ok {
		return true
	}
	return false
}

// GetAgent gets a session agent. Returns nil if session does not exists
func (s *Session) GetAgent(id string) *Agent {
	if s.HasSession(id) {
		return s.agents[id]
	}
	return nil
}

// SendOp sends an operation to a session agent
func (s *Session) SendOp(id string, op Op) error {
	agent := s.GetAgent(id)
	if agent == nil {
		return fmt.Errorf("session not found")
	}
	agent.msgChan <- op
	return nil
}

// Stop all active sessions
func (s *Session) Stop() {
	if s.db != nil {
		s.db.Close()
	}
}

// CreateSession creates a new session agent
func (s *Session) CreateSession(id string) string {

	if id == "" {
		id = util.UUID4()
	}

	msgChan := make(chan Op)
	agent := NewAgent(s.db, msgChan)
	go agent.Begin()
	s.agents[id] = agent
	return id
}
