package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"gregaf/order-tracking-go/internal/dto"
	repository "gregaf/order-tracking-go/internal/repository/dynamodb"
	"gregaf/order-tracking-go/internal/service/core"
	transport "gregaf/order-tracking-go/internal/transport/http"
	"gregaf/order-tracking-go/internal/transport/http/middleware"
	"gregaf/order-tracking-go/internal/user"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Request = events.APIGatewayV2HTTPRequest
type Response = events.APIGatewayV2HTTPResponse

type handler struct {
	userSvc user.UserService
	logger  *slog.Logger
}

func (h *handler) handleRequest(ctx context.Context, r Request) (Response, error) {
	h.logger.Info("Income Request Data", "event", r)

	data := []byte(r.Body)

	userDto := &dto.CreateUserDTO{}
	err := json.Unmarshal(data, userDto)
	if err != nil {
		return transport.Failure(transport.ErrorResponse{
			Message: "Invalid JSON submitted",
			Code:    "INVALID_JSON",
			Details: map[string]string{
				"error": err.Error(),
			},
		}, http.StatusBadRequest)
	}

	authCtx, err := middleware.GetAuthContext(r)
	if err != nil {
		return transport.Failure(transport.ErrorResponse{
			Message: "User not authorized",
			Code:    "UNAUTHORIZED",
			Details: map[string]string{
				"requestorID": authCtx.RequestorID,
			},
		}, http.StatusUnauthorized)
	}

	user, err := h.userSvc.CreateUser(ctx, *authCtx, *userDto)
	if err != nil {
		return transport.Failure(transport.ErrorResponse{
			Message: "Internal Server Error",
			Code:    "UNKNOWN_ERROR",
			Details: map[string]string{
				"error": err.Error(),
			},
		}, http.StatusInternalServerError)
	}

	h.logger.Debug("Created User", "user", user)

	return transport.Success(map[string]string{
		"Message": "User successfully created",
	}, http.StatusCreated)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbEndpoint := os.Getenv("DB_ENDPOINT")
	// TODO: Pass table name
	region := os.Getenv("AWS_REGION")

	logger.Info("Loaded environment variables", "dbEndpoint", dbEndpoint, "region", region)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// logger.Info("Loaded configuration", "config", cfg)
	repo := repository.NewDynamoDbUserRepository(logger, cfg, func(o *dynamodb.Options) {
		// o.BaseEndpoint = &dbEndpoint
		o.Region = region
	})

	// Initializing persistent connections, etc...
	h := handler{
		userSvc: core.NewUserService(logger, repo),
		logger:  logger,
	}

	lambda.Start(h.handleRequest)
}
