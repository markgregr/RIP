package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/app/ds"
)

//методы для таблицы baggage
func (h *Handler) GetBaggages(c *gin.Context) {
	baggages, err := h.Repo.GetBaggages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	searchQuery := c.DefaultQuery("searchQuery", "")
	var foundBaggages []ds.Baggage
	for _, baggage := range baggages {
		if strings.HasPrefix(strings.ToLower(baggage.BaggageCode), strings.ToLower(searchQuery)) {
			foundBaggages = append(foundBaggages, baggage)
		}
	}
	//c.JSON(http.StatusOK, gin.H{"baggages": foundBaggages})
	data := gin.H{"baggages": foundBaggages}
	c.HTML(http.StatusOK,"index.tmpl", data)
}

func (h *Handler) GetBaggageByID(c *gin.Context) {
	baggageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	baggage, err := h.Repo.GetBaggageByID(baggageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GetBaggageByID"})
		return
	}
	c.HTML(http.StatusOK, "card.tmpl", baggage)
}
func (h *Handler) CreateBaggage(c *gin.Context) {
	var baggage ds.Baggage
	if err := c.BindJSON(&baggage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Repo.CreateBaggage(&baggage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Baggage created successfully"})
}
func (h *Handler) DeleteBaggage(c *gin.Context) {
	baggageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	err = h.Repo.DeleteBaggage(baggageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Baggage deleted successfully"})
}
func (h *Handler) UpdateBaggage(c *gin.Context) {
	baggageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}

	var updatedBaggage ds.Baggage
	if err := c.BindJSON(&updatedBaggage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Repo.UpdateBaggage(baggageID, &updatedBaggage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Baggage updated successfully"})
}
