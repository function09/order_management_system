package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/function09/order_management_system/server/config"
	"github.com/function09/order_management_system/server/db"
	"github.com/function09/order_management_system/server/internal/auth"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}

	cfg := config.LoadConfig()

	database := db.ConnectToDB(cfg.DatabaseURL)
	defer database.Close()

	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    cfg.Port,
	}

	store := &auth.Store{DB: database}

	mux.HandleFunc("POST /auth/register", auth.RegisterUserHandler(store))
	mux.HandleFunc("POST /auth/login", auth.LoginUserHandler(store, cfg.JWTSecret))
	mux.HandleFunc("POST /auth/logout", auth.LogOutHandler)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	log.Print("Shutdown complete")

}
