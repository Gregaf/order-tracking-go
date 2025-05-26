package dto

import (
	"gregaf/order-tracking-go/internal/models"
	"gregaf/order-tracking-go/internal/types"
)

type CreateProductDTO struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Category    string       `json:"category"`
	Price       types.USD    `json:"price"`
	Weight      types.Weight `json:"weight"`
}

type GetProductsResponse struct {
	Products      []models.Product `json:"products"`
	NextPageToken string           `json:"nextPageToken,omitempty"`
}
