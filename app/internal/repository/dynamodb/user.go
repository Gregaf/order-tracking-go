package dynamodbrepository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gregaf/order-tracking-go/internal/models"
)

type DynamoDbUserRepository struct {
	db     *dynamodb.Client
	logger *slog.Logger
}

func NewDynamoDbUserRepository(logger *slog.Logger, cfg aws.Config, optFns ...func(*dynamodb.Options)) *DynamoDbUserRepository {
	svc := dynamodb.NewFromConfig(cfg, optFns...)

	return &DynamoDbUserRepository{
		db:     svc,
		logger: logger,
	}
}

type DynamoDbUser struct {
	Pk string `dynamodbav:"Pk"`
	Sk string `dynamodbav:"Sk"`
	models.User
}

func (du *DynamoDbUser) ToUser() *models.User {
	return &du.User
}

func (d *DynamoDbUserRepository) CreateUser(ctx context.Context, user models.User) error {
	r := DynamoDbUser{
		Pk:   fmt.Sprintf("USER#%s", user.ID),
		Sk:   "METADATA",
		User: user,
	}
	av, err := attributevalue.MarshalMap(r)
	if err != nil {
		return fmt.Errorf("failed to marshal dynamodb record, %w", err)
	}
	d.logger.Info("Marshalled record", "record", av)
	input := &dynamodb.PutItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Item:      av,
	}

	d.logger.Info("PutItemInput", "input", input)

	res, err := d.db.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in dynamodb, %w", err)
	}
	d.logger.Info("Successfully created user", "response", res)

	return nil
}

// DeleteUser implements repository.UserRepository.
func (d *DynamoDbUserRepository) DeleteUser(ctx context.Context, ID string) error {
	av, err := attributevalue.MarshalMap(DynamoDbCompositeKey{Pk: fmt.Sprintf("USER#%s", ID), Sk: "METADATA"})
	if err != nil {
		return fmt.Errorf("failed to marshal dynamodb record, %w", err)
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Key:       av,
	}

	res, err := d.db.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete item in dynamodb, %w", err)
	}

	d.logger.Info("Successfully deleted user", "response", res)

	return nil
}

// GetUserByID implements repository.UserRepository.
func (d *DynamoDbUserRepository) GetUserByID(ctx context.Context, ID string) (*models.User, error) {
	av, err := attributevalue.MarshalMap(DynamoDbCompositeKey{Pk: fmt.Sprintf("USER#%s", ID), Sk: "METADATA"})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dynamodb record, %w", err)
	}

	d.logger.Info("Marshalled record", "record", av)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Key:       av,
	}

	res, err := d.db.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item in dynamodb, %w", err)
	}

	if res.Item == nil {
		return nil, nil
	}

	d.logger.Info("Successfully fetched user", "response", res)

	var fetchedUser DynamoDbUser
	err = attributevalue.UnmarshalMap(res.Item, &fetchedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dynamodb record, %w", err)
	}

	return fetchedUser.ToUser(), nil
}
