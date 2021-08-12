package main

import (
	"github.com/gin-gonic/gin"
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	"github.com/timoth-y/chainmetric-contracts/shared/utils"
	"github.com/timoth-y/chainmetric-contracts/src/users/api"
	"github.com/timoth-y/chainmetric-contracts/src/users/usecase/identity"
)

func init() {
	core.InitCore()
	core.InitMongoDB()
	utils.MustExecute(identity.Init, "failed to initialize identity package")
}

func main() {
	engine := gin.Default()
	api.Setup(engine)

	_ = engine.Run(":8080")
}
