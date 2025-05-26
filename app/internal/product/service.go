package product

import (
	"context"

	"gregaf/order-tracking-go/internal/dto"
	"gregaf/order-tracking-go/internal/models"

	"gregaf/order-tracking-go/internal/transport/http/middleware"
)

type ProductService interface {
	GetProductByID(ctx context.Context, authCtx middleware.AuthContext, ID string) (*models.Product, error)
	CreateProduct(ctx context.Context, authCtx middleware.AuthContext, product dto.CreateProductDTO) (*models.Product, error)
	GetProducts(ctx context.Context, authCtx middleware.AuthContext, options models.GetResourceOptions) (*GetProductsData, error)
	DeleteProductByID(ctx context.Context, authCtx middleware.AuthContext, id string) error
}

type GetProductsData struct {
	Products      []models.Product `json:"products"`
	Count         int              `json:"count"`
	NextPageToken string           `json:"nextPageToken,omitempty"`
}
