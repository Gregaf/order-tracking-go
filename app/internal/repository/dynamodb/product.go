package dynamodbrepository

import (
	"context"
	"fmt"
	"log/slog"

	"gregaf/order-tracking-go/internal/models"
	"gregaf/order-tracking-go/internal/product"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDbProductRepository struct {
	db     *dynamodb.Client
	logger *slog.Logger
}

func NewDynamoDbProductRepository(logger *slog.Logger, cfg aws.Config, optFns ...func(*dynamodb.Options)) *DynamoDbProductRepository {
	svc := dynamodb.NewFromConfig(cfg, optFns...)

	return &DynamoDbProductRepository{
		db:     svc,
		logger: logger,
	}
}

type DynamoDbProduct struct {
	Pk string `dynamodbav:"Pk"`
	Sk string `dynamodbav:"Sk"`
	models.Product
}

func (d *DynamoDbProduct) ToProduct() *models.Product {
	return &d.Product
}

func (d *DynamoDbProductRepository) CreateProduct(ctx context.Context, product models.Product) error {
	r := DynamoDbProduct{
		Pk:      "METADATA",
		Sk:      fmt.Sprintf("PRODUCT#%s", product.ID),
		Product: product,
	}
	av, err := attributevalue.MarshalMap(r)
	if err != nil {
		return fmt.Errorf("failed to marshal dynamodb record, %w", err)
	}
	d.logger.Info("Marshalled record", "record", av)
	input := &dynamodb.PutItemInput{
		TableName: aws.String(PRODUCT_TABLE_NAME),
		Item:      av,
	}

	d.logger.Info("PutItemInput", "input", input)

	res, err := d.db.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item in dynamodb, %w", err)
	}
	d.logger.Info("Successfully created product", "response", res)

	return nil
}

func (d *DynamoDbProductRepository) GetProductByID(ctx context.Context, ID string) (*models.Product, error) {
	av, err := attributevalue.MarshalMap(DynamoDbCompositeKey{Pk: fmt.Sprintf("PRODUCT#%s", ID), Sk: "METADATA"})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dynamodb record, %w", err)
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(PRODUCT_TABLE_NAME),
		Key:       av,
	}

	res, err := d.db.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item in dynamodb, %w", err)
	}

	if res.Item == nil {
		return nil, nil
	}

	var fetchedProduct DynamoDbProduct
	err = attributevalue.UnmarshalMap(res.Item, &fetchedProduct)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dynamodb record, %w", err)
	}

	return fetchedProduct.ToProduct(), nil
}

func (d *DynamoDbProductRepository) GetAllProducts(ctx context.Context, options models.GetResourceOptions) (*product.GetProductsData, error) {

	d.logger.Info("filter criterias", "options", options)

	builder := expression.NewBuilder()
	keyCondition := expression.Key("Pk").Equal(expression.Value("METADATA")).And(expression.Key("Sk").BeginsWith("PRODUCT#"))

	filter, err := buildFilterExpression(options.FilterCriteria)
	if err != nil {
		return nil, err
	}

	if filter != nil {
		builder = builder.WithFilter(*filter)
	}

	expr, err := builder.WithKeyCondition(keyCondition).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build expression, %w", err)
	}
	d.logger.Info("out", "names", expr.Names(), "values", expr.Values())

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(PRODUCT_TABLE_NAME),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
	}

	res, err := d.db.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item in dynamodb, %w", err)
	}

	var fetchedProducts []DynamoDbProduct
	err = attributevalue.UnmarshalListOfMaps(res.Items, &fetchedProducts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dynamodb record, %w", err)
	}

	d.logger.Info("Successfully fetched products", "response", res)

	products := make([]models.Product, len(fetchedProducts))
	for i, fetchedProduct := range fetchedProducts {
		products[i] = *fetchedProduct.ToProduct()
	}

	return &product.GetProductsData{Products: products, Count: len(products)}, nil
}

// amazonq-ignore-next-line
func (d *DynamoDbProductRepository) DeleteProductByID(ctx context.Context, id string) error {
	panic("not implemented") // TODO: Implement
}
