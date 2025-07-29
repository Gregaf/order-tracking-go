package dynamodbrepository

import (
	"context"
	"fmt"
	"log/slog"

	"gregaf/order-tracking-go/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func (d *DynamoDbUserRepository) SyncUser(ctx context.Context, user models.User) error {
	update := expression.Set(
		expression.Name("FirstName"), expression.Value(user.FirstName),
	).Set(
		expression.Name("Email"), expression.Value(user.Email),
	).Set(
		expression.Name("UpdatedAtDate"), expression.Value(user.UpdatedAtDate),
	).Set(
		expression.Name("DisplayID"), expression.IfNotExists(
			expression.Name("DisplayID"),
			expression.Value(user.DisplayID),
		),
	).Set(
		expression.Name("CreatedAtDate"), expression.IfNotExists(
			expression.Name("CreatedAtDate"),
			expression.Value(user.CreatedAtDate),
		),
	)

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression: %w", err)
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"Pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", user.ID)},
			"Sk": &types.AttributeValueMemberS{Value: "METADATA"},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	res, err := d.db.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to update item in dynamodb: %w", err)
	}

	d.logger.Info("Successfully synced user", "response", res)
	return nil
}
