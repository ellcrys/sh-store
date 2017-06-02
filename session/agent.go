package session

import (
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ncodes/patchain"
	"github.com/ncodes/patchain/cockroach/tables"
	"github.com/ncodes/patchain/object"
	"github.com/ncodes/safehold/config"
	"github.com/pkg/errors"
)

// OpType represents a db connection operation
type OpType int

var (
	logAgent = config.MakeLogger("session.agent")

	// OpCommit represents a commit operation
	OpCommit OpType = 1

	// OpRollback represents a rollback operation
	OpRollback OpType = 2

	// OpPutObjects represents an object creation operation.
	OpPutObjects OpType = 3

	// OpGetObjects represents a get object operation
	OpGetObjects OpType = 4

	// OpCountObjects represents a count operation
	OpCountObjects OpType = 5

	// MaxSessionIdleTime is the maximum duration a session can be idle before stopping
	MaxSessionIdleTime = 10 * time.Minute

	// ErrAgentBusy represents an agent in busy state
	ErrAgentBusy = fmt.Errorf("agent is busy")
)

// Op represents an agent operation
type Op struct {
	OpType  OpType
	Data    interface{}
	Out     interface{}
	Done    chan struct{}
	Limit   int
	OrderBy string
	Error   error
}

// Agent defines a structure for an agent
// that holds a database session
type Agent struct {
	opChan     chan *Op
	curOp      *Op
	db         patchain.DB
	Error      error
	stop       bool
	o          *object.Object
	busy       bool
	tx         patchain.DB
	txFinished bool
	began      bool
	debug      bool
}

// NewAgent creates a new agent
func NewAgent(db patchain.DB, opChan chan *Op) *Agent {
	return &Agent{db: db, opChan: opChan, o: object.NewObject(db)}
}

// Reset resets the agent
func (a *Agent) Reset() {
	a.Error = nil
	a.stop = false
	a.busy = false
	a.txFinished = false
	a.newTx()
}

// Stop stops the current session
func (a *Agent) Stop() {
	if a.curOp != nil {
		a.closeDoneChan(a.curOp)
	}
	a.stop = true
}

// put an object
func (a *Agent) put() error {

	if !a.began {
		return fmt.Errorf("agent has not started. Did you call Begin()?")
	}

	var dbOp = patchain.UseDBOption{DB: a.tx}
	return a.o.Put(a.curOp.Data, &dbOp)
}

// get gets objects
func (a *Agent) get() error {

	if !a.began {
		return fmt.Errorf("agent has not started. Did you call Begin()?")
	}

	if a.curOp == nil {
		return fmt.Errorf("operation object is required")
	}

	queryJSON, ok := a.curOp.Data.(string)
	if !ok || !govalidator.IsJSON(queryJSON) {
		return fmt.Errorf("query must be a json string")
	}

	jsq := a.tx.NewQuery()
	if err := jsq.Parse(queryJSON); err != nil {
		return errors.Wrap(err, "failed to parse query")
	}

	sql, args, err := jsq.ToSQL()
	if err != nil {
		return errors.Wrap(err, "failed to generate SQL")
	}

	if a.debug {
		logAgent.Debugf("jsq: %s %v, order: %s, limit: %d", sql, args, a.curOp.OrderBy, a.curOp.Limit)
	}

	err = a.tx.GetAll(&tables.Object{
		QueryParams: patchain.QueryParams{
			Limit:   a.curOp.Limit,
			OrderBy: a.curOp.OrderBy,
			Expr: patchain.Expr{
				Expr: sql,
				Args: args,
			},
		},
	}, a.curOp.Out)

	if err != nil {
		return err
	}

	return nil
}

// counts object that match a query
func (a *Agent) count() error {

	if !a.began {
		return fmt.Errorf("agent has not started. Did you call Begin()?")
	}

	if a.curOp == nil {
		return fmt.Errorf("operation object is required")
	}

	queryJSON, ok := a.curOp.Data.(string)
	if !ok || !govalidator.IsJSON(queryJSON) {
		return fmt.Errorf("query must be a json string")
	}

	jsq := a.tx.NewQuery()
	if err := jsq.Parse(queryJSON); err != nil {
		return errors.Wrap(err, "failed to parse query")
	}

	sql, args, err := jsq.ToSQL()
	if err != nil {
		return errors.Wrap(err, "failed to generate SQL")
	}

	if a.debug {
		logAgent.Debugf("jsq: %s %v, order: %s, limit: %d", sql, args, a.curOp.OrderBy, a.curOp.Limit)
	}

	err = a.tx.Count(&tables.Object{
		QueryParams: patchain.QueryParams{
			Limit:   a.curOp.Limit,
			OrderBy: a.curOp.OrderBy,
			Expr: patchain.Expr{
				Expr: sql,
				Args: args,
			},
		},
	}, a.curOp.Out)

	if err != nil {
		return err
	}

	return nil
}

// Debug turns on logging
func (a *Agent) Debug() {
	a.debug = true
}

// newTx creates a new transaction
func (a *Agent) newTx() {
	a.tx = a.db.Begin()
	a.txFinished = false
}

// commit commits the current transaction
func (a *Agent) commit() {
	if a.tx != nil && !a.txFinished {
		a.Error = a.tx.Commit()
		a.txFinished = true
	}
}

// rollback rolls back the current transaction
func (a *Agent) rollback() {
	if a.tx != nil && !a.txFinished {
		a.Error = a.tx.Rollback()
		a.txFinished = true
	}
}

// closeDoneChan closes the Op done channel and sets it to nil
func (a *Agent) closeDoneChan(op *Op) {
	if op == nil {
		return
	}
	if op.Done != nil && !a.stop {
		close(op.Done)
		op.Done = nil
	}
}

// Begin starts a database transaction session.
// Only a single operation is allowed to run at any point in time.
// It closes the operations done/wait channel when it completes
// and passes any error to the operation
func (a *Agent) Begin(endCb func()) {

	a.tx = a.db.Begin()
	a.began = true

	for !a.stop {
		select {
		case op := <-a.opChan:

			// if agent is busy, do not accept new operation
			if a.busy {
				op.Error = ErrAgentBusy
				a.closeDoneChan(op)
				continue
			}

			a.busy = true
			a.curOp = op

			// process operations
			switch op.OpType {
			case OpCommit:
				a.commit()
				a.newTx()
			case OpRollback:
				a.rollback()
				a.newTx()
			case OpPutObjects:
				op.Error = a.put()
				if op.Error != nil {
					a.rollback()
				}
			case OpGetObjects:
				op.Error = a.get()
				if op.Error != nil {
					a.rollback()
				}
			case OpCountObjects:
				op.Error = a.count()
				if op.Error != nil {
					a.rollback()
				}
			}

			a.busy = false
			a.closeDoneChan(op)

		case <-time.After(MaxSessionIdleTime):
			a.rollback()
			a.Stop()
		}
	}

	a.began = false

	if endCb != nil {
		endCb()
	}
}
