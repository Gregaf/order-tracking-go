package dynamodbrepository

import (
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gregaf/order-tracking/internal/models"
	"github.com/gregaf/order-tracking/internal/repository"
)

type DynamoDbUserRepository struct {
	db     *dynamodb.Client
	logger *slog.Logger
}

// CreateUser implements repository.UserRepository.
func (d *DynamoDbUserRepository) CreateUser(user *models.User) error {
	panic("unimplemented")
}

// DeleteUser implements repository.UserRepository.
func (d *DynamoDbUserRepository) DeleteUser(id string) error {
	panic("unimplemented")
}

// GetUserByID implements repository.UserRepository.
func (d *DynamoDbUserRepository) GetUserByID(id string) (*models.User, error) {
	panic("unimplemented")
}

// UpdateUser implements repository.UserRepository.
func (d *DynamoDbUserRepository) UpdateUser(user *models.User) error {
	panic("unimplemented")
}

func NewDynamoDbUserRepository(cfg aws.Config) repository.UserRepository {
	svc := dynamodb.NewFromConfig(cfg, nil)

	return &DynamoDbUserRepository{
		db: svc,
	}
}
