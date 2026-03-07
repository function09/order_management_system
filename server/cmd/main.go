package main

import (
	"log"
	"net/http"

	"github.com/function09/order_management_system/server/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	cfg := config.LoadConfig()

	mux := http.NewServeMux()

	log.Printf("Listening on port %s...", cfg.Port)

	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		log.Fatal(err)
	}

}
