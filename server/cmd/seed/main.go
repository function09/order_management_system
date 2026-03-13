package main

import (
	"log"

	"github.com/function09/order_management_system/server/config"
	"github.com/function09/order_management_system/server/db"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal(err)
	}

	cfg := config.LoadConfig()

	database := db.ConnectToDB(cfg.DatabaseURL)
	defer database.Close()

	categories := []string{"Electronics", "Clothing", "Food & Beverage", "Home & Garden", "Sports"}

	_, err := database.Exec("TRUNCATE TABLE categories RESTART IDENTITY CASCADE")

	if err != nil {
		log.Printf("Error clearing table: %s", err)
	}
	for _, e := range categories {
		_, err := database.Exec("INSERT INTO categories (category) VALUES ($1);", e)

		if err != nil {
			log.Printf("Error inserting categories %s", err)
		}
	}
}
