package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/timoth-y/chainmetric-contracts/src/users/api/docs"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/routes/auth"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/routes/users"
)

func Setup(engine *gin.Engine) {
	engine.GET("", apiIndex)
	engine.GET("/health", healthCheck)

	auth.Routes(engine.Group("/users"))
	users.Routes(engine.Group("/users"))

	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func apiIndex(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Chainmetric Identity API")
}

func healthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Healthy")
}
