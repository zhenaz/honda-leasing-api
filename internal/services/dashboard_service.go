package service

import (
	"context"

	"honda-leasing-api/internal/domain/query"

	"gorm.io/gorm"
)

type DashboardService interface {
	GetSummary(ctx context.Context) (*DashboardSummaryResponse, error)
}

type DashboardSummaryResponse struct {
	TotalContract    int64   `json:"total_contract"`
	ActiveContract   int64   `json:"active_contract"`
	ClosedContract   int64   `json:"closed_contract"`
	TotalOutstanding float64 `json:"total_outstanding"`
	TotalPaid        float64 `json:"total_paid"`
}

type dashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) DashboardService {
	return &dashboardService{db: db}
}

func (s *dashboardService) GetSummary(ctx context.Context) (*DashboardSummaryResponse, error) {

	q := query.Use(s.db)

	// Count contract
	totalContract, _ := q.LeasingContract.WithContext(ctx).Count()
	activeContract, _ := q.LeasingContract.WithContext(ctx).
		Where(q.LeasingContract.Status.Eq("active")).
		Count()
	closedContract, _ := q.LeasingContract.WithContext(ctx).
		Where(q.LeasingContract.Status.Eq("closed")).
		Count()

	// Total Outstanding (unpaid)
	unpaidSchedules, _ := q.PaymentSchedule.WithContext(ctx).
		Where(q.PaymentSchedule.StatusPembayaran.Eq("unpaid")).
		Find()

	totalOutstanding := 0.0
	for _, sc := range unpaidSchedules {
		totalOutstanding += sc.TotalTagihan
	}

	// Total Paid
	payments, _ := q.Payment.WithContext(ctx).Find()

	totalPaid := 0.0
	for _, p := range payments {
		totalPaid += p.JumlahBayar
	}

	return &DashboardSummaryResponse{
		TotalContract:    totalContract,
		ActiveContract:   activeContract,
		ClosedContract:   closedContract,
		TotalOutstanding: totalOutstanding,
		TotalPaid:        totalPaid,
	}, nil
}
