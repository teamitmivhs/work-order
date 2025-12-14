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

		// Set default values untuk field yang tidak ada di DB tapi diperlukan frontend
		wo.Executors = []int{}
		wo.SafetyChecklist = []string{}
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

	// 2. Store executors in a JSON format in Problem field (since no separate table exists)
	if len(req.Executors) > 0 {
		executorsJSON := fmt.Sprintf("Executors: %v", req.Executors)
		_, err = tx.Exec("UPDATE orders SET Problem = CONCAT(Problem, '\n', ?) WHERE ID = ?", executorsJSON, orderID)
		if err != nil {
			return fmt.Errorf("failed to store executors: %w", err)
		}
	}

	// 3. Store safety checklist items in Location field as notes
	if len(req.SafetyChecklistItems) > 0 {
		checklistJSON := fmt.Sprintf("Safety Checklist: %v", req.SafetyChecklistItems)
		_, err = tx.Exec("UPDATE orders SET Location = CONCAT(Location, '\n', ?) WHERE ID = ?", checklistJSON, orderID)
		if err != nil {
			return fmt.Errorf("failed to store safety checklist: %w", err)
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

	return tx.Commit()
}

// Implementasi DeleteOrder - menghapus work order
func (r *workOrderRepository) DeleteOrder(orderID int64) error {
	// Simple delete from orders table since no separate tables exist
	_, err := r.db.Exec("DELETE FROM orders WHERE ID = ?", orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
