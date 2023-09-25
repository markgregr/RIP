package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)


func StartServer() {
	log.Println("Server start up")

	file, err := os.Open("resources/data/baggage.json")
	if err != nil {
		log.Println("Ошибка при открытии JSON файла:", err)
		return
	}
	defer file.Close()

	var baggages []Baggage
	log.Println(baggages)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&baggages); err != nil {
		log.Println("Ошибка при декодировании JSON данных:", err)
		return
	}
	
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
		searchQuery := c.DefaultQuery("q", "")
		var foundBaggages []Baggage
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
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			// Обработка ошибки
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		baggage := baggages[id-1]
		c.HTML(http.StatusOK, "card.tmpl", baggage)
	})

	r.Run()

	log.Println("Server down")
}
