package core

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gregaf/order-tracking-go/internal/dto"
	"github.com/gregaf/order-tracking-go/internal/models"
	"github.com/gregaf/order-tracking-go/internal/product"
	"github.com/gregaf/order-tracking-go/internal/transport/http/middleware"
	"github.com/gregaf/order-tracking-go/internal/util"
)

type ProductServiceCore struct {
	productRepo product.ProductRepository
	logger      *slog.Logger
}

func NewProductServiceCore(logger *slog.Logger, productRepo product.ProductRepository) *ProductServiceCore {
	return &ProductServiceCore{productRepo: productRepo, logger: logger}
}

func (p *ProductServiceCore) GetProductByID(ctx context.Context, authCtx middleware.AuthContext, ID string) (*models.Product, error) {
	// if middleware.HasPermission(authCtx.Permissions, authCtx.RequestorID, ID, "product", "read") {
	// 	return nil, middleware.ErrNotAuthorized
	// }

	foundProduct, err := p.productRepo.GetProductByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if foundProduct == nil {
		return nil, product.ErrProductNotFound
	}

	return foundProduct, nil
}

var productFilterFields = []string{"name", "category", "description", "price"}

func (p *ProductServiceCore) GetProducts(ctx context.Context, authCtx middleware.AuthContext, options models.GetResourceOptions) (*product.GetProductsData, error) {

	for _, filterCriteria := range options.FilterCriteria {
		err := filterCriteria.Validate()
		if err != nil {
			return nil, err
		}

		isValidField := util.IsOneOf(filterCriteria.Field, productFilterFields...)

		if !isValidField {
			return nil, errors.New("invalid filter field")
		}
	}

	return p.productRepo.GetAllProducts(ctx, options)
}

func (p *ProductServiceCore) CreateProduct(ctx context.Context, authCtx middleware.AuthContext, product dto.CreateProductDTO) (*models.Product, error) {
	// Check auth

	// Validate product

	productID, err := util.GenerateRandomSegment(10)
	if err != nil {
		return nil, err
	}

	createdAtDate := time.Now().UnixMilli()

	newProduct := &models.Product{
		Name:          product.Name,
		Category:      product.Category,
		Description:   product.Description,
		Price:         product.Price,
		Weight:        product.Weight,
		ID:            productID,
		CreatedAtDate: createdAtDate,
		UpdatedAtDate: createdAtDate,
	}

	err = p.productRepo.CreateProduct(ctx, *newProduct)
	if err != nil {
		return nil, err
	}

	return newProduct, nil
}

func (p *ProductServiceCore) DeleteProductByID(ctx context.Context, authCtx middleware.AuthContext, id string) error {
	panic("not implemented") // TODO: Implement
}
