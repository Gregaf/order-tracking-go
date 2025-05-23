package user

import (
	"context"

	"github.com/gregaf/order-tracking-go/internal/dto"
	"github.com/gregaf/order-tracking-go/internal/models"
	"github.com/gregaf/order-tracking-go/internal/transport/http/middleware"
)

type UserService interface {
	CreateUser(ctx context.Context, authCtx middleware.AuthContext, userDto dto.CreateUserDTO) (*models.User, error)
	GetUserByID(ctx context.Context, authCtx middleware.AuthContext, userID string) (*models.User, error)
	DeleteUserByID(ctx context.Context, authCtx middleware.AuthContext, id string) error
}
