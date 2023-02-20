package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	fmt.Println(server.ListenAndServe())
}
