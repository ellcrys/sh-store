package db

import (
	"time"

	"fmt"

	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach/tables"
)

// OpType represents a db connection operation
type OpType int

var (
	// OpCommit represents a commit operation
	OpCommit OpType = 1

	// OpRollback represents a rollback operation
	OpRollback OpType = 2

	// OpPutObject represents an object creation operation.
	OpPutObject OpType = 3
)

// Op represents an agent operation
type Op struct {
	OpType OpType
	Data   interface{}
}

// Agent defines a structure for an agent
// that holds a database session
type Agent struct {
	msgChan chan Op
	db      patchain.DB
	Error   error
	stop    bool
}

// NewAgent creates a new agent
func NewAgent(db patchain.DB, msgChan chan Op) *Agent {
	return &Agent{db: db, msgChan: msgChan}
}

// Reset resets the agent
func (a *Agent) Reset() {
	a.Error = nil
	a.stop = false
}

// Stop stops the current session
func (a *Agent) Stop() {
	a.stop = true
}

// Begin starts a database session
func (a *Agent) Begin() {
	tx := a.db.Begin()
	for !a.stop {
		select {
		case msg := <-a.msgChan:
			switch msg.OpType {
			case OpCommit:
				a.Error = tx.Commit()
			case OpRollback:
				a.Error = tx.Rollback()
				// case OpPutObject:

			}
		case <-time.After(10 * time.Second):
			num := int64(0)
			err := tx.Count(&tables.Object{}, &num)
			fmt.Println(num, err)
			tx.Rollback()
			a.Stop()
		}
	}
}
