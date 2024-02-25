package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/aivaki/calculator/internal/types"
)

const (
	httpStatus = `http://localhost:8080/status`
)

func (h handler) Status(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		statusPath,
		headerPath,
		headPath,
	)
	if err != nil {
		log.Println(`error loading worker status template:`, err)
		return
	}

	response, err := http.Get(httpStatus)
	if err != nil {
		log.Println(`failed to get workers information:`, err)
		return
	}
	defer response.Body.Close()

	var ctx types.StatusData
	if err := json.NewDecoder(response.Body).Decode(&ctx); err != nil {
		log.Println(`failed to decode:`, err)
		return
	}

	t.ExecuteTemplate(w, `status`, ctx)
}