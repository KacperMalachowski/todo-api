package db_test

import (
	"testing"

	"github.com/KacperMalachowski/todo-api/internal/db"
	"github.com/KacperMalachowski/todo-api/internal/todos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//Test Adding a task with an empty ID
func TestAddTaskWithEmptyID(t *testing.T) {
	task := &todos.Task{
		Title: "Test Task",
	}
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.Error(t, err)
	assert.Equal(t, db.ErrInvalidTask, err)
}

// Test Adding a task with an empty title
func TestAddTaskWithEmptyTitle(t *testing.T) {
	task := &todos.Task{
		ID: "task1",
	}
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.Error(t, err)
	assert.Equal(t, db.ErrInvalidTask, err)
}

// Test Adding a nil task
func TestAddNilTask(t *testing.T) {
	var task *todos.Task
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.Error(t, err)
	assert.Equal(t, db.ErrInvalidTask, err)
}

// Test Adding a valid task
func TestAddValidTask(t *testing.T) {
	task := todos.NewTask("task1", "Test Task")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.NoError(t, err)
}

// Test Adding a duplicate task
func TestAddDuplicateTask(t *testing.T) {
	task := todos.NewTask("task1", "Test Task")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.NoError(t, err)

	err = database.Add(task)
	assert.Error(t, err)
	assert.Equal(t, db.ErrDuplicateTask, err)
}

func TestGetNonExistentTask(t *testing.T) {
	database := db.NewInMemory()
	require.NotNil(t, database)

	task, err := database.Get("nonexistent")
	assert.Nil(t, task)
	assert.Error(t, err)
	assert.Equal(t, db.ErrNotFound, err)
}

func TestGetExistingTask(t *testing.T) {
	task := todos.NewTask("task1", "Test Task")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.NoError(t, err)

	id := task.ID
	retrievedTask, err := database.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, task, retrievedTask)
}

func TestUpdateNonExistentTask(t *testing.T) {
	task := todos.NewTask("task1", "Test Task")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Update("nonexistent", task)
	assert.Error(t, err)
	assert.Equal(t, db.ErrNotFound, err)
}

func TestUpdateExistingTask(t *testing.T) {
	task := todos.NewTask("task1", "Test Task")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.NoError(t, err)

	updatedTask := todos.NewTask("task1", "Updated Task")
	updatedTask.ID = task.ID

	err = database.Update(task.ID, updatedTask)
	assert.NoError(t, err)

	retrievedTask, err := database.Get(task.ID)
	assert.NoError(t, err)
	assert.Equal(t, updatedTask.Title, retrievedTask.Title)
}

func TestUpdateTaskWithEmptyID(t *testing.T) {
	task := &todos.Task{
		Title: "Test Task",
	}
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Update("", task)
	assert.Error(t, err)
	assert.Equal(t, db.ErrInvalidTask, err)
}

func TestDeleteNonExistentTask(t *testing.T) {
	database := db.NewInMemory()
	require.NotNil(t, database)

	database.Delete("nonexistent")
}

func TestDeleteExistingTask(t *testing.T) {
	task := todos.NewTask("task1", "Test Task")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.NoError(t, err)

	id := task.ID
	database.Delete(id)

	retrievedTask, err := database.Get(id)
	assert.Nil(t, retrievedTask)
	assert.Error(t, err)
	assert.Equal(t, db.ErrNotFound, err)
}

func TestListTasksWhenEmpty(t *testing.T) {
	database := db.NewInMemory()
	require.NotNil(t, database)

	tasks, err := database.List()
	assert.Nil(t, tasks)
	assert.Error(t, err)
	assert.Equal(t, db.ErrNotFound, err)
}

func TestListTasksWhenNotEmpty(t *testing.T) {
	task1 := todos.NewTask("task1", "Test Task 1")
	task2 := todos.NewTask("task2", "Test Task 2")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task1)
	assert.NoError(t, err)

	err = database.Add(task2)
	assert.NoError(t, err)

	tasks, err := database.List()
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Contains(t, tasks, task1)
	assert.Contains(t, tasks, task2)
}

func TestListTasksWithSingleTask(t *testing.T) {
	task := todos.NewTask("task1", "Test Task")
	database := db.NewInMemory()
	require.NotNil(t, database)

	err := database.Add(task)
	assert.NoError(t, err)

	tasks, err := database.List()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, task, tasks[0])
}
