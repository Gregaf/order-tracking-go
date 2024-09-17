package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gregaf/order-tracking/internal/models"
	dynamodbrepository "github.com/gregaf/order-tracking/internal/repository/dynamodb_repository"
	"github.com/gregaf/order-tracking/internal/transport"
	userservice "github.com/gregaf/order-tracking/internal/user_service"
)

type Request = events.APIGatewayV2HTTPRequest
type Response = events.APIGatewayV2HTTPResponse

type handler struct {
	userSvc *userservice.UserService
	logger  *slog.Logger
}

type CreateUserRequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (h *handler) handleRequest(ctx context.Context, r Request) (Response, error) {
	h.logger.Info("Income Request Data", "event", r)

	data := []byte(r.Body)
	ok := json.Valid(data)
	if !ok {
		return transport.Failure(400, fmt.Errorf("Invalid JSON Body, valid JSON is required in the Body"))
	}

	userDto := &CreateUserRequestBody{}
	err := json.Unmarshal(data, userDto)
	if err != nil {
		return transport.Failure(400, err)
	}

	err = h.userSvc.CreateUser(&models.User{
		FirstName: userDto.FirstName,
		LastName:  userDto.LastName,
		Email:     userDto.Email,
	})
	if err != nil {
		return transport.Failure(400, err)
	}

	return transport.Success("User successfully created", http.StatusCreated)
}

func main() {
	dbEndpoint := os.Getenv("DB_ENDPOINT")
	region := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	repo := dynamodbrepository.NewDynamoDbUserRepository(cfg, func(o *dynamodb.Options) {
		o.EndpointResolverV2.ResolveEndpoint(context.TODO(), dynamodb.EndpointParameters{
			Endpoint: &dbEndpoint,
			Region:   &region,
		})
	})

	// Initializing persistent connections, etc...
	h := handler{
		userSvc: userservice.NewUserService(repo),
		logger:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	lambda.Start(h.handleRequest)
}
