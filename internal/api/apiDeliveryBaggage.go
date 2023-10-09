package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
func (h *Handler) AddBaggageToDelivery(c *gin.Context) {
    // Получаем данные из JSON-запроса
    var payload struct {
        BaggageID  uint `json:"baggage_id"`
        DeliveryID uint `json:"delivery_id"`
    }

    if err := c.BindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Попытка обновления связи между багажом и доставкой в репозитории
    err := h.Repo.AddBaggageToDelivery(payload.BaggageID, payload.DeliveryID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Baggage delivery updated successfully"})
}
func (h *Handler) RemoveBaggageFromDelivery(c *gin.Context) {
    // Получаем данные из JSON-запроса
    var payload struct {
        BaggageID  uint `json:"baggage_id"`
        DeliveryID uint `json:"delivery_id"`
    }

    if err := c.BindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Попытка удаления связи между багажом и доставкой в репозитории
    err := h.Repo.RemoveBaggageFromDelivery(payload.BaggageID, payload.DeliveryID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Baggage removed from delivery successfully"})
}


