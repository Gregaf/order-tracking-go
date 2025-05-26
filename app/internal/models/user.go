package models

import "gregaf/order-tracking-go/internal/types"

type User struct {
	ID            string      `json:"id" dynamodbav:"ID"`
	DisplayID     string      `json:"displayID" dynamodbav:"UID"`
	FirstName     string      `json:"firstName" dynamodbav:"FirstName"`
	LastName      string      `json:"lastName" dynamodbav:"LastName"`
	Email         types.Email `json:"email" dynamodbav:"Email"`
	PrimaryPhone  string      `json:"primaryPhone" dynamodbav:"PrimaryPhone"`
	CreatedAtDate int64       `json:"createdAtDate" dynamodbav:"CreatedAtDate"`
	UpdatedAtDate int64       `json:"updatedAtDate" dynamodbav:"UpdatedAtDate"`
}
