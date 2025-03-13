package route

import (
	"github.com/gin-gonic/gin"
	"github.com/jeremyseow/rates-service/server/handler"
)

func SetupRoutes(router *gin.Engine, handlers *handler.Handlers) {
	router.GET("/rates", handlers.GetRates)
	router.GET("/user-agent", handlers.GetUserAgent)
	router.POST("/files/:content", handlers.PostFile)
	router.POST("/webhook", handlers.PostWebHook)
}
