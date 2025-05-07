package entity

import (
	"errors"
	"sync"
	"time"
)

type TodoItem struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	TodoID      int        `json:"todo_id"`
	UserID      int        `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type TodoItemModel struct {
	sync.RWMutex
	items  map[int]*TodoItem
	nextID int
}

func NewTodoItemModel() *TodoItemModel {
	return &TodoItemModel{
		items:  make(map[int]*TodoItem),
		nextID: 1,
	}
}

func (m *TodoItemModel) Create(item *TodoItem) error {
	m.Lock()
	defer m.Unlock()

	item.ID = m.nextID
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	item.DeletedAt = nil

	m.items[item.ID] = item
	m.nextID++

	return nil
}

func (m *TodoItemModel) GetByID(id int) (*TodoItem, error) {
	m.RLock()
	defer m.RUnlock()

	item, exists := m.items[id]
	if !exists || item.DeletedAt != nil {
		return nil, errors.New("item not found")
	}

	return item, nil
}

func (m *TodoItemModel) GetByIDWithDeleted(id int) (*TodoItem, error) {
	m.RLock()
	defer m.RUnlock()

	item, exists := m.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	return item, nil
}

func (m *TodoItemModel) GetByTodoID(todoID int) []*TodoItem {
	m.RLock()
	defer m.RUnlock()

	var todoItems []*TodoItem
	for _, item := range m.items {
		if item.TodoID == todoID && item.DeletedAt == nil {
			todoItems = append(todoItems, item)
		}
	}

	return todoItems
}

func (m *TodoItemModel) GetByTodoIDWithDeleted(todoID int) []*TodoItem {
	m.RLock()
	defer m.RUnlock()

	var todoItems []*TodoItem
	for _, item := range m.items {
		if item.TodoID == todoID {
			todoItems = append(todoItems, item)
		}
	}

	return todoItems
}

func (m *TodoItemModel) GetAll() []*TodoItem {
	m.RLock()
	defer m.RUnlock()

	var activeItems []*TodoItem
	for _, item := range m.items {
		if item.DeletedAt == nil {
			activeItems = append(activeItems, item)
		}
	}

	return activeItems
}

func (m *TodoItemModel) GetAllWithDeleted() []*TodoItem {
	m.RLock()
	defer m.RUnlock()

	var allItems []*TodoItem
	for _, item := range m.items {
		allItems = append(allItems, item)
	}

	return allItems
}

func (m *TodoItemModel) Update(item *TodoItem) error {
	m.Lock()
	defer m.Unlock()

	existing, exists := m.items[item.ID]
	if !exists || existing.DeletedAt != nil {
		return errors.New("item not found")
	}

	item.UpdatedAt = time.Now()
	m.items[item.ID] = item

	return nil
}

func (m *TodoItemModel) Delete(id int) error {
	m.Lock()
	defer m.Unlock()

	item, exists := m.items[id]
	if !exists || item.DeletedAt != nil {
		return errors.New("item not found")
	}

	now := time.Now()
	item.DeletedAt = &now
	return nil
}
