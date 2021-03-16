package ports

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func validateWithTags(ko keyOpts) error {
	err := validate.Struct(ko)

	return err
}

func validateWithFunc(ko keyOpts) error {
	err := validate.Var(ko.Scope, "required,min=1,max=50")
	err = validate.Var(ko.Expiration, "required")

	return err
}

func validateOzzo(ko keyOpts) error {
	err := validation.Validate(ko.Scope,
		validation.Required,
		validation.Length(1, 50),
	)
	err = validation.Validate(ko.Expiration,
		validation.Required,
		validation.Date(time.RFC3339),
	)

	return err
}

func validateHomeSolution(ko keyOpts) error {
	if len(ko.Scope) < 1 {
		return errors.New("scope must be longer than 0")
	}
	if len(ko.Scope) > 50 {
		return errors.New("scope must be shorter than 50")
	}

	return nil
}
