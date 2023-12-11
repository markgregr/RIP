package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/markgregr/RIP/internal/http/repository"
)

func Authenticate(redisClient *redis.Client, jwtSecretKey []byte, r *repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessTokenStr := ctx.GetHeader("Authorization")
		if accessTokenStr == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Для продолжения необходимо авторизоваться"})
			ctx.Abort()
			return
		}

		userID, err := ParseAndValidateToken(accessTokenStr, jwtSecretKey)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		if err := r.CheckTokenInRedis(userID, accessTokenStr); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		log.Println(userID)
		ctx.Next()
	}
}

func Guest(redisClient *redis.Client, jwtSecretKey []byte, r *repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessTokenStr := ctx.GetHeader("Authorization")
		if accessTokenStr == "" {
			ctx.Set("userID", uint(0))
			return
		}
		
		userID, err := ParseAndValidateToken(accessTokenStr, jwtSecretKey)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Set("userID", uint(0))
			ctx.Abort()
			return 
		}

		if err := r.CheckTokenInRedis(userID, accessTokenStr); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Set("userID", uint(0))
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		log.Println(userID)
		ctx.Next()
	}
}



