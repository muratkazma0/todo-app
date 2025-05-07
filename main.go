package main

import (
	"log"
	"todoapp/controllers"
	"todoapp/entity"
	"todoapp/routes"
)

func createDefaultUser(userModel *entity.UserModel, username, password, role string) (*entity.User, error) {
	user, err := userModel.GetByUsername(username)
	if err != nil {
		user = &entity.User{
			Username: username,
			Role:     role,
			Password: password,
		}
		if err := userModel.Create(user); err != nil {
			return nil, err
		}
		log.Printf("Default %s user created successfully", role)
	}
	return user, nil
}

func createDefaultTodo(todoModel *entity.TodoModel, title, description string, userID int) (*entity.Todo, error) {
	todo := &entity.Todo{
		Title:       title,
		Description: description,
		UserID:      userID,
	}
	if err := todoModel.Create(todo); err != nil {
		return nil, err
	}
	log.Printf("Todo '%s' created successfully", title)
	return todo, nil
}

func createDefaultTodoItem(todoItemModel *entity.TodoItemModel, title, description string, todoID, userID int) error {
	todoItem := &entity.TodoItem{
		Title:       title,
		Description: description,
		TodoID:      todoID,
		UserID:      userID,
	}
	if err := todoItemModel.Create(todoItem); err != nil {
		return err
	}
	log.Printf("Todo item '%s' created successfully", title)
	return nil
}

func initializeDefaultData(userModel *entity.UserModel, todoModel *entity.TodoModel, todoItemModel *entity.TodoItemModel) error {
	adminUser, err := createDefaultUser(userModel, "admin", "admin123", "admin")
	if err != nil {
		return err
	}

	adminTodo, err := createDefaultTodo(todoModel, "Admin Todo", "This is admin's todo", adminUser.ID)
	if err != nil {
		return err
	}

	if err := createDefaultTodoItem(todoItemModel, "Admin Todo Item", "This is admin's todo item", adminTodo.ID, adminUser.ID); err != nil {
		return err
	}

	normalUser, err := createDefaultUser(userModel, "user", "user123", "user")
	if err != nil {
		return err
	}

	normalTodo, err := createDefaultTodo(todoModel, "User Todo", "This is normal user's todo", normalUser.ID)
	if err != nil {
		return err
	}

	if err := createDefaultTodoItem(todoItemModel, "User Todo Item", "This is normal user's todo item", normalTodo.ID, normalUser.ID); err != nil {
		return err
	}

	return nil
}

func main() {
	userModel := entity.NewUserModel()
	todoModel := entity.NewTodoModel()
	todoItemModel := entity.NewTodoItemModel()

	if err := initializeDefaultData(userModel, todoModel, todoItemModel); err != nil {
		log.Fatalf("Failed to initialize default data: %v", err)
	}

	authController := controllers.NewAuthController(userModel)
	userController := controllers.NewUserController(userModel)
	todoController := controllers.NewTodoController(todoModel)
	todoItemController := controllers.NewTodoItemController(todoItemModel, todoModel)

	r := routes.SetupRoutes(
		authController,
		userController,
		todoController,
		todoItemController,
	)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}