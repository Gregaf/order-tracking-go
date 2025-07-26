package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Event = events.CognitoEventUserPoolsPreTokenGenV2_0
type Response = events.CognitoEventUserPoolsPreTokenGenResponseV2_0

type handler struct {
	logger *slog.Logger
}

func (h *handler) handleRequest(ctx context.Context, event Event) (Event, error) {
	h.logger.Info("Income Request Data", "event", event)

	customClaims := make(map[string]interface{})

	customClaims["role"] = "user"
	customClaims["permissions"] = []string{"own:user:read", "own:user:write", "own:user:delete", "own:product:read", "own:product:write"}

	event.Response.ClaimsAndScopeOverrideDetails.AccessTokenGeneration.ClaimsToAddOrOverride = customClaims

	return event, nil
}

func main() {
	h := handler{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
	lambda.Start(h.handleRequest)
}
