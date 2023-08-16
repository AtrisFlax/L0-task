package api

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (item Item) Validate() error {
	if err := validate.Struct(item); err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return err
		}
		return err
	}
	if err := validate.Var(
		item.DateCreated.String(),
		"required,datetime=2006-01-02 15:04:05.999999999 -0700 MST"); err != nil {
		return err
	}

	return nil
}
