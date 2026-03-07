package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	port := ":8000"

	log.Printf("Listening on port %s...", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}

}
