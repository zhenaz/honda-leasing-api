package repository

import (
	"context"
	"honda-leasing-api/internal/domain/model"
	"honda-leasing-api/internal/domain/query"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
}

type paymentRepository struct {
	q *query.Query
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{
		q: query.Use(db),
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	return r.q.Payment.WithContext(ctx).
		Create(payment)
}
