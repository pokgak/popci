package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

type Payload struct {
	Message string `json:"message"`
}

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

	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		s.logger.Error("Error decoding JSON", "error", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	s.logger.Info("Received payload", "payload", payload)
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