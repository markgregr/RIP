package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/model"
	"github.com/markgregr/RIP/internal/pkg/middleware"
)

// GetDeliveries godoc
// @Summary Получение списка доставок
// @Description Возвращает список всех не удаленных доставок
// @Tags Доставка
// @Produce json
// @Param searchFlightNumber query string false "Номер рейса" Format(email)
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param deliveryStatus query string false "Статус доставки" Format(email)
// @Success 200 {object} model.DeliveryRequest "Список доставок"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /delivery [get]
func (h *Handler) GetDeliveries(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
    startFormationDate := c.DefaultQuery("startFormationDate", "")
    endFormationDate := c.DefaultQuery("endFormationDate", "")
    deliveryStatus := c.DefaultQuery("deliveryStatus", "")

    var deliveries []model.DeliveryRequest
    var err error

    if middleware.ModeratorOnly(h.UseCase.Repository, c){
        deliveries, err = h.UseCase.GetDeliveriesModerator(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus)  
    } else {
        deliveries, err = h.UseCase.GetDeliveriesUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus, userID)
    }
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
}

// GetDeliveryByID godoc
// @Summary Получение доставки по идентификатору
// @Description Возвращает информацию о доставке по её идентификатору
// @Tags Доставка
// @Produce json
// @Param delivery_id path int true "Идентификатор доставки"
// @Success 200 {object} model.DeliveryGetResponse "Информация о доставке"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /delivery/{delivery_id} [get]
func (h *Handler) GetDeliveryByID(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    var delivery model.DeliveryGetResponse

    if middleware.ModeratorOnly(h.UseCase.Repository, c){
        delivery, err = h.UseCase.GetDeliveryByIDModerator(uint(deliveryID)) 
    } else {
        delivery, err = h.UseCase.GetDeliveryByIDUser(uint(deliveryID), userID)
    }
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"delivery" : delivery})
}

// DeleteDelivery godoc
// @Summary Удаление доставки
// @Description Удаляет доставку по её идентификатору
// @Tags Доставка
// @Produce json
// @Param delivery_id path int true "Идентификатор доставки"
// @Param searchFlightNumber query string false "Номер рейса" Format(email)
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param deliveryStatus query string false "Статус доставки" Format(email)
// @Success 200 {object} model.DeliveryRequest "Список багажей"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /delivery/{delivery_id} [delete]
func (h *Handler) DeleteDelivery(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
    startFormationDate := c.DefaultQuery("startFormationDate", "")
    endFormationDate := c.DefaultQuery("endFormationDate", "")
    deliveryStatus := c.DefaultQuery("deliveryStatus", "")
    deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    if middleware.ModeratorOnly(h.UseCase.Repository, c){
        err = h.UseCase.DeleteDeliveryUser(uint(deliveryID), userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        deliveries, err := h.UseCase.GetDeliveriesModerator(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
    } else {
        err = h.UseCase.DeleteDeliveryUser(uint(deliveryID), userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        deliveries, err := h.UseCase.GetDeliveriesUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus, userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
    }
}

// UpdateDeliveryFlightNumber godoc
// @Summary Обновление номера рейса доставки
// @Description Обновляет номер рейса для доставки по её идентификатору
// @Tags Доставка
// @Produce json
// @Param delivery_id path int true "Идентификатор доставки"
// @Param flightNumber body model.DeliveryUpdateFlightNumberRequest true "Новый номер рейса"
// @Success 200 {object} model.DeliveryGetResponse "Информация о доставке"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /delivery/{delivery_id} [put]
func (h *Handler) UpdateDeliveryFlightNumber(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    var flightNumber model.DeliveryUpdateFlightNumberRequest
    if err := c.BindJSON(&flightNumber); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения JSON объекта"})
        return
    }

    if middleware.ModeratorOnly(h.UseCase.Repository, c){
        err = h.UseCase.UpdateFlightNumberUser(uint(deliveryID), userID, flightNumber)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        delivery, err := h.UseCase.GetDeliveryByIDModerator(uint(deliveryID))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
        }

        c.JSON(http.StatusOK, gin.H{"delivery": delivery})
    } else {
        err = h.UseCase.UpdateFlightNumberUser(uint(deliveryID), userID, flightNumber)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        delivery, err := h.UseCase.GetDeliveryByIDUser(uint(deliveryID), userID)
        if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
        }

        c.JSON(http.StatusOK, gin.H{"delivery": delivery})
    }
}

// UpdateDeliveryStatusUser godoc
// @Summary Обновление статуса доставки для пользователя
// @Description Обновляет статус доставки для пользователя по идентификатору доставки
// @Tags Доставка
// @Produce json
// @Param delivery_id path int true "Идентификатор доставки"
// @Success 200 {object} model.DeliveryGetResponse "Информация о доставке"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /delivery/{delivery_id}/status/user [put]
func (h *Handler) UpdateDeliveryStatusUser(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недоупстимый ИД доставки"})
        return
    }

    if middleware.ModeratorOnly(h.UseCase.Repository, c) {
        err = h.UseCase.UpdateDeliveryStatusUser(uint(deliveryID), userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        delivery, err := h.UseCase.GetDeliveryByIDModerator(uint(deliveryID))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
        }

        c.JSON(http.StatusOK, gin.H{"delivery": delivery})
    
    } else {
        err = h.UseCase.UpdateDeliveryStatusUser(uint(deliveryID), userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        delivery, err := h.UseCase.GetDeliveryByIDUser(uint(deliveryID), userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
        }

        c.JSON(http.StatusOK, gin.H{"delivery": delivery})
    }
}

// UpdateDeliveryStatusModerator godoc
// @Summary Обновление статуса доставки для модератора
// @Description Обновляет статус доставки для модератора по идентификатору доставки
// @Tags Доставка
// @Produce json
// @Param delivery_id path int true "Идентификатор доставки"
// @Param deliveryStatus body model.DeliveryUpdateStatusRequest true "Новый статус доставки"
// @Success 200 {object} model.DeliveryGetResponse "Информация о доставке"
// @Failure 400 {object} model.ErrorResponse "Обработанная ошибка сервера"
// @Failure 401 {object} model.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /delivery/{delivery_id}/status/moderator [put]
func (h *Handler) UpdateDeliveryStatusModerator(c *gin.Context) {
    ctxUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Идентификатор пользователя отсутствует в контексте"})
		return
	}
	userID := ctxUserID.(uint)

    deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    var deliveryStatus model.DeliveryUpdateStatusRequest
    if err := c.BindJSON(&deliveryStatus); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if middleware.ModeratorOnly(h.UseCase.Repository, c){
        err = h.UseCase.UpdateDeliveryStatusModerator(uint(deliveryID), userID, deliveryStatus)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        delivery, err := h.UseCase.GetDeliveryByIDUser(uint(deliveryID), userID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
        }

        c.JSON(http.StatusOK, gin.H{"delivery": delivery})
    } else {
        c.JSON(http.StatusForbidden, gin.H{"error": "данный запрос доступен только модератору"})
        return
    }
}
