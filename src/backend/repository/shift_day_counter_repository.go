package repository

import (
	"database/sql"
	"time"

	"teamitmivhs/work-order-backend/config"
)

type ShiftDayCounterSnapshot struct {
	LastDate               string `json:"lastDate"`
	LastDay                int    `json:"lastDay"`
	LastMonth              int    `json:"lastMonth"`
	LastYear               int    `json:"lastYear"`
	RolloverCount          int    `json:"rolloverCount"`
	MovedToStandby         int64  `json:"movedToStandby"`
	ClearedScheduleEntries int64  `json:"clearedScheduleEntries"`
}

func crossedSunday(last, current time.Time) bool {
	last = time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, last.Location())
	current = time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
	if !last.Before(current) {
		return false
	}
	daysUntilSunday := (7 - int(last.Weekday())) % 7
	if daysUntilSunday == 0 {
		daysUntilSunday = 7
	}
	return !last.AddDate(0, 0, daysUntilSunday).After(current)
}

func RunShiftDayRollover() (*ShiftDayCounterSnapshot, error) {
	if config.DB == nil {
		return nil, sql.ErrConnDone
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	today := now.Format("2006-01-02")

	tx, err := config.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	insertResult, err := tx.Exec(`
		INSERT INTO shift_day_counter (ID, LastDate, LastDay, LastMonth, LastYear, RolloverCount)
		VALUES (1, ?, ?, ?, ?, 0)
		ON DUPLICATE KEY UPDATE ID = ID
	`, today, now.Day(), int(now.Month()), now.Year())
	if err != nil {
		return nil, err
	}
	inserted, _ := insertResult.RowsAffected()

	var lastDate time.Time
	var lastDay, lastMonth, lastYear, rolloverCount int
	if err := tx.QueryRow(`
		SELECT LastDate, LastDay, LastMonth, LastYear, RolloverCount
		FROM shift_day_counter
		WHERE ID = 1
		FOR UPDATE
	`).Scan(&lastDate, &lastDay, &lastMonth, &lastYear, &rolloverCount); err != nil {
		return nil, err
	}

	moved := int64(0)
	clearedScheduleEntries := int64(0)
	shouldClearSchedule := (inserted > 0 && now.Weekday() == time.Sunday) ||
		crossedSunday(lastDate.In(loc), now)
	if lastDate.In(loc).Format("2006-01-02") != today {
		result, err := tx.Exec("UPDATE members SET Status = 'standby' WHERE Status = 'nextshift'")
		if err != nil {
			return nil, err
		}
		moved, _ = result.RowsAffected()

		rolloverCount++
		lastDay = now.Day()
		lastMonth = int(now.Month())
		lastYear = now.Year()
		lastDate = now
		if _, err := tx.Exec(`
			UPDATE shift_day_counter
			SET LastDate = ?, LastDay = ?, LastMonth = ?, LastYear = ?, RolloverCount = ?, UpdatedAt = NOW()
			WHERE ID = 1
		`, today, lastDay, lastMonth, lastYear, rolloverCount); err != nil {
			return nil, err
		}
		if moved > 0 {
			if err := enqueueEvent(tx, "member.shift_rolled_over", "member", 0); err != nil {
				return nil, err
			}
		}
	}
	if shouldClearSchedule {
		result, err := tx.Exec("DELETE FROM weekly_shift_schedule")
		if err != nil {
			return nil, err
		}
		clearedScheduleEntries, _ = result.RowsAffected()
		if clearedScheduleEntries > 0 {
			if err := enqueueEvent(tx, "member.shift_schedule_updated", "shift_schedule", 0); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &ShiftDayCounterSnapshot{
		LastDate:               lastDate.In(loc).Format("2006-01-02"),
		LastDay:                lastDay,
		LastMonth:              lastMonth,
		LastYear:               lastYear,
		RolloverCount:          rolloverCount,
		MovedToStandby:         moved,
		ClearedScheduleEntries: clearedScheduleEntries,
	}, nil
}
