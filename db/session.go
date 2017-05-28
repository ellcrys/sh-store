package db

import (
	"fmt"

	"sync"

	"github.com/ellcrys/util"
	"github.com/ncodes/patchain"
)

var (
	maxOpenConnection = 10
	maxIdleConnection = 5
)

// Session defines a structure for a partition chain session manager
type Session struct {
	sync.Mutex
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
	s.Lock()
	defer s.Unlock()
	return len(s.agents)
}

// HasSession checks whether a session exists
func (s *Session) HasSession(id string) bool {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.agents[id]; ok {
		return true
	}
	return false
}

// GetAgent gets a session agent. Returns nil if session does not exists
func (s *Session) GetAgent(id string) *Agent {
	if s.HasSession(id) {
		s.Lock()
		defer s.Unlock()
		return s.agents[id]
	}
	return nil
}

// SendOp sends an operation to a session agent
func (s *Session) SendOp(id string, op *Op) error {
	agent := s.GetAgent(id)
	if agent == nil {
		return fmt.Errorf("session not found")
	}
	if agent.busy {
		return fmt.Errorf("session is busy")
	}

	agent.opChan <- op
	if op.Done != nil {
		<-op.Done
	}

	if op.Error != nil {
		return op.Error
	}

	return agent.Error
}

// Stop all active sessions
func (s *Session) Stop() {
	if s.db != nil {
		s.db.Close()
	}

	s.Lock()
	defer s.Unlock()
	for _, agent := range s.agents {
		agent.Stop()
	}
}

// remove an agent reference
func (s *Session) removeAgentRef(id string) {
	s.Lock()
	defer s.Unlock()
	delete(s.agents, id)
}

// RemoveAgent stops and removes an agent
func (s *Session) RemoveAgent(id string) {
	if s.HasSession(id) {
		s.GetAgent(id).Stop()
		s.removeAgentRef(id)
	}
}

// CommitEnd commits the current transaction of an agent and
// removes the session/agent
func (s *Session) CommitEnd(id string) {
	if s.HasSession(id) {
		s.GetAgent(id).commit()
		s.RemoveAgent(id)
	}
}

// RollbackEnd rolls back the current transaction of an agent and
// removes the session/agent
func (s *Session) RollbackEnd(id string) {
	if s.HasSession(id) {
		s.GetAgent(id).commit()
		s.RemoveAgent(id)
	}
}

// CreateSession creates a new session agent
func (s *Session) CreateSession(id string) string {

	if id == "" {
		id = util.UUID4()
	}

	if s.db == nil {
		panic("not connected to database")
	}

	// create the new session agent, start it on a goroutine
	// and save a reference to it.
	msgChan := make(chan *Op)
	agent := NewAgent(s.db, msgChan)
	go agent.Begin(func() {
		s.RemoveAgent(id)
	})

	s.Lock()
	s.agents[id] = agent
	s.Unlock()

	return id
}
