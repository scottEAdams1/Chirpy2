package main

import (
	"fmt"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", handlerReadiness)

	fmt.Printf("Serving on port: %s\n", port)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
