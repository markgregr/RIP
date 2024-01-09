package middleware

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/markgregr/RIP/internal/model"
)

func GenerateJWTAccessToken(userID uint) (model.Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour*2).Unix() 

	accessToken, err := token.SignedString([]byte("AccessSecretKey"))
	if err != nil {
		return model.Token{}, err
	}
	
	return model.Token{AccessToken:accessToken}, nil
}

 func ParseAndValidateToken(tokenStr string, jwtSecretKey []byte) (uint, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return jwtSecretKey, nil
    })
	if err != nil{
		return 0, errors.New("ошибка парсинга акцесс токена")
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
