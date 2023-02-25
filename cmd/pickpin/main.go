package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/users"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	_, err := db.New(logger)
	if err != nil {
		return
	}

	mux := httprouter.New()

	users.RegisterHandlers(mux, logger)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	logger.Println("Starting server...")
	err = server.ListenAndServe()
	if err != nil {
		logger.Printf("Failed to start server, %s", err.Error())
	}
}
