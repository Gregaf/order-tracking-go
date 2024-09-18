package models

type User struct {
	Uid           string `json:"uid" dynamodbav:"Uid"`
	Sub           string `json:"sub" dynamodbav:"Sub"`
	FirstName     string `json:"firstName" dynamodbav:"FirstName"`
	LastName      string `json:"lastName" dynamodbav:"LastName"`
	Email         Email  `json:"email" dynamodbav:"Email"`
	PrimaryPhone  string `json:"primaryPhone" dynamodbav:"PrimaryPhone"`
	CreatedAtDate int64  `json:"createdAtDate" dynamodbav:"CreatedAtDate"`
	UpdatedAtDate int64  `json:"updatedAtDate" dynamodbav:"UpdatedAtDate"`
}
