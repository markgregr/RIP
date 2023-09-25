package app

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/app/ds"
)

func (a *Application) StartServer() {
	
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/css", "./resources/css")
	r.Static("/data", "./resources/data")
	r.Static("/images", "./resources/images")
	r.Static("/fonts", "./resources/fonts")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		baggages, err := a.Repository.GetAllBaggage()
		if err != nil {
			log.Println("Error Repository method GetAll:", err)
			return
		}
		log.Println(baggages)
		searchQuery := c.DefaultQuery("q", "")
		
		var foundBaggages []ds.Baggage
		for _, baggage := range baggages {
			if strings.HasPrefix(strings.ToLower(baggage.BaggageCode), strings.ToLower(searchQuery)) {
				foundBaggages = append(foundBaggages, baggage)
			}
		}
		data := gin.H{
			"baggages": foundBaggages,
		}
		c.HTML(http.StatusOK, "index.tmpl", data)
	})

	r.GET("/baggage/:id", func(c *gin.Context) {
		baggage_id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			// Обработка ошибки
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		baggage, err := a.Repository.GetBaggageByID(baggage_id)
		if err != nil {
			// Обработка ошибки
			c.JSON(http.StatusBadRequest, gin.H{"error": "GetBaggageByID"})
			return
		}

		c.HTML(http.StatusOK, "card.tmpl", baggage)
	})

	r.GET("/baggage/:id/delete", func(c *gin.Context) {
		baggage_id, err := strconv.Atoi(c.Param("id"))
		log.Println(baggage_id)
		if err != nil {
			// Обработка ошибки
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		a.Repository.DeleteBaggage(baggage_id)

		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.Run()

	log.Println("Server down")
}
