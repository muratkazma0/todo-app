package controllers

import (
	"net/http"
	"strconv"

	"todoapp/entity"

	"github.com/gin-gonic/gin"
)

type TodoItemController struct {
	todoItemModel *entity.TodoItemModel
	todoModel     *entity.TodoModel
}

func NewTodoItemController(todoItemModel *entity.TodoItemModel, todoModel *entity.TodoModel) *TodoItemController {
	return &TodoItemController{
		todoItemModel: todoItemModel,
		todoModel:     todoModel,
	}
}

type CreateTodoItemRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateTodoItemRequest struct {
	Completed bool `json:"completed"`
}

func (c *TodoItemController) Create(ctx *gin.Context) {
	todoID, err := strconv.Atoi(ctx.Param("todo_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
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
		todo, err2 = c.todoModel.GetByIDWithDeleted(todoID)
	} else {
		todo, err2 = c.todoModel.GetByID(todoID)
	}

	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	if todo.UserID != userID.(int) && userRole != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req CreateTodoItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := &entity.TodoItem{
		Title:       req.Title,
		Description: req.Description,
		TodoID:      todoID,
		UserID:      userID.(int),
	}

	if err := c.todoItemModel.Create(item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update todo completion percentage
	if err := c.todoModel.UpdateCompletionPct(todoID, c.todoItemModel); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo completion percentage"})
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

func (c *TodoItemController) GetByTodoID(ctx *gin.Context) {
	todoID, err := strconv.Atoi(ctx.Param("todo_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
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
		todo, err2 = c.todoModel.GetByIDWithDeleted(todoID)
	} else {
		todo, err2 = c.todoModel.GetByID(todoID)
	}

	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	if todo.UserID != userID.(int) && userRole != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var items []*entity.TodoItem
	if userRole == "admin" {
		items = c.todoItemModel.GetByTodoIDWithDeleted(todoID)
	} else {
		items = c.todoItemModel.GetByTodoID(todoID)
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *TodoItemController) Update(ctx *gin.Context) {
	todoID, err := strconv.Atoi(ctx.Param("todo_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("item_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
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
		todo, err2 = c.todoModel.GetByIDWithDeleted(todoID)
	} else {
		todo, err2 = c.todoModel.GetByID(todoID)
	}

	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	if todo.UserID != userID.(int) && userRole != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var item *entity.TodoItem
	var err3 error
	if userRole == "admin" {
		item, err3 = c.todoItemModel.GetByIDWithDeleted(itemID)
	} else {
		item, err3 = c.todoItemModel.GetByID(itemID)
	}

	if err3 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo item not found"})
		return
	}

	var req UpdateTodoItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.Completed = req.Completed

	if err := c.todoItemModel.Update(item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.todoModel.UpdateCompletionPct(todoID, c.todoItemModel); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo completion percentage"})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *TodoItemController) Delete(ctx *gin.Context) {
	todoID, err := strconv.Atoi(ctx.Param("todo_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("item_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
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
		todo, err2 = c.todoModel.GetByIDWithDeleted(todoID)
	} else {
		todo, err2 = c.todoModel.GetByID(todoID)
	}

	if err2 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	// Check if user is the owner or admin
	if todo.UserID != userID.(int) && userRole != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var err3 error
	if userRole == "admin" {
		_, err3 = c.todoItemModel.GetByIDWithDeleted(itemID)
	} else {
		_, err3 = c.todoItemModel.GetByID(itemID)
	}

	if err3 != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo item not found"})
		return
	}

	if err := c.todoItemModel.Delete(itemID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.todoModel.UpdateCompletionPct(todoID, c.todoItemModel); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update todo completion percentage"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "todo item deleted"})
}
