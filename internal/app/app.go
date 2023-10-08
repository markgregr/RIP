package app

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/internal/api"
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
func (app *Application) Run() {
    handler := api.NewHandler(app.Repository)
    r := gin.Default()

    r.LoadHTMLGlob("templates/*")
    r.Static("/css", "./resources/css")
    r.Static("/data", "./resources/data")
    r.Static("/images", "./resources/images")
    r.Static("/fonts", "./resources/fonts")

    // Группа запросов для багажа
    baggageGroup := r.Group("/baggage")
    {
        baggageGroup.GET("/", handler.GetBaggages)
        baggageGroup.GET("/:id", handler.GetBaggageByID)
        baggageGroup.DELETE("/:id/delete", handler.DeleteBaggage)
        baggageGroup.POST("/create", handler.CreateBaggage)
        baggageGroup.PUT("/:id/update", handler.UpdateBaggage)
    }

    // Группа запросов для заявок
    deliveryGroup := r.Group("/deliveries")
    {
        deliveryGroup.GET("/", handler.GetDeliveries)
        deliveryGroup.GET("/:id", handler.GetDeliveryByID)
        deliveryGroup.DELETE("/:id/delete", handler.DeleteDelivery)
        deliveryGroup.PUT("/:id/update", handler.UpdateDelivery)
    }
    

    addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
    r.Run(addr)
    log.Println("Server down")
}

