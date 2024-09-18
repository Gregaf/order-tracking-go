package models

import (
	"errors"
	"testing"

	"github.com/gregaf/order-tracking/internal/validation"
)

func TestEmailValidate(t *testing.T) {
	testCases := []struct {
		email    Email
		expected []error
	}{
		{email: "not-a-valid-email", expected: []error{&validation.ErrValidationInvalidEmailFormat{}}},
	}

	for _, tc := range testCases {
		errs := tc.email.Validate()
		if len(errs) != len(tc.expected) {
			t.Errorf("Expected '%d' errors, got '%d'", len(tc.expected), len(errs))
		}
		for i, err := range errs {
			if !errors.Is(err, tc.expected[i]) {
				t.Errorf("Expected error '%v', got '%v'", tc.expected[i], err)
			}
		}
	}
}
