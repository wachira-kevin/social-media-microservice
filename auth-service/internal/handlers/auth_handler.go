package handlers

import (
	"auth-service/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthHandler struct {
	KeycloakService *services.KeycloakService
}

func NewAuthHandler(ks *services.KeycloakService) *AuthHandler {
	return &AuthHandler{KeycloakService: ks}
}

func (ah *AuthHandler) Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	tokens, err := ah.KeycloakService.Login(credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (ah *AuthHandler) Logout(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	if err := ah.KeycloakService.Logout(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (ah *AuthHandler) RefreshToken(c *gin.Context) {
	token := strings.Split(c.GetHeader("Authorization"), " ")[1]
	tokens, err := ah.KeycloakService.RefreshToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tokens)
}
