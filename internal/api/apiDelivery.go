package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/app/ds"
)

//методы для таблицы delivery
func (h *Handler) GetDeliveries(c *gin.Context) {
	searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
    deliveryStatus := c.DefaultQuery("deliveryStatus", "")
	deliveries, err := h.Repo.GetDeliveries(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deliveries": deliveries})
}
func (h *Handler) GetDeliveryByID(c *gin.Context) {
	deliveryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	delivery, err := h.Repo.GetDeliveryByID(deliveryID)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }
	c.JSON(http.StatusOK, delivery)
}
func (h *Handler) DeleteDelivery(c *gin.Context) {
	searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
    deliveryStatus := c.DefaultQuery("deliveryStatus", "")
	deliveryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err = h.Repo.DeleteDelivery(deliveryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Получаем обновленный список багажей
	deliveries, err := h.Repo.GetDeliveries(searchFlightNumber, startFormationDate, endFormationDate, deliveryStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery deleted successfully", "deliveries": deliveries})
}
func (h *Handler) UpdateDelivery(c *gin.Context) {
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
        return
    }

    var updatedDeliveryRequest ds.Delivery
    if err := c.BindJSON(&updatedDeliveryRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Попытка обновления доставки в репозитории
    if err := h.Repo.UpdateDelivery(deliveryID, &updatedDeliveryRequest); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    // Получаем обновленный объект доставки
    updatedDelivery, err := h.Repo.GetDeliveryByID(deliveryID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    // Здесь предполагается, что метод UpdateDelivery вернет ошибку, если обновление не удалось
    // В противном случае, вы не сможете получить обновленный объект

    c.JSON(http.StatusOK, gin.H{"message": "Delivery updated successfully", "delivery": updatedDelivery})
}
func (h *Handler) UpdateDeliveryStatusForUser(c *gin.Context) {
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
        return
    }

    var updateRequest ds.Delivery
    if err := c.BindJSON(&updateRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.Repo.UpdateDeliveryStatusForUser(deliveryID, 1, &updateRequest) // Здесь предполагается, что userID передается в репозиторий (в данном случае равен 0)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    updatedDelivery, err := h.Repo.GetDeliveryByID(deliveryID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Delivery updated successfully", "delivery": updatedDelivery})
}
func (h *Handler) UpdateDeliveryStatusForModerator(c *gin.Context) {
    deliveryID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
        return
    }

    var updateRequest ds.Delivery
    if err := c.BindJSON(&updateRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.Repo.UpdateDeliveryStatusForModerator(deliveryID, 1, &updateRequest) // Здесь предполагается, что moderatorID передается в репозиторий (в данном случае равен 1)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    updatedDelivery, err := h.Repo.GetDeliveryByID(deliveryID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Delivery updated successfully", "delivery": updatedDelivery})
}
