package main

import (
	"fmt"
	"net/http"
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

	fmt.Println(server.ListenAndServe())
}
