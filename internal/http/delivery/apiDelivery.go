package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/auth"
	"github.com/markgregr/RIP/internal/model"
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
// @Success 200 {object} model.DeliveryRequest "Список багажей"
// @Failure 500 {object} model.DeliveryRequest "Ошибка сервера"
// @Router /delivery [get]
func (h *Handler) GetDeliveries(c *gin.Context) {
    authInstance := auth.GetAuthInstance()
    searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
    startFormationDate := c.DefaultQuery("startFormationDate", "")
    endFormationDate := c.DefaultQuery("endFormationDate", "")
    deliveryStatus := c.DefaultQuery("deliveryStatus", "")

    var deliveries []model.DeliveryRequest
    var err error

    if authInstance.Role == "moderator" {
        deliveries, err = h.UseCase.GetDeliveriesModerator(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus, authInstance.UserID)
    } else {
        deliveries, err = h.UseCase.GetDeliveriesUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus, authInstance.UserID)
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
}

// GetDeliveryByID godoc
// @Summary Получение доставки по идентификатору
// @Description Возвращает информацию о доставке по её идентификатору
// @Tags Доставка
// @Produce json
// @Param id path int true "Идентификатор доставки"
// @Success 200 {object} model.DeliveryGetResponse "Информация о доставке"
// @Failure 400 {object} model.DeliveryGetResponse "Недопустимый идентификатор доставки"
// @Failure 500 {object} model.DeliveryGetResponse "Ошибка сервера"
// @Router /delivery/{id} [get]
func (h *Handler) GetDeliveryByID(c *gin.Context) {
    authInstance := auth.GetAuthInstance()
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    var delivery model.DeliveryGetResponse

    if authInstance.Role == "moderator"{
        delivery, err = h.UseCase.GetDeliveryByIDModerator(uint(deliveryID), authInstance.UserID)
    } else {
        // Получение доставки для пользователя
        delivery, err = h.UseCase.GetDeliveryByIDUser(uint(deliveryID), authInstance.UserID)
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, delivery)
}

// DeleteDelivery godoc
// @Summary Удаление доставки
// @Description Удаляет доставку по её идентификатору
// @Tags Доставка
// @Produce json
// @Param id path int true "Идентификатор доставки"
// @Param searchFlightNumber query string false "Номер рейса" Format(email)
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param deliveryStatus query string false "Статус доставки" Format(email)
// @Success 200 {object} model.DeliveryRequest "Список багажей"
// @Failure 400 {object} model.DeliveryRequest "Недопустимый идентификатор доставки"
// @Failure 500 {object} model.DeliveryRequest "Ошибка сервера"
// @Router /delivery/{id}/delete [delete]
func (h *Handler) DeleteDelivery(c *gin.Context) {
    authInstance := auth.GetAuthInstance()
    searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
    startFormationDate := c.DefaultQuery("startFormationDate", "")
    endFormationDate := c.DefaultQuery("endFormationDate", "")
    deliveryStatus := c.DefaultQuery("deliveryStatus", "")
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    if authInstance.Role == "moderator" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос недоступен для модератора"})
        return
    }

    err = h.UseCase.DeleteDeliveryUser(uint(deliveryID), authInstance.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    deliveries, err := h.UseCase.GetDeliveriesUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus, authInstance.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
}

// UpdateDeliveryFlightNumber godoc
// @Summary Обновление номера рейса доставки
// @Description Обновляет номер рейса для доставки по её идентификатору
// @Tags Доставка
// @Produce json
// @Param id path int true "Идентификатор доставки"
// @Param flightNumber body model.DeliveryUpdateFlightNumberRequest true "Новый номер рейса"
// @Success 200 {object} model.DeliveryGetResponse "Информация о доставке"
// @Failure 400 {object} model.DeliveryGetResponse "Недопустимый идентификатор доставки или ошибка чтения JSON объекта"
// @Failure 500 {object} model.DeliveryGetResponse "Ошибка сервера"
// @Router /delivery/{id}/update [put]
func (h *Handler) UpdateDeliveryFlightNumber(c *gin.Context) {
    authInstance := auth.GetAuthInstance()
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    var flightNumber model.DeliveryUpdateFlightNumberRequest
    if err := c.BindJSON(&flightNumber); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения JSON объекта"})
        return
    }

    if authInstance.Role == "moderator" {
        err = h.UseCase.UpdateFlightNumberModerator(uint(deliveryID), authInstance.UserID, flightNumber)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        delivery, err := h.UseCase.GetDeliveryByIDModerator(uint(deliveryID), authInstance.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
        }

        c.JSON(http.StatusOK, gin.H{"delivery": delivery})
    } else {
        err = h.UseCase.UpdateFlightNumberUser(uint(deliveryID), authInstance.UserID, flightNumber)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        delivery, err := h.UseCase.GetDeliveryByIDUser(uint(deliveryID),authInstance.UserID)
        if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Param id path int true "Идентификатор доставки"
// @Success 200 {object} model.DeliveryGetResponse "Информация о доставке"
// @Failure 400 {object} model.DeliveryGetResponse "Недопустимый идентификатор доставки"
// @Failure 500 {object} model.DeliveryGetResponse "Ошибка сервера"
// @Router /delivery/{id}/user [put]
func (h *Handler) UpdateDeliveryStatusUser(c *gin.Context) {
    authInstance := auth.GetAuthInstance()
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недоупстимый ИД доставки"})
        return
    }

    if authInstance.Role == "user" {
        err = h.UseCase.UpdateDeliveryStatusUser(uint(deliveryID), authInstance.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        delivery, err := h.UseCase.GetDeliveryByIDUser(uint(deliveryID), authInstance.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
        }

        c.JSON(http.StatusOK, gin.H{"delivery": delivery})
    } else if authInstance.Role == "moderator" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос доступен только пользователю"})
        return
    }
}

// UpdateDeliveryStatusModerator godoc
// @Summary Обновление статуса доставки для модератора
// @Description Обновляет статус доставки для модератора по идентификатору доставки
// @Tags Доставка
// @Produce json
// @Param id path int true "Идентификатор доставки"
// @Param deliveryStatus body model.DeliveryUpdateStatusRequest true "Новый статус доставки"
// @Success 200 {object} model.DeliveryGetResponse "Информация о доставке"
// @Failure 400 {object} model.DeliveryGetResponse "Недопустимый идентификатор доставки или ошибка чтения JSON объекта"
// @Failure 500 {object} model.DeliveryGetResponse "Ошибка сервера"
// @Router /delivery/{id}/status [put]
func (h *Handler) UpdateDeliveryStatusModerator(c *gin.Context) {
    // Получение экземпляра singleton для аутентификации
    authInstance := auth.GetAuthInstance()

    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    var deliveryStatus model.DeliveryUpdateStatusRequest
    if err := c.BindJSON(&deliveryStatus); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if authInstance.Role == "moderator" {
        err = h.UseCase.UpdateDeliveryStatusModerator(uint(deliveryID), authInstance.UserID, deliveryStatus)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        delivery, err := h.UseCase.GetDeliveryByIDUser(uint(deliveryID), authInstance.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
        }

        c.JSON(http.StatusOK, gin.H{"delivery": delivery})
    } else if authInstance.Role == "user" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос доступен только модератору"})
        return
    }
}
