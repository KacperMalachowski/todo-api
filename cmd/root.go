package cmd

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/KacperMalachowski/todo-api/internal/db"
	"github.com/KacperMalachowski/todo-api/internal/todos"
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
	if err != nil {
		s.logger.Printf("Error marshalling tasks", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		s.logger.Printf("Error writing response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *server) AddTask(w http.ResponseWriter, r *http.Request) {
	var task todos.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		s.logger.Printf("Error decoding task: %v", err)
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	if err := s.database.Add(&task); err != nil {
		s.logger.Printf("Error adding task: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *server) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	task, err := s.database.Get(id)
	if err != nil {
		s.logger.Printf("Error getting task: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		s.logger.Printf("Error marshalling task: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		s.logger.Printf("Error writing response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *server) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var task todos.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		s.logger.Printf("Error decoding task: %v", err)
		http.Error(w, "Invalid task data", http.StatusBadRequest)
		return
	}

	if err := s.database.Update(id, &task); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			s.logger.Printf("Task not found: %v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		s.logger.Printf("Error updating task: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	s.database.Delete(id)

	w.WriteHeader(http.StatusNoContent)
}

func run(cmd *cobra.Command, args []string) error {
	r := mux.NewRouter()
	database := db.NewInMemory()
	logger := log.New(log.Writer(), "todo-api: ", log.LstdFlags)
	server := NewServer(logger, database)

	r.HandleFunc("/tasks", server.AllTasks).Methods("GET")
	r.HandleFunc("/tasks", server.AddTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", server.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", server.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", server.DeleteTask).Methods("DELETE")

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
