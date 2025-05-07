package controllers

import (
	"net/http"
	"strconv"
	"todoapp/entity"

	"github.com/gin-gonic/gin"
)

type TodoController struct {
	todoModel *entity.TodoModel
}

func NewTodoController(todoModel *entity.TodoModel) *TodoController {
	return &TodoController{
		todoModel: todoModel,
	}
}

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (c *TodoController) Create(ctx *gin.Context) {
	var req CreateTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	todo := &entity.Todo{
		Title:         req.Title,
		Description:   req.Description,
		UserID:        userID.(int),
		CompletionPct: 0,
	}

	if err := c.todoModel.Create(todo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, todo)
}

func (c *TodoController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRole, _ := ctx.Get("user_role")
	var todo *entity.Todo
	var err2 error

	if userRole == "admin" {
		todo, err2 = c.todoModel.GetByIDWithDeleted(id)
	} else {
		todo, err2 = c.todoModel.GetByID(id)
	}

	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	if todo.UserID != userID.(int) && userRole != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (c *TodoController) GetAll(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRole, _ := ctx.Get("user_role")

	var todos []*entity.Todo
	if userRole == "admin" {
		todos = c.todoModel.GetAllWithDeleted()
	} else {
		todos = c.todoModel.GetByUserID(userID.(int))
	}

	ctx.JSON(http.StatusOK, todos)
}

func (c *TodoController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRole, _ := ctx.Get("user_role")

	var todo *entity.Todo
	var err2 error
	if userRole == "admin" {
		todo, err2 = c.todoModel.GetByIDWithDeleted(id)
	} else {
		todo, err2 = c.todoModel.GetByID(id)
	}

	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	if todo.UserID != userID.(int) && userRole != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req UpdateTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.Title = req.Title
	todo.Description = req.Description

	if err := c.todoModel.Update(todo); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (c *TodoController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRole, _ := ctx.Get("user_role")

	var todo *entity.Todo
	var err2 error
	if userRole == "admin" {
		todo, err2 = c.todoModel.GetByIDWithDeleted(id)
	} else {
		todo, err2 = c.todoModel.GetByID(id)
	}

	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	if todo.UserID != userID.(int) && userRole != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := c.todoModel.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "todo deleted"})
}
