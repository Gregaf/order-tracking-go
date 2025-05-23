package transport

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
}

func Success(data interface{}, statusCode int) (events.APIGatewayV2HTTPResponse, error) {
	// TODO: Validate that statusCode is not 4xx or 5xx

	body, _ := json.Marshal(data)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func Failure(errRes ErrorResponse, statusCode int) (events.APIGatewayV2HTTPResponse, error) {
	errorBody, _ := json.Marshal(errRes)
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(errorBody),
	}, nil
}
