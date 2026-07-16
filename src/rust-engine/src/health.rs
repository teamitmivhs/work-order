use axum::{extract::State, http::StatusCode, Json};
use serde_json::{json, Value};
use std::sync::atomic::Ordering;

use crate::state::AppState;

pub async fn health_check(
    State(state): State<AppState>,
) -> Result<Json<Value>, (StatusCode, Json<Value>)> {
    let pending =
        sqlx::query_scalar::<_, i64>("SELECT COUNT(*) FROM event_outbox WHERE ProcessedAt IS NULL")
            .fetch_one(&state.pool)
            .await;

    match pending {
        Ok(pending_events) => Ok(Json(json!({
            "status": "healthy",
            "service": "event-engine",
            "pending_events": pending_events,
            "processed_events": state.processed_events.load(Ordering::Relaxed),
            "subscribers": state.events.receiver_count()
        }))),
        Err(error) => {
            tracing::error!(%error, "Event engine database health check failed");
            Err((
                StatusCode::SERVICE_UNAVAILABLE,
                Json(json!({
                    "status": "unhealthy",
                    "service": "event-engine",
                    "error": error.to_string()
                })),
            ))
        }
    }
}
