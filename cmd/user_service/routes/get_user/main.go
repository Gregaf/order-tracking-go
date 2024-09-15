package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
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

func (h *handler) handleRequest(ctx context.Context, r Request) (Response, error) {
	h.logger.Info("Income Request Data", "event", r)

	userId := r.PathParameters["user_id"]
	if userId == "" {
		return transport.Failure(http.StatusBadRequest, errors.New("Invalid path parameter, 'user_id' path paramter is required"))
	}

	user, err := h.userSvc.GetUserByID(userId)
	if err != nil {
		return transport.Failure(http.StatusInternalServerError, err)
	}

	if user == nil {
		return transport.Success(nil, http.StatusNotFound)
	}

	return transport.Success(user, http.StatusOK)
}

func main() {
	// Fetch ENV vars...
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	repo := dynamodbrepository.NewDynamoDbUserRepository(cfg)

	// Initializing persistent connections, etc...
	h := handler{
		userSvc: userservice.NewUserService(repo),
		logger:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	lambda.Start(h.handleRequest)
}
