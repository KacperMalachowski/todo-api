package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KacperMalachowski/todo-api/internal/db"
	"github.com/KacperMalachowski/todo-api/internal/todos"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var logger = log.New(io.Discard, "", 0)

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
	router.HandleFunc("/tasks", NewServer(logger, database).AllTasks).Methods("GET")

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

func TestListAllTasksEmpty(t *testing.T) {
	database := db.NewInMemory()
	require.NotNil(t, database)

	req, err := http.NewRequest("GET", "/tasks", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks", NewServer(logger, database).AllTasks).Methods("GET")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	responseBody, err := io.ReadAll(rr.Body)
	require.NoError(t, err)

	assert.Equal(t, "task not found\n", string(responseBody))
}

// Test adding a new task
func TestAddNewTask(t *testing.T) {
	task := todos.NewTask("task1", "Test Task")
	database := db.NewInMemory()
	require.NotNil(t, database)

	taskJSON, err := json.Marshal(task)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/tasks", io.NopCloser(bytes.NewBuffer(taskJSON)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks", NewServer(logger, database).AddTask).Methods("POST")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestUpdateTask(t *testing.T) {
	task := todos.NewTask("task1", "Test Task")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.NoError(t, err)

	updatedTask := todos.NewTask("task1", "Updated Task")
	updatedTaskJSON, err := json.Marshal(updatedTask)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", fmt.Sprintf("/tasks/%s", task.ID), bytes.NewReader(updatedTaskJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", NewServer(logger, database).UpdateTask).Methods("PUT")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify the task was updated
	retrievedTask, err := database.Get(task.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Task", retrievedTask.Description)
}

func TestUpdateNonExistentTask(t *testing.T) {
	updatedTask := todos.NewTask("task1", "Updated Task")
	updatedTaskJSON, err := json.Marshal(updatedTask)
	require.NoError(t, err)

	database := db.NewInMemory()
	require.NotNil(t, database)

	req, err := http.NewRequest("PUT", "/tasks/nonexistent", io.NopCloser(bytes.NewBuffer(updatedTaskJSON)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", NewServer(logger, database).UpdateTask).Methods("PUT")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	responseBody, err := io.ReadAll(rr.Body)
	require.NoError(t, err)

	assert.Equal(t, "task not found\n", string(responseBody))
}

func TestUpdateTaskWithInvalidData(t *testing.T) {
	invalidTask := `{"id": "", "title": ""}` // Invalid task with empty ID and title
	database := db.NewInMemory()
	require.NotNil(t, database)

	req, err := http.NewRequest("PUT", "/tasks/task1", io.NopCloser(bytes.NewBufferString(invalidTask)))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", NewServer(logger, database).UpdateTask).Methods("PUT")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	responseBody, err := io.ReadAll(rr.Body)
	require.NoError(t, err)

	assert.Equal(t, db.ErrInvalidTask.Error(), strings.Trim(string(responseBody), "\n"))
}

func TestDeleteTask(t *testing.T) {
	task := todos.NewTask("task1", "Test Task")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.NoError(t, err)

	req, err := http.NewRequest("DELETE", "/tasks/task1", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/tasks/{id}", NewServer(logger, database).DeleteTask).Methods("DELETE")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify the task was deleted
	_, err = database.Get("task1")
	assert.Error(t, err)
	assert.Equal(t, db.ErrNotFound, err)
}

