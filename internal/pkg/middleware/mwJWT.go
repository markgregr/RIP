package middleware

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/markgregr/RIP/internal/http/repository"
	"github.com/markgregr/RIP/internal/model"
)

func GenerateJWTAccessToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Minute*2).Unix() 

	accessToken, err := token.SignedString([]byte("AccessSecretKey"))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func GenerateJWTRefreshToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", errors.New("ошибка генерации случайной последовательности")
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix() 

	refreshToken, err := token.SignedString([]byte("RefreshSecretKey"))
	if err != nil {
		return "", errors.New("ошибка генерации рефреш токена")
	}

	return refreshToken, nil
}


func GenerateJWTTokenPair(userID uint) (model.TokenPair, error) {
	accessToken, err := GenerateJWTAccessToken(userID)
	if err != nil {
	   return model.TokenPair{}, err
	}
 
	refreshToken, err := GenerateJWTRefreshToken()
	if err != nil {
	   return model.TokenPair{}, err
	}
 
	return model.TokenPair{
	   AccessToken:  accessToken,
	   RefreshToken: refreshToken,
	}, nil
 }

 func ParseAndValidateToken(tokenString string, jwtSecretKey []byte) (uint, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecretKey, nil
    })

    if err != nil {
        if ve, ok := err.(*jwt.ValidationError); ok && ve.Errors == jwt.ValidationErrorExpired {
            // Токен истек, но мы все равно можем получить userID
            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                return 0, errors.New("ошибка возвращения claims")
            }

            userIDFloat, ok := claims["userID"].(float64)
            if !ok {
                return 0, errors.New("ошибка перевода типа userID")
            }

            userID := uint(userIDFloat)
            return userID, nil
        }
        return 0, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, errors.New("ошибка возвращения claims")
    }

    userIDFloat, ok := claims["userID"].(float64)
    if !ok {
        return 0, errors.New("ошибка перевода типа userID")
    }

    userID := uint(userIDFloat)

    return userID, nil
}


func RefreshTokens(oldToken string, userID uint, r *repository.Repository, jwtSecretKey []byte) (model.TokenPair, error) {
	newTokens, err := GenerateJWTTokenPair(userID)
	if err != nil {
		return model.TokenPair{}, err
	}
	if err := r.DeleteJWTTokenPair(userID); err != nil {
		return model.TokenPair{}, err
	}
	if err := r.SaveJWTTokenPair(userID, newTokens.AccessToken, newTokens.RefreshToken); err != nil {
		return model.TokenPair{}, err
	}

	return newTokens, nil
}