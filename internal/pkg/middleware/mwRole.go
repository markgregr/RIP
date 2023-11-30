package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/http/repository"
	"github.com/markgregr/RIP/internal/model"
)

func ModeratorOnly(r *repository.Repository, c *gin.Context) bool {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется аутентификация"})
		c.Abort()
	}

	userID, ok := ctxUserID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при преобразовании идентификатора пользователя"})
		c.Abort()
	}

	role, err := r.GetUserRoleByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
	}	

	if role == model.ModeratorRole {
		return true
	}
	return false
}
