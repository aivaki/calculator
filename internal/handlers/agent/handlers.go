package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/aivaki/calculator/internal/storage"
	"github.com/aivaki/calculator/internal/types"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

const (
	configPath = `../config/.env`

	operationsExecutionTime = `SELECT execution_time FROM Operations`

	unfinishedExpressions = `SELECT key, expression FROM Expressions WHERE status = 'in processing'`
)

type handler struct {
	m  sync.Mutex
	ch chan *types.Job
	u *websocket.Upgrader
}

type Handlers interface {
	Status(w http.ResponseWriter, r *http.Request)
	Connection(w http.ResponseWriter, r *http.Request)
}

func New() Handlers {
	h := new(handler)

	err := godotenv.Load(configPath)
	if err != nil {
		log.Println(`error loading configuration file:`, err)
	}

	h.upgrader()
	h.workers()
	h.previousExpressions()

	return h
}

func (h *handler) workers() {
	countWorkers := os.Getenv(countWorkersEnv)
	
	maxWorkers, err := strconv.Atoi(countWorkers)
	if err != nil {
		maxWorkers = 5
	}

	h.u = new(websocket.Upgrader)
	h.u.CheckOrigin = func(r *http.Request) bool { return true }

	h.ch = make(chan *types.Job, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go h.worker(h.ch)
	}
}

func (h *handler) upgrader() {
	h.u = new(websocket.Upgrader)
	h.u.CheckOrigin = func(r *http.Request) bool { return true }
}

func (h *handler) worker(ch <-chan *types.Job) {
	db, err := sql.Open(storage.DriverName, storage.DataSourceName)
	if err != nil {
		log.Println(`database connection error in worker:`, err)
		return
	}
	defer db.Close()

	for {
		job := <-ch

		result, err := h.evaluateExpression(job.Expression, job.ID)
		if err != nil {
			h.m.Lock()
			_, err = db.Exec(fmt.Sprintf(`UPDATE Expressions SET completed_at = CURRENT_TIMESTAMP, status = '%s' WHERE key = '%s'`, `error`, job.ID))
			if err != nil {
				log.Println(`error adding an entry to the database about an invalid expression:`, err)
				return
			}
			h.m.Unlock()
			continue
		}

		h.m.Lock()
		_, err = db.Exec(fmt.Sprintf(`UPDATE Expressions SET completed_at = CURRENT_TIMESTAMP, result = '%s', status = '%s' WHERE key = '%s'`, fmt.Sprintf(`%v`, result), `completed`, job.ID))
		if err != nil {
			log.Println(`error adding a record to the database about the executed expression:`, err)
			return
		}
		h.m.Unlock()
	}
}


func (h *handler) evaluateExpression(expression, id string) (result any, err error) {
	db, err := sql.Open(storage.DriverName, storage.DataSourceName)
	if err != nil {
		err = errors.New(`database connection error in calculation: ` + err.Error())
		return
	}
	defer db.Close()

	data, err := db.Query(operationsExecutionTime)
	if err != nil {
		err = errors.New(`failed to request operation times from the database: ` + err.Error())
		return
	}

	operations := []int{}
	for data.Next() {
		var time int
		if err = data.Scan(&time); err != nil {
			err = errors.New(`error loading operations from database: ` + err.Error())
			return
		}
		operations = append(operations, time)
	}

	plusTime := operations[0] * strings.Count(expression, `+`)
	minusTime := operations[1] * strings.Count(expression, `-`)
	multipleTime := operations[2] * strings.Count(expression, `*`)
	divisionTime := operations[3] * strings.Count(expression, `/`)

	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		err = errors.New(`failed to convert to expression: ` + err.Error())
		return
	}
	result, err = expr.Evaluate(nil)
	if err != nil {
		err = errors.New(`failed to calculate the expression: ` + err.Error())
		return
	}

	timing := time.Duration(plusTime+minusTime+multipleTime+divisionTime) * time.Millisecond
	if timing < time.Millisecond {
		timing = time.Millisecond
	}

	<-time.After(timing)
	return
}

func (h *handler) previousExpressions() {
	db, err := sql.Open(storage.DriverName, storage.DataSourceName)
	if err != nil {
		log.Println(`database connection error when checking expressions:`, err)
		return
	}
	defer db.Close()

	data, err := db.Query(unfinishedExpressions)
	if err != nil {
		log.Println(`error getting information about unfinished expressions:`, err)
	}

	for data.Next() {
		var id, expression string
		if err := data.Scan(&id, &expression); err != nil {
			log.Println(`error loading unfinished expression:`, err)
		}

		j := new(types.Job)
		j.ID = id
		j.Expression = expression

		h.ch <- j
	}
	if err != nil {
		log.Println(`unfinished expression error:`, err)
		return
	}
}