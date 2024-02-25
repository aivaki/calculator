package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	"github.com/aivaki/calculator/internal/storage"
)

const (
	operationsExecutionTime = `SELECT execution_time FROM Operations`
)

func (h handler) Settings(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		settingsPath,
		headerPath,
		headPath,
	)
	if err != nil {
		log.Println(`template connection error in settings:`, err)
		return
	}

	db, err := sql.Open(storage.DriverName, storage.DataSourceName)
	if err != nil {
		log.Println(`database connection error in settings:`, err)
		return
	}
	defer db.Close()

	data, err := db.Query(operationsExecutionTime)
	if err != nil {
		log.Println(`error when retrieving information about operations from the database:`, err)
		return
	}

	operations := []int{}
	for data.Next() {
		var time int
		if err := data.Scan(&time); err != nil {
			log.Println(`error loading information:`, err)
		}
		operations = append(operations, time)
	}

	ctx := struct {
		Plus,
		Minus,
		Multiple,
		Division int
	}{
		Plus:     operations[0],
		Minus:    operations[1],
		Multiple: operations[2],
		Division: operations[3],
	}

	t.ExecuteTemplate(w, `settings`, ctx)
}