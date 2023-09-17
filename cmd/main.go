package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Baggage struct {
	BaggageID     int
	OwnerName     string
	Destination   string
	Weight        float64
	Status        string
	Src           string
	Href          string
	BaggageType   string
	Width, Height int
	Length        int
}

func main() {
	log.Println("Server start up")

	file, err := os.Open("resources/data/baggage.json")
	if err != nil {
		log.Println("Ошибка при открытии JSON файла:", err)
		return
	}
	defer file.Close()

	var baggages []Baggage

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&baggages); err != nil {
		log.Println("Ошибка при декодировании JSON данных:", err)
		return
	}
	log.Println(baggages)
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/css", "./resources/css")
	r.Static("/data", "./resources/data")
	r.Static("/images", "./resources/images")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		data := gin.H{
			"baggages": baggages,
		}
		c.HTML(http.StatusOK, "index.tmpl", data)
	})

	r.Run()

	log.Println("Server down")
}
