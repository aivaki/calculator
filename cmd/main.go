package main

import (
	"log"

	"github.com/aivaki/calculator/internal/storage"
	"github.com/aivaki/calculator/pkg/agent"
	"github.com/aivaki/calculator/pkg/orchestrator"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := storage.New()
	if err != nil {
		log.Fatalln(err)
		return
	}
	db.Close()
	
	a := agent.New()
	a.Run()

	o := orchestrator.New()
	o.Run()
}