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

/*
====================================
STATUS CONSTANTS
====================================
*/

// Schedule status
const (
	ScheduleUnpaid  = "unpaid"
	SchedulePaid    = "paid"
	ScheduleOverdue = "overdue"
)

// Contract status (sesuai CHECK constraint DB)
const (
	ContractActive  = "active"
	ContractLate    = "late"
	ContractPaidOff = "paid_off"
)

/*
====================================
INTERFACE
====================================
*/

type PaymentService interface {
	PayInstallment(ctx context.Context, req *CreatePaymentRequest) error
	CheckOverdue(ctx context.Context) error
}

/*
====================================
REQUEST STRUCT
====================================
*/

type CreatePaymentRequest struct {
	ScheduleID       int64   `json:"schedule_id"`
	JumlahBayar      float64 `json:"jumlah_bayar"`
	MetodePembayaran string  `json:"metode_pembayaran"`
	Provider         string  `json:"provider"`
}

/*
====================================
SERVICE IMPLEMENTATION
====================================
*/

type paymentService struct {
	db *gorm.DB
}

func NewPaymentService(db *gorm.DB) PaymentService {
	return &paymentService{
		db: db,
	}
}

/*
====================================
PAY INSTALLMENT
====================================
*/

func (s *paymentService) PayInstallment(
	ctx context.Context,
	req *CreatePaymentRequest,
) error {

	return s.db.Transaction(func(tx *gorm.DB) error {

		q := query.Use(tx)

		// 1️⃣ Ambil schedule
		schedule, err := q.PaymentSchedule.WithContext(ctx).
			Where(q.PaymentSchedule.ScheduleID.Eq(req.ScheduleID)).
			First()

		if err != nil {
			return errors.New("schedule not found")
		}

		// 2️⃣ Cek sudah dibayar?
		if schedule.StatusPembayaran != nil &&
			*schedule.StatusPembayaran == SchedulePaid {
			return errors.New("installment already paid")
		}

		now := time.Now()

		// 3️⃣ Insert payment
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

		// 4️⃣ Update schedule → paid
		statusPaid := SchedulePaid
		schedule.StatusPembayaran = &statusPaid
		schedule.TanggalBayar = &now

		if _, err := q.PaymentSchedule.WithContext(ctx).
			Where(q.PaymentSchedule.ScheduleID.Eq(schedule.ScheduleID)).
			Updates(schedule); err != nil {
			return err
		}

		// 5️⃣ Cek apakah semua schedule sudah paid
		allSchedules, err := q.PaymentSchedule.WithContext(ctx).
			Where(q.PaymentSchedule.ContractID.Eq(schedule.ContractID)).
			Find()

		if err != nil {
			return err
		}

		allPaid := true
		for _, sc := range allSchedules {
			if sc.StatusPembayaran == nil ||
				*sc.StatusPembayaran != SchedulePaid {
				allPaid = false
				break
			}
		}

		// 6️⃣ Jika semua paid → contract = paid_off
		if allPaid {
			if _, err := q.LeasingContract.WithContext(ctx).
				Where(q.LeasingContract.ContractID.Eq(schedule.ContractID)).
				Update(q.LeasingContract.Status, ContractPaidOff); err != nil {
				return err
			}
		}

		return nil
	})
}

/*
====================================
CHECK OVERDUE
====================================
*/

func (s *paymentService) CheckOverdue(ctx context.Context) error {

	return s.db.Transaction(func(tx *gorm.DB) error {

		q := query.Use(tx)
		now := time.Now()

		// Ambil semua schedule unpaid
		schedules, err := q.PaymentSchedule.WithContext(ctx).
			Where(q.PaymentSchedule.StatusPembayaran.Eq(ScheduleUnpaid)).
			Find()

		if err != nil {
			return err
		}

		for _, sc := range schedules {

			// Jika sudah lewat jatuh tempo
			if now.After(sc.JatuhTempo) {

				statusOverdue := ScheduleOverdue

				// Update schedule → overdue
				if _, err := q.PaymentSchedule.WithContext(ctx).
					Where(q.PaymentSchedule.ScheduleID.Eq(sc.ScheduleID)).
					Update(q.PaymentSchedule.StatusPembayaran, statusOverdue); err != nil {
					return err
				}

				// Update contract → late
				if _, err := q.LeasingContract.WithContext(ctx).
					Where(q.LeasingContract.ContractID.Eq(sc.ContractID)).
					Update(q.LeasingContract.Status, ContractLate); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
