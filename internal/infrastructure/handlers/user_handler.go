package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ngodeyuk-core/internal/domain/dtos"
	"ngodeyuk-core/internal/services"
)

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	GetByUsername(ctx *gin.Context)
}

type userHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) UserHandler {
	return &userHandler{service}
}

func (handler *userHandler) Register(ctx *gin.Context) {
	var input dtos.RegisterDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.service.Register(&input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"name":     input.Name,
			"username": input.Username,
		},
	})
}

func (handler *userHandler) Login(ctx *gin.Context) {
	var input dtos.LoginDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := handler.service.Login(&input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (handler *userHandler) ChangePassword(ctx *gin.Context) {
	var input dtos.ChangePasswordDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	input.Username = username.(string)
	err := handler.service.ChangePassword(&input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "change password successfully"})
}

func (handler *userHandler) Update(ctx *gin.Context) {
	var input dtos.UpdateDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := handler.service.Update(username.(string), &input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"name":  input.Name,
			"point": input.Point,
			"heart": input.Heart,
		},
	})
}

func (handler *userHandler) GetAll(ctx *gin.Context) {
	users, err := handler.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var response []dtos.UserDTO
	for _, user := range users {
		response = append(response, dtos.UserDTO{
			UserId:   user.UserId,
			Name:     user.Name,
			Username: user.Username,
			ImgURL:   user.ImgURL,
			Heart:    user.Heart,
			Points:   user.Points,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (handler *userHandler) GetByUsername(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized"})
		return
	}
	user, err := handler.service.GetByUsername(username.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	response := dtos.UserDTO{
		UserId:   user.UserId,
		Name:     user.Name,
		ImgURL:   user.ImgURL,
		Username: user.Username,
		Heart:    user.Heart,
		Points:   user.Points,
	}
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}
