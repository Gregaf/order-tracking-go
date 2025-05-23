package dto

import "github.com/gregaf/order-tracking-go/internal/types"

type CreateUserDTO struct {
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Email     types.Email `json:"email"`
}

func (cu *CreateUserDTO) Validate() error {
	return nil
}
