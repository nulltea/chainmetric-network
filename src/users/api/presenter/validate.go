package presenter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// ValidationErrorsResponse stores validation errors by struct fields.
//
// swagger:model
type ValidationErrorsResponse map[string]string

// PresentValidation handles validation error and sends its according representation to the API user.
func PresentValidation(ctx *gin.Context, err error) {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "failed to validate request"))
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		errors := make(ValidationErrorsResponse)
		for _, fe := range errs {
			errors[fe.Field()] = fe.Error()
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, errors)
	}
}
