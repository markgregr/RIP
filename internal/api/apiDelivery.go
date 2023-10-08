package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/app/ds"
)

//методы для таблицы delivery
func (h *Handler) GetDeliveries(c *gin.Context) {
	deliveries, err := h.Repo.GetDeliveries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	searchQuery := c.DefaultQuery("searchQuery", "")
	var foundDeliveries []ds.Delivery
	for _, delivery := range deliveries {
		if strings.HasPrefix(strings.ToLower(delivery.FlightNumber), strings.ToLower(searchQuery)) {
			foundDeliveries = append(foundDeliveries, delivery)
		}
	}
	data := gin.H{"deliveries": foundDeliveries}
	c.JSON(http.StatusOK, data)
}
func (h *Handler) GetDeliveryByID(c *gin.Context) {
	deliveryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	delivery, err := h.Repo.GetDeliveryByID(deliveryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GetDeliveryByID"})
		return
	}
	c.JSON(http.StatusOK, delivery)
}
func (h *Handler) DeleteDelivery(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "Delivery deleted successfully"})
}
func (h *Handler) UpdateDelivery(c *gin.Context) {
	deliveryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var updatedDelivery ds.Delivery
	if err := c.ShouldBindJSON(&updatedDelivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Repo.UpdateDelivery(deliveryID, &updatedDelivery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery updated successfully"})
}
