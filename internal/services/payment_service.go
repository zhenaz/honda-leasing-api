package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"honda-leasing-api/internal/domain/model"
	"honda-leasing-api/internal/domain/query"

	"gorm.io/gorm"
)

type PaymentService interface {
	PayInstallment(ctx context.Context, req *CreatePaymentRequest) error
}

type CreatePaymentRequest struct {
	ScheduleID       int64   `json:"schedule_id"`
	JumlahBayar      float64 `json:"jumlah_bayar"`
	MetodePembayaran string  `json:"metode_pembayaran"`
	Provider         string  `json:"provider"`
}

type paymentService struct {
	db *gorm.DB
}

func NewPaymentService(db *gorm.DB) PaymentService {
	return &paymentService{
		db: db,
	}
}

func (s *paymentService) PayInstallment(
	ctx context.Context,
	req *CreatePaymentRequest,
) error {

	return s.db.Transaction(func(tx *gorm.DB) error {

		q := query.Use(tx)

		// =========================
		// 1️⃣ Ambil schedule
		// =========================
		schedule, err := q.PaymentSchedule.WithContext(ctx).
			Where(q.PaymentSchedule.ScheduleID.Eq(req.ScheduleID)).
			First()

		if err != nil {
			return errors.New("schedule not found")
		}

		// =========================
		// 2️⃣ Cek sudah dibayar?
		// =========================
		if schedule.StatusPembayaran != nil &&
			*schedule.StatusPembayaran == "paid" {
			return errors.New("installment already paid")
		}

		// =========================
		// 3️⃣ Insert payment
		// =========================
		now := time.Now()
		statusPaid := "paid"

		nomorBukti := fmt.Sprintf(
			"PAY-%d-%d",
			schedule.ScheduleID,
			now.UnixNano(),
		)

		payment := &model.Payment{
			NomorBukti:       nomorBukti,
			JumlahBayar:      req.JumlahBayar,
			TanggalBayar:     now,
			MetodePembayaran: req.MetodePembayaran,
			Provider:         &req.Provider,
			ContractID:       schedule.ContractID,
			ScheduleID:       &schedule.ScheduleID,
		}

		if err := q.Payment.WithContext(ctx).Create(payment); err != nil {
			return err
		}

		// =========================
		// 4️⃣ Update schedule
		// =========================
		schedule.StatusPembayaran = &statusPaid
		schedule.TanggalBayar = &now

		if _, err := q.PaymentSchedule.WithContext(ctx).
			Where(q.PaymentSchedule.ScheduleID.Eq(schedule.ScheduleID)).
			Updates(schedule); err != nil {
			return err
		}

		// =========================
		// 5️⃣ Cek apakah semua paid
		// =========================
		allSchedules, err := q.PaymentSchedule.WithContext(ctx).
			Where(q.PaymentSchedule.ContractID.Eq(schedule.ContractID)).
			Find()
		if err != nil {
			return err
		}

		allPaid := true
		for _, sc := range allSchedules {
			if sc.StatusPembayaran == nil ||
				*sc.StatusPembayaran != "paid" {
				allPaid = false
				break
			}
		}

		if allPaid {
			if _, err := q.LeasingContract.WithContext(ctx).
				Where(q.LeasingContract.ContractID.Eq(schedule.ContractID)).
				Update(q.LeasingContract.Status, "paid_off"); err != nil {
				return err
			}
		}

		return nil
	})
}
