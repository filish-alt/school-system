package main

import (
	"log"
	"net/http"
	"school-exam/internal/initiator"
)

func main() {
	app, err := initiator.Initiate()
	if err != nil {
		log.Fatal(err)
	}
	s := app.Server
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
