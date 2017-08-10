package session

import (
	"fmt"
	"time"

	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/ellcrys/elldb/config"
	"github.com/ellcrys/elldb/core/servers/db"
	"github.com/jinzhu/gorm"
	"github.com/ncodes/jsq"
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

// QueryError represents an error with a query object or expression
type QueryError struct {
	Message string
}

func (e *QueryError) Error() string {
	return e.Message
}

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

	// opChan is the channel where new operation is received from
	opChan chan *Op

	// curOp represents an operation
	curOp *Op

	// db is the database connection
	db *gorm.DB

	// Error indicates the latest error that occurred while executing an operation
	Error error

	// stop causes the agent to exit
	stop bool

	// busy indicates that an operation is being executed
	busy bool

	// tx represents the current db transaction
	tx *gorm.DB

	// txFinished indicates that the current transaction has completed
	txFinished bool

	// indicates that the agent has started or Start() has been called
	began bool

	// debug turns on jsq logs
	debug bool
}

// NewAgent creates a new agent
func NewAgent(db *gorm.DB, opChan chan *Op) *Agent {
	return &Agent{db: db, opChan: opChan, txFinished: true}
}

// Reset resets the agent
func (a *Agent) Reset() {
	a.Lock()
	a.Error = nil
	a.stop = false
	a.busy = false
	a.txFinished = true
	a.curOp = nil
	a.Unlock()
	a.newTx()
}

// Stop stops the agent and rolls back existing transaction
func (a *Agent) Stop() {
	a.Lock()
	defer a.Unlock()
	if !a.stop {
		if a.curOp != nil {
			a.Unlock()
			a.closeDoneChan(a.curOp)
			a.Lock()
		}
		if !a.txFinished && a.tx != nil {
			a.tx.Rollback()
		}
		a.stop = true
	}
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
			return &QueryError{Message: "query must be a json string"}
		}

		jsq := jsq.NewJSQ(db.GetValidObjectFields())
		if err := jsq.Parse(a.curOp.QueryWithJSQ); err != nil {
			return &QueryError{Message: err.Error()}
		}

		sql, args, err := jsq.ToSQL()
		if err != nil {
			return &QueryError{Message: err.Error()}
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
			return &QueryError{Message: "query must be a json string"}
		}

		jsq := jsq.NewJSQ(db.GetValidObjectFields())
		if err := jsq.Parse(a.curOp.QueryWithJSQ); err != nil {
			return &QueryError{Message: err.Error()}
		}

		sql, args, err := jsq.ToSQL()
		if err != nil {
			return &QueryError{Message: err.Error()}
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

	a.Lock()
	defer a.Unlock()
	var err error

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
			return &QueryError{Message: "query must be a json string"}
		}

		jsq := jsq.NewJSQ(db.GetValidObjectFields())
		if err := jsq.Parse(a.curOp.QueryWithJSQ); err != nil {
			return &QueryError{Message: err.Error()}
		}

		sql, args, err := jsq.ToSQL()
		if err != nil {
			return &QueryError{Message: err.Error()}
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

	a.Lock()
	defer a.Unlock()
	var err error

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
			return &QueryError{Message: "query must be a json string"}
		}

		jsq := jsq.NewJSQ(db.GetValidObjectFields())
		if err := jsq.Parse(a.curOp.QueryWithJSQ); err != nil {
			return &QueryError{Message: err.Error()}
		}

		sql, args, err := jsq.ToSQL()
		if err != nil {
			return &QueryError{Message: err.Error()}
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
	a.Lock()
	if a.tx != nil && !a.txFinished {
		a.Error = a.tx.Commit().Error
		a.txFinished = true
		a.Unlock()
		a.newTx()
	} else {
		a.Unlock()
	}
}

// rollback rolls back the current transaction
func (a *Agent) rollback() {
	a.Lock()
	if a.tx != nil && !a.txFinished {
		a.Error = a.tx.Rollback().Error
		a.txFinished = true
		a.Unlock()
		a.newTx()
	} else {
		a.Unlock()
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

	if !a.stop {
		a.began = true
		a.newTx()
	}

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
			case OpGetObjects:
				op.Error = a.get()
			case OpCountObjects:
				op.Error = a.count()
			case OpUpdateObjects:
				op.Error = a.update()
			case OpDeleteObjects:
				op.Error = a.delete()
			}

			a.closeDoneChan(op)
			a.busy = false

		case <-time.After(MaxSessionIdleTime):
			a.Stop()
		}
	}

	a.began = false
	if endCb != nil {
		endCb()
	}
}
