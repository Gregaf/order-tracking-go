package models

import (
	"net/mail"

	"github.com/gregaf/order-tracking/internal/validation"
)

type Email string

func (e Email) Validate() []error {
	var errors = make([]error, 0)

	if len(e) == 0 {
		errors = append(errors, &validation.ErrValidationEmptyString{})
	}

	_, err := mail.ParseAddress(string(e))
	if err != nil {
		errors = append(errors, &validation.ErrValidationInvalidEmailFormat{})
	}

	return errors
}
