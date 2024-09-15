package models

type User struct {
	Id            string `json:"id"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	PrimaryPhone  string `json:"primaryPhone"`
	CreatedAtDate int    `json:"createdAtDate"`
	UpdatedAtDate int    `json:"updatedAtDate"`
}
