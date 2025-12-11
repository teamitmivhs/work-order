package controllers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"teamitmivhs/work-order-backend/models"

	"github.com/gin-gonic/gin"
)

var (
	mu         sync.Mutex
	workOrders = []models.WorkOrder{}
	nextID     = 1
)

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
			if payload.WorkingHours != "" {
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
	c.JSON(http.StatusOK, gin.H{"total": total, "pending": pending, "progress": progress, "completed": completed})
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
	rating := "Needs Improvement"
	suggestion := "Investigate bottlenecks."
	if completionRate >= 80 {
		rating = "Excellent"
		suggestion = "Keep up the good work."
	} else if completionRate >= 60 {
		rating = "Good"
		suggestion = "Focus on reducing pending items."
	} else if completionRate >= 40 {
		rating = "Fair"
		suggestion = "Consider process improvements."
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
