package repository

import (
	"database/sql"
	"fmt"
)

func enqueueEvent(tx *sql.Tx, eventType, aggregateType string, aggregateID int64) error {
	_, err := tx.Exec(
		"INSERT INTO event_outbox (EventType, AggregateType, AggregateID) VALUES (?, ?, ?)",
		eventType, aggregateType, aggregateID,
	)
	if err != nil {
		return fmt.Errorf("enqueue %s event: %w", eventType, err)
	}
	return nil
}
