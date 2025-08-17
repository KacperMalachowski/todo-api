package todos_test

import (
	"testing"

	"github.com/KacperMalachowski/todo-api/internal/todos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test MarkAsCompleted method
func TestMarkAsCompleted(t *testing.T) {
	task := todos.NewTask("Test Task", "This is a test task")
	require.NotNil(t, task)
	require.False(t, task.IsCompleted(), "Task should not be completed initially")

	task.MarkAsCompleted()
	assert.True(t, task.IsCompleted(), "Task should be marked as completed")
}
