package validation

type ErrValidationInvalidEmailFormat struct{}

func (e *ErrValidationInvalidEmailFormat) Error() string {
	return "invalid email format"
}
