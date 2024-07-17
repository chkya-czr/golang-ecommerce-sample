package service

import (
	"context"
	"product_service/internal/domain/product/model"
)

type IProductService interface {
	Get(ctx context.Context, productId int) (*model.ChildProductModel, error)
}

type ProductService struct {
}

func New() *ProductService {
	return &ProductService{}
}

func (svc *ProductService) Get(ctx context.Context, productId int) (*model.ChildProductModel, error) {
	return nil, nil
}
