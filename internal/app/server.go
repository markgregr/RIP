package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/markgregr/RIP/docs"
	"github.com/markgregr/RIP/internal/pkg/middleware"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func (app *Application) Run() {
    r := gin.Default()  
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
        BaggageGroup.GET("/", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetBaggages)
        BaggageGroup.GET("/:baggage_id", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetBaggageByID)
        BaggageGroup.DELETE("/:baggage_id", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.DeleteBaggage)
        BaggageGroup.POST("/", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.CreateBaggage)
        BaggageGroup.PUT("/:baggage_id", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.UpdateBaggage)
        BaggageGroup.POST("/:baggage_id/delivery", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.AddBaggageToDelivery)
        BaggageGroup.DELETE("/:baggage_id/delivery", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.RemoveBaggageFromDelivery)
        BaggageGroup.POST("/:baggage_id/image", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.AddBaggageImage)
    }
    

    // Группа запросов для доставки
    DeliveryGroup := r.Group("/delivery").Use(middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository))
    {
        DeliveryGroup.GET("/", app.Handler.GetDeliveries)
        DeliveryGroup.GET("/:delivery_id", app.Handler.GetDeliveryByID)
        DeliveryGroup.DELETE("/:delivery_id", app.Handler.DeleteDelivery)
        DeliveryGroup.PUT("/:delivery_id", app.Handler.UpdateDeliveryFlightNumber)
        DeliveryGroup.PUT("/:delivery_id/status/user", app.Handler.UpdateDeliveryStatusUser)  // Новый маршрут для обновления статуса доставки пользователем
        DeliveryGroup.PUT("/:delivery_id/status/moderator", app.Handler.UpdateDeliveryStatusModerator)  // Новый маршрут для обновления статуса доставки модератором
    }

    UserGroup := r.Group("/user")
    {
        UserGroup.POST("/register", app.Handler.Register)
        UserGroup.POST("/login", app.Handler.Login)
        UserGroup.POST("/logout", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.Logout)
    }
    addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
    r.Run(addr)
    log.Println("Server down")
}