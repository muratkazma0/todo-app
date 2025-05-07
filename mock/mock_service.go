package mock

import (
	"todoapp/entity"
)

type MockService struct {
	todoModel     *entity.TodoModel
	userModel     *entity.UserModel
	todoItemModel *entity.TodoItemModel
}

func NewMockService() *MockService {
	todoModel := entity.NewTodoModel()
	userModel := entity.NewUserModel()
	todoItemModel := entity.NewTodoItemModel()

	service := &MockService{
		todoModel:     todoModel,
		userModel:     userModel,
		todoItemModel: todoItemModel,
	}

	service.createMockData()

	return service
}

func (s *MockService) GetTodoModel() *entity.TodoModel {
	return s.todoModel
}

func (s *MockService) GetUserModel() *entity.UserModel {
	return s.userModel
}

func (s *MockService) GetTodoItemModel() *entity.TodoItemModel {
	return s.todoItemModel
}

func (s *MockService) createMockData() {
	adminUser := &entity.User{
		Username: "admin",
		Password: "admin123",
		Role:     "admin",
	}
	s.userModel.Create(adminUser)

	normalUser := &entity.User{
		Username: "user",
		Password: "user123",
		Role:     "user",
	}
	s.userModel.Create(normalUser)

	adminTodo1 := &entity.Todo{
		Title:       "Admin Todo 1",
		Description: "Admin's first todo",
		UserID:      1,
	}
	s.todoModel.Create(adminTodo1)

	adminTodo2 := &entity.Todo{
		Title:       "Admin Todo 2",
		Description: "Admin's second todo",
		UserID:      1,
	}
	s.todoModel.Create(adminTodo2)

	userTodo1 := &entity.Todo{
		Title:       "User Todo 1",
		Description: "User's first todo",
		UserID:      2,
	}
	s.todoModel.Create(userTodo1)

	userTodo2 := &entity.Todo{
		Title:       "User Todo 2",
		Description: "User's second todo",
		UserID:      2,
	}
	s.todoModel.Create(userTodo2)

	adminTodoItem1 := &entity.TodoItem{
		Title:       "Admin Todo Item 1",
		Description: "First item for admin's first todo",
		TodoID:      1,
		UserID:      1,
	}
	s.todoItemModel.Create(adminTodoItem1)

	adminTodoItem2 := &entity.TodoItem{
		Title:       "Admin Todo Item 2",
		Description: "Second item for admin's first todo",
		TodoID:      1,
		UserID:      1,
	}
	s.todoItemModel.Create(adminTodoItem2)

	userTodoItem1 := &entity.TodoItem{
		Title:       "User Todo Item 1",
		Description: "First item for user's first todo",
		TodoID:      3,
		UserID:      2,
	}
	s.todoItemModel.Create(userTodoItem1)

	userTodoItem2 := &entity.TodoItem{
		Title:       "User Todo Item 2",
		Description: "Second item for user's first todo",
		TodoID:      3,
		UserID:      2,
	}
	s.todoItemModel.Create(userTodoItem2)
}
