package repository

import (
	"context"
	"fmt"
	"honda-leasing-api/internal/domain/model"
	"honda-leasing-api/internal/domain/query"

	"gorm.io/gorm"
)

type LeasingContractRepository interface {
	Create(ctx context.Context, contract *model.LeasingContract) error
	FindByID(ctx context.Context, id int64) (*model.LeasingContract, error)
	Update(ctx context.Context, contract *model.LeasingContract) error
}

type leasingContractRepository struct {
	q *query.Query
}

func NewLeasingContractRepository(db *gorm.DB) LeasingContractRepository {
	return &leasingContractRepository{
		q: query.Use(db),
	}
}

// --- CREATE ---
func (r *leasingContractRepository) Create(ctx context.Context, contract *model.LeasingContract) error {
	return r.q.LeasingContract.WithContext(ctx).
		Create(contract)
}

// --- FIND BY ID ---
func (r *leasingContractRepository) FindByID(ctx context.Context, id int64) (*model.LeasingContract, error) {
	contract, err := r.q.LeasingContract.WithContext(ctx).
		Where(r.q.LeasingContract.ContractID.Eq(id)).
		First()

	if err != nil {
		return nil, fmt.Errorf("leasing contract not found: %w", err)
	}

	return contract, nil
}

// --- UPDATE ---
func (r *leasingContractRepository) Update(ctx context.Context, contract *model.LeasingContract) error {
	_, err := r.q.LeasingContract.WithContext(ctx).
		Where(r.q.LeasingContract.ContractID.Eq(contract.ContractID)).
		Updates(contract)

	if err != nil {
		return fmt.Errorf("failed to update contract: %w", err)
	}

	return nil
}
