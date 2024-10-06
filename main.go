package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Server struct {
	logger *slog.Logger
}

func (s *Server) webhookHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Received request", "method", r.Method, "url", r.URL.Path)
	
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	err := HandlePayload(r.Body)
	if err != nil {
		s.logger.Error("Failed to handle payload", "error", err.Error())
		return
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	server := &Server{logger: logger}

	http.HandleFunc("/webhook", server.webhookHandler)
	logger.Info("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if (err != nil) {
		logger.Error(err.Error())
	}
}