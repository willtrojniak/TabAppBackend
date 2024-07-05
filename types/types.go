package types

import (
	"log/slog"
	"reflect"
	"strings"

	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateData(data interface{}, logger *slog.Logger) error {
	err := Validate.Struct(data)

	if err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			logger.Error("Error while attempting to validate user")
			return services.NewInternalServiceError(err)
		}

		errors := make(services.ValidationErrors)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = services.ValidationError{Value: err.Value(), Error: err.Tag()}
		}
		logger.Debug("User validation failed...", "errors", errors)
		return services.NewValidationServiceError(err, errors)

	}

	return nil
}
