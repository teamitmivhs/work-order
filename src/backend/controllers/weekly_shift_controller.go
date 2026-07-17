package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"teamitmivhs/work-order-backend/middleware"
	"teamitmivhs/work-order-backend/repository"
	"teamitmivhs/work-order-backend/services"
	"teamitmivhs/work-order-backend/utils"

	"github.com/gin-gonic/gin"
)

const weeklyShiftSize = 3

var shiftDayNames = map[int]string{
	1: "Senin",
	2: "Selasa",
	3: "Rabu",
	4: "Kamis",
	5: "Jumat",
}

type weeklyShiftRequest struct {
	MemberIDs []int `json:"memberIds" binding:"required"`
}

func normalizeWeeklyShiftMembers(memberIDs []int) ([]int, bool) {
	if len(memberIDs) != weeklyShiftSize {
		return nil, false
	}
	seen := make(map[int]bool, weeklyShiftSize)
	for _, id := range memberIDs {
		if id <= 0 || seen[id] {
			return nil, false
		}
		seen[id] = true
	}
	return memberIDs, true
}

func shiftDayName(day int) string {
	return shiftDayNames[day]
}

func canManageShift(c *gin.Context, memberRepo repository.MemberRepository) bool {
	role, _ := middleware.GetUserRoleFromContext(c)
	if utils.IsAdminRole(role) {
		return true
	}
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok || userID == 0 {
		return false
	}
	member, err := memberRepo.GetMemberByID(userID)
	return err == nil && strings.EqualFold(strings.TrimSpace(member.Division), "Data Analyst")
}

func GetWeeklyShiftScheduleHandler(c *gin.Context) {
	entries, err := repository.GetWeeklyShiftSchedule()
	if err != nil {
		utils.InternalServerError(c, "Failed to retrieve weekly shift schedule", err)
		return
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
	}
	now := time.Now().In(location)
	today := int(now.Weekday())
	if today < 1 || today > 5 {
		today = 0
	}
	utils.RespondSuccess(c, http.StatusOK, gin.H{
		"todayDay":      today,
		"todayName":     shiftDayName(today),
		"todayDate":     now.Format("2006-01-02"),
		"entries":       entries,
		"canManage":     canManageShift(c, repository.NewMemberRepository()),
		"membersPerDay": weeklyShiftSize,
	})
}

func UpdateWeeklyShiftScheduleHandler(c *gin.Context) {
	day, err := strconv.Atoi(c.Param("day"))
	if err != nil || shiftDayName(day) == "" {
		utils.BadRequest(c, "Shift day must be Monday through Friday")
		return
	}
	memberRepo := repository.NewMemberRepository()
	if !canManageShift(c, memberRepo) {
		utils.Forbidden(c, "Only admins, teachers, or Data Analyst can manage shifts")
		return
	}

	var req weeklyShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request payload", err.Error())
		return
	}
	memberIDs, valid := normalizeWeeklyShiftMembers(req.MemberIDs)
	if !valid {
		utils.BadRequest(c, "Each weekday shift must contain exactly 3 different members")
		return
	}
	updatedBy, _ := middleware.GetUserIDFromContext(c)
	newlyAssigned, err := repository.ReplaceWeeklyShift(day, memberIDs, updatedBy)
	if errors.Is(err, repository.ErrInvalidWeeklyShiftMembers) {
		utils.BadRequest(c, "Shift members must be active staff and cannot be teachers")
		return
	}
	if err != nil {
		utils.InternalServerError(c, "Failed to update weekly shift schedule", err)
		return
	}

	if len(newlyAssigned) > 0 {
		dayName := shiftDayName(day)
		go services.NotifyWorkOrderUsers(
			newlyAssigned,
			"Jadwal Shift • "+dayName,
			fmt.Sprintf("Kamu dijadwalkan shift hari %s. Buka untuk melihat tim shift.", dayName),
			0,
			map[string]any{"url": "/shift", "tag": fmt.Sprintf("weekly-shift-%d", day)},
		)
	}

	utils.RespondWithMessage(c, http.StatusOK, "Weekly shift schedule updated", gin.H{
		"dayOfWeek": day,
		"memberIds": memberIDs,
	})
}
