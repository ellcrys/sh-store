package session

import (
	"fmt"
	"time"

	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/elldb/config"
	"github.com/ellcrys/elldb/servers/db"
	"github.com/jinzhu/gorm"
	"github.com/ncodes/jsq"
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

	// OpUpdateObjects represents an update operation
	OpUpdateObjects OpType = 4

	// OpDeleteObjects represents a delete operation
	OpDeleteObjects OpType = 5

	// MaxSessionIdleTime is the maximum duration a session can be idle before stopping
	MaxSessionIdleTime = 10 * time.Minute

	// ErrAgentBusy represents an agent in busy state
	ErrAgentBusy = fmt.Errorf("agent is busy")
)

// Op represents an agent operation
type Op struct {

	// Operation type
	OpType OpType

	// Bucket name to query
	Bucket string

	// PutObjects holds the list of objects to put in the Bucket
	PutObjects []*db.Object

	// QueryWithJSQ is the query expression for retrieving objects
	QueryWithJSQ string

	// QueryWithObject is the query object for retrieving objects
	QueryWithObject *db.Object

	// UpdateObject is the object containing updates
	UpdateObject interface{}

	// NumAffectedObjects indicates the number of objects affected in an update/delete operation
	NumAffectedObjects int64

	// Out will be populated with queried objects
	Out interface{}

	// Done is closed when the operation is completed
	Done chan struct{}

	// Limit determines the number of objects fetched
	Limit int

	// OrderBy determines the ordering of the fetched objects
	OrderBy string

	// Error contains the last error that occurred during the operation
	Error error
}

// Agent defines a structure for an agent
// that holds a database session
type Agent struct {
	sync.Mutex
	opChan     chan *Op
	curOp      *Op
	db         *gorm.DB
	Error      error
	stop       bool
	busy       bool
	tx         *gorm.DB
	txFinished bool
	began      bool
	debug      bool
}

// NewAgent creates a new agent
func NewAgent(db *gorm.DB, opChan chan *Op) *Agent {
	return &Agent{db: db, opChan: opChan}
}

// Reset resets the agent
func (a *Agent) Reset() {
	a.Error = nil
	a.stop = false
	a.busy = false
	a.txFinished = false
	a.curOp = nil
	a.newTx()
}

// Stop stops the current session
func (a *Agent) Stop() {
	if a.curOp != nil {
		a.closeDoneChan(a.curOp)
	}
	a.stop = true
}

// put one or more objects
func (a *Agent) put() error {
	a.Lock()
	defer a.Unlock()

	if !a.began {
		return fmt.Errorf("agent has not started. Did you call Start()?")
	}

	return db.CreateObjects(a.tx, a.curOp.Bucket, a.curOp.PutObjects)
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

	var query db.Object

	// Use JSQ enabled query if op.QueryWithJSQ is set
	if len(a.curOp.QueryWithJSQ) > 0 {

		if !govalidator.IsJSON(a.curOp.QueryWithJSQ) {
			return fmt.Errorf("query must be a json string")
		}

		jsq := jsq.NewJSQ(db.GetValidObjectFields())
		if err := jsq.Parse(a.curOp.QueryWithJSQ); err != nil {
			return fmt.Errorf("parser:%s", err)
		}

		sql, args, err := jsq.ToSQL()
		if err != nil {
			return errors.Wrap(err, "failed to generate SQL")
		}

		if a.debug {
			logAgent.Debugf("jsq: %s %v, order: %s, limit: %d", sql, args, a.curOp.OrderBy, a.curOp.Limit)
		}

		query = db.Object{
			QueryParams: db.QueryParams{
				Limit:   a.curOp.Limit,
				OrderBy: a.curOp.OrderBy,
				Expr: db.Expr{
					Expr: sql,
					Args: args,
				},
			},
		}
	} else {
		query = *a.curOp.QueryWithObject
		query.QueryParams.Limit = a.curOp.Limit
		query.QueryParams.OrderBy = a.curOp.OrderBy
	}

	return db.QueryObjects(a.tx, &query, a.curOp.Out)
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

	var query db.Object

	// Use JSQ enabled query if op.QueryWithJSQ is set
	if len(a.curOp.QueryWithJSQ) > 0 {

		if !govalidator.IsJSON(a.curOp.QueryWithJSQ) {
			return fmt.Errorf("query must be a json string")
		}

		jsq := jsq.NewJSQ(db.GetValidObjectFields())
		if err := jsq.Parse(a.curOp.QueryWithJSQ); err != nil {
			return fmt.Errorf("parser:%s", err)
		}

		sql, args, err := jsq.ToSQL()
		if err != nil {
			return errors.Wrap(err, "failed to generate SQL")
		}

		if a.debug {
			logAgent.Debugf("jsq: %s %v, order: %s, limit: %d", sql, args, a.curOp.OrderBy, a.curOp.Limit)
		}

		query = db.Object{
			QueryParams: db.QueryParams{
				Expr: db.Expr{
					Expr: sql,
					Args: args,
				},
			},
		}
	} else {
		query = *a.curOp.QueryWithObject
	}

	return db.CountObjects(a.tx, &query, a.curOp.Out)
}

// update objects matching a query
func (a *Agent) update() error {

	var err error

	a.Lock()
	defer a.Unlock()

	if !a.began {
		return fmt.Errorf("agent has not started. Did you call Start()?")
	}

	if a.curOp == nil {
		return fmt.Errorf("operation object is required")
	}

	var query db.Object

	// Use JSQ enabled query if op.QueryWithJSQ is set
	if len(a.curOp.QueryWithJSQ) > 0 {

		if !govalidator.IsJSON(a.curOp.QueryWithJSQ) {
			return fmt.Errorf("query must be a json string")
		}

		jsq := jsq.NewJSQ(db.GetValidObjectFields())
		if err := jsq.Parse(a.curOp.QueryWithJSQ); err != nil {
			return fmt.Errorf("parser:%s", err)
		}

		sql, args, err := jsq.ToSQL()
		if err != nil {
			return errors.Wrap(err, "failed to generate SQL")
		}

		if a.debug {
			logAgent.Debugf("jsq: %s %v, order: %s, limit: %d", sql, args, a.curOp.OrderBy, a.curOp.Limit)
		}

		query = db.Object{
			QueryParams: db.QueryParams{
				Expr: db.Expr{
					Expr: sql,
					Args: args,
				},
			},
		}
	} else {
		query = *a.curOp.QueryWithObject
	}

	a.curOp.NumAffectedObjects, err = db.UpdateObjects(a.tx, &query, a.curOp.UpdateObject)
	return err
}

// delete objects matching a query
func (a *Agent) delete() error {

	var err error

	a.Lock()
	defer a.Unlock()

	if !a.began {
		return fmt.Errorf("agent has not started. Did you call Start()?")
	}

	if a.curOp == nil {
		return fmt.Errorf("operation object is required")
	}

	var query db.Object

	// Use JSQ enabled query if op.QueryWithJSQ is set
	if len(a.curOp.QueryWithJSQ) > 0 {

		if !govalidator.IsJSON(a.curOp.QueryWithJSQ) {
			return fmt.Errorf("query must be a json string")
		}

		jsq := jsq.NewJSQ(db.GetValidObjectFields())
		if err := jsq.Parse(a.curOp.QueryWithJSQ); err != nil {
			return fmt.Errorf("parser:%s", err)
		}

		sql, args, err := jsq.ToSQL()
		if err != nil {
			return errors.Wrap(err, "failed to generate SQL")
		}

		if a.debug {
			logAgent.Debugf("jsq: %s %v, order: %s, limit: %d", sql, args, a.curOp.OrderBy, a.curOp.Limit)
		}

		query = db.Object{
			QueryParams: db.QueryParams{
				Expr: db.Expr{
					Expr: sql,
					Args: args,
				},
			},
		}
	} else {
		query = *a.curOp.QueryWithObject
	}

	a.curOp.NumAffectedObjects, err = db.DeleteObjects(a.tx, &query)
	return err
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
		a.Error = a.tx.Commit().Error
		a.txFinished = true
		a.Unlock()
		a.newTx()
	}
}

// rollback rolls back the current transaction
func (a *Agent) rollback() {
	if a.tx != nil && !a.txFinished {
		a.Lock()
		a.Error = a.tx.Rollback().Error
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
				if op.Error = a.put(); op.Error != nil {
					a.rollback()
				}
			case OpGetObjects:
				if op.Error = a.get(); op.Error != nil {
					a.rollback()
				}
			case OpCountObjects:
				if op.Error = a.count(); op.Error != nil {
					a.rollback()
				}
			case OpUpdateObjects:
				if op.Error = a.update(); op.Error != nil {
					a.rollback()
				}
			case OpDeleteObjects:
				if op.Error = a.delete(); op.Error != nil {
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
