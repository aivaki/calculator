package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aivaki/calculator/internal/storage"
	"github.com/aivaki/calculator/internal/types"
	"github.com/google/uuid"
)

func (h *handler) Connection(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(storage.DriverName, storage.DataSourceName)
	if err != nil {
		log.Println(`database connection error in web sockets:`, err)
		return
	}
	defer db.Close()

	connection, err := h.u.Upgrade(w, r, nil)
	if err != nil {
		log.Println(`socket connection error:`, err)
		return
	}
	defer connection.Close()

	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			return
		}

		data := make(map[string]string)
		if err = json.Unmarshal(message, &data); err != nil {
			log.Println(`failed to load data:`, err)
			continue
		}

		expression, ok1 := data[`expression`]
		keyForResult, ok2 := data[`getresult`]
		if !ok1 && !ok2 {
			log.Println(`information about the expression was not found in the request:`, err)
			continue
		}

		if ok2 {
			var result, expression string

			data := db.QueryRow(`SELECT result, expression FROM Expressions WHERE key = ?`, keyForResult)
			data.Scan(&result, &expression)

			if len(result) > 0 {
				response := map[string]any{
					`result`:     result,
					`expression`: expression,
					`id`:         keyForResult,
				}
				if err = connection.WriteJSON(response); err != nil {
					log.Println(`failed to send expression response data:`, err)
				}
			}
			continue
		}

		idForExpression := uuid.New().String()
		h.ch <- &types.Job{ID: idForExpression, Expression: expression}

		query := fmt.Sprintf(`INSERT INTO Expressions (key, expression, status, error_message) VALUES ('%s', '%s', 'in processing', 'nil')`, idForExpression, expression)

		if _, err = db.Exec(query); err != nil {
			log.Println(`failed to add a new expression to the database:`, err)
			return
		}

		response := map[string]any{
			`status`:     `in processing`,
			`expression`: expression,
			`id`:         idForExpression,
		}
		if err = connection.WriteJSON(response); err != nil {
			log.Println(err)
			continue
		}
	}
}