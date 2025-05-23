package models

import "github.com/gregaf/order-tracking-go/internal/types"

type Product struct {
	ID            string       `json:"id" dynamodbav:"ID"`
	Name          string       `json:"name" dynamodbav:"Name"`
	Description   string       `json:"description" dynamodbav:"Description"`
	Category      string       `json:"category" dynamodbav:"Category"`
	Price         types.USD    `json:"price" dynamodbav:"Price"`
	Weight        types.Weight `json:"weight" dynamodbav:"Weight"`
	CreatedAtDate int64        `json:"createdAtDate" dynamodbav:"CreatedAtDate"`
	UpdatedAtDate int64        `json:"updatedAtDate" dynamodbav:"UpdatedAtDate"`
}
