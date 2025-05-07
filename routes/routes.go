package routes

import (
	"todoapp/controllers"
	"todoapp/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	authController *controllers.AuthController,
	userController *controllers.UserController,
	todoController *controllers.TodoController,
	todoItemController *controllers.TodoItemController,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	r.POST("/login", authController.Login)

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("", userController.Create)
			users.GET("", middleware.AuthMiddleware(), middleware.AdminOnly(), userController.GetAll)
			users.GET("/:id", middleware.AuthMiddleware(), userController.GetByID)
			users.GET("/username/:username", middleware.AuthMiddleware(), userController.GetByUsername)
			users.PUT("/:id", middleware.AuthMiddleware(), userController.Update)
			users.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminOnly(), userController.Delete)
		}

		// Todo routes
		todos := api.Group("/todos")
		todos.Use(middleware.AuthMiddleware())
		{
			// Todo item routes
			items := todos.Group("/items")
			{
				items.POST("/:todo_id", todoItemController.Create)
				items.GET("/:todo_id", todoItemController.GetByTodoID)
				items.PUT("/:todo_id/:item_id", todoItemController.Update)
				items.DELETE("/:todo_id/:item_id", todoItemController.Delete)
			}

			// Todo routes
			todos.POST("", todoController.Create)
			todos.GET("", todoController.GetAll)
			todos.GET("/:id", todoController.GetByID)
			todos.PUT("/:id", todoController.Update)
			todos.DELETE("/:id", todoController.Delete)
		}
	}

	return r
}
