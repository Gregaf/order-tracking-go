package userservice

import (
	"errors"
	"time"

	"github.com/gregaf/order-tracking/internal/models"
	"github.com/gregaf/order-tracking/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user *models.User) error {
	// validate user data...

	user.CreatedAtDate = time.Now().UnixMilli()
	user.UpdatedAtDate = user.CreatedAtDate

	return s.userRepo.CreateUser(user)
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	// validate user data...

	return s.userRepo.UpdateUser(user)
}

func (s *UserService) DeleteUser(id string) error {
	return s.userRepo.DeleteUser(id)
}
