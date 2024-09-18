package validation

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Errors map[string][]error
}

func (ve *ValidationError) Add(field string, err error) {
	if ve.Errors == nil {
		ve.Errors = make(map[string][]error)
	}
	ve.Errors[field] = append(ve.Errors[field], err)
}

func (ve *ValidationError) Error() string {
	if len(ve.Errors) == 0 {
		return "no validation errors"
	}

	var sb strings.Builder
	for field, errors := range ve.Errors {
		sb.WriteString(fmt.Sprintf("%s: ", field))
		for _, err := range errors {
			sb.WriteString(err.Error() + "; ")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (ve *ValidationError) HasErrors() bool {
	return len(ve.Errors) > 0
}

type ErrValidationEmptyString struct{}

func (e *ErrValidationEmptyString) Error() string {
	return "field required, must contain a value"
}
