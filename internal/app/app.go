package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/app/config"
	"github.com/markgregr/RIP/internal/app/dsn"
	"github.com/markgregr/RIP/internal/app/repository"
)

// Application представляет основное приложение.
type Application struct {
    Config       *config.Config
    Repository   *repository.Repository
    RequestLimit int
}

// New создает новый объект Application и настраивает его.
func New(ctx context.Context) (*Application, error) {
    // Инициализируйте конфигурацию
    cfg, err := config.NewConfig(ctx)
    if err != nil {
        return nil, err
    }

    // Инициализируйте подключение к базе данных (DB)
    repo, err := repository.New(dsn.FromEnv())
    if err != nil {
        return nil, err
    }

    // Инициализируйте и настройте объект Application
    app := &Application{
        Config: cfg,
        Repository: repo,
        // Установите другие параметры вашего приложения, если необходимо
    }

    return app, nil
}

// Run запускает приложение.
func (app *Application) Run(){
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
		
		searchCode := c.DefaultQuery("searchCode", "")

		baggages, err := app.Repository.GetBaggages(searchCode)
		if err != nil {
			log.Println("Error Repository method GetAll:", err)
			return
		}
		data := gin.H{
			"baggages": baggages,
			"searchCode": searchCode,
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

		baggage, err := app.Repository.GetBaggageByID(baggage_id)
		if err != nil {
			// Обработка ошибки
			c.JSON(http.StatusBadRequest, gin.H{"error": "GetBaggageByID"})
			return
		}

		c.HTML(http.StatusOK, "card.tmpl", baggage)
	})

	r.POST("/baggage/:id/delete", func(c *gin.Context) {

		baggage_id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			// Обработка ошибки
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		app.Repository.DeleteBaggage(baggage_id)
		c.Redirect(http.StatusMovedPermanently, "/")
	})
    
	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
    r.Run(addr)
	log.Println("Server down")
}
