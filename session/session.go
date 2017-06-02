package session

import (
	"fmt"
	"net"
	"strconv"

	"sync"

	"github.com/ellcrys/util"
	"github.com/ncodes/patchain"
	"github.com/pkg/errors"
)

var (
	maxOpenConnection = 10
	maxIdleConnection = 5

	// nodePort is the port to the running RPC service
	nodePort int
)

func init() {
}

func getNodeAddr() (string, int, error) {
	nodeAddr := "127.0.0.1"
	_, nodePort, err := net.SplitHostPort(util.Env("SH_RPC_ADDR", "localhost:9002"))
	if err != nil {
		return "", 0, errors.Wrap(err, "failed to split rpc addr")
	}
	port, _ := strconv.Atoi(nodePort)
	return nodeAddr, port, err
}

// Session defines a structure for a partition chain session manager
type Session struct {
	sync.Mutex
	db         patchain.DB
	agents     map[string]*Agent
	sessionReg SessionRegistry
}

// NewSession creates a new partition chain session
func NewSession(sessionReg SessionRegistry) *Session {
	return &Session{
		agents:     make(map[string]*Agent),
		sessionReg: sessionReg,
	}
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

// GetAgents returns all agents
func (s *Session) GetAgents() map[string]*Agent {
	return s.agents
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
	for sid, agent := range s.agents {
		agent.Stop()
		s.End(sid)
	}
}

// removeAgent an agent reference
func (s *Session) removeAgent(id string) {
	s.Lock()
	defer s.Unlock()
	delete(s.agents, id)
}

// End stops and removes an agent from memory and registry
func (s *Session) End(id string) {
	if s.HasSession(id) {
		s.GetAgent(id).Stop()
		s.removeAgent(id)
		s.sessionReg.Del(id) // delete from session registry
	}
}

// CommitEnd commits the current transaction of an agent and
// removes the session/agent
func (s *Session) CommitEnd(id string) {
	if s.HasSession(id) {
		s.GetAgent(id).commit()
		s.End(id)
	}
}

// RollbackEnd rolls back the current transaction of an agent and
// removes the session/agent
func (s *Session) RollbackEnd(id string) {
	if s.HasSession(id) {
		s.GetAgent(id).rollback()
		s.End(id)
	}
}

// CreateSession creates a new session agent
func (s *Session) CreateSession(id, identityID string) (string, error) {
	return s.createSession(id, identityID, true)
}

// CreateUnregisteredSession creates a session that is not registered
func (s *Session) CreateUnregisteredSession(id, identityID string) (sid string) {
	sid, _ = s.createSession(id, identityID, false)
	return
}

// createSession creates a new session. It will generate an id if
// id param is not empty and will only register the session on the registry
// if registered is true.
func (s *Session) createSession(id, identityID string, registered bool) (string, error) {

	if id == "" {
		id = util.UUID4()
	}

	// create the new session agent, start it on a goroutine
	// and save a reference to it.
	msgChan := make(chan *Op)
	agent := NewAgent(s.db.NewDB(), msgChan)
	go agent.Begin(func() { // on stop
		s.End(id) // remove agent
	})

	s.Lock()
	s.agents[id] = agent
	s.Unlock()

	if registered {
		if err := s.registerSession(id, identityID); err != nil {
			s.End(id)
			return "", fmt.Errorf("failed to register session")
		}
	}

	return id, nil
}

// registerSessionService register the session with the session registry
func (s *Session) registerSession(sid string, identityID string) error {

	addr, port, err := getNodeAddr()
	if err != nil {
		return err
	}

	return s.sessionReg.Add(RegItem{
		Address: addr,
		Port:    port,
		SID:     sid,
		Meta: map[string]interface{}{
			"identity": identityID,
		},
	})
}

// SendQueryOp sends a query operation
func SendQueryOp(ses *Session, query string, limit int, order string, out interface{}) error {
	sid := ses.CreateUnregisteredSession(util.UUID4(), util.UUID4())
	if err := ses.SendOp(sid, &Op{
		OpType:  OpGetObjects,
		Data:    query,
		Done:    make(chan struct{}),
		Out:     out,
		Limit:   limit,
		OrderBy: order,
	}); err != nil {
		return err
	}
	return nil
}

// SendQueryOpWithSession sends a query operation using an existing session id
func SendQueryOpWithSession(ses *Session, sid, query string, limit int, order string, out interface{}) error {
	agent := ses.GetAgent(sid)
	if agent == nil {
		return fmt.Errorf("session not found")
	}
	if err := ses.SendOp(sid, &Op{
		OpType:  OpGetObjects,
		Data:    query,
		Done:    make(chan struct{}),
		Out:     out,
		Limit:   limit,
		OrderBy: order,
	}); err != nil {
		return err
	}
	return nil
}

// SendCountOpWithSession sends a query operation using an existing session id
func SendCountOpWithSession(ses *Session, sid, query string, out interface{}) error {
	agent := ses.GetAgent(sid)
	if agent == nil {
		return fmt.Errorf("session not found")
	}
	if err := ses.SendOp(sid, &Op{
		OpType: OpCountObjects,
		Data:   query,
		Done:   make(chan struct{}),
		Out:    out,
	}); err != nil {
		return err
	}
	return nil
}
