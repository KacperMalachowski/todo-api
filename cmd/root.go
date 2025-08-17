package cmd

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

type Logger interface {
	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
}

var rootCmd = &cobra.Command{
	Use: "todo-api",
	Short: "A simple Todo API",
	Long: `A simple Todo API built with Go and Cobra.
This API allows you to manage your todo items with basic CRUD operations.`,
	RunE: run,
}



func run(cmd *cobra.Command, args []string) error {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Welcome to the Todo API!"))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}).Methods("GET")

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
