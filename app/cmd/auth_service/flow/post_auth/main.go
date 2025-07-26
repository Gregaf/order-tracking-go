package main

import (
	"context"
	"gregaf/order-tracking-go/internal/dto"
	dynamodbrepository "gregaf/order-tracking-go/internal/repository/dynamodb"
	"gregaf/order-tracking-go/internal/service/core"
	"gregaf/order-tracking-go/internal/transport/http/middleware"
	"gregaf/order-tracking-go/internal/types"
	"gregaf/order-tracking-go/internal/user"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Event = events.CognitoEventUserPoolsPostAuthentication

type handler struct {
	userSvc user.UserService
	logger  *slog.Logger
}

func (h *handler) handleRequest(ctx context.Context, event Event) (Event, error) {
	h.logger.Info("Income Request Data", "event", event)

	upsertDto := dto.UpsertUserDTO{
		ID:        event.UserName,
		FirstName: event.Request.UserAttributes["given_name"],
		LastName:  event.Request.UserAttributes["family_name"],
		Email:     types.Email(event.Request.UserAttributes["email"]),
	}

	res, err := h.userSvc.Upsert(ctx, middleware.AuthContext{Role: "Fake", RequestorID: "Fake", Permissions: []string{"Fake"}}, upsertDto)
	if err != nil {
		h.logger.Error("Failed to upsert user", "error", err)
		return event, err
	}
	h.logger.Info("User upserted successfully", "user", res)

	return event, nil
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

	repo := dynamodbrepository.NewDynamoDbUserRepository(logger, cfg, func(o *dynamodb.Options) {
		o.Region = region
	})

	// Initializing persistent connections, etc...
	h := handler{
		userSvc: core.NewUserService(logger, repo),
		logger:  logger,
	}

	lambda.Start(h.handleRequest)
}
