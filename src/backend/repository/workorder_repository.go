package repository

import (
	"database/sql"
	"fmt"

	"teamitmivhs/work-order-backend/models"
)

type WorkOrderRepository interface {
	CreateTask(task models.WorkOrderRequest) (int64, error)
	TakeOrder(orderID int64, req models.TakeWorkOrder) error
	CompleteOrder(orderID int64, req models.CompleteWorkOrder) error
	DeleteOrder(orderID int64) error
	GetAllTasks() ([]models.WorkOrder, error)
}

type workOrderRepository struct {
	db *sql.DB
}

func NewWorkOrderRepository(db *sql.DB) WorkOrderRepository {
	return &workOrderRepository{db: db}
}

// Implementasi GetAllTasks - mengambil semua work orders dari database
func (r *workOrderRepository) GetAllTasks() ([]models.WorkOrder, error) {
	query := `
        SELECT ID, Priority, TimeDisplay, Requester, Location, Device, Problem, WorkingHours, Status, CompletedAt
        FROM orders
        ORDER BY TimeSort DESC
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("querying orders failed: %w", err)
	}
	defer rows.Close()

	var workOrders []models.WorkOrder

	for rows.Next() {
		var wo models.WorkOrder
		var priority, timeDisplay, requester, location, device, problem, workingHours, status, completedAt sql.NullString

		// Urutan Scan HARUS sama persis dengan urutan SELECT
		err := rows.Scan(
			&wo.ID, &priority, &timeDisplay, &requester, &location,
			&device, &problem, &workingHours, &status, &completedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning order row failed: %w", err)
		}

		// Cek .Valid sebelum mapping ke struct untuk menghindari error NULL
		if priority.Valid {
			wo.Priority = priority.String
		}
		if timeDisplay.Valid {
			wo.Time = timeDisplay.String // Mapping TimeDisplay ke Time untuk JavaScript
		}
		if requester.Valid {
			wo.Requester = requester.String
		}
		if location.Valid {
			wo.Location = location.String
		}
		if device.Valid {
			wo.Device = device.String
		}
		if problem.Valid {
			wo.Problem = problem.String
		}
		if status.Valid {
			wo.Status = status.String
		}
		if completedAt.Valid && completedAt.String != "" {
			wo.CompletedAt = completedAt.String
		}
		if workingHours.Valid {
			var hours int
			if _, err := fmt.Sscanf(workingHours.String, "%d", &hours); err == nil {
				wo.WorkingHours = &hours
			}
		}

		// Ambil data executors untuk setiap work order
		executorsQuery := "SELECT Executors FROM executors WHERE ID = ?"
		execRows, err := r.db.Query(executorsQuery, wo.ID)
		if err != nil {
			return nil, fmt.Errorf("querying executors for order %d failed: %w", wo.ID, err)
		}
		defer execRows.Close()

		var executors []int
		for execRows.Next() {
			var executorID int
			if err := execRows.Scan(&executorID); err != nil {
				return nil, fmt.Errorf("scanning executor for order %d failed: %w", wo.ID, err)
			}
			executors = append(executors, executorID)
		}
		wo.Executors = executors

		// Ambil data safety checklist untuk setiap work order
		checklistQuery := "SELECT SafetyChecklist FROM safetychecklist WHERE ID = ?"
		checkRows, err := r.db.Query(checklistQuery, wo.ID)
		if err != nil {
			return nil, fmt.Errorf("querying safety checklist for order %d failed: %w", wo.ID, err)
		}
		defer checkRows.Close()

		var checklist []string
		for checkRows.Next() {
			var item string
			if err := checkRows.Scan(&item); err != nil {
				return nil, fmt.Errorf("scanning safety checklist item for order %d failed: %w", wo.ID, err)
			}
			checklist = append(checklist, item)
		}
		wo.SafetyChecklist = checklist

		workOrders = append(workOrders, wo)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration over order rows failed: %w", err)
	}

	return workOrders, nil
}

// Implementasi CreateTask - membuat work order baru
func (r *workOrderRepository) CreateTask(task models.WorkOrderRequest) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	taskInsertQuery := `
        INSERT INTO orders 
        (Priority, TimeDisplay, TimeSort, Requester, Location, Device, Problem, WorkingHours, Status) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	result, err := tx.Exec(
		taskInsertQuery,
		task.Priority, task.TimeDisplay, task.TimeSort, task.Requester,
		task.Location, task.Device, task.Problem, task.WorkingHours, task.Status,
	)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

// Implementasi TakeOrder - mengambil work order untuk diproses
func (r *workOrderRepository) TakeOrder(orderID int64, req models.TakeWorkOrder) error {
	// Cek status executor sebelum memulai transaksi
	for _, executorID := range req.Executors {
		var status string
		err := r.db.QueryRow("SELECT Status FROM members WHERE ID = ?", executorID).Scan(&status)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("executor with ID %d not found", executorID)
			}
			return fmt.Errorf("failed to query executor status: %w", err)
		}
		if status == "onjob" {
			return fmt.Errorf("executor with ID %d is already on another job", executorID)
		}
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. UPDATE status task
	_, err = tx.Exec("UPDATE orders SET Status = ?, TimeSort = NOW() WHERE ID = ?", req.Status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// 2. Hapus executors lama sebelum insert yang baru
	_, err = tx.Exec("DELETE FROM executors WHERE ID = ?", orderID)
	if err != nil {
		return fmt.Errorf("failed to delete old executors: %w", err)
	}

	// 3. Insert executors baru
	for _, executorID := range req.Executors {
		_, err = tx.Exec("INSERT INTO executors (ID, Executors) VALUES (?, ?)", orderID, executorID)
		if err != nil {
			return fmt.Errorf("failed to insert executor: %w", err)
		}
		// Update status member menjadi 'onjob'
		_, err = tx.Exec("UPDATE members SET Status = 'onjob' WHERE ID = ?", executorID)
		if err != nil {
			return fmt.Errorf("failed to update member status: %w", err)
		}
	}

	// 4. Hapus safety checklist lama sebelum insert yang baru
	_, err = tx.Exec("DELETE FROM safetychecklist WHERE ID = ?", orderID)
	if err != nil {
		return fmt.Errorf("failed to delete old safety checklist: %w", err)
	}

	// 5. Insert safety checklist baru
	for _, item := range req.SafetyChecklistItems {
		_, err = tx.Exec("INSERT INTO safetychecklist (ID, SafetyChecklist) VALUES (?, ?)", orderID, item)
		if err != nil {
			return fmt.Errorf("failed to insert safety checklist item: %w", err)
		}
	}

	return tx.Commit()
}

// Implementasi CompleteOrder - menandai work order sebagai selesai
func (r *workOrderRepository) CompleteOrder(orderID int64, req models.CompleteWorkOrder) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. UPDATE status dan waktu selesai task
	_, err = tx.Exec("UPDATE orders SET Status = ?, CompletedAt = ? WHERE ID = ?",
		req.Status, req.CompletedAtDisplay, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order completion: %w", err)
	}

	// 2. Ambil semua executor dari order ini
	rows, err := tx.Query("SELECT Executors FROM executors WHERE ID = ?", orderID)
	if err != nil {
		return fmt.Errorf("failed to query executors for order completion: %w", err)
	}
	defer rows.Close()

	var executorIDs []int
	for rows.Next() {
		var executorID int
		if err := rows.Scan(&executorID); err != nil {
			return fmt.Errorf("failed to scan executor ID: %w", err)
		}
		executorIDs = append(executorIDs, executorID)
	}

	// 3. Update status semua executor menjadi 'standby'
	for _, executorID := range executorIDs {
		_, err := tx.Exec("UPDATE members SET Status = 'standby' WHERE ID = ?", executorID)
		if err != nil {
			return fmt.Errorf("failed to update member status to standby: %w", err)
		}
	}

	return tx.Commit()
}

// Implementasi DeleteOrder - menghapus work order
func (r *workOrderRepository) DeleteOrder(orderID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Hapus dari tabel child dulu untuk menghindari masalah foreign key
	if _, err := tx.Exec("DELETE FROM executors WHERE ID = ?", orderID); err != nil {
		return fmt.Errorf("failed to delete from executors: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM safetychecklist WHERE ID = ?", orderID); err != nil {
		return fmt.Errorf("failed to delete from safetychecklist: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM orders WHERE ID = ?", orderID); err != nil {
		return fmt.Errorf("failed to delete from orders: %w", err)
	}

	return tx.Commit()
}
