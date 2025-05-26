package models

import "gregaf/order-tracking-go/internal/types"

type Order struct {
	ID         string      `json:"id" dynamodbav:"ID"`
	Items      []OrderItem `json:"items" dynamodbav:"Items"`
	TotalItems int         `json:"totalItems" dynamodbav:"TotalItems"`
	Status     OrderStatus `json:"status" dynamodbav:"Status"`
}

type OrderItem struct {
	ID            string          `json:"id" dynamodbav:"ID"`
	Name          string          `json:"name" dynamodbav:"Name"`
	Description   string          `json:"description" dynamodbav:"Description"`
	Category      string          `json:"category" dynamodbav:"Category"`
	Price         types.USD       `json:"price" dynamodbav:"Price"`
	Weight        types.Weight    `json:"weight" dynamodbav:"Weight"`
	Quantity      int             `json:"quantity" dynamodbav:"Quantity"`
	Status        OrderItemStatus `json:"status" dynamodbav:"Status"`
	CreatedAtDate int64           `json:"createdAtDate" dynamodbav:"CreatedAtDate"`
	UpdatedAtDate int64           `json:"updatedAtDate" dynamodbav:"UpdatedAtDate"`
}

type OrderStatus int

const (
	PROCESSING OrderStatus = iota
	SHIPPED
	DELIVERED
	CANCELLED
)

var orderStatusNames = map[OrderStatus]string{
	PROCESSING: "PROCESSING",
	SHIPPED:    "SHIPPED",
	DELIVERED:  "DELIVERED",
	CANCELLED:  "CANCELLED",
}

func (o OrderStatus) String() string {
	return orderStatusNames[o]
}

type OrderItemStatus int

const (
	AWAITING_ARRIVAL OrderItemStatus = iota
	RECEVIED
)

var orderItemStatusNames = map[OrderItemStatus]string{
	AWAITING_ARRIVAL: "AWAITING_ARRIVAL",
	RECEVIED:         "RECEVIED",
}

func (o OrderItemStatus) String() string {
	return orderItemStatusNames[o]
}
