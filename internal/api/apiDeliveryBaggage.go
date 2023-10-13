package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)
func (h *Handler) AddBaggageToDelivery(c *gin.Context) {
    searchCode := c.DefaultQuery("searchCode", "")
    // Получаем параметры из URL
    baggageID, err := strconv.Atoi(c.Param("baggage_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid baggage_id"})
        return
    }

    deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery_id"})
        return
    }

    // Попытка обновления связи между багажом и доставкой в репозитории
    err = h.Repo.AddBaggageToDelivery(uint(baggageID), uint(deliveryID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    // Получаем обновленный список багажей
	baggages, err := h.Repo.GetBaggages(searchCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Baggage delivery update successfully", "baggages": baggages})
}

func (h *Handler) RemoveBaggageFromDelivery(c *gin.Context) {
    searchCode := c.DefaultQuery("searchCode", "")
    var err error  // Объявляем переменную здесь

    // Получаем параметры из URL
    baggageID, err := strconv.Atoi(c.Param("baggage_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid baggage_id"})
        return
    }
    deliveryID, err := strconv.Atoi(c.Param("delivery_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery_id"})
        return
    }

    // Попытка удаления связи между багажом и доставкой в репозитории
    err = h.Repo.RemoveBaggageFromDelivery(uint(baggageID), uint(deliveryID))  // Используем объявленную переменную err
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    baggages, err := h.Repo.GetBaggages(searchCode)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Baggage removed from delivery successfully", "baggages": baggages})
}
