use sqlx::Row;
use std::{sync::atomic::Ordering, time::Duration};
use tokio::time;

use crate::{models::OutboxEvent, state::AppState};

pub async fn run(state: AppState) {
    let mut poll = time::interval(Duration::from_millis(500));
    let mut cleanup = time::interval(Duration::from_secs(60 * 60));
    cleanup.tick().await;

    loop {
        tokio::select! {
            _ = poll.tick() => {
                if let Err(error) = publish_pending(&state).await {
                    tracing::error!(%error, "Failed to consume event outbox");
                }
            }
            _ = cleanup.tick() => {
                if let Err(error) = cleanup_processed(&state).await {
                    tracing::warn!(%error, "Failed to clean processed events");
                }
            }
        }
    }
}

async fn publish_pending(state: &AppState) -> Result<(), sqlx::Error> {
    // ponytail: one engine replica; add row claiming only when horizontal scaling is needed.
    let rows = sqlx::query(
        "SELECT ID, EventType, AggregateType, AggregateID \
         FROM event_outbox WHERE ProcessedAt IS NULL ORDER BY ID LIMIT 100",
    )
    .fetch_all(&state.pool)
    .await?;

    for row in rows {
        let event = OutboxEvent {
            id: row.try_get("ID")?,
            event_type: row.try_get("EventType")?,
            aggregate_type: row.try_get("AggregateType")?,
            aggregate_id: row.try_get("AggregateID")?,
        };

        // No subscribers is fine: every page loads current state before opening SSE.
        let _ = state.events.send(event.clone());
        let result = sqlx::query(
            "UPDATE event_outbox SET ProcessedAt = NOW(3) WHERE ID = ? AND ProcessedAt IS NULL",
        )
        .bind(event.id)
        .execute(&state.pool)
        .await?;

        if result.rows_affected() == 1 {
            state.processed_events.fetch_add(1, Ordering::Relaxed);
        }
    }

    Ok(())
}

async fn cleanup_processed(state: &AppState) -> Result<(), sqlx::Error> {
    sqlx::query(
        "DELETE FROM event_outbox \
         WHERE ProcessedAt IS NOT NULL AND ProcessedAt < NOW() - INTERVAL 7 DAY \
         LIMIT 1000",
    )
    .execute(&state.pool)
    .await?;
    Ok(())
}
