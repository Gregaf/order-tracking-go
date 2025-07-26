package user

import (
	"context"

	"gregaf/order-tracking-go/internal/dto"
	"gregaf/order-tracking-go/internal/models"
	"gregaf/order-tracking-go/internal/transport/http/middleware"
)

type UserService interface {
	Upsert(ctx context.Context, authCtx middleware.AuthContext, userDto dto.UpsertUserDTO) (*models.User, error)
	CreateUser(ctx context.Context, authCtx middleware.AuthContext, userDto dto.CreateUserDTO) (*models.User, error)
	GetUserByID(ctx context.Context, authCtx middleware.AuthContext, userID string) (*models.User, error)
	DeleteUserByID(ctx context.Context, authCtx middleware.AuthContext, id string) error
}
