package repository

import (
	"database/sql"
	"fmt"
	"log"

	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/models"
)

type WorkOrderRepository interface {
	CreateTask(task models.WorkOrderRequest) (int64, error)
	TakeOrder(orderID int64, req models.TakeWorkOrder) error
	CompleteOrder(orderID int64, req models.CompleteWorkOrder) error
	DeleteOrder(orderID int64) error
	GetAllTasks() ([]models.WorkOrder, error)
	// ...
}

// Implementasi GetAllTasks
func (r *workOrderRepository) GetAllTasks() ([]models.WorkOrder, error) {
	query := `
        SELECT ID, Priority, TimeDisplay, Requester, Location, Device, Problem, WorkingHours, Status
        FROM orders
        ORDER BY TimeSort DESC
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("querying orders failed: %w", err)
	}
	defer rows.Close()

	taskMap := make(map[int]*models.WorkOrder)
	var allTaskIDs []interface{}

	for rows.Next() {
		var wo models.WorkOrder
		var priority, timeDisplay, requester, location, device, problem, workingHours, status sql.NullString

		// Urutan Scan HARUS sama persis dengan urutan SELECT
		err := rows.Scan(
			&wo.ID, &priority, &timeDisplay, &requester, &location,
			&device, &problem, &workingHours, &status,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning order row failed: %w", err)
		}

		// Cek .Valid sebelum mapping ke struct untuk menghindari error NULL
		if priority.Valid {
			wo.Priority = priority.String
		}
		if timeDisplay.Valid {
			wo.Time = timeDisplay.String
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
		if workingHours.Valid {
			var hours int
			if _, err := fmt.Sscanf(workingHours.String, "%d", &hours); err == nil {
				wo.WorkingHours = &hours
			}
		}

		wo.Executors = []int{}
		taskMap[wo.ID] = &wo
		allTaskIDs = append(allTaskIDs, wo.ID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration over order rows failed: %w", err)
	}
	
	if len(allTaskIDs) == 0 {
		return []models.WorkOrder{}, nil
	}

	// Ambil semua executor dalam satu query yang efisien
	placeholders := ""
	for i := 0; i < len(allTaskIDs); i++ {
		placeholders += "?"
		if i < len(allTaskIDs)-1 {
			placeholders += ","
		}
	}
	execQuery := fmt.Sprintf("SELECT ID, Executors FROM executors WHERE ID IN (%s)", placeholders)

	execRows, err := r.db.Query(execQuery, allTaskIDs...)
	if err != nil {
		return nil, fmt.Errorf("querying executors failed: %w", err)
	}
	defer execRows.Close()

	for execRows.Next() {
		var taskID, memberID int
		if err := execRows.Scan(&taskID, &memberID); err != nil {
			return nil, fmt.Errorf("scanning executor row failed: %w", err)
		}
		if task, ok := taskMap[taskID]; ok {
			task.Executors = append(task.Executors, memberID)
		}
	}

	// Konversi map kembali ke slice untuk menjaga urutan asli jika diperlukan
	// (saat ini tidak menjaga urutan, tapi cukup untuk JSON)
	finalWorkOrders := make([]models.WorkOrder, 0, len(taskMap))
	for _, task := range taskMap {
		finalWorkOrders = append(finalWorkOrders, *task)
	}

	return finalWorkOrders, nil
}

// Implementasi TakeOrder
func (r *workOrderRepository) TakeOrder(orderID int64, req models.TakeWorkOrder) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. UPDATE status task
	_, err = tx.Exec("UPDATE orders SET Status = ?, TimeSort = NOW() WHERE ID = ?", req.Status, orderID)
	if err != nil {
		return err
	}

	// 2. INSERT ke executors (Loop melalui semua executor IDs)
	for _, execID := range req.Executors {
		// Asumsi: Anda memiliki tabel users/members dengan kolom status
		_, err = tx.Exec("INSERT INTO executors (ID, Executors) VALUES (?, ?)", orderID, execID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE members SET Status = 'onjob' WHERE ID = ?", execID)
		if err != nil {
			return err
		}
	}

	// 3. INSERT ke safetychecklist (Loop melalui semua checklist IDs)
	for _, itemID := range req.SafetyChecklistItems {
		// Asumsi: itemID adalah ID dari tabel master safety_checklist_items
		_, err = tx.Exec("INSERT INTO safetychecklist (ID, SafetyChecklist) VALUES (?, ?)", orderID, itemID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Implementasi CompleteOrder
func (r *workOrderRepository) CompleteOrder(orderID int64, req models.CompleteWorkOrder) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. UPDATE status dan waktu selesai task
	_, err = tx.Exec("UPDATE orders SET Status = ?, CompletedAt = ? WHERE ID = ?",
		req.Status, req.CompletedAtDisplay, orderID)
	if err != nil {
		return err
	}

	// 2. Dapatkan semua executor ID yang terkait dengan task ini (dari executors)
	rows, err := tx.Query("SELECT Executors FROM executors WHERE ID = ?", orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var executorIDs []int64
	for rows.Next() {
		var execID int64
		if err := rows.Scan(&execID); err != nil {
			return err
		}
		executorIDs = append(executorIDs, execID)
	}

	// 3. UPDATE status semua executor menjadi 'standby'
	for _, execID := range executorIDs {
		// Hanya update member yang statusnya masih 'onjob'
		_, err = tx.Exec("UPDATE members SET Status = 'standby' WHERE ID = ? AND Status = 'onjob'", execID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Implementasi DeleteOrder
func (r *workOrderRepository) DeleteOrder(orderID int64) error {
	// Transaksi diperlukan untuk menghapus dari tabel relasi terlebih dahulu
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM executors WHERE ID = ?", orderID)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM safetychecklist WHERE ID = ?", orderID)
	if err != nil {
		return err
	}

	// 2. Hapus task utama
	_, err = tx.Exec("DELETE FROM orders WHERE ID = ?", orderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

type workOrderRepository struct {
	db *sql.DB
}

func NewWorkOrderRepository(db *sql.DB) WorkOrderRepository {
	return &workOrderRepository{db: db}
}

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

func GetAllMembers() ([]models.Member, error) {
	rows, err := config.DB.Query("SELECT ID, Name, Role, Status, Avatar FROM members")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		var m models.Member
		if err := rows.Scan(&m.ID, &m.Name, &m.Role, &m.Status, &m.Avatar); err != nil {
			log.Printf("Error scanning member row: %v", err)
			return nil, err
		}
		members = append(members, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func CreateMember(member *models.Member) error {
	_, err := config.DB.Exec("INSERT INTO members (Name, Password, Role, Status, Avatar) VALUES (?, ?, ?, ?, ?)", member.Name, member.Password, member.Role, member.Status, member.Avatar)
	return err
}

func GetMemberByName(name string) (*models.Member, error) {
	row := config.DB.QueryRow("SELECT ID, Name, Password, Role, Status, Avatar FROM members WHERE Name = ?", name)

	var m models.Member
	if err := row.Scan(&m.ID, &m.Name, &m.Password, &m.Role, &m.Status, &m.Avatar); err != nil {
		return nil, err
	}
	return &m, nil
}
