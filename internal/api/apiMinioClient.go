package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddBaggageImage(c *gin.Context) {
	baggageID, err := strconv.Atoi(c.Param("baggage_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопсутимый ИД багажа"})
		return
	}

	// Чтение изображения из запроса
    image, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимое изображение"})
        return
    }

    // Чтение содержимого изображения в байтах
    file, err := image.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось открыть изображение"})
        return
    }
    defer file.Close()

    imageBytes, err := io.ReadAll(file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось прочитать изображение в байтах"})
        return
    }
	// Получение Content-Type из заголовков запроса
	contentType := image.Header.Get("Content-Type")

	// Вызов функции репозитория для добавления изображения
	err = h.Repo.AddBaggageImage(baggageID, imageBytes, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Изображение усспешно загружено"})
}
