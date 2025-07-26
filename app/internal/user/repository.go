package user

import (
	"context"
	"gregaf/order-tracking-go/internal/models"
)

type UserRepository interface {
	UpsertUser(ctx context.Context, user models.User) error
	CreateUser(ctx context.Context, user models.User) error
	GetUserByID(ctx context.Context, ID string) (*models.User, error)
	// UpdateUser(user *models.User) error
	DeleteUser(ctx context.Context, id string) error
}
