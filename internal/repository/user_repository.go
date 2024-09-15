package repository

import "github.com/gregaf/order-tracking/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
}
