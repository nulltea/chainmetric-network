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

// handleRegister ...
// @Summary User registration
// @Description Request user initial registration.
// @Tags users
// @Accept json
// @Produce json
// @Param register body user.RegistrationRequest true "Request to register new user"
// @Success 200 {object} user.User
// @Failure 400 {object} presenter.HTTPError
// @Failure 500 {object} presenter.HTTPError
// @Router /users/register [get]
func handleRegister(ctx *gin.Context) {
	var req model.RegistrationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		presenter.AbortWithError(ctx, http.StatusBadRequest, errors.Wrap(err, "invalid request structure"))
	}

	if err := validator.New().Struct(req); err != nil {
		presenter.PresentValidation(ctx, err)
	}

	user, err := identity.Register(req)
	if err != nil {
		presenter.AbortWithError(ctx, http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, user)
}
