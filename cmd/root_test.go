package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KacperMalachowski/todo-api/internal/db"
	"github.com/KacperMalachowski/todo-api/internal/todos"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test listing all tasks
func TestListAllTasks(t *testing.T) {
	task1 := todos.NewTask("task1", "Test Task 1")
	task2 := todos.NewTask("task2", "Test Task 2")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task1)
	assert.NoError(t, err)

	err = database.Add(task2)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/tasks", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks", NewServer(nil, database).AllTasks).Methods("GET")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	responseBody, err := io.ReadAll(rr.Body)
	require.NoError(t, err)

	var tasks []todos.Task
	err = json.Unmarshal(responseBody, &tasks)
	require.NoError(t, err)

	expectedTasks := []todos.Task{*task1, *task2}
	assert.ElementsMatch(t, expectedTasks, tasks)
}

