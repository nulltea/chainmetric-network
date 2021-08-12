package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	model "github.com/timoth-y/chainmetric-contracts/model/user"
	"github.com/timoth-y/chainmetric-contracts/src/users/api/presenter"
	"github.com/timoth-y/chainmetric-contracts/src/users/usecase/identity"
)

// handleEnroll ...
// @Summary Enroll new user
// @Description Generates signing cryptographic identity for user and confirms one.
// @Tags users
// @Accept json
// @Produce json
// @Param register body user.EnrollmentRequest true "Request to enroll new user"
// @Success 200
// @Failure 400 {object} presenter.HTTPError
// @Failure 500 {object} presenter.HTTPError
// @Router /users/enroll [get]
func handleEnroll(ctx *gin.Context) {
	var req model.EnrollmentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.AbortWithError(ctx, http.StatusBadRequest, errors.Wrap(err, "invalid request structure"))
		return
	}

	if err := validator.New().Struct(req); err != nil {
		presenter.PresentValidation(ctx, err)
		return
	}

	if err := identity.Enroll(req); err != nil {
		presenter.AbortWithError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
