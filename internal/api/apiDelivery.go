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
	deliveries, err := h.Repo.GetDeliveries(searchFlightNumber)
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
	deliveries, err := h.Repo.GetDeliveries(searchFlightNumber)
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

    var updatedDelivery ds.Delivery
    if err := c.BindJSON(&updatedDelivery); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Попытка обновления багажа в репозитории
    err = h.Repo.UpdateDelivery(deliveryID, &updatedDelivery)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    // Получаем обновленный объект багажа (указатель на Baggage)
    updatedDeliveryPtr, err := h.Repo.GetDeliveryByID(deliveryID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    // Преобразуем указатель в значение
    updatedDelivery = *updatedDeliveryPtr

    c.JSON(http.StatusOK, gin.H{"message": "Delivery updated successfully", "delivery": updatedDelivery})
}
