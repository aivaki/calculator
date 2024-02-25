package orchestrator

import (
	"net/http"

	handlers "github.com/aivaki/calculator/internal/handlers/orchestrator"
)

type Orchestrator struct {
	handlers handlers.Handlers
	mux *http.ServeMux
}

func New() (o *Orchestrator) {
	o = new(Orchestrator)
	o.handlers = handlers.New()
	o.mux = http.NewServeMux()
	return 
}

func (o *Orchestrator) Run() {
	o.newHandlers()
	o.listenAndServe()
}

func (o *Orchestrator) newHandlers() {
	o.mux.HandleFunc(`/`, o.handlers.Main)
	o.mux.HandleFunc(`/history/`, o.handlers.History)
	o.mux.HandleFunc(`/saving/`, o.handlers.Saving)
	o.mux.HandleFunc(`/settings/`, o.handlers.Settings)
	o.mux.HandleFunc(`/status/`, o.handlers.Status)
}

func (o *Orchestrator) listenAndServe() {
	http.ListenAndServe(`:5500`, o.mux)
}