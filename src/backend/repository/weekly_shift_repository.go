package repository

import (
	"errors"
	"fmt"

	"teamitmivhs/work-order-backend/config"
)

var ErrInvalidWeeklyShiftMembers = errors.New("weekly shift members must be three active staff")

type WeeklyShiftEntry struct {
	DayOfWeek int    `json:"dayOfWeek"`
	Position  int    `json:"position"`
	MemberID  int    `json:"memberId"`
	Name      string `json:"name"`
	Division  string `json:"division"`
	Avatar    string `json:"avatar"`
}

func GetWeeklyShiftSchedule() ([]WeeklyShiftEntry, error) {
	rows, err := config.DB.Query(`
		SELECT s.DayOfWeek, s.Position, m.ID, m.Name, COALESCE(m.Division, ''), m.Avatar
		FROM weekly_shift_schedule s
		JOIN members m ON m.ID = s.MemberID
		WHERE m.AccountStatus = 'active'
		  AND m.MembershipStatus = 'active'
		  AND m.CanHandleWorkOrder = 1
		  AND LOWER(COALESCE(m.Role, '')) NOT IN ('guru', 'guest')
		ORDER BY s.DayOfWeek, s.Position
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]WeeklyShiftEntry, 0, 15)
	for rows.Next() {
		var entry WeeklyShiftEntry
		if err := rows.Scan(
			&entry.DayOfWeek,
			&entry.Position,
			&entry.MemberID,
			&entry.Name,
			&entry.Division,
			&entry.Avatar,
		); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func ReplaceWeeklyShift(day int, memberIDs []int, updatedBy int) ([]int, error) {
	if day < 1 || day > 5 || len(memberIDs) != 3 {
		return nil, ErrInvalidWeeklyShiftMembers
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(`
		SELECT ID
		FROM members
		WHERE ID IN (?, ?, ?)
		  AND AccountStatus = 'active'
		  AND MembershipStatus = 'active'
		  AND CanHandleWorkOrder = 1
		  AND LOWER(COALESCE(Role, '')) NOT IN ('guru', 'guest')
		FOR UPDATE
	`, memberIDs[0], memberIDs[1], memberIDs[2])
	if err != nil {
		return nil, err
	}
	validMembers := make(map[int]bool, 3)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		validMembers[id] = true
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if len(validMembers) != 3 {
		return nil, ErrInvalidWeeklyShiftMembers
	}

	oldRows, err := tx.Query(
		"SELECT MemberID FROM weekly_shift_schedule WHERE DayOfWeek = ? FOR UPDATE",
		day,
	)
	if err != nil {
		return nil, err
	}
	oldMembers := make(map[int]bool, 3)
	for oldRows.Next() {
		var id int
		if err := oldRows.Scan(&id); err != nil {
			oldRows.Close()
			return nil, err
		}
		oldMembers[id] = true
	}
	if err := oldRows.Err(); err != nil {
		oldRows.Close()
		return nil, err
	}
	if err := oldRows.Close(); err != nil {
		return nil, err
	}

	if _, err := tx.Exec("DELETE FROM weekly_shift_schedule WHERE DayOfWeek = ?", day); err != nil {
		return nil, err
	}
	for position, memberID := range memberIDs {
		if _, err := tx.Exec(`
			INSERT INTO weekly_shift_schedule (DayOfWeek, Position, MemberID, UpdatedBy)
			VALUES (?, ?, ?, ?)
		`, day, position+1, memberID, updatedBy); err != nil {
			return nil, fmt.Errorf("save weekly shift member: %w", err)
		}
	}
	if err := enqueueEvent(tx, "member.shift_schedule_updated", "shift_schedule", int64(day)); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	newlyAssigned := make([]int, 0, 3)
	for _, memberID := range memberIDs {
		if !oldMembers[memberID] {
			newlyAssigned = append(newlyAssigned, memberID)
		}
	}
	return newlyAssigned, nil
}
