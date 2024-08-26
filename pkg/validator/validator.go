package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

const (
	usernameMinLength = 3
	usernameMaxLength = 32
	defaultLimit      = 40
)

var (
	usernameRegex = regexp.MustCompile(fmt.Sprintf(`^[a-z0-9_]{%d,%d}$`, usernameMinLength, usernameMaxLength))
)

type Validator interface {
	Validate(i interface{}) error
}

type valid struct {
	v *validator.Validate
}

func NewValidator() (Validator, error) {
	v := validator.New()

	if err := v.RegisterValidation("username", usernameValidate); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("value", valueValidate); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("limit", limitValidate); err != nil {
		return nil, err
	}

	return &valid{v: v}, nil
}

func (v *valid) Validate(i interface{}) error {
	if err := v.v.Struct(i); err != nil {
		return validateError(err.(validator.ValidationErrors)[0])
	}
	return nil
}

func validateError(err validator.FieldError) error {
	switch err.Tag() {
	case "username":
		return fmt.Errorf("field username can only consist of lower Latin characters, numbers and underscore symbol. Min length is 3, max: 32, your input: %s", err.Value())
	case "value":
		return fmt.Errorf("field value is incorrect. Value must be >= 0, got %d", err.Value())
	case "limit":
		return fmt.Errorf("field limit is incorrect. Value must be >= 0 and <= 20, got %d", err.Value())
	default:
		return fmt.Errorf("field %s is required", err.Field())
	}
}

func usernameValidate(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.String {
		return false
	}
	return usernameRegex.MatchString(fl.Field().String())
}

func valueValidate(fl validator.FieldLevel) bool {
	return fl.Field().Int() >= 0
}

func limitValidate(fl validator.FieldLevel) bool {
	v := fl.Field().Int()
	return v > 0 && v <= defaultLimit
}
