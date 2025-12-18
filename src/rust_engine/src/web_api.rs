use axum::{
    extract::{Path, State},
    http::StatusCode,
    Json,
};
use std::sync::Arc;

use crate::{
    models::*,
    tracker::TimeTracker,
};

pub async fn start_timer(
    State(tracker): State<Arc<TimeTracker>>,
    Json(req): Json<StartTimeRequest>,
) -> Result<Json<StartTimeResponse>, (StatusCode, Json<ErrorResponse>)> {
    tracker
        .start(req.work_order_id, req.executor_id)
        .map(|started_at| {
            Json(StartTimeResponse {
                work_order_id: req.work_order_id,
                started_at,
            })
        })
        .map_err(|e| {
            (
                StatusCode::BAD_REQUEST,
                Json(ErrorResponse { message: e }),
            )
        })
}

pub async fn stop_timer(
    State(tracker): State<Arc<TimeTracker>>,
    Json(req): Json<StopTimeRequest>,
) -> Result<Json<StopTimeResponse>, (StatusCode, Json<ErrorResponse>)> {
    tracker
        .stop(req.work_order_id)
        .map(|(started_at, duration)| {
            let stopped_at = started_at + duration;

            Json(StopTimeResponse {
                work_order_id: req.work_order_id,
                started_at,
                stopped_at,
                duration_seconds: duration,
            })
        })
        .map_err(|e| {
            (
                StatusCode::BAD_REQUEST,
                Json(ErrorResponse { message: e }),
            )
        })
}

pub async fn timer_status(
    State(tracker): State<Arc<TimeTracker>>,
    Path(work_order_id): Path<u64>,
) -> Json<TimerStatusResponse> {
    match tracker.status(work_order_id) {
        Some((started_at, elapsed)) => Json(TimerStatusResponse {
            work_order_id,
            is_running: true,
            started_at: Some(started_at),
            elapsed_seconds: Some(elapsed),
        }),
        None => Json(TimerStatusResponse {
            work_order_id,
            is_running: false,
            started_at: None,
            elapsed_seconds: None,
        }),
    }
}
