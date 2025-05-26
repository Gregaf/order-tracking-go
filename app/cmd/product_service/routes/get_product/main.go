package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"gregaf/order-tracking-go/internal/product"
	repository "gregaf/order-tracking-go/internal/repository/dynamodb"
	"gregaf/order-tracking-go/internal/service/core"
	transport "gregaf/order-tracking-go/internal/transport/http"
	"gregaf/order-tracking-go/internal/transport/http/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Request = events.APIGatewayV2HTTPRequest
type Response = events.APIGatewayV2HTTPResponse

type handler struct {
	productSvc product.ProductService
	logger     *slog.Logger
}

func (h *handler) handleRequest(ctx context.Context, r Request) (Response, error) {
	h.logger.Info("Income Request Data", "event", r)

	productID := r.PathParameters["productID"]

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

	foundProduct, err := h.productSvc.GetProductByID(ctx, *authCtx, productID)
	if err != nil {
		return transport.Failure(transport.ErrorResponse{
			Message: "Internal Server Error",
			Code:    "UNKNOWN_ERROR",
			Details: map[string]string{
				"error": err.Error(),
			},
		}, http.StatusInternalServerError)
	}

	return transport.Success(foundProduct, http.StatusOK)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbEndpoint := os.Getenv("DB_ENDPOINT")
	region := os.Getenv("AWS_REGION")

	logger.Info("Loaded environment variables", "dbEndpoint", dbEndpoint, "region", region)

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	repo := repository.NewDynamoDbProductRepository(logger, cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = &dbEndpoint
		o.Region = region
	})

	// Initializing persistent connections, etc...
	h := handler{
		productSvc: core.NewProductServiceCore(logger, repo),
		logger:     logger,
	}

	lambda.Start(h.handleRequest)
}
