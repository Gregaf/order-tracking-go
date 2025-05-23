package product

import (
	"context"

	"github.com/gregaf/order-tracking-go/internal/models"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product models.Product) error
	GetProductByID(ctx context.Context, ID string) (*models.Product, error)
	GetAllProducts(ctx context.Context, options models.GetResourceOptions) (*GetProductsData, error)
	DeleteProductByID(ctx context.Context, id string) error
}
