package db

import (
	"fmt"

	"github.com/KacperMalachowski/todo-api/internal/todos"
)

var (
	ErrNotFound = fmt.Errorf("task not found")
	ErrInvalidID = fmt.Errorf("invalid task ID")
	ErrInvalidTask = fmt.Errorf("invalid task data")
	ErrInternal = fmt.Errorf("internal server error")
	ErrDuplicateTask = fmt.Errorf("task with the same ID already exists")
)

type inMemory struct {
	tasks map[string]*todos.Task
}

func NewInMemory() *inMemory {
	return &inMemory{
		tasks: make(map[string]*todos.Task),
	}
}

func (db *inMemory) Get(id string) (*todos.Task, error) {
	task, exists := db.tasks[id]
	if !exists {
		return nil,	ErrNotFound
	}
	return task, nil
}

func (db *inMemory) Add(task *todos.Task) error {
	if task == nil || task.ID == "" || task.Title == "" {
		return	ErrInvalidTask
	}
	if _, exists := db.tasks[task.ID]; exists {
		return	ErrDuplicateTask
	}
	db.tasks[task.ID] = task
	return nil
}

func (db *inMemory) Update(id string, task *todos.Task) error {
	if task == nil || task.ID == "" || task.Title == "" {
		return	ErrInvalidTask
	}
	if _, exists := db.tasks[id]; !exists {
		return	ErrNotFound
	}
	db.tasks[id] = task
	return nil
}

func (db *inMemory) Delete(id string) {
	// Treat non-existent ID as an valid operation - task is deleted
	delete(db.tasks, id)
}
