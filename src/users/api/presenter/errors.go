package presenter

import "github.com/gin-gonic/gin"

type HTTPError struct {
	Code   int    `json:"code" example:"400"`
	Reason string `json:"reason"`
}

type HTTPMultiError struct {
	Code   int         `json:"code" example:"400"`
	Reason interface{} `json:"reason"`
}

func AbortWithError(ctx *gin.Context, status int, reason interface{}) {
	switch r := reason.(type) {
	case error:
		ctx.AbortWithStatusJSON(status, HTTPError{
			Code:    status,
			Reason: r.Error(),
		})
	default:
		ctx.AbortWithStatusJSON(status, HTTPMultiError{
			Code:    status,
			Reason: r,
		})
	}
}
