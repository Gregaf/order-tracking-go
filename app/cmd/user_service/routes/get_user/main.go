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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	repository "github.com/gregaf/order-tracking-go/internal/repository/dynamodb"
	"github.com/gregaf/order-tracking-go/internal/service/core"
	transport "github.com/gregaf/order-tracking-go/internal/transport/http"
	"github.com/gregaf/order-tracking-go/internal/transport/http/middleware"
	"github.com/gregaf/order-tracking-go/internal/user"
)

type Request = events.APIGatewayV2HTTPRequest
type Response = events.APIGatewayV2HTTPResponse

type handler struct {
	userSvc user.UserService
	logger  *slog.Logger
}

func (h *handler) handleRequest(ctx context.Context, r Request) (Response, error) {
	h.logger.Debug("Income Request Data", "event", r)

	userID := r.PathParameters["userID"]

	authCtx, err := middleware.GetAuthContext(r)
	if err != nil {
		return transport.Failure(transport.ErrorResponse{
			Message: "User not authorized",
			Code:    "UNAUTHORIZED",
			Details: map[string]string{
				"userID":      userID,
				"requestorID": authCtx.RequestorID,
			},
		}, http.StatusUnauthorized)
	}

	foundUser, err := h.userSvc.GetUserByID(ctx, *authCtx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return transport.Failure(transport.ErrorResponse{
				Message: "User not found",
				Code:    "NOT_FOUND",
				Details: map[string]string{
					"userID": userID,
				},
			}, http.StatusNotFound)
		} else if errors.Is(err, middleware.ErrNotAuthorized) {
			return transport.Failure(transport.ErrorResponse{
				Message: "User not authorized",
				Code:    "UNAUTHORIZED",
				Details: map[string]string{
					"userID":      userID,
					"requestorID": authCtx.RequestorID,
				},
			}, http.StatusUnauthorized)
		}

		return transport.Failure(transport.ErrorResponse{
			Message: "Internal Server Error",
			Code:    "UNKNOWN_ERROR",
			Details: map[string]string{
				"error": err.Error(),
			},
		}, http.StatusInternalServerError)
	}

	return transport.Success(foundUser, http.StatusOK)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbEndpoint := os.Getenv("DB_ENDPOINT")
	region := os.Getenv("AWS_REGION")

	logger.Debug("Loaded environment variables", "dbEndpoint", dbEndpoint, "region", region)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	repo := repository.NewDynamoDbUserRepository(logger, cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = &dbEndpoint
		o.Region = region
	})

	h := handler{
		userSvc: core.NewUserService(logger, repo),
		logger:  logger,
	}

	lambda.Start(h.handleRequest)
}
