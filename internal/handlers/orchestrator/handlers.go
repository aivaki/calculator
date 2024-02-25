package handlers

import (
	"log"
	"net/http"
	"text/template"
)

const (
	frontendPath = `../frontend/`
	includesPath = frontendPath + `includes/`

	mainPath = frontendPath + `main.html`
	historyPath = frontendPath + `history.html`
	statusPath = frontendPath + `status.html`
	settingsPath = frontendPath + `settings.html`

	headerPath = includesPath + `header.html`
	headPath = includesPath + `head.html`
)

type handler uintptr

type Handlers interface {
	Main(w http.ResponseWriter, r *http.Request)
	Settings(w http.ResponseWriter, r *http.Request)
	History(w http.ResponseWriter, r *http.Request)
	Saving(w http.ResponseWriter, r *http.Request)
	Status(w http.ResponseWriter, r *http.Request)
}

func New() Handlers {
	h := new(handler)
	return h
}

func (h handler) Main(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		mainPath,
		headerPath,
		headPath,
	)
	if err != nil {
		log.Println(`error loading main template:`, err)
		return
	}
	t.ExecuteTemplate(w, `main`, nil)
}