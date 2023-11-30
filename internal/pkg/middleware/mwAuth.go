package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/markgregr/RIP/internal/http/repository"
)

// Функция для проверки и обновления токена
func AuthenticateAndRefresh(redisClient *redis.Client, jwtSecretKey []byte, r *repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Для продолжения необходимо авторизоваться"})
			ctx.Abort()
			return
		}

		userID, err := ParseAndValidateToken(tokenString, jwtSecretKey)
		if err != nil {
			// Если ошибка связана с истекшим токеном, попробуем обновить токены и распарсить новый access token
			if strings.Contains(err.Error(), "Token is expired") {
				newTokens, newErr := RefreshTokens(tokenString, userID, r, jwtSecretKey)
				if newErr != nil {
					ctx.JSON(http.StatusUnauthorized, gin.H{"error": newErr.Error()})
					return
				}

				// Токены успешно обновлены, попробуем распарсить новый access token
				userID, err = ParseAndValidateToken(newTokens.AccessToken, jwtSecretKey)
				if err != nil {
					ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
					return
				}
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
		}

		if err := r.CheckTokenInRedis(userID, tokenString); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Set("userID", userID)
		log.Println(userID)
		ctx.Next()
	}
}


