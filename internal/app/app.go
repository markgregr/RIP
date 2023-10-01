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
func (app *Application) Run(){

	handler := api.NewHandler(app.Repository);


    r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/css", "./resources/css")
	r.Static("/data", "./resources/data")
	r.Static("/images", "./resources/images")
	r.Static("/fonts", "./resources/fonts")
    
	
    //методы для багажа
	r.GET("/", handler.GetBaggages)
	r.GET("/baggage/:id", handler.GetBaggageByID)
	r.DELETE("/baggage/:id/delete", handler.DeleteBaggage)
	r.POST("/create", handler.CreateBaggage)
    r.PUT("/baggage/:id/update", handler.UpdateBaggage)

    //методы для заявок 
    r.GET("/deliveries", handler.GetDeliveries)
    r.GET("/delivery/:id", handler.GetDeliveryByID)
    r.DELETE("/delivery/:id/delete", handler.DeleteDelivery)
    r.POST("/deliveries/create", handler.CreateDelivery)
    r.PUT("/delivery/:id/update", handler.UpdateDelivery)

    //методы для пользователей
    r.GET("/users", handler.GetUsers)
    r.GET("/user/:id", handler.GetUserByID)
    r.DELETE("/user/:id/delete", handler.DeleteUser)
    r.POST("/users/create", handler.CreateUser)
    r.PUT("/user/:id/update", handler.UpdateUser)

    
    addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
    r.Run(addr)
	log.Println("Server down")
}
