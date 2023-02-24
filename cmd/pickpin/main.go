package main

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/db"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World!"))
		if err != nil {
			fmt.Println(err)
		}
	})

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	db.Connect()

	fmt.Println(server.ListenAndServe())
}
