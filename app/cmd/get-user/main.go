package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

type MyResponse struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, event MyEvent) (MyResponse, error) {
	msg := fmt.Sprintf("Hello, %s!", event.Name)
	return MyResponse{Message: msg}, nil
}

func main() {
	lambda.Start(handler)
}
