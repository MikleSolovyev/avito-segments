package validator

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	v *validator.Validate
}

func New() *CustomValidator {
	v := validator.New(validator.WithRequiredStructEnabled())
	return &CustomValidator{
		v: v,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.v.Struct(i)
}
