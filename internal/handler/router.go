package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/yuraaaam1/testTask/docs"
)

func NewRouter(sh *SubscriptionHandler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")

	api.GET("/subscriptions/cost", sh.CalculateTotalCostHandler)

	subscriptions := api.Group("/subscriptions")
	subscriptions.POST("", sh.CreateHandler)
	subscriptions.GET("/:id", sh.GetByIDHandler)
	subscriptions.GET("", sh.ListHandler)
	subscriptions.PUT("/:id", sh.UpdateHandler)
	subscriptions.DELETE("/:id", sh.DeleteHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
