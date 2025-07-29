package dto

import "gregaf/order-tracking-go/internal/types"

type SyncUserDTO struct {
	ID        string
	FirstName string
	LastName  string
	Email     types.Email
}

type CreateUserDTO struct {
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Email     types.Email `json:"email"`
}

func (cu *CreateUserDTO) Validate() error {
	return nil
}
