package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"unicode"

	"gregaf/order-tracking-go/internal/dto"
	"gregaf/order-tracking-go/internal/models"
	"gregaf/order-tracking-go/internal/transport/http/middleware"
	"gregaf/order-tracking-go/internal/user"
	"gregaf/order-tracking-go/internal/util"
)

type UserServiceCore struct {
	userRepo user.UserRepository
	logger   *slog.Logger
}

func NewUserService(logger *slog.Logger, userRepo user.UserRepository) *UserServiceCore {
	return &UserServiceCore{userRepo: userRepo, logger: logger}
}

func (s *UserServiceCore) Upsert(ctx context.Context, authCtx middleware.AuthContext, userDto dto.UpsertUserDTO) (*models.User, error) {
	if userDto.ID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	updatedAtDate := time.Now().UnixMilli()

	user := models.User{
		ID:            userDto.ID,
		FirstName:     userDto.FirstName,
		LastName:      userDto.LastName,
		Email:         userDto.Email,
		UpdatedAtDate: updatedAtDate,
	}

	err := s.userRepo.UpsertUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserServiceCore) CreateUser(ctx context.Context, authCtx middleware.AuthContext, userDto dto.CreateUserDTO) (*models.User, error) {
	errValidation := userDto.Validate()
	if errValidation != nil {
		return nil, errValidation
	}

	displayID, err := generateDisplayID(userDto.FirstName, userDto.LastName)
	if err != nil {
		return nil, err
	}

	createdAtDate := time.Now().UnixMilli()

	newUser := models.User{
		ID:            authCtx.RequestorID,
		DisplayID:     displayID,
		FirstName:     userDto.FirstName,
		LastName:      userDto.LastName,
		Email:         userDto.Email,
		CreatedAtDate: createdAtDate,
		UpdatedAtDate: createdAtDate,
	}

	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *UserServiceCore) GetUserByID(ctx context.Context, authCtx middleware.AuthContext, userID string) (*models.User, error) {
	if !middleware.HasPermission(authCtx.Permissions, authCtx.RequestorID, userID, "user", "read") {
		return nil, middleware.ErrNotAuthorized
	}

	s.logger.Debug("User authorized", "userID", userID, "Permissions", authCtx.Permissions)

	foundUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if foundUser == nil {
		return nil, user.ErrUserNotFound
	}

	return foundUser, nil
}

func (s *UserServiceCore) DeleteUserByID(ctx context.Context, authCtx middleware.AuthContext, id string) error {
	if !middleware.HasPermission(authCtx.Permissions, authCtx.RequestorID, id, "user", "delete") {
		return middleware.ErrNotAuthorized
	}

	return s.userRepo.DeleteUser(ctx, id)
}

func generateDisplayID(firstName, lastName string) (string, error) {
	if len(firstName) == 0 || len(lastName) == 0 {
		return "", errors.New("first name and last name cannot be empty")
	}

	// Extract initials (first letter of first name and last name)
	firstInitial := string(unicode.ToUpper(rune(firstName[0])))
	lastInitial := string(unicode.ToUpper(rune(lastName[0])))

	// Generate two random 4-character segments
	randomSegment1, err := util.GenerateRandomSegment(4)
	if err != nil {
		return "", fmt.Errorf("failed to generate random segment 1: %w", err)
	}

	randomSegment2, err := util.GenerateRandomSegment(4)
	if err != nil {
		return "", fmt.Errorf("failed to generate random segment 2: %w", err)
	}

	customerID := fmt.Sprintf("%s%s-%s-%s", firstInitial, lastInitial, randomSegment1, randomSegment2)

	return customerID, nil
}
