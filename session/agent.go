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

	// MaxSessionIdleTime is the maximum duration a session can be idle before stopping
	MaxSessionIdleTime = 10 * time.Minute
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
	a.stop = true
}

// put an object
func (a *Agent) put(d interface{}) error {

	if !a.began {
		return fmt.Errorf("agent has not started. Did you can Begin()?")
	}

	var dbOp = patchain.UseDBOption{DB: a.tx}
	return a.o.Put(d, &dbOp)
}

// get gets objects
func (a *Agent) get(op *Op, out interface{}) error {

	if !a.began {
		return fmt.Errorf("agent has not started. Did you can Begin()?")
	}

	if op == nil {
		return fmt.Errorf("operation object is required")
	}

	queryJSON, ok := op.Data.(string)
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
		logAgent.Debugf("jsq: %s %v, order: %s, limit: %d", sql, args, op.OrderBy, op.Limit)
	}

	err = a.tx.GetAll(&tables.Object{
		QueryParams: patchain.QueryParams{
			Limit:   op.Limit,
			OrderBy: op.OrderBy,
			Expr: patchain.Expr{
				Expr: sql,
				Args: args,
			},
		},
	}, out)

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

// closeDoneChan closes the Op done channel
func (a *Agent) closeDoneChan(op *Op) {
	if op.Done != nil {
		close(op.Done)
	}
}

// Begin starts a database session
func (a *Agent) Begin(endCb func()) {
	a.tx = a.db.Begin()
	a.began = true
	for !a.stop {
		select {
		case op := <-a.opChan:
			a.busy = true
			switch op.OpType {
			case OpCommit:
				a.commit()
				a.newTx()
			case OpRollback:
				a.rollback()
				a.newTx()
			case OpPutObjects:
				op.Error = a.put(op.Data)
				if op.Error != nil {
					a.rollback()
				}
			case OpGetObjects:
				op.Error = a.get(op, op.Out)
				if op.Error != nil {
					a.rollback()
				}
			}
			a.closeDoneChan(op)
		case <-time.After(MaxSessionIdleTime):
			a.rollback()
			a.Stop()
		}
		a.busy = false
	}
	a.began = false

	if endCb != nil {
		endCb()
	}
}
