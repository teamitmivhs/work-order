use axum::{
    extract::{State, Path},
    Json,
    http::StatusCode,
};
use std::sync::Arc;

use crate::time_tracker::{
    models::*,
    tracker::TimeTracker,
};

/// POST /timer/start
pub async fn start_timer(
    State(tracker): State<Arc<TimeTracker>>,
    Json(req): Json<StartTimeRequest>,
) -> Result<Json<StartTimeResponse>, (StatusCode, Json<ErrorResponse>)> {
    match tracker.start(req.workorder_id, req.executor_id) {
        Ok(started_at) => Ok(Json(StartTimeResponse {
            work_order_id: req.work_order_id,
            started_at,
        })),
        Err(e) => Err((
            StatusCode::BAD_REQUEST,
            Json(ErrorResponse { message: e }),
        )),
    }
}

/// POST /timer/stop
pub async fn stop_timer(
    State(tracker): State<Arc<TimeTracker>>,
    Json(req): Json<StopTimeRequest>,
) -> Result<Json<StopTimeResponse>, (StatusCode, Json<ErrorResponse>)> {
    match tracker.stop(req.workorder_id) {
        Ok((started_at, duration)) => {
            let stopped_at = started_at + duration;

            Ok(Json(StopTimeResponse {
                work_order_id: req.workorder_id,
                started_at,
                stopped_at,
                duration_seconds: duration,
            }))
        }
        Err(e) => Err((
            StatusCode::BAD_REQUEST,
            Json(ErrorResponse { message: e }),
        )),
    }
}

/// GET /timer/:workorder_id
pub async fn timer_status(
    State(tracker): State<Arc<TimeTracker>>,
    Path(workorder_id): Path<u64>,
) -> Json<TimerStatusResponse> {
    match tracker.status(workorder_id) {
        Some((started_at, elapsed)) => Json(TimerStatusResponse {
            work_order_id,
            is_running: true,
            started_at: Some(started_at),
            elapsed_seconds: Some(elapsed),
        }),
        None => Json(TimerStatusResponse {
            workorder_id,
            is_running: false,
            started_at: None,
            elapsed_seconds: None,
        }),
    }
}
