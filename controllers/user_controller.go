package controllers

import (
	"net/http"
	"strconv"

	"todoapp/entity"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userModel *entity.UserModel
}

func NewUserController(userModel *entity.UserModel) *UserController {
	return &UserController{
		userModel: userModel,
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Role     string `json:"role" binding:"required"`
}

func (c *UserController) Create(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := c.userModel.GetByUsername(req.Username); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}

	user := &entity.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := c.userModel.Create(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := c.userModel.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) GetByUsername(ctx *gin.Context) {
	username := ctx.Param("username")
	user, err := c.userModel.GetByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) GetAll(ctx *gin.Context) {
	users := c.userModel.GetAll()
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := c.userModel.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username != user.Username {
		if _, err := c.userModel.GetByUsername(req.Username); err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
			return
		}
	}

	user.Username = req.Username
	if req.Password != "" {
		user.Password = req.Password
	}
	user.Role = req.Role

	if err := c.userModel.Update(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.userModel.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
