package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	model "github.com/timoth-y/chainmetric-contracts/model/user"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/usecase/identity"
)

func register(ctx *gin.Context) {
	var (
		req model.RegistrationRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "invalid request structure"))
	}

	user, err := identity.Register(req)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, user)
}
