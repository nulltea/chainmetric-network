package main

import (
	"github.com/gin-gonic/gin"
	"github.com/timoth-y/chainmetric-contracts/shared"
	"github.com/timoth-y/chainmetric-contracts/src/users/api"
	"github.com/timoth-y/chainmetric-contracts/src/users/usecase/identity"
)

func init() {
	shared.InitCore()
	shared.InitMongoDB()
	identity.Init()
}

// @title Chainmetric Identity API
// @version 1.0
// @description Chainmetric Identity API users to authenticate and save preferences.

// @contact.name Timothy Yalugin
// @contact.url https://github.com/timoth-y
// @contact.email timauthx@gmail.com

// @license.name Apache 2.0
// @license.url https://raw.githubusercontent.com/timoth-y/chainmetric-contracts/main/LICENSE

// @host identity.chainmetric.network
// @BasePath /
func main() {
	engine := gin.Default()
	api.Setup(engine)

	engine.Run(":8080")
}
