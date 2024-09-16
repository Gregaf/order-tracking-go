package dynamodbrepository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gregaf/order-tracking/internal/models"
	"github.com/gregaf/order-tracking/internal/repository"
)

type DynamoDbUserRepository struct {
	db     *dynamodb.Client
	logger *slog.Logger
}

type DynamoDbUser struct {
	Pk string `dynamodbav:"Pk"`
	Sk string `dynamodbav:"Sk"`
	models.User
}

// CreateUser implements repository.UserRepository.
func (d *DynamoDbUserRepository) CreateUser(user *models.User) error {
	r := DynamoDbUser{
		Pk:   fmt.Sprintf("USER#%s", user.Sub),
		Sk:   "METADATA",
		User: *user,
	}
	av, err := attributevalue.MarshalMap(r)
	if err != nil {
		return fmt.Errorf("failed to marshal dynamodb record, %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("UserServiceTable"),
		Item:      av,
	}

	res, err := d.db.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to put item in dynamodb, %w", err)
	}
	d.logger.Info("Successfully created user", "response", res)

	return nil
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

func NewDynamoDbUserRepository(cfg aws.Config, optFns ...func(*dynamodb.Options)) repository.UserRepository {
	svc := dynamodb.NewFromConfig(cfg, optFns...)

	return &DynamoDbUserRepository{
		db: svc,
	}
}
