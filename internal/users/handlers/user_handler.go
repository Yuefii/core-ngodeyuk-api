package handlers

import (
	"net/http"
	"ngodeyuk-core/internal/users/services"
	"ngodeyuk-core/pkg/dto"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	RegisterUser(c *gin.Context)
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var registerDTO dto.RegisterDTO
	if err := c.ShouldBindJSON(&registerDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.RegisterUser(registerDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "register successfully",
		"data": gin.H{
			"Name":     user.Name,
			"Username": user.Username,
		},
	})
}
