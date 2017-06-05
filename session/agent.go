package session

import (
	"fmt"
	"time"

	"sync"

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

	// OpPutObjects represents an object creation operation.
	OpPutObjects OpType = 1

	// OpGetObjects represents a get object operation
	OpGetObjects OpType = 2

	// OpCountObjects represents a count operation
	OpCountObjects OpType = 3

	// MaxSessionIdleTime is the maximum duration a session can be idle before stopping
	MaxSessionIdleTime = 10 * time.Minute

	// ErrAgentBusy represents an agent in busy state
	ErrAgentBusy = fmt.Errorf("agent is busy")
)

// Op represents an agent operation
type Op struct {
	OpType          OpType
	PutData         interface{}
	QueryWithJSQ    string
	QueryWithObject *tables.Object
	Out             interface{}
	Done            chan struct{}
	Limit           int
	OrderBy         string
	Error           error
}

// Agent defines a structure for an agent
// that holds a database session
type Agent struct {
	sync.Mutex
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
	a.Lock()
	defer a.Unlock()

	if !a.began {
		return fmt.Errorf("agent has not started. Did you call Start()?")
	}

	var dbOp = patchain.UseDBOption{DB: a.tx}
	return a.o.Put(a.curOp.PutData, &dbOp)
}

// get gets objects by query with JSQ or a putchain.Object.
// If op.QueryWithJSQ is set, it is used, otherwise op.QueryWithObject is used.
func (a *Agent) get() error {
	a.Lock()
	defer a.Unlock()

	if !a.began {
		return fmt.Errorf("agent has not started. Did you call Start()?")
	}

	if a.curOp == nil {
		return fmt.Errorf("operation object is required")
	}

	var finalQueryObj *tables.Object

	// Use JSQ enabled query if op.QueryWithJSQ is set
	if a.curOp.QueryWithJSQ != "" {

		if !govalidator.IsJSON(a.curOp.QueryWithJSQ) {
			return fmt.Errorf("query must be a json string")
		}

		jsq := a.tx.NewQuery()
		if err := jsq.Parse(a.curOp.QueryWithJSQ); err != nil {
			return err
		}

		sql, args, err := jsq.ToSQL()
		if err != nil {
			return errors.Wrap(err, "failed to generate SQL")
		}

		if a.debug {
			logAgent.Debugf("jsq: %s %v, order: %s, limit: %d", sql, args, a.curOp.OrderBy, a.curOp.Limit)
		}

		finalQueryObj = &tables.Object{
			QueryParams: patchain.QueryParams{
				Limit:   a.curOp.Limit,
				OrderBy: a.curOp.OrderBy,
				Expr: patchain.Expr{
					Expr: sql,
					Args: args,
				},
			},
		}
	} else {
		finalQueryObj = a.curOp.QueryWithObject
		finalQueryObj.QueryParams.Limit = a.curOp.Limit
		finalQueryObj.QueryParams.OrderBy = a.curOp.OrderBy
	}

	err := a.tx.GetAll(finalQueryObj, a.curOp.Out)
	if err != nil {
		return err
	}

	return nil
}

// counts object that match a query
func (a *Agent) count() error {
	a.Lock()
	defer a.Unlock()

	if !a.began {
		return fmt.Errorf("agent has not started. Did you call Start()?")
	}

	if a.curOp == nil {
		return fmt.Errorf("operation object is required")
	}

	var finalQueryObj *tables.Object

	// Use JSQ enabled query if op.QueryWithJSQ is set
	if a.curOp.QueryWithJSQ != "" {

		if !govalidator.IsJSON(a.curOp.QueryWithJSQ) {
			return fmt.Errorf("query must be a json string")
		}

		jsq := a.tx.NewQuery()
		if err := jsq.Parse(a.curOp.QueryWithJSQ); err != nil {
			return err
		}

		sql, args, err := jsq.ToSQL()
		if err != nil {
			return errors.Wrap(err, "failed to generate SQL")
		}

		if a.debug {
			logAgent.Debugf("jsq: %s %v, order: %s, limit: %d", sql, args, a.curOp.OrderBy, a.curOp.Limit)
		}

		finalQueryObj = &tables.Object{
			QueryParams: patchain.QueryParams{
				Limit:   a.curOp.Limit,
				OrderBy: a.curOp.OrderBy,
				Expr: patchain.Expr{
					Expr: sql,
					Args: args,
				},
			},
		}
	} else {
		finalQueryObj = a.curOp.QueryWithObject
		finalQueryObj.QueryParams.Limit = a.curOp.Limit
		finalQueryObj.QueryParams.OrderBy = a.curOp.OrderBy
	}

	err := a.tx.Count(finalQueryObj, a.curOp.Out)
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
	a.Lock()
	defer a.Unlock()

	if a.txFinished {
		a.tx = a.db.Begin()
		a.txFinished = false
	}
}

// commit commits the current transaction
func (a *Agent) commit() {
	if a.tx != nil && !a.txFinished {
		a.Lock()
		a.Error = a.tx.Commit()
		a.txFinished = true
		a.Unlock()
		a.newTx()
	}
}

// rollback rolls back the current transaction
func (a *Agent) rollback() {
	if a.tx != nil && !a.txFinished {
		a.Lock()
		a.Error = a.tx.Rollback()
		a.txFinished = true
		a.Unlock()
		a.newTx()
	}
}

// closeDoneChan closes the Op done channel and sets it to nil
func (a *Agent) closeDoneChan(op *Op) {
	a.Lock()
	defer a.Unlock()
	if op == nil {
		return
	}
	if op.Done != nil && !a.stop {
		close(op.Done)
		op.Done = nil
	}
}

// Start starts a database transaction session.
// Only a single operation is allowed to run at any point in time.
// It closes the operations done/wait channel when it completes
// and passes any error to the operation
func (a *Agent) Start(endCb func()) {

	a.Lock()
	a.tx = a.db.Begin()
	a.began = true
	a.Unlock()

	for !a.stop {
		select {
		case op := <-a.opChan:

			// if agent is busy, do not accept new operations
			if a.busy {
				op.Error = ErrAgentBusy
				a.closeDoneChan(op)
				continue
			}

			a.busy = true
			a.curOp = op

			// process operations
			switch op.OpType {
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
