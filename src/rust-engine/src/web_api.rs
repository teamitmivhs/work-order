use axum::{
    extract::{Path, State},
    http::{HeaderMap, StatusCode},
    Json,
};
use std::{env, sync::Arc};

use crate::{models::*, time_tracker::TimeTracker};

/// Helper: validasi internal API key dari header X-Internal-Key
/// FIX: semua endpoint dilindungi shared secret agar tidak bisa diakses sembarangan
fn validate_internal_key(headers: &HeaderMap) -> bool {
    let expected = env::var("INTERNAL_API_KEY").unwrap_or_else(|_| "changeme-internal-key".to_string());
    headers
        .get("X-Internal-Key")
        .and_then(|v| v.to_str().ok())
        .map(|v| v == expected)
        .unwrap_or(false)
}

fn unauthorized() -> (StatusCode, Json<ErrorResponse>) {
    (
        StatusCode::UNAUTHORIZED,
        Json(ErrorResponse {
            message: "Missing or invalid internal API key".to_string(),
        }),
    )
}

/// POST /timer/start — mulai timer untuk work order
pub async fn start_timer(
    State(tracker): State<Arc<TimeTracker>>,
    headers: HeaderMap,
    Json(req): Json<StartTimeRequest>,
) -> Result<Json<StartTimeResponse>, (StatusCode, Json<ErrorResponse>)> {
    // FIX: validasi internal key
    if !validate_internal_key(&headers) {
        return Err(unauthorized());
    }

    tracker
        .start(req.work_order_id, req.executor_id)
        .await
        .map(|started_at| {
            Json(StartTimeResponse {
                work_order_id: req.work_order_id,
                started_at,
            })
        })
        .map_err(|e| {
            tracing::warn!(error = %e, "Failed to start timer");
            (StatusCode::BAD_REQUEST, Json(ErrorResponse { message: e }))
        })
}

/// POST /timer/stop — hentikan timer dan dapatkan durasi
pub async fn stop_timer(
    State(tracker): State<Arc<TimeTracker>>,
    headers: HeaderMap,
    Json(req): Json<StopTimeRequest>,
) -> Result<Json<StopTimeResponse>, (StatusCode, Json<ErrorResponse>)> {
    // FIX: validasi internal key
    if !validate_internal_key(&headers) {
        return Err(unauthorized());
    }

    tracker
        .stop(req.work_order_id)
        .await
        // FIX: stop() sekarang return (started_at, stopped_at, duration)
        // stopped_at adalah actual timestamp, bukan rekonstruksi started_at + duration
        .map(|(started_at, stopped_at, duration)| {
            Json(StopTimeResponse {
                work_order_id: req.work_order_id,
                started_at,
                stopped_at, // actual timestamp saat stop
                duration_seconds: duration,
            })
        })
        .map_err(|e| {
            tracing::warn!(error = %e, "Failed to stop timer");
            (StatusCode::BAD_REQUEST, Json(ErrorResponse { message: e }))
        })
}

/// GET /timer/:work_order_id — status satu timer
pub async fn timer_status(
    State(tracker): State<Arc<TimeTracker>>,
    headers: HeaderMap,
    Path(work_order_id): Path<u64>,
) -> Result<Json<TimerStatusResponse>, (StatusCode, Json<ErrorResponse>)> {
    // FIX: validasi internal key
    if !validate_internal_key(&headers) {
        return Err(unauthorized());
    }

    match tracker.status(work_order_id).await {
        Ok(Some((started_at, elapsed))) => Ok(Json(TimerStatusResponse {
            work_order_id,
            is_running: true,
            started_at: Some(started_at),
            elapsed_seconds: Some(elapsed),
        })),
        Ok(None) => Ok(Json(TimerStatusResponse {
            work_order_id,
            is_running: false,
            started_at: None,
            elapsed_seconds: None,
        })),
        Err(e) => {
            tracing::warn!(error = %e, "Failed to get timer status");
            Err((
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(ErrorResponse { message: e }),
            ))
        }
    }
}

/// GET /timers — list semua timer aktif
/// FIX: endpoint baru untuk monitoring dan debugging
pub async fn list_timers(
    State(tracker): State<Arc<TimeTracker>>,
    headers: HeaderMap,
) -> Result<Json<Vec<ActiveTimerInfo>>, (StatusCode, Json<ErrorResponse>)> {
    if !validate_internal_key(&headers) {
        return Err(unauthorized());
    }

    tracker
        .list_active()
        .await
        .map(|list| {
            Json(
                list.into_iter()
                    .map(|(work_order_id, executor_id, started_at, elapsed)| ActiveTimerInfo {
                        work_order_id,
                        executor_id,
                        started_at,
                        elapsed_seconds: elapsed,
                    })
                    .collect(),
            )
        })
        .map_err(|e| {
            tracing::warn!(error = %e, "Failed to list timers");
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(ErrorResponse { message: e }),
            )
        })
}