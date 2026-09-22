package handler

import (
	"7solutions-challenge/auth/domain"
	"7solutions-challenge/auth/service"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type authHandler struct {
	service service.AuthService
}

var validate = validator.New()

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{
		service: service,
	}
}

func (h *authHandler) Register(c *gin.Context) {
	var req domain.UserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Register(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
	})
}

func (h *authHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.LoginResponse{
		Token: token,
	})
}
