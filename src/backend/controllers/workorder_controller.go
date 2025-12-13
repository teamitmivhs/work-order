package controllers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"teamitmivhs/work-order-backend/models"
	"teamitmivhs/work-order-backend/repository"

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
	completed := 0
	executorMembersMap := make(map[int]models.Member) // To store unique executor members

	allMembers, err := repository.GetAllMembers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve members for summary"})
		return
	}

	memberIDToMember := make(map[int]models.Member)
	for _, member := range allMembers {
		memberIDToMember[member.ID] = member
	}

	for _, o := range workOrders {
		if o.Status == "pending" {
			pending++
		} else if o.Status == "progress" {
			progress++
		} else if o.Status == "completed" {
			completed++
			for _, executorID := range o.Executors {
				if member, ok := memberIDToMember[executorID]; ok {
					executorMembersMap[executorID] = member
				}
			}
		}
	}

	var uniqueExecutorMembers []models.Member
	for _, member := range executorMembersMap {
		uniqueExecutorMembers = append(uniqueExecutorMembers, member)
	}

	c.JSON(http.StatusOK, gin.H{"total": total, "pending": pending, "progress": progress, "completed": completed, "executors": uniqueExecutorMembers})
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
