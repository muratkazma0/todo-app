package entity

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Todo struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	UserID        int        `json:"user_id"`
	CompletionPct float64    `json:"completion_pct"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
}

type TodoModel struct {
	sync.RWMutex
	todos  map[int]*Todo
	nextID int
}

func NewTodoModel() *TodoModel {
	return &TodoModel{
		todos:  make(map[int]*Todo),
		nextID: 1,
	}
}

func (m *TodoModel) Create(todo *Todo) error {
	m.Lock()
	defer m.Unlock()

	todo.ID = m.nextID
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	todo.DeletedAt = nil

	m.todos[todo.ID] = todo
	m.nextID++

	return nil
}

func (m *TodoModel) GetByID(id int) (*Todo, error) {
	m.RLock()
	defer m.RUnlock()

	todo, exists := m.todos[id]
	if !exists || todo.DeletedAt != nil {
		return nil, errors.New("todo not found")
	}

	return todo, nil
}

func (m *TodoModel) GetByIDWithDeleted(id int) (*Todo, error) {
	m.RLock()
	defer m.RUnlock()

	todo, exists := m.todos[id]
	if !exists {
		return nil, errors.New("todo not found")
	}

	return todo, nil
}

func (m *TodoModel) GetAll() []*Todo {
	m.RLock()
	defer m.RUnlock()

	var activeTodos []*Todo
	for _, todo := range m.todos {
		if todo.DeletedAt == nil {
			activeTodos = append(activeTodos, todo)
		}
	}

	return activeTodos
}

func (m *TodoModel) GetAllWithDeleted() []*Todo {
	m.RLock()
	defer m.RUnlock()

	var allTodos []*Todo
	fmt.Printf("Total todos in map: %d\n", len(m.todos))
	for id, todo := range m.todos {
		fmt.Printf("Todo ID: %d, Title: %s, DeletedAt: %v\n", id, todo.Title, todo.DeletedAt)
		allTodos = append(allTodos, todo)
	}
	fmt.Printf("Returning %d todos\n", len(allTodos))
	return allTodos
}

func (m *TodoModel) GetByUserID(userID int) []*Todo {
	m.RLock()
	defer m.RUnlock()

	var userTodos []*Todo
	for _, todo := range m.todos {
		if todo.UserID == userID && todo.DeletedAt == nil {
			userTodos = append(userTodos, todo)
		}
	}

	return userTodos
}

func (m *TodoModel) GetByUserIDWithDeleted(userID int) []*Todo {
	m.RLock()
	defer m.RUnlock()

	var userTodos []*Todo
	for _, todo := range m.todos {
		if todo.UserID == userID {
			userTodos = append(userTodos, todo)
		}
	}

	return userTodos
}

func (m *TodoModel) Update(todo *Todo) error {
	m.Lock()
	defer m.Unlock()

	existing, exists := m.todos[todo.ID]
	if !exists || existing.DeletedAt != nil {
		return errors.New("todo not found")
	}

	todo.UpdatedAt = time.Now()
	m.todos[todo.ID] = todo

	return nil
}

func (m *TodoModel) Delete(id int) error {
	m.Lock()
	defer m.Unlock()

	todo, exists := m.todos[id]
	if !exists || todo.DeletedAt != nil {
		return errors.New("todo not found")
	}

	now := time.Now()
	todo.DeletedAt = &now
	return nil
}

func (m *TodoModel) UpdateCompletionPct(todoID int, todoItemModel *TodoItemModel) error {
	m.Lock()
	defer m.Unlock()

	todo, exists := m.todos[todoID]
	if !exists {
		return errors.New("todo not found")
	}

	// Get all items for this todo
	items := todoItemModel.GetByTodoID(todoID)
	if len(items) == 0 {
		todo.CompletionPct = 0
		return nil
	}

	completedCount := 0
	for _, item := range items {
		if item.Completed {
			completedCount++
		}
	}

	todo.CompletionPct = float64(completedCount) / float64(len(items)) * 100
	todo.UpdatedAt = time.Now()

	return nil
}
