package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	"github.com/aivaki/calculator/internal/storage"
	"github.com/aivaki/calculator/internal/types"
)

const (
	processRows = `SELECT status, expression, created_at FROM Expressions WHERE status = 'in processing'`

	finishedRows = `SELECT status, result, expression, created_at, completed_at FROM Expressions WHERE status = 'completed'`

	errorRows = `SELECT status, expression, created_at, completed_at FROM Expressions WHERE status = 'error'`
)

func (h handler) History(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		historyPath,
		headerPath,
		headPath,
	)
	if err != nil {
		log.Println(`error parsing expression history template:`, err)
		return
	}

	db, err := sql.Open(storage.DriverName, storage.DataSourceName)
	if err != nil {
		log.Println(`database connection error in expression history:`, err)
		return
	}
	defer db.Close()

	var expressions []types.Expression

	pr, err := db.Query(processRows)
	if err != nil {
		log.Println(`error retrieving unprocessed expressions from the database:`, err)
		return
	}
	for pr.Next() {
		var expression types.Expression
		if err := pr.Scan(&expression.Status, &expression.Expression, &expression.CreatedAt); err != nil {
			log.Println(`error loading unprocessed expression:`, err)
			return
		}
		expressions = append(expressions, expression)
	}
	if err := pr.Err(); err != nil {
		log.Println(`unprocessed expression error:`, err)
		return
	}
	pr.Close()

	fr, err := db.Query(finishedRows)
	if err != nil {
		log.Println(`error retrieving executed expressions from the database:`, err)
		return
	}
	for fr.Next() {
		var expression types.Expression
		if err := fr.Scan(&expression.Status, &expression.Result, &expression.Expression, &expression.CreatedAt, &expression.CompletedAt); err != nil {
			log.Println(`error loading executed expression:`, err)
			return
		}
		expressions = append(expressions, expression)
	}
	if err := fr.Err(); err != nil {
		log.Println(`expression executed error:`, err)
		return
	}
	fr.Close()

	er, err := db.Query(errorRows)
	if err != nil {
		log.Println(`error retrieving invalid expressions from the database:`, err)
		return
	}
	for er.Next() {
		var expression types.Expression
		if err := er.Scan(&expression.Status, &expression.Expression, &expression.CreatedAt, &expression.CompletedAt); err != nil {
			log.Println(`invalid expression error:`, err)
			return
		}
		expressions = append(expressions, expression)
	}
	if err := er.Err(); err != nil {
		log.Println(`invalid expression error:`, err)
		return
	}
	er.Close()

	t.ExecuteTemplate(w, `history`, expressions)
}