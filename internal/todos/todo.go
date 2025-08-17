package todos

import "github.com/google/uuid"

type Task struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Completed bool `json:"completed"`
}

func (t *Task) IsCompleted() bool {
	return t.Completed
}

func (t *Task) MarkAsCompleted() {
	t.Completed = true
}

func NewTask(title, description string) *Task {
	id := uuid.New().String() // Assuming you have a way to generate unique IDs
	return &Task{
		ID: id,
		Title: title,
		Description: description,
		Completed: false,
	}
}

