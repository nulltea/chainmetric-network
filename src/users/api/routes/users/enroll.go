package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/presenter"
)

func enroll(ctx *gin.Context) {
	var (
		req presenter.RegistrationRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "invalid request structure"))
	}


}
