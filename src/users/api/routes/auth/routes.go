package auth

import "github.com/gin-gonic/gin"

// Routes defines routes for users authentication.
func Routes(routes *gin.RouterGroup) {
	routes.POST("/", authenticate)
}
