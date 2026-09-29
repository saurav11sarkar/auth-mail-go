package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/saurav11sarkar/go/internal/app"
	"github.com/saurav11sarkar/go/internal/config"
	"github.com/saurav11sarkar/go/internal/database"
)

func main() {
	_ = godotenv.Load()
	cfg := config.MustLoad()

	db, err := database.Connect(context.Background(), cfg.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	handler, err := app.NewHandler(db, cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		log.Println("Starting server")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server shutdown:", err)
	}
}
