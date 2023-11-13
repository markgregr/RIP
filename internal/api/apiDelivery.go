package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/app/ds"
	"github.com/markgregr/RIP/internal/auth"
)

func (h *Handler) GetDeliveries(c *gin.Context) {
    // Получение экземпляра singleton для аутентификации
    authInstance := auth.GetAuthInstance()

    searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
    startFormationDate := c.DefaultQuery("startFormationDate", "")
    endFormationDate := c.DefaultQuery("endFormationDate", "")
    deliveryStatus := c.DefaultQuery("deliveryStatus", "")

    // Выбор соответствующего метода репозитория в зависимости от роли пользователя
    var deliveries map[string]interface{}
    var err error
    if authInstance.Role == "moderator" {
        // Получение доставок для модератора
        deliveries, err = h.Repo.GetDeliveriesForModerator(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus, authInstance.UserID)
    } else {
        // Получение доставок для пользователя
        deliveries, err = h.Repo.GetDeliveriesForUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus, authInstance.UserID)
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
}

func (h *Handler) GetDeliveryByID(c *gin.Context) {
    // Получение экземпляра singleton для аутентификации
    authInstance := auth.GetAuthInstance()

    // Получение идентификатора доставки из параметров запроса
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    // Получение информации о доставке в зависимости от роли пользователя
    var delivery map[string]interface{}
    var repoErr error
    if authInstance.Role == "moderator"{
        // Получение доставки для модератора
        delivery, repoErr = h.Repo.GetDeliveryByIDForModerator(deliveryID, authInstance.UserID)
    } else {
        // Получение доставки для пользователя
        delivery, repoErr = h.Repo.GetDeliveryByIDForUser(deliveryID, authInstance.UserID)
    }

    if repoErr != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": repoErr.Error()})
        return
    }
    // Возвращение информации о доставке
    c.JSON(http.StatusOK, delivery)
}

func (h *Handler) DeleteDelivery(c *gin.Context) {
    // Получение экземпляра singleton для аутентификации
    authInstance := auth.GetAuthInstance()

    // Получение идентификатора доставки из параметров запроса
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }
    // Возвращение сообщения об успешном удалении и обновленного списка доставок
    searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
    startFormationDate := c.DefaultQuery("startFormationDate", "")
    endFormationDate := c.DefaultQuery("endFormationDate", "")
    deliveryStatus := c.DefaultQuery("deliveryStatus", "")

    // Проверка, является ли текущий пользователь пользователем (не модератором)
    if authInstance.Role == "moderator" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос недоступен для модератора"})
        return
    }

    // Удаление доставки только если она принадлежит текущему пользователю
    err = h.Repo.DeleteDeliveryForUser(deliveryID, authInstance.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Получаем обновленный список доставок
    deliveries, err := h.Repo.GetDeliveriesForUser(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus, authInstance.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Доставка успешно удалена", "deliveries": deliveries})
}

func (h *Handler) UpdateDelivery(c *gin.Context) {
    // Получение идентификатора доставки из параметров запроса
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    // Привязка JSON-запроса к структуре Delivery
    var updatedDeliveryRequest ds.Delivery
    if err := c.BindJSON(&updatedDeliveryRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Получение экземпляра singleton для аутентификации
    authInstance := auth.GetAuthInstance()

    // Проверка, является ли пользователь авторизованным и имеет ли права на обновление доставки
    var repoErr error
    if authInstance.Role == "moderator" {
        // Обновление доставки для модератора
        repoErr = h.Repo.UpdateDeliveryForModerator(deliveryID, authInstance.UserID, &updatedDeliveryRequest)
        if repoErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": repoErr.Error()})
            return
        }
        // Получение обновленного объекта доставки
        updatedDelivery, err := h.Repo.GetDeliveryByIDForUser(deliveryID, authInstance.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Доставка успешно изменена", "delivery": updatedDelivery})
    } else {
        // Обновление доставки для пользователя
        repoErr = h.Repo.UpdateDeliveryForUser(deliveryID, authInstance.UserID, &updatedDeliveryRequest)
        if repoErr != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": repoErr.Error()})
            return
        }
        updatedDelivery, err := h.Repo.GetDeliveryByIDForUser(deliveryID,authInstance.UserID)
        if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Доставка успешно изменена", "delivery": updatedDelivery})
    }
}

func (h *Handler) UpdateDeliveryStatusForUser(c *gin.Context) {
    // Получение экземпляра singleton для аутентификации
    authInstance := auth.GetAuthInstance()

    // Получение идентификатора доставки из параметров запроса
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недоупстимый ИД доставки"})
        return
    }

    // Проверка роли пользователя
    if authInstance.Role == "user" {
        // Пользователь может обновлять только свои доставки
        err = h.Repo.UpdateDeliveryStatusForUser(deliveryID, authInstance.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        // Получение обновленного объекта доставки
        updatedDelivery, err := h.Repo.GetDeliveryByIDForUser(deliveryID, authInstance.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Доставка успешно обновлена", "delivery": updatedDelivery})
    } else if authInstance.Role == "moderator" {
        // Модератор не имеет права обновлять статус доставок пользователя
        c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос доступен только пользователю"})
        return
    }
}

func (h *Handler) UpdateDeliveryStatusForModerator(c *gin.Context) {
    // Получение экземпляра singleton для аутентификации
    authInstance := auth.GetAuthInstance()

    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
        return
    }

    var updateRequest ds.Delivery
    if err := c.BindJSON(&updateRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

   // Проверка роли пользователя
    if authInstance.Role == "moderator" {
        // Пользователь может обновлять только свои доставки
        err = h.Repo.UpdateDeliveryStatusForModerator(deliveryID, authInstance.UserID, &updateRequest)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        // Получение обновленного объекта доставки
        updatedDelivery, err := h.Repo.GetDeliveryByIDForUser(deliveryID, authInstance.UserID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
        }
        c.JSON(http.StatusOK, gin.H{"message": "Доставка успешно обновлена", "delivery": updatedDelivery})
    } else if authInstance.Role == "user" {
        // Модератор не имеет права обновлять статус доставок пользователя
        c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос доступен только модератору"})
        return
    }
}
