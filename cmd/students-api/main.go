package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ashok2025-eng/students-api/internal/config"
	"github.com/Ashok2025-eng/students-api/internal/http/handlers/student"
	"github.com/Ashok2025-eng/students-api/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.MustLoad()
	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))
	// Setup router
	router := http.NewServeMux() // creates a router(like a traffic controller)

	router.HandleFunc("/api/students", student.New(storage)) //- “Whenever someone visits /api/students, run the student.New() handler.”

	// Setup server
	//creates the actual hhtp server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("Server started ", slog.String("address", cfg.HTTPServer.Addr))
	// ✅ Print server started message
	// fmt.Printf("Server started at http://%s\n", cfg.HTTPServer.Addr)
	done := make(chan os.Signal, 1) // handles shutdown signaals

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//Strt server in Background
	go func() {
		err := server.ListenAndServe() //makes the server start listening for requests.
		if err != nil {
			log.Fatal("failed to start server:", err)
		}
	}()
	<-done

	slog.Info("sutting down the server")

	//gracefull shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
	// Start server

}
