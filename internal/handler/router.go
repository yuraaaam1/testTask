package handler

import "github.com/gin-gonic/gin"

func NewRouter(sh *SubscriptionHandler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")

	subscriptions := api.Group("/subscriptions")

	subscriptions.POST("", sh.CreateHandler)
	subscriptions.GET("/:id", sh.GetByIDHandler)
	subscriptions.GET("", sh.ListHandler)
	subscriptions.PUT("/:id", sh.UpdateHandler)
	subscriptions.DELETE("/:id", sh.DeleteHandler)

	return router
}
