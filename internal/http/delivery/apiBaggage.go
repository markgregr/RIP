package delivery

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/model"
	"github.com/markgregr/RIP/internal/pkg/middleware"
)

// @Summary Получение списка багажа
// @Description Возращает список всех активных багажей
// @Tags Багаж
// @Produce json
// @Param searchCode query string false "Код багажа" Format(email)
// @Success 200 {object} model.BaggagesGetResponse "Список багажей"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /baggage [get]
func (h *Handler) GetBaggages(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте пп"})
		return
	}
	userID := ctxUserID.(uint)
    searchCode := c.DefaultQuery("searchCode", "")

    baggages, err := h.UseCase.GetBaggages(searchCode,userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"baggages": baggages.Baggages, "deliveryID":baggages.DeliveryID})
}

// @Summary Получение багажа по ID
// @Description Возвращает информацию о багаже по его ID
// @Tags Багаж
// @Produce json
// @Param baggage_id path int true "ID багажа"
// @Success 200 {object} model.Baggage "Информация о багаже"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /baggage/{baggage_id} [get]
func (h *Handler) GetBaggageByID(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    baggageID, err := strconv.Atoi(c.Param("baggage_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД багажа"})
        return
    }

    baggage, err := h.UseCase.GetBaggageByID(uint(baggageID), userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"baggage": baggage})
}

// @Summary Создание нового багажа
// @Description Создает новый багаж с предоставленными данными
// @Tags Багаж
// @Accept json
// @Produce json
// @Param searchCode query string false "Код багажа" Format(email)
// @Param baggage body model.BaggageRequest true "Пользовательский объект в формате JSON"
// @Success 200 {object} model.BaggagesGetResponse "Список багажей"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} model.ErrorResponse "У пользователя нет прав для этого запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /baggage/create [post]
func (h *Handler) CreateBaggage(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    searchCode := c.DefaultQuery("searchCode", "")

	var baggage model.BaggageRequest

	if err := c.BindJSON(&baggage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать JSON"})
		return
	}

    if middleware.ModeratorOnly(h.UseCase.Repository, c){
        err := h.UseCase.CreateBaggage(userID, baggage)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    
        baggages, err := h.UseCase.GetBaggages(searchCode, userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    
        c.JSON(http.StatusOK, gin.H{"baggages": baggages.Baggages, "deliveryID":baggages.DeliveryID})
    } else {
        c.JSON(http.StatusForbidden, gin.H{"error": "данный запрос доступен только модератору"})
        return
    }
}

// @Summary Удаление багажа
// @Description Удаляет багаж по его ID
// @Tags Багаж
// @Produce json
// @Param baggage_id path int true "ID багажа"
// @Param searchCode query string false "Код багажа" Format(email)
// @Success 200 {object} model.BaggagesGetResponse "Список багажей"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} model.ErrorResponse "У пользователя нет прав для этого запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /baggage/{baggage_id}/delete [delete]
func (h *Handler) DeleteBaggage(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)
	
    searchCode := c.DefaultQuery("searchCode", "")

	baggageID, err := strconv.Atoi(c.Param("baggage_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД багажа"})
		return
	}

    if middleware.ModeratorOnly(h.UseCase.Repository, c){
        err = h.UseCase.DeleteBaggage(uint(baggageID), userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    
        baggages, err := h.UseCase.GetBaggages(searchCode, userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
    
        c.JSON(http.StatusOK, gin.H{"baggages": baggages.Baggages, "deliveryID":baggages.DeliveryID})
    } else {
        c.JSON(http.StatusForbidden, gin.H{"error": "данный запрос доступен только модератору"})
        return
    }
}

// @Summary Обновление информации о багаже
// @Description Обновляет информацию о багаже по его ID
// @Tags Багаж
// @Accept json
// @Produce json
// @Param baggage_id path int true "ID багажа"
// @Success 200 {object} model.Baggage "Информация о багаже"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} model.ErrorResponse "У пользователя нет прав для этого запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /baggage/{baggage_id}/update [put]
func (h *Handler) UpdateBaggage(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    baggageID, err := strconv.Atoi(c.Param("baggage_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"error": "недопустимый ИД багажа"}})
        return
    }

    var baggage model.BaggageRequest
    if err := c.BindJSON(&baggage); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать JSON"})
        return
    }
    
    if middleware.ModeratorOnly(h.UseCase.Repository, c){
        err = h.UseCase.UpdateBaggage(uint(baggageID),uint(userID), baggage)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        updatedBaggage, err := h.UseCase.GetBaggageByID(uint(baggageID), uint(userID))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"baggage": updatedBaggage})
    } else {
        c.JSON(http.StatusForbidden, gin.H{"error": "данный запрос дсотупен только модератору"})
        return
    }
}

// @Summary Добавление багажа к доставке
// @Description Добавляет багаж к доставке по его ID
// @Tags Багаж
// @Produce json
// @Param baggage_id path int true "ID багажа"
// @Param searchCode query string false "Код багажа" Format(email)
// @Success 200 {object} model.BaggagesGetResponse  "Список багажей"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /baggage/{baggage_id}/delivery [post]
func (h *Handler) AddBaggageToDelivery(c *gin.Context) {
	ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    searchCode := c.DefaultQuery("searchCode", "")

    baggageID, err := strconv.Atoi(c.Param("baggage_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД багажа"})
        return
    }

    err = h.UseCase.AddBaggageToDelivery(uint(baggageID), uint(userID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	baggages, err := h.UseCase.GetBaggages(searchCode, uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

    c.JSON(http.StatusOK, gin.H{"baggages": baggages.Baggages, "deliveryID":baggages.DeliveryID})
}

// @Summary Удаление багажа из доставки
// @Description Удаляет багаж из доставки по его ID
// @Tags Багаж
// @Produce json
// @Param baggage_id path int true "ID багажа"
// @Param searchCode query string false "Код багажа" Format(email)
// @Success 200 {object} model.BaggagesGetResponse "Список багажей"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /baggage/{baggage_id}/delivery [delete]
func (h *Handler) RemoveBaggageFromDelivery(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    searchCode := c.DefaultQuery("searchCode", "")

    baggageID, err := strconv.Atoi(c.Param("baggage_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД багажа"})
        return
    }
   
    err = h.UseCase.RemoveBaggageFromDelivery(uint(baggageID), uint(userID))  
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    baggages, err := h.UseCase.GetBaggages(searchCode, uint(userID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"baggages": baggages.Baggages, "deliveryID":baggages.DeliveryID})
}

// @Summary Добавление изображения к багажу
// @Description Добавляет изображение к багажу по его ID
// @Tags Багаж
// @Accept mpfd
// @Produce json
// @Param baggage_id path int true "ID багажа"
// @Param image formData file true "Изображение багажа"
// @Success 200 {object} model.Baggage "Информация о багаже с изображением"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 403 {object} model.ErrorResponse "У пользователя нет прав для этого запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /baggage/{baggage_id}/image [post]
func (h* Handler) AddBaggageImage(c* gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    baggageID, err := strconv.Atoi(c.Param("baggage_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД багажа"})
        return
    }

    image, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимое изображение"})
        return
    }

    file, err := image.Open()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось открыть изображение"})
        return
    }
    defer file.Close()

    imageBytes, err := io.ReadAll(file)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать изображение в байтах"})
        return
    }

	contentType := image.Header.Get("Content-Type")
    
    if middleware.ModeratorOnly(h.UseCase.Repository, c){
        err = h.UseCase.AddBaggageImage(uint(baggageID), uint(userID),imageBytes, contentType)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        baggage, err := h.UseCase.GetBaggageByID(uint(baggageID),uint(userID))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"baggage": baggage})
    } else {
        c.JSON(http.StatusForbidden, gin.H{"error": "данный запрос доступен только модератору"})
        return
    }
}



