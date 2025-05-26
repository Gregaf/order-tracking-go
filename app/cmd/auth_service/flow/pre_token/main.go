package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request = events.CognitoEventUserPoolsPreTokenGenRequestV2_0
type Response = events.CognitoEventUserPoolsPreTokenGenResponseV2_0

type handler struct {
	logger *slog.Logger
}

func (h *handler) handleRequest(ctx context.Context, r Request) (Response, error) {
	h.logger.Info("Income Request Data", "event", r)

	customClaims := make(map[string]interface{})

	customClaims["role"] = "user"
	customClaims["permissions"] = []string{"own:user:read", "own:user:write", "own:user:delete", "own:product:read", "own:product:write"}

	overrideDetails := events.ClaimsAndScopeOverrideDetailsV2_0{
		AccessTokenGeneration: events.AccessTokenGenerationV2_0{
			ClaimsToAddOrOverride: customClaims,
		},
	}

	return Response{ClaimsAndScopeOverrideDetails: overrideDetails}, nil
}

func main() {
	h := handler{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
	lambda.Start(h.handleRequest)
}
