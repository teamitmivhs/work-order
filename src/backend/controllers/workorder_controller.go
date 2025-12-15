package controllers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"teamitmivhs/work-order-backend/config"
	"teamitmivhs/work-order-backend/models"
	"teamitmivhs/work-order-backend/repository"

	"github.com/gin-gonic/gin"
)

type WorkOrderController struct {
	Repo repository.WorkOrderRepository
}

var (
	mu         sync.Mutex
	workOrders = []models.WorkOrder{}
	nextID     = 1
)

func NewWorkOrderController(repo repository.WorkOrderRepository) *WorkOrderController {
	return &WorkOrderController{Repo: repo}
}

// GetTaskListHandler menangani request GET /api/workorders
func (ctrl *WorkOrderController) GetTaskListHandler(c *gin.Context) {
	tasks, err := ctrl.Repo.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks from database", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// CreateTaskHandler menangani request POST /api/workorders
func (ctrl *WorkOrderController) CreateTaskHandler(c *gin.Context) {
	var req models.WorkOrderRequest

	// 1. Baca dan Bind Payload dari JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// 2. Validasi sederhana
	if req.Priority == "" || req.Requester == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Priority and Requester are required"})
		return
	}

	// Pastikan Status di set ke pending jika tidak ada atau salah
	if req.Status == "" {
		req.Status = "pending"
	}

	// 3. Panggil Repository untuk menyimpan data
	newID, err := ctrl.Repo.CreateTask(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save task to database"})
		return
	}

	// 4. Kirim Respon Sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "Work Order created successfully",
		"id":      newID, // Mengembalikan ID yang baru dibuat
	})
}

// TakeOrderHandler: POST /api/workorders/{id}/take
func (ctrl *WorkOrderController) TakeOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
		return
	}

	var req models.TakeWorkOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.Status != "progress" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'progress' to take order"})
		return
	}

	err = ctrl.Repo.TakeOrder(orderID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to take order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order taken successfully", "order_id": orderID})
}

// CompleteOrderHandler: PATCH /api/workorders/{id}/complete
func (ctrl *WorkOrderController) CompleteOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	// ... (Error handling untuk orderID) ...

	var req models.CompleteWorkOrder
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if req.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'completed'"})
		return
	}

	err = ctrl.Repo.CompleteOrder(orderID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order completed successfully", "order_id": orderID})
}

// DeleteOrderHandler: DELETE /api/workorders/{id}
func (ctrl *WorkOrderController) DeleteOrderHandler(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	// ... (Error handling untuk orderID) ...

	err = ctrl.Repo.DeleteOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully", "order_id": orderID})
}

func GetWorkOrders(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, workOrders)
}

func CreateWorkOrder(c *gin.Context) {
	var w models.WorkOrder
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	w.ID = nextID
	nextID++
	if w.Status == "" {
		w.Status = "pending"
	}
	if w.Time == "" {
		w.Time = time.Now().Format("15:04")
	}
	workOrders = append(workOrders, w)
	mu.Unlock()

	c.JSON(http.StatusCreated, w)
}

func UpdateWorkOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var payload models.WorkOrder
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i := range workOrders {
		if workOrders[i].ID == id {
			if payload.Executors != nil {
				workOrders[i].Executors = payload.Executors
			}
			if payload.Status != "" {
				workOrders[i].Status = payload.Status
			}
			if payload.SafetyChecklist != nil {
				workOrders[i].SafetyChecklist = payload.SafetyChecklist
			}
			if payload.WorkingHours != nil {
				workOrders[i].WorkingHours = payload.WorkingHours
			}
			c.JSON(http.StatusOK, workOrders[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func DeleteWorkOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i := range workOrders {
		if workOrders[i].ID == id {
			workOrders = append(workOrders[:i], workOrders[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"deleted": id})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}

func GetSummary(c *gin.Context) {
	// Initialize repository
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}
	repo := repository.NewWorkOrderRepository(db)

	// Fetch all tasks from the database
	allTasks, err := repo.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks for summary"})
		return
	}

	// Calculate summary counts
	total := len(allTasks)
	pending := 0
	progress := 0
	completed := 0

	for _, task := range allTasks {
		switch task.Status {
		case "pending":
			pending++
		case "progress":
			progress++
		case "completed":
			completed++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"pending":   pending,
		"progress":  progress,
		"completed": completed,
	})
}

func GetKaizen(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	total := len(workOrders)
	pending := 0
	progress := 0
	for _, o := range workOrders {
		if o.Status == "pending" {
			pending++
		} else if o.Status == "progress" {
			progress++
		}
	}
	completed := total - pending - progress
	completionRate := 0
	if total > 0 {
		completionRate = (completed * 100) / total
	}
	rating := "Perlu perbaikan"
	suggestion := "Investigasi hambatan kerja."
	if completionRate >= 80 {
		rating = "Sempurna"
		suggestion = "Pertahankan kinerja tim."
	} else if completionRate >= 60 {
		rating = "Baik"
		suggestion = "Performa baik, namun masih ada ruang untuk peningkatan."
	} else if completionRate >= 40 {
		rating = "Cukup"
		suggestion = "Perbaikan perlu dipertimbangkan."
	}
	c.JSON(http.StatusOK, gin.H{
		"total":          total,
		"pending":        pending,
		"progress":       progress,
		"completed":      completed,
		"completionRate": completionRate,
		"rating":         rating,
		"suggestion":     suggestion,
	})
}

func GetMembersHandler(c *gin.Context) {
	members, err := repository.GetAllMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve members"})
		return
	}

	c.JSON(http.StatusOK, members)
}
