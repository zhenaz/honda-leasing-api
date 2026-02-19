package repository

import (
	"context"
	"fmt"
	"honda-leasing-api/internal/domain/model"
	"honda-leasing-api/internal/domain/query"

	"gorm.io/gorm"
)

type LeasingProductRepository interface {
	FindByID(ctx context.Context, id int64) (*model.LeasingProduct, error)
}

type leasingProductRepository struct {
	q *query.Query
}

func NewLeasingProductRepository(db *gorm.DB) LeasingProductRepository {
	return &leasingProductRepository{
		q: query.Use(db),
	}
}

func (r *leasingProductRepository) FindByID(ctx context.Context, id int64) (*model.LeasingProduct, error) {
	product, err := r.q.LeasingProduct.WithContext(ctx).
		Where(r.q.LeasingProduct.ProductID.Eq(id)).
		First()

	if err != nil {
		return nil, fmt.Errorf("leasing product not found: %w", err)
	}

	return product, nil
}
