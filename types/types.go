package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/WilliamTrojniak/TabAppBackend/services"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

	Validate.RegisterStructValidation(TabUpdateStructLevelValidation, TabUpdate{})
	Validate.RegisterValidation("future", dateFutureValidation)
}

func dateFutureValidation(f1 validator.FieldLevel) bool {
	date := f1.Field().Interface().(Date)
	return date.After(civil.DateOf(time.Now().AddDate(0, 0, -2)))
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
		logger.Debug("Data validation failed", "errors", errors)
		return services.NewValidationServiceError(err, errors)

	}

	return nil
}

func ReadRequestJson(r *http.Request, dest interface{}) error {
	mediaType := getMediaType(r)
	if mediaType != "application/json" {
		return services.NewServiceError(nil, http.StatusUnsupportedMediaType, nil)
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dest)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxError):
			return services.NewServiceError(err, http.StatusBadRequest, fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset))
		case errors.As(err, &unmarshalTypeError):
			return services.NewServiceError(err, http.StatusBadRequest, fmt.Sprintf("Request body contains an invalid value for the %q field of type %v (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Offset))
		case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
			return services.NewServiceError(err, http.StatusBadRequest, fmt.Sprintf("Unexpected EOF in request body"))
		case uuid.IsInvalidLengthError(err):
			return services.NewServiceError(err, http.StatusBadRequest, fmt.Sprintf("Invalid UUID value"))
		default:
			slog.Warn("Issue reading JSON from request body", "err", err)
			return services.NewInternalServiceError(err)
		}
	}

	return nil
}

func getMediaType(r *http.Request) string {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return ""
	}

	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	return mediaType

}
