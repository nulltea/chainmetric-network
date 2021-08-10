package api

import (
	"github.com/gin-gonic/gin"
	"github.com/timoth-y/chainmetric-contracts/shared"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/routes"
)

func init() {
	shared.InitCore()
	shared.InitMongoDB()
}

func main() {
	engine := gin.Default()
	routes.Setup(engine)

	engine.Run(":8080")
}
