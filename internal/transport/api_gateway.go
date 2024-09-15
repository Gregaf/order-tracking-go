package transport

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponse struct {
	Error string `json:"error"`
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

func Failure(statusCode int, err error) (events.APIGatewayV2HTTPResponse, error) {
	errorBody, _ := json.Marshal(ErrorResponse{Error: err.Error()})
	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(errorBody),
	}, nil
}
