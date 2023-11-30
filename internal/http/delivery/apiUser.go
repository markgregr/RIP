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
// @Param user body model.User true "Пользовательский объект в формате JSON"
// @Success 201 {object} []model.User "Успешно зарегистрированный пользователь"
// @Router /user/register [post]
func (h *Handler) Register(c *gin.Context) {
	var user model.UserRegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.UseCase.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	users, err := h.UseCase.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"users": users})
}

// @BasePath /user/login
// @Summary Вход пользователя
// @Description Авторизация пользователя и генерация JWT-токена
// @Tags Пользователь
// @Accept json
// @Produce json
// @Param body body model.UserLoginRequest true "Данные для входа"
// @Success 200 {object} map[string]string "Успешный ответ"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /user/login [post]
func (h *Handler) Login(c *gin.Context) {
	var user model.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := h.UseCase.LoginUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("access_token", tokenPair.AccessToken, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": tokenPair.AccessToken})
}

// @BasePath /api/user
// @Summary Получить пользователя по идентификатору
// @Description Получение данных пользователя по его идентификатору
// @Produce json
// @Success 200 {object} []model.User "Успешный ответ"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /user/ [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	cUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}

	userID, ok := cUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при преобразовании идентификатора пользователя"})
		return
	}

	user, err := h.UseCase.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// @BasePath /user/logout
// @Summary Выход пользователя
// @Description Выход пользователя из системы и удаление токена из куков
// @Tags Пользователь
// @Produce json
// @Success 200 {object} map[string]string "Успешный ответ"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /user/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	cUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}

	userID, ok := cUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при преобразовании идентификатора пользователя"})
		return
	}

	err := h.UseCase.LogoutUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно вышел из системы"})
}