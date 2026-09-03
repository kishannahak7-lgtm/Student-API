package main

import (
	"GO-API/internal/config"
	"GO-API/internal/config/http/handler/student"
	"GO-API/internal/storage/sqlite"
	"context"

	//"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//load setup
	cfg := config.MustLoad()

	// database setups
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	//setup router
	router := http.NewServeMux()

	router.Handle("/api/students", student.New(storage))
	router.Handle("/api/students/{id}", student.GetByID(storage))
	//setup server

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server", err)
		}
	}()

	<-done

	//gracefull shutdown
	slog.Info("shutting down the server ")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("shutdown successfully")

}
