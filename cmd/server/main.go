package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pc-06/golangsmtp/internal/config"
	"github.com/pc-06/golangsmtp/internal/handler"
	"github.com/pc-06/golangsmtp/internal/middleware"
	"github.com/pc-06/golangsmtp/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("warning: could not load .env file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	mailSvc := service.NewMailService(cfg)
	mailHandler := handler.NewMailHandler(mailSvc)

	rateLimiter := middleware.NewRateLimiter(5, time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/send-email", rateLimiter.Limit(mailHandler.SendEmail))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("server starting on %s", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
