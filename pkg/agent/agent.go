package agent

import (
	"net/http"

	handlers "github.com/aivaki/calculator/internal/handlers/agent"
	"github.com/rs/cors"
)

var (

)

type Agent struct {
	mux http.Handler
	handlers handlers.Handlers
}

func New() (a *Agent) {
	a = new(Agent)
	a.handlers = handlers.New()
	
	c := cors.Default()
	a.mux = c.Handler(http.DefaultServeMux)
	return
}

func (a *Agent) Run() {
	a.newHandlers()
	a.listenAndServe()
}

func (a *Agent) newHandlers() {
	http.HandleFunc(`/status`, a.handlers.Status)
	http.HandleFunc(`/ws`, a.handlers.Connection)
}

func (a *Agent) listenAndServe() {
	go http.ListenAndServe(`:8080`, a.mux)
}