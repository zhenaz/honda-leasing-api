package repository

import (
	"context"
	"fmt"
	"honda-leasing-api/internal/domain/model"
	"honda-leasing-api/internal/domain/query"

	"gorm.io/gorm"
)

type PaymentScheduleRepository interface {
	Create(ctx context.Context, schedule *model.PaymentSchedule) error
	FindByContractID(ctx context.Context, contractID int64) ([]*model.PaymentSchedule, error)
	Update(ctx context.Context, schedule *model.PaymentSchedule) error
	FindByID(ctx context.Context, id int64) (*model.PaymentSchedule, error)
}

type paymentScheduleRepository struct {
	q *query.Query
}

func NewPaymentScheduleRepository(db *gorm.DB) PaymentScheduleRepository {
	return &paymentScheduleRepository{
		q: query.Use(db),
	}
}

func (r *paymentScheduleRepository) FindByID(ctx context.Context, id int64) (*model.PaymentSchedule, error) {
	return r.q.PaymentSchedule.WithContext(ctx).
		Where(r.q.PaymentSchedule.ScheduleID.Eq(id)).
		First()
}

// --- CREATE ---
func (r *paymentScheduleRepository) Create(ctx context.Context, schedule *model.PaymentSchedule) error {
	return r.q.PaymentSchedule.WithContext(ctx).
		Create(schedule)
}

// --- FIND BY CONTRACT ID ---
func (r *paymentScheduleRepository) FindByContractID(ctx context.Context, contractID int64) ([]*model.PaymentSchedule, error) {
	schedules, err := r.q.PaymentSchedule.WithContext(ctx).
		Where(r.q.PaymentSchedule.ContractID.Eq(contractID)).
		Order(r.q.PaymentSchedule.AngsuranKe.Asc()).
		Find()

	if err != nil {
		return nil, fmt.Errorf("failed to find schedules: %w", err)
	}

	return schedules, nil
}

// --- UPDATE ---
func (r *paymentScheduleRepository) Update(ctx context.Context, schedule *model.PaymentSchedule) error {
	_, err := r.q.PaymentSchedule.WithContext(ctx).
		Where(r.q.PaymentSchedule.ScheduleID.Eq(schedule.ScheduleID)).
		Updates(schedule)

	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	return nil
}
