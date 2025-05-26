package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
)

type Request = events.APIGatewayV2CustomAuthorizerV2Request
type Response = events.APIGatewayV2CustomAuthorizerSimpleResponse

type handler struct {
	logger *slog.Logger
}

type CustomClaims struct {
	Permissions []string `json:"permissions"`
	Role        string   `json:"role"`
	jwt.RegisteredClaims
}

func (c *CustomClaims) GetPermissions() []string {
	if c.Permissions == nil {
		return []string{}
	}

	return c.Permissions
}

func (c *CustomClaims) GetRole() string {
	if len(c.Role) <= 0 {
		return "Anonymous"
	}

	return c.Role
}

func (h *handler) handleRequest(ctx context.Context, r Request) (Response, error) {
	h.logger.Info("Income Request Data", "event", r)

	tokenString := strings.Split(r.Headers["authorization"], " ")[1]
	if tokenString == "" {
		return Response{IsAuthorized: false}, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		h.logger.Error("Failed to parse JWT", "error", err)
		return Response{IsAuthorized: false}, nil
	}

	if claims, ok := token.Claims.(*CustomClaims); !ok {
		h.logger.Error("Failed to parse JWT claims", "error", err)
		return Response{IsAuthorized: false}, nil
	} else {

		return Response{IsAuthorized: true, Context: map[string]interface{}{
			"permissions": claims.GetPermissions(),
			"role":        claims.GetRole(),
			"requestorID": claims.Subject,
		}}, nil
	}
}

func main() {
	h := handler{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
	lambda.Start(h.handleRequest)
}
