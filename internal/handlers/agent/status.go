package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aivaki/calculator/internal/storage"
	"github.com/aivaki/calculator/internal/types"
)

const (
	countWorkersEnv = `COUNT_WORKERS`
)

func (h *handler) Status(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(storage.DriverName, storage.DataSourceName)
	if err != nil {
		log.Println(`database connection error in worker status request:`, err)
		return
	}
	defer db.Close()

	h.m.Lock()
	data, err := db.Query(`SELECT expression FROM Expressions WHERE status = 'in processing'`)
	if err != nil {
		log.Println(`error in obtaining information for worker status:`, err)
		return
	}
	h.m.Unlock()

	inProcessWorkers := []string{}

	for data.Next() {
		var expression string
		err := data.Scan(&expression)
		if err != nil {
			log.Println(`error loading information for worker status:`, err)
		}
		inProcessWorkers = append(inProcessWorkers, expression)
	}

	countWorkers := os.Getenv(countWorkersEnv)
	maxWorkers, err := strconv.Atoi(countWorkers)
	if err != nil {
		maxWorkers = 5
	}
	freeWorkers := maxWorkers - len(inProcessWorkers)
	if freeWorkers < 0 {
		freeWorkers = 0
	}

	statusData := types.StatusData{
		FreeWorkers: freeWorkers,
		MaxWorkers:  maxWorkers,
		Expressions: inProcessWorkers,
	}

	jsonData, _ := json.Marshal(statusData)

	w.Header().Set(`Content-Type`, `application/json`)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println(`error writing json about workers:`, err)
		return
	}
}