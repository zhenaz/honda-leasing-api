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

type LeasingService interface {
	CreateContract(ctx context.Context, req *CreateContractRequest) (*model.LeasingContract, error)
	GetContractDetail(ctx context.Context, id int64) (*ContractDetailResponse, error)
	ApproveContract(ctx context.Context, id int64) error
	// ActivateContract(ctx context.Context, id int64) error
	ListContracts(ctx context.Context, page int, limit int, status string) ([]*model.LeasingContract, int64, error)
	CompleteTask(ctx context.Context, taskID int64) error
}

type CreateContractRequest struct {
	CustomerID     int64   `json:"customer_id"`
	MotorID        int64   `json:"motor_id"`
	ProductID      int64   `json:"product_id"`
	NilaiKendaraan float64 `json:"nilai_kendaraan"`
	DpDibayar      float64 `json:"dp_dibayar"`
}

type ContractDetailResponse struct {
	Contract  *model.LeasingContract   `json:"contract"`
	Schedules []*model.PaymentSchedule `json:"schedules"`
	Summary   ContractSummary          `json:"summary"`
}

type ContractSummary struct {
	TotalCicilan    int     `json:"total_cicilan"`
	SudahDibayar    int     `json:"sudah_dibayar"`
	SisaCicilan     int     `json:"sisa_cicilan"`
	ProgressPercent float64 `json:"progress_percent"`
}

type leasingService struct {
	db *gorm.DB
}

func NewLeasingService(db *gorm.DB) LeasingService {
	return &leasingService{db: db}
}

func (s *leasingService) CreateContract(
	ctx context.Context,
	req *CreateContractRequest,
) (*model.LeasingContract, error) {

	var createdContract *model.LeasingContract

	err := s.db.Transaction(func(tx *gorm.DB) error {

		q := query.Use(tx)

		// =========================
		// 1️⃣ Ambil Product
		// =========================
		product, err := q.LeasingProduct.WithContext(ctx).
			Where(q.LeasingProduct.ProductID.Eq(req.ProductID)).
			First()
		if err != nil {
			return errors.New("product not found")
		}

		// =========================
		// 2️⃣ Validasi DP
		// =========================
		dpPercent := (req.DpDibayar / req.NilaiKendaraan) * 100
		if dpPercent < product.DpPersenMin || dpPercent > product.DpPersenMax {
			return errors.New("DP di luar batas product")
		}

		// =========================
		// 3️⃣ Hitung Leasing (Flat Rate / bulan)
		// =========================
		pokok := req.NilaiKendaraan - req.DpDibayar
		totalBunga := pokok * (product.BungaFlat / 100) * float64(product.TenorBulan)
		totalPinjaman := pokok + totalBunga
		cicilan := totalPinjaman / float64(product.TenorBulan)

		now := time.Now()
		contractNumber := fmt.Sprintf(
			"KTR-%d-%d",
			now.Year(),
			now.UnixNano(),
		)

		// =========================
		// 4️⃣ Insert Contract
		// =========================
		contract := &model.LeasingContract{
			ContractNumber:  &contractNumber,
			RequestDate:     now,
			TenorBulan:      product.TenorBulan,
			NilaiKendaraan:  req.NilaiKendaraan,
			DpDibayar:       req.DpDibayar,
			PokokPinjaman:   pokok,
			TotalPinjaman:   totalPinjaman,
			CicilanPerBulan: cicilan,
			Status:          "draft",
			CustomerID:      req.CustomerID,
			MotorID:         req.MotorID,
			ProductID:       req.ProductID,
		}

		if err := q.LeasingContract.WithContext(ctx).Create(contract); err != nil {
			return err
		}

		// =========================
		// 5️⃣ Generate Payment Schedule
		// =========================
		pokokBulanan := pokok / float64(product.TenorBulan)
		marginBulanan := totalBunga / float64(product.TenorBulan)

		for i := 1; i <= int(product.TenorBulan); i++ {

			jatuhTempo := now.AddDate(0, i-1, 0)
			status := "unpaid"

			schedule := &model.PaymentSchedule{
				AngsuranKe:       int16(i),
				JatuhTempo:       jatuhTempo,
				Pokok:            pokokBulanan,
				Margin:           marginBulanan,
				TotalTagihan:     pokokBulanan + marginBulanan,
				StatusPembayaran: &status,
				ContractID:       contract.ContractID,
			}

			if err := q.PaymentSchedule.WithContext(ctx).Create(schedule); err != nil {
				return err
			}
		}

		// =========================
		// 6️⃣ Generate Leasing Tasks (Workflow)
		// =========================
		statusInProgress := "inprogress"

		tasksTemplate := []struct {
			Name     string
			RoleName string
			Sequence int16
		}{
			{"Input Pengajuan & Upload Dokumen", "SALES", 1},
			{"Auto Scoring & Review", "ADMIN_CABANG", 2},
			{"Survey Lapangan", "SURVEYOR", 3},
			{"Approval Final", "ADMIN_CABANG", 4},
			{"Akad & Tanda Tangan", "SALES", 5},
			{"Monitoring Pembayaran", "COLLECTION", 6},
		}

		// Map role_name → role_id
		roleMap := make(map[string]int64)
		roles, err := q.Role.WithContext(ctx).Find()
		if err != nil {
			return err
		}

		for _, r := range roles {
			roleMap[r.RoleName] = r.RoleID
		}

		for _, t := range tasksTemplate {

			roleID, ok := roleMap[t.RoleName]
			if !ok {
				return errors.New("role not found: " + t.RoleName)
			}

			task := &model.LeasingTask{
				TaskName:   t.Name,
				Startdate:  &now,
				Status:     &statusInProgress,
				ContractID: contract.ContractID,
				RoleID:     roleID,
				SequenceNo: &t.Sequence,
			}

			if err := q.LeasingTask.WithContext(ctx).Create(task); err != nil {
				return err
			}
		}

		createdContract = contract
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdContract, nil
}

func (s *leasingService) GetContractDetail(
	ctx context.Context,
	id int64,
) (*ContractDetailResponse, error) {

	q := query.Use(s.db)

	contract, err := q.LeasingContract.WithContext(ctx).
		Where(q.LeasingContract.ContractID.Eq(id)).
		First()
	if err != nil {
		return nil, errors.New("contract not found")
	}

	schedules, err := q.PaymentSchedule.WithContext(ctx).
		Where(q.PaymentSchedule.ContractID.Eq(id)).
		Order(q.PaymentSchedule.AngsuranKe.Asc()).
		Find()
	if err != nil {
		return nil, err
	}

	total := len(schedules)
	paid := 0

	for _, sc := range schedules {
		if sc.StatusPembayaran != nil &&
			*sc.StatusPembayaran == "paid" {
			paid++
		}
	}

	progress := 0.0
	if total > 0 {
		progress = (float64(paid) / float64(total)) * 100
	}

	return &ContractDetailResponse{
		Contract:  contract,
		Schedules: schedules,
		Summary: ContractSummary{
			TotalCicilan:    total,
			SudahDibayar:    paid,
			SisaCicilan:     total - paid,
			ProgressPercent: progress,
		},
	}, nil
}

const (
	ContractStatusDraft    = "draft"
	ContractStatusApproved = "approved"
	ContractStatusActive   = "active"
)

func (s *leasingService) ApproveContract(
	ctx context.Context,
	id int64,
) error {

	q := query.Use(s.db)

	contract, err := q.LeasingContract.WithContext(ctx).
		Where(q.LeasingContract.ContractID.Eq(id)).
		First()

	if err != nil {
		return errors.New("contract not found")
	}

	if contract.Status != ContractStatusApproved {
		return errors.New("only approved contract can be activated")
	}

	_, err = q.LeasingContract.WithContext(ctx).
		Where(q.LeasingContract.ContractID.Eq(id)).
		Update(q.LeasingContract.Status, ContractStatusActive)

	return err
}

// func (s *leasingService) ActivateContract(
// 	ctx context.Context,
// 	id int64,
// ) error {

// 	q := query.Use(s.db)

// 	contract, err := q.LeasingContract.WithContext(ctx).
// 		Where(q.LeasingContract.ContractID.Eq(id)).
// 		First()

// 	if err != nil {
// 		return errors.New("contract not found")
// 	}

// 	if contract.Status != ContractStatusApproved {
// 		return errors.New("only approved contract can be activated")
// 	}

// 	_, err = q.LeasingContract.WithContext(ctx).
// 		Where(q.LeasingContract.ContractID.Eq(id)).
// 		Update(q.LeasingContract.Status, ContractStatusActive)

// 	return err
// }

func (s *leasingService) ListContracts(
	ctx context.Context,
	page int,
	limit int,
	status string,
) ([]*model.LeasingContract, int64, error) {

	q := query.Use(s.db)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	baseQuery := q.LeasingContract.WithContext(ctx)

	if status != "" {
		baseQuery = baseQuery.Where(q.LeasingContract.Status.Eq(status))
	}

	total, err := baseQuery.Count()
	if err != nil {
		return nil, 0, err
	}

	contracts, err := baseQuery.
		Order(q.LeasingContract.CreatedAt.Desc()).
		Limit(limit).
		Offset(offset).
		Find()

	if err != nil {
		return nil, 0, err
	}

	return contracts, total, nil
}

func (s *leasingService) CompleteTask(
	ctx context.Context,
	taskID int64,
) error {

	return s.db.Transaction(func(tx *gorm.DB) error {

		q := query.Use(tx)

		// 1️⃣ Ambil task
		task, err := q.LeasingTask.WithContext(ctx).
			Where(q.LeasingTask.TaskID.Eq(taskID)).
			First()

		if err != nil {
			return errors.New("task not found")
		}

		// 2️⃣ Update task status
		statusCompleted := "completed"

		_, err = q.LeasingTask.WithContext(ctx).
			Where(q.LeasingTask.TaskID.Eq(taskID)).
			Update(q.LeasingTask.Status, statusCompleted)

		if err != nil {
			return err
		}

		// 3️⃣ Ambil semua task contract
		tasks, err := q.LeasingTask.WithContext(ctx).
			Where(q.LeasingTask.ContractID.Eq(task.ContractID)).
			Find()

		if err != nil {
			return err
		}

		allCompleted := true

		for _, t := range tasks {
			if t.Status == nil || *t.Status != "completed" {
				allCompleted = false
				break
			}
		}

		if allCompleted {

			contract, err := q.LeasingContract.WithContext(ctx).
				Where(q.LeasingContract.ContractID.Eq(task.ContractID)).
				First()

			if err != nil {
				return err
			}

			// Hanya dari draft → approved
			if contract.Status == ContractStatusDraft {

				_, err = q.LeasingContract.WithContext(ctx).
					Where(q.LeasingContract.ContractID.Eq(contract.ContractID)).
					Update(q.LeasingContract.Status, ContractStatusApproved)

				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
