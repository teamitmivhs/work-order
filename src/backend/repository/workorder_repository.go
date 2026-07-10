package repository

import (
	"database/sql"
	"fmt"
	"log"

	"teamitmivhs/work-order-backend/models"
	"teamitmivhs/work-order-backend/services"
)

type WorkOrderRepository interface {
	CreateTask(task models.WorkOrderRequest) (int64, error)
	JoinPendingOrder(orderID int64, memberID int) error
	TakeOrder(orderID int64, req models.TakeWorkOrder) error
	CompleteOrder(orderID int64, req models.CompleteWorkOrder) error
	DeleteOrder(orderID int64) error
	UpdateOrderExecutors(orderID int64, req models.UpdateWorkOrderRequest) error
	UpdateOrderNotes(orderID int64, adminNotes string, rating *int, notesQuality *int) error
	UpdateOrderNoteText(orderID int64, notes string) error
	UpdateOrderDocumentationPhoto(orderID int64, filename string) error
	HasDocumentationPhoto(orderID int64) (bool, error)
	GetAllTasks() ([]models.WorkOrder, error)
	GetTasksForOperator(operatorID int) ([]models.WorkOrder, error)
	GetTaskByTrackingCode(trackingCode string) (models.WorkOrder, error)
	GetSafetyChecklist(orderID int64) ([]string, error)
	UpdateSafetyChecklist(orderID int64, items []string) error
	IsSafetyChecklistFulfilled(orderID int64) (bool, error)
	GetKaizenMetrics() (models.Kaizen, error)
}

type workOrderRepository struct {
	db *sql.DB
}

func NewWorkOrderRepository(db *sql.DB) WorkOrderRepository {
	return &workOrderRepository{db: db}
}

const progressTimerSelect = `
		       CASE
		         WHEN o.Status = 'progress' AND o.StartedAt IS NOT NULL
		           THEN DATE_FORMAT(o.StartedAt, '%Y-%m-%d %H:%i:%s')
		         WHEN o.Status = 'progress' AND o.TimeSort IS NOT NULL
		           THEN DATE_FORMAT(TIMESTAMP(CURDATE(), o.TimeSort), '%Y-%m-%d %H:%i:%s')
		         ELSE NULL
		       END AS StartedAt,
		       CASE
		         WHEN o.Status = 'progress' AND o.StartedAt IS NOT NULL
		           THEN GREATEST(TIMESTAMPDIFF(SECOND, o.StartedAt, NOW()), 0)
		         WHEN o.Status = 'progress' AND o.TimeSort IS NOT NULL THEN
		           CASE
		             WHEN TIME_TO_SEC(TIMEDIFF(CURTIME(), o.TimeSort)) < 0
		               THEN TIME_TO_SEC(TIMEDIFF(CURTIME(), o.TimeSort)) + 86400
		             ELSE TIME_TO_SEC(TIMEDIFF(CURTIME(), o.TimeSort))
		           END
		         ELSE NULL
		       END AS ProgressSeconds`

// scanWorkOrderRow adalah helper untuk menghindari duplikasi scan di endpoint list work order.
func scanWorkOrderRow(rows *sql.Rows) (models.WorkOrder, error) {
	var wo models.WorkOrder
	var priority, timeDisplay, createdAt, startedAt, trackingCode, requester, location, device, problem, workingHours, status, completedAt, notes, adminNotes, documentationPhoto sql.NullString
	var progressSeconds, rating, notesQuality sql.NullInt64

	err := rows.Scan(
		&wo.ID, &priority, &timeDisplay, &createdAt, &startedAt, &progressSeconds, &trackingCode, &requester, &location,
		&device, &problem, &workingHours, &status, &completedAt, &notes, &adminNotes, &rating, &notesQuality, &documentationPhoto,
	)
	if err != nil {
		return wo, err
	}

	if priority.Valid {
		wo.Priority = priority.String
	}
	if timeDisplay.Valid {
		wo.Time = timeDisplay.String
	}
	if createdAt.Valid {
		wo.CreatedAt = createdAt.String
	}
	if startedAt.Valid {
		wo.StartedAt = startedAt.String
	}
	if progressSeconds.Valid {
		value := int(progressSeconds.Int64)
		wo.ProgressSeconds = &value
	}
	if trackingCode.Valid {
		wo.TrackingCode = trackingCode.String
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
	if notes.Valid {
		wo.Notes = notes.String
	}
	if adminNotes.Valid {
		wo.AdminNotes = adminNotes.String
	}
	if rating.Valid {
		value := int(rating.Int64)
		wo.Rating = &value
	}
	if notesQuality.Valid {
		value := int(notesQuality.Int64)
		wo.NotesQuality = &value
	}
	if documentationPhoto.Valid {
		wo.DocumentationPhoto = documentationPhoto.String
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

	return wo, nil
}

// GetAllTasks mengambil semua work orders dari database
func (r *workOrderRepository) GetAllTasks() ([]models.WorkOrder, error) {
	query := `
		SELECT DISTINCT o.ID, o.Priority, o.TimeDisplay,
		       DATE_FORMAT(COALESCE(o.CreatedAt, TIMESTAMP(CURDATE(), o.TimeSort)), '%Y-%m-%d %H:%i:%s') AS CreatedAt,` + progressTimerSelect + `,
		       o.TrackingCode,
		       o.Requester, o.Location, o.Device, o.Problem, o.WorkingHours, o.Status,
		       o.CompletedAt, o.Notes, o.AdminNotes, o.Rating, o.NotesQuality, o.DocumentationPhoto
		FROM orders o
		ORDER BY o.ID DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("Error querying orders: %v", err)
		return nil, fmt.Errorf("querying orders failed: %w", err)
	}
	defer rows.Close()

	// FIX: inisialisasi sebagai empty slice agar JSON response [] bukan null
	workOrders := make([]models.WorkOrder, 0)

	for rows.Next() {
		wo, err := scanWorkOrderRow(rows)
		if err != nil {
			log.Printf("Error scanning order row: %v", err)
			return nil, fmt.Errorf("scanning order row failed: %w", err)
		}

		executors, err := r.getOrderExecutors(wo.ID)
		if err != nil {
			log.Printf("Error getting executors for order %d: %v", wo.ID, err)
			return nil, err
		}
		wo.Executors = executors
		workOrders = append(workOrders, wo)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	return workOrders, nil
}

// GetTasksForOperator mengambil order yang bisa dikerjakan operator:
// - semua order pending agar operator bisa mengambil pekerjaan baru
// - order progress/completed yang memang ditugaskan ke operator tersebut
func (r *workOrderRepository) GetTasksForOperator(operatorID int) ([]models.WorkOrder, error) {
	query := `
		SELECT DISTINCT o.ID, o.Priority, o.TimeDisplay,
		       DATE_FORMAT(COALESCE(o.CreatedAt, TIMESTAMP(CURDATE(), o.TimeSort)), '%Y-%m-%d %H:%i:%s') AS CreatedAt,` + progressTimerSelect + `,
		       o.TrackingCode,
		       o.Requester, o.Location, o.Device, o.Problem, o.WorkingHours, o.Status,
		       o.CompletedAt, o.Notes, o.AdminNotes, o.Rating, o.NotesQuality, o.DocumentationPhoto
		FROM orders o
		WHERE o.Status = 'pending'
		   OR EXISTS (
		       SELECT 1
		       FROM executors e
		       WHERE e.order_id = o.ID
		         AND e.member_id = ?
		   )
		ORDER BY o.ID DESC
	`
	rows, err := r.db.Query(query, operatorID)
	if err != nil {
		log.Printf("Error querying orders for operator: %v", err)
		return nil, fmt.Errorf("querying orders failed: %w", err)
	}
	defer rows.Close()

	workOrders := make([]models.WorkOrder, 0)

	for rows.Next() {
		wo, err := scanWorkOrderRow(rows)
		if err != nil {
			log.Printf("Error scanning order row: %v", err)
			return nil, err
		}

		executors, err := r.getOrderExecutors(wo.ID)
		if err != nil {
			return nil, err
		}
		wo.Executors = executors
		workOrders = append(workOrders, wo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workOrders, nil
}

func (r *workOrderRepository) GetTaskByTrackingCode(trackingCode string) (models.WorkOrder, error) {
	query := `
		SELECT DISTINCT o.ID, o.Priority, o.TimeDisplay,
		       DATE_FORMAT(COALESCE(o.CreatedAt, TIMESTAMP(CURDATE(), o.TimeSort)), '%Y-%m-%d %H:%i:%s') AS CreatedAt,` + progressTimerSelect + `,
		       o.TrackingCode,
		       o.Requester, o.Location, o.Device, o.Problem, o.WorkingHours, o.Status,
		       o.CompletedAt, o.Notes, o.AdminNotes, o.Rating, o.NotesQuality, o.DocumentationPhoto
		FROM orders o
		WHERE UPPER(o.TrackingCode) = UPPER(?)
		LIMIT 1
	`
	rows, err := r.db.Query(query, trackingCode)
	if err != nil {
		return models.WorkOrder{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return models.WorkOrder{}, sql.ErrNoRows
	}

	wo, err := scanWorkOrderRow(rows)
	if err != nil {
		return wo, err
	}
	wo.Executors, err = r.getOrderExecutors(wo.ID)
	if err != nil {
		return wo, err
	}
	return wo, rows.Err()
}

// getOrderExecutors adalah helper untuk mengambil list executor ID sebuah order
func (r *workOrderRepository) getOrderExecutors(orderID int) ([]int, error) {
	rows, err := r.db.Query("SELECT member_id FROM executors WHERE order_id = ?", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executors := make([]int, 0)
	for rows.Next() {
		var executorID int
		if err := rows.Scan(&executorID); err != nil {
			return nil, err
		}
		executors = append(executors, executorID)
	}
	return executors, nil
}

// UpdateOrderNotes menyimpan catatan admin tanpa menimpa catatan guest.
func (r *workOrderRepository) UpdateOrderNotes(orderID int64, adminNotes string, rating *int, notesQuality *int) error {
	_, err := r.db.Exec(
		"UPDATE orders SET AdminNotes = ?, Rating = ?, NotesQuality = ? WHERE ID = ?",
		adminNotes, nullableInt(rating), nullableInt(notesQuality), orderID,
	)
	if err != nil {
		return fmt.Errorf("failed to update notes for order %d: %w", orderID, err)
	}
	return nil
}

func (r *workOrderRepository) UpdateOrderNoteText(orderID int64, notes string) error {
	_, err := r.db.Exec("UPDATE orders SET Notes = ? WHERE ID = ?", notes, orderID)
	if err != nil {
		return fmt.Errorf("failed to update note text for order %d: %w", orderID, err)
	}
	return nil
}

func nullableInt(value *int) interface{} {
	if value == nil {
		return nil
	}
	return *value
}

func (r *workOrderRepository) UpdateOrderDocumentationPhoto(orderID int64, filename string) error {
	_, err := r.db.Exec(
		"UPDATE orders SET DocumentationPhoto = ? WHERE ID = ?",
		filename, orderID,
	)
	if err != nil {
		return fmt.Errorf("failed to update documentation photo for order %d: %w", orderID, err)
	}
	return nil
}

func (r *workOrderRepository) HasDocumentationPhoto(orderID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM orders WHERE ID = ? AND DocumentationPhoto IS NOT NULL AND DocumentationPhoto <> ''",
		orderID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetSafetyChecklist mengambil item checklist untuk sebuah order
func (r *workOrderRepository) GetSafetyChecklist(orderID int64) ([]string, error) {
	rows, err := r.db.Query(
		"SELECT SafetyChecklist FROM safetychecklist WHERE order_id = ? ORDER BY SafetyChecklist ASC",
		orderID,
	)
	if err != nil {
		log.Printf("Error querying safety checklist: %v", err)
		return nil, err
	}
	defer rows.Close()

	// FIX: inisialisasi sebagai empty slice
	checklist := make([]string, 0)
	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			log.Printf("Error scanning safety checklist item: %v", err)
			return nil, err
		}
		checklist = append(checklist, item)
	}
	return checklist, nil
}

// UpdateSafetyChecklist mengganti semua checklist item untuk sebuah order
func (r *workOrderRepository) UpdateSafetyChecklist(orderID int64, items []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM safetychecklist WHERE order_id = ?", orderID); err != nil {
		return fmt.Errorf("failed to delete existing checklist: %w", err)
	}

	for _, item := range items {
		if _, err := tx.Exec(
			"INSERT INTO safetychecklist (order_id, SafetyChecklist) VALUES (?, ?)",
			orderID, item,
		); err != nil {
			return fmt.Errorf("failed to insert checklist item: %w", err)
		}
	}

	return tx.Commit()
}

// IsSafetyChecklistFulfilled mengecek apakah checklist sudah diisi untuk order ini
// FIX: return true jika tidak ada item checklist sama sekali (checklist opsional)
// Sebelumnya selalu return false jika tidak ada item → complete order selalu blocked
func (r *workOrderRepository) IsSafetyChecklistFulfilled(orderID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM safetychecklist WHERE order_id = ?",
		orderID,
	).Scan(&count)
	if err != nil {
		log.Printf("Error checking safety checklist count: %v", err)
		return false, err
	}
	// Jika tidak ada item checklist → dianggap fulfilled (checklist opsional)
	// Jika ada item → dianggap fulfilled (sudah diisi saat take order)
	return true, nil
}

// GetKaizenMetrics mengambil metrik performa work orders
func (r *workOrderRepository) GetKaizenMetrics() (models.Kaizen, error) {
	var metrics models.Kaizen

	err := r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE Status = 'completed'").
		Scan(&metrics.ImplementedKaizens)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying completed orders: %v", err)
		return metrics, err
	}

	err = r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE Status = 'pending'").
		Scan(&metrics.PendingKaizens)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying pending orders: %v", err)
		return metrics, err
	}

	// FIX: hitung juga order yang sedang progress agar TotalKaizens akurat
	var progressCount int
	err = r.db.QueryRow("SELECT COUNT(*) FROM orders WHERE Status = 'progress'").
		Scan(&progressCount)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying progress orders: %v", err)
		return metrics, err
	}

	metrics.TotalKaizens = metrics.ImplementedKaizens + metrics.PendingKaizens + progressCount
	return metrics, nil
}

// CreateTask membuat work order baru di database
func (r *workOrderRepository) CreateTask(task models.WorkOrderRequest) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		`INSERT INTO orders
		(CreatedAt, Priority, TimeDisplay, TimeSort, TrackingCode, Requester, Location, Device, Problem, WorkingHours, Status)
		VALUES (NOW(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		task.Priority, task.TimeDisplay, task.TimeSort, task.TrackingCode, task.Requester,
		task.Location, task.Device, task.Problem, task.WorkingHours, task.Status,
	)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, tx.Commit()
}

// TakeOrder mengambil work order untuk diproses oleh executor
func (r *workOrderRepository) JoinPendingOrder(orderID int64, memberID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var orderStatus string
	if err := tx.QueryRow("SELECT Status FROM orders WHERE ID = ? FOR UPDATE", orderID).Scan(&orderStatus); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("order not found")
		}
		return fmt.Errorf("failed to query order status: %w", err)
	}
	if orderStatus != "pending" {
		return fmt.Errorf("only pending orders can be joined")
	}

	var memberStatus string
	if err := tx.QueryRow("SELECT Status FROM members WHERE ID = ? FOR UPDATE", memberID).Scan(&memberStatus); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("member not found")
		}
		return fmt.Errorf("failed to query member status: %w", err)
	}
	if memberStatus == "onjob" {
		return fmt.Errorf("member is already on another job")
	}

	if _, err := tx.Exec(
		"INSERT IGNORE INTO executors (order_id, member_id) VALUES (?, ?)",
		orderID, memberID,
	); err != nil {
		return fmt.Errorf("failed to join pending order: %w", err)
	}

	return tx.Commit()
}

// TakeOrder mengambil work order untuk diproses oleh executor
func (r *workOrderRepository) TakeOrder(orderID int64, req models.TakeWorkOrder) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// FIX: cek status executor DI DALAM transaksi dengan SELECT FOR UPDATE
	// untuk mencegah race condition (sebelumnya cek di luar transaksi)
	for _, executorID := range req.Executors {
		var status string
		err := tx.QueryRow(
			"SELECT Status FROM members WHERE ID = ? FOR UPDATE",
			executorID,
		).Scan(&status)
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

	// 1. Update status order dan reset waktu mulai ke waktu take.
	// StartedAt menyimpan timestamp penuh untuk durasi lintas hari; TimeSort tetap
	// diisi untuk kompatibilitas UI/schema lama.
	if _, err = tx.Exec(
		"UPDATE orders SET Status = ?, StartedAt = NOW(), TimeSort = CURTIME(), TimeDisplay = DATE_FORMAT(CURTIME(), '%H:%i') WHERE ID = ?",
		req.Status, orderID,
	); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// 2. Hapus executor lama
	if _, err = tx.Exec("DELETE FROM executors WHERE order_id = ?", orderID); err != nil {
		return fmt.Errorf("failed to delete old executors: %w", err)
	}

	// 3. Insert executor baru dan update status member
	for _, executorID := range req.Executors {
		if _, err = tx.Exec(
			"INSERT INTO executors (order_id, member_id) VALUES (?, ?)",
			orderID, executorID,
		); err != nil {
			return fmt.Errorf("failed to insert executor: %w", err)
		}
		if _, err = tx.Exec(
			"UPDATE members SET Status = 'onjob' WHERE ID = ?",
			executorID,
		); err != nil {
			return fmt.Errorf("failed to update member status: %w", err)
		}
	}

	// 4. Replace safety checklist
	if _, err = tx.Exec("DELETE FROM safetychecklist WHERE order_id = ?", orderID); err != nil {
		return fmt.Errorf("failed to delete old safety checklist: %w", err)
	}
	for _, item := range req.SafetyChecklistItems {
		if _, err = tx.Exec(
			"INSERT INTO safetychecklist (order_id, SafetyChecklist) VALUES (?, ?)",
			orderID, item,
		); err != nil {
			return fmt.Errorf("failed to insert safety checklist item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// FIX: panggil StartTimer ke Rust engine SETELAH commit berhasil
	// Sebelumnya StartTimer tidak pernah dipanggil sama sekali sehingga
	// WorkingHours tidak pernah terisi otomatis.
	// Gunakan executor pertama sebagai referensi timer (satu timer per order)
	if len(req.Executors) > 0 {
		if err := services.StartTimer(uint64(orderID), uint64(req.Executors[0])); err != nil {
			// Timer gagal start tidak membatalkan take order — log saja
			// agar WorkingHours bisa diisi manual jika perlu
			log.Printf("[WARNING] Failed to start timer for order %d: %v", orderID, err)
		}
	}

	return nil
}

// CompleteOrder menandai work order sebagai selesai dan reset status executor
func (r *workOrderRepository) CompleteOrder(orderID int64, req models.CompleteWorkOrder) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Update status dan waktu selesai order
	if _, err = tx.Exec(
		"UPDATE orders SET Status = ?, CompletedAt = ? WHERE ID = ?",
		req.Status, req.CompletedAtDisplay, orderID,
	); err != nil {
		return fmt.Errorf("failed to update order completion: %w", err)
	}

	// 2. Ambil semua executor order ini
	rows, err := tx.Query("SELECT member_id FROM executors WHERE order_id = ?", orderID)
	if err != nil {
		return fmt.Errorf("failed to query executors for order completion: %w", err)
	}

	var executorIDs []int
	for rows.Next() {
		var executorID int
		if err := rows.Scan(&executorID); err != nil {
			rows.Close() // FIX: close eksplisit sebelum lanjut ke UPDATE
			return fmt.Errorf("failed to scan executor ID: %w", err)
		}
		executorIDs = append(executorIDs, executorID)
	}
	// FIX: close rows SEBELUM tx.Exec berikutnya untuk mencegah deadlock transaksi
	rows.Close()

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating executor rows: %w", err)
	}

	// 3. Reset status semua executor menjadi standby
	for _, executorID := range executorIDs {
		if _, err := tx.Exec(
			"UPDATE members SET Status = 'standby' WHERE ID = ?",
			executorID,
		); err != nil {
			return fmt.Errorf("failed to update member status to standby: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// FIX: panggil StopTimer ke Rust engine SETELAH commit berhasil
	// Durasi yang dikembalikan disimpan ke kolom WorkingHours di database
	durationSeconds, err := services.StopTimer(uint64(orderID))
	if err != nil {
		log.Printf("[WARNING] Failed to stop timer for order %d: %v", orderID, err)
		// Timer gagal stop tidak membatalkan complete order — WorkingHours dikosongkan
		return nil
	}

	// Simpan menit numerik agar frontend bisa format konsisten:
	// 0 = Less than 1 minute, 54 = 54 minutes, 94 = 1 hour 34 minutes.
	durationMinutes := durationSeconds / 60
	workingHoursStr := fmt.Sprintf("%d minutes", durationMinutes)
	if _, err := r.db.Exec(
		"UPDATE orders SET WorkingHours = ? WHERE ID = ?",
		workingHoursStr, orderID,
	); err != nil {
		log.Printf("[WARNING] Failed to update WorkingHours for order %d: %v", orderID, err)
	}

	return nil
}

// DeleteOrder menghapus work order beserta data terkait
// FIX: juga reset status member yang sedang onjob untuk order ini
func (r *workOrderRepository) DeleteOrder(orderID int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// FIX: ambil executor yang sedang onjob untuk order ini sebelum dihapus
	rows, err := tx.Query(
		"SELECT member_id FROM executors WHERE order_id = ?",
		orderID,
	)
	if err != nil {
		return fmt.Errorf("failed to query executors before delete: %w", err)
	}

	var executorIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return fmt.Errorf("failed to scan executor ID: %w", err)
		}
		executorIDs = append(executorIDs, id)
	}
	rows.Close() // FIX: close sebelum lanjut

	// FIX: reset status member yang sedang onjob karena order-nya dihapus
	for _, executorID := range executorIDs {
		if _, err := tx.Exec(
			"UPDATE members SET Status = 'standby' WHERE ID = ? AND Status = 'onjob'",
			executorID,
		); err != nil {
			return fmt.Errorf("failed to reset member status on delete: %w", err)
		}
	}

	// Hapus child tables dulu (foreign key constraint)
	if _, err := tx.Exec("DELETE FROM executors WHERE order_id = ?", orderID); err != nil {
		return fmt.Errorf("failed to delete from executors: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM safetychecklist WHERE order_id = ?", orderID); err != nil {
		return fmt.Errorf("failed to delete from safetychecklist: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM orders WHERE ID = ?", orderID); err != nil {
		return fmt.Errorf("failed to delete from orders: %w", err)
	}

	return tx.Commit()
}

// UpdateOrderExecutors memperbarui daftar executor untuk order yang masih pending
// FIX: method baru untuk endpoint PATCH /api/workorders/:id yang dipanggil frontend
func (r *workOrderRepository) UpdateOrderExecutors(orderID int64, req models.UpdateWorkOrderRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Ambil executor lama untuk reset status jika dihapus
	oldRows, err := tx.Query("SELECT member_id FROM executors WHERE order_id = ?", orderID)
	if err != nil {
		return fmt.Errorf("failed to query old executors: %w", err)
	}

	oldExecutors := make(map[int]bool)
	for oldRows.Next() {
		var id int
		if err := oldRows.Scan(&id); err != nil {
			oldRows.Close()
			return err
		}
		oldExecutors[id] = true
	}
	oldRows.Close()

	newExecutors := make(map[int]bool)
	for _, id := range req.Executors {
		newExecutors[id] = true
	}

	// Reset status member yang dihapus dari executor list
	for id := range oldExecutors {
		if !newExecutors[id] {
			if _, err := tx.Exec(
				"UPDATE members SET Status = 'standby' WHERE ID = ? AND Status = 'onjob'",
				id,
			); err != nil {
				return fmt.Errorf("failed to reset removed executor status: %w", err)
			}
		}
	}

	// Ganti executor list
	if _, err := tx.Exec("DELETE FROM executors WHERE order_id = ?", orderID); err != nil {
		return fmt.Errorf("failed to delete old executors: %w", err)
	}
	for _, executorID := range req.Executors {
		if _, err := tx.Exec(
			"INSERT INTO executors (order_id, member_id) VALUES (?, ?)",
			orderID, executorID,
		); err != nil {
			return fmt.Errorf("failed to insert executor: %w", err)
		}
	}

	// Update status order jika ada
	if req.Status != nil {
		if _, err := tx.Exec(
			"UPDATE orders SET Status = ? WHERE ID = ?",
			*req.Status, orderID,
		); err != nil {
			return fmt.Errorf("failed to update order status: %w", err)
		}
	}

	return tx.Commit()
}
