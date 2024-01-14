package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/model"
)

// @BasePath /user/register
// @Summary Регистрация нового пользователя.
// @Description Регистрация нового пользователя с предоставленной информацией.
// @Tags Пользователь
// @Accept json
// @Produce json
// @Param user body model.UserRegisterRequest true "Пользовательский объект в формате JSON"
// @Success 200 {object} model.UserLoginResponse "Токен"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /user/register [post]
func (h *Handler) Register(c *gin.Context) {
	var user model.UserRegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginResponse, err := h.UseCase.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": loginResponse.AccessToken, "full_name":loginResponse.FullName})
}

// @BasePath /user/login
// @Summary Вход пользователя
// @Description Авторизация пользователя и генерация JWT-токена
// @Tags Пользователь
// @Accept json
// @Produce json
// @Param body body model.UserLoginRequest true "Данные для входа"
// @Success 200 {object} model.UserLoginResponse "Токен"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /user/login [post]
func (h *Handler) Login(c *gin.Context) {
	var user model.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginResponse, err := h.UseCase.LoginUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": loginResponse.AccessToken, "full_name":loginResponse.FullName})

}

// @BasePath /user/logout
// @Summary Выход пользователя
// @Description Выход пользователя из системы и удаление токена из куков
// @Tags Пользователь
// @Produce json
// @Success 200 {string} string "Успешный ответ"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /user/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	cUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := cUserID.(uint)
	
	err := h.UseCase.LogoutUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно вышел из системы"})
}
