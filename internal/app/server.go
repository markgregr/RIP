package app

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/docs"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// Run запускает приложение.
func (app *Application) Run() {
    r := gin.Default()
    r.Use(cors.Default())  
    // Это нужно для автоматического создания папки "docs" в вашем проекте
    docs.SwaggerInfo.Title = "BagTracker RestAPI"
    docs.SwaggerInfo.Description = "API server for BagTracker application"
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.Host = "localhost:8081"
    docs.SwaggerInfo.BasePath = "/"
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    // Группа запросов для багажа
    BaggageGroup := r.Group("/baggage")
    {   
        BaggageGroup.GET("/", app.Handler.GetBaggages)
        BaggageGroup.GET("/:baggage_id", app.Handler.GetBaggageByID) 
        BaggageGroup.DELETE("/:baggage_id/delete", app.Handler.DeleteBaggage) 
        BaggageGroup.POST("/create", app.Handler.CreateBaggage)
        BaggageGroup.PUT("/:baggage_id/update", app.Handler.UpdateBaggage) 
        BaggageGroup.POST("/:baggage_id/delivery", app.Handler.AddBaggageToDelivery) 
        BaggageGroup.DELETE("/:baggage_id/delivery/delete", app.Handler.RemoveBaggageFromDelivery)
        BaggageGroup.POST("/:baggage_id/image",app.Handler.AddBaggageImage)
    }
    
    // Группа запросов для доставки
    DeliveryGroup := r.Group("/delivery")
    {
        DeliveryGroup.GET("/", app.Handler.GetDeliveries)
        DeliveryGroup.GET("/:id", app.Handler.GetDeliveryByID)
        DeliveryGroup.DELETE("/:id/delete", app.Handler.DeleteDelivery)
        DeliveryGroup.PUT("/:id/update", app.Handler.UpdateDeliveryFlightNumber)
        DeliveryGroup.PUT("/:id/status/user", app.Handler.UpdateDeliveryStatusUser)  // Новый маршрут для обновления статуса доставки пользователем
        DeliveryGroup.PUT("/:id/status/moderator", app.Handler.UpdateDeliveryStatusModerator)  // Новый маршрут для обновления статуса доставки модератором
    }

    addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
    r.Run(addr)
    log.Println("Server down")
}