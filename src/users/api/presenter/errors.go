package presenter

import "github.com/gin-gonic/gin"

type HTTPError struct {
	Code   int         `json:"code" example:"400"`
	Reason interface{} `json:"reason"`
}

func AbortWithError(ctx *gin.Context, status int, reason interface{}) {
	ctx.AbortWithStatusJSON(status, HTTPError{
		Code:    status,
		Reason: reason,
	})
}
