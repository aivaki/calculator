package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aivaki/calculator/internal/storage"
)

const (
	updateOperations = `
	UPDATE Operations 
	SET execution_time = ? 
	WHERE operation_type = ?
	`
)

func (h handler) Saving(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(storage.DriverName, storage.DataSourceName)
	if err != nil {
		log.Println(`database connection error when saving operations:`, err)
		return
	}
	defer db.Close()

	operations := map[string]string{
		r.FormValue(`plus`): `+`,
		r.FormValue(`minus`): `-`,
		r.FormValue(`multiple`): `*`,
		r.FormValue(`division`): `/`,
	}
	for k, v := range operations {
		if 	_, err = db.Exec(updateOperations, k, v); err != nil {
			log.Printf("update error (%v): %v\n", v, err)
			return
		}
	}

	http.Redirect(w, r, `/settings`, http.StatusMovedPermanently)
}