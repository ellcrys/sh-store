package db

import "github.com/ncodes/patchain"

var (
	maxOpenConnection = 10
	maxIdleConnection = 5
)

// Session defines a structure for a partition chain session manager
type Session struct {
	ConnectionString string
	db               patchain.DB
}

// NewSession creates a new partition chain session
func NewSession(conStr string) *Session {
	return &Session{ConnectionString: conStr}
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

// Stop all active sessions
func (s *Session) Stop() {
	if s.db != nil {
		s.db.Close()
	}
}
