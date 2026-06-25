use axum::{
    extract::State,
    http::StatusCode,
    Json,
};
use serde_json::{json, Value};
use std::sync::Arc;

use crate::time_tracker::TimeTracker;

/// GET /health — health check dengan info timer aktif
/// FIX: active_count() kini diekspos lewat sini sehingga tidak mubazir
pub async fn health_check(
    State(tracker): State<Arc<TimeTracker>>,
) -> Result<Json<Value>, (StatusCode, Json<Value>)> {
    match tracker.active_count().await {
        Ok(count) => Ok(Json(json!({
            "status": "healthy",
            "service": "time-tracker",
            "active_timers": count 
        }))),
        Err(e) => {
            tracing::error!(error = %e, "Health check failed to get timer count");
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(json!({
                    "status": "unhealthy",
                    "service": "time-tracker",
                    "error": e
                })),
            ))
        }
    }
}