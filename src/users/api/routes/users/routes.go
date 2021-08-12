package users

import "github.com/gin-gonic/gin"

// Routes defines routes for user registration and management.
func Routes(routes *gin.RouterGroup) {
	routes.GET("/", handleGetUser)
	routes.POST("/register", handleRegister)
	routes.POST("/enroll", handleEnroll)
}


