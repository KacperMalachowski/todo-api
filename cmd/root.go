package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/KacperMalachowski/todo-api/internal/db"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

var rootCmd = &cobra.Command{
	Use: "todo-api",
	Short: "A simple Todo API",
	Long: `A simple Todo API built with Go and Cobra.
This API allows you to manage your todo items with basic CRUD operations.`,
	RunE: run,
}

type server struct {
	logger Logger
	database *db.InMemory
}

func NewServer(logger Logger, database *db.InMemory) *server {
	return &server{
		logger:   logger,
		database: database,
	}
}

func (s *server) AllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.database.List()
	if err != nil {
		s.logger.Printf("Error listing tasks", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(tasks) == 0 {
		http.Error(w, "No tasks found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(tasks)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		s.logger.Printf("Error writing response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func run(cmd *cobra.Command, args []string) error {
	r := mux.NewRouter()
	database := db.NewInMemory()
	logger := log.New(log.Writer(), "todo-api: ", log.LstdFlags)
	server := NewServer(logger, database)

	r.HandleFunc("/tasks", server.AllTasks).Methods("GET")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error executing root command: %w", err)
	}
}
