package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	model "github.com/timoth-y/chainmetric-contracts/model/user"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/presenter"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/usecase/identity"
)

func handleEnroll(ctx *gin.Context) {
	var req model.EnrollmentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.Wrap(err, "invalid request structure"))
	}

	if err := validator.New().Struct(req); err != nil {
		presenter.PresentValidation(ctx, err)
	}

	if err := identity.Enroll(req); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.Status(http.StatusOK)
}
