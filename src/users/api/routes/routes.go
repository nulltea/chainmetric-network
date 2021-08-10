package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine) {
	engine.GET("", apiIndex)
	engine.GET("/health", healthCheck)

	Routes(engine.Group("/auth"))
	Routes(engine.Group("/users"))
}

func apiIndex(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Chainmetric Identity API")
}

func healthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Healthy")
}
