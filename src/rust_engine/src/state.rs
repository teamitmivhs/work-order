use axum::{
    extract::State,
    Json,
    http::StatusCode,
};

use crate::time_tracker::{
    models::*,
    tracker::TimeTracker,
};
use std::sync::Arc;

pub async fn start_timer(
    State(tracker): State<Arc<TimeTracker>>,
    Json(req): Json<StartTimeRequest>,
) -> Result<Json<StartTimeResponse>, (StatusCode, Json<ErrorResponse>)> {
    match tracker.start(req.work_order_id, req.executor_id) {
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

pub async fn stop_timer(
    State(tracker): State<Arc<TimeTracker>>,
    Json(req): Json<StopTimeRequest>,
) -> Result<Json<StopTimeResponse>, (StatusCode, Json<ErrorResponse>)> {
    match tracker.stop(req.work_order_id) {
        Ok(result) => Ok(Json(StopTimeResponse {
            work_order_id: req.work_order_id,
            started_at: result.started_at,
            stopped_at: result.stopped_at,
            duration_seconds: result.duration_seconds,
        })),
        Err(e) => Err((
            StatusCode::BAD_REQUEST,
            Json(ErrorResponse { message: e }),
        )),
    }
}

pub async fn timer_status(
    State(tracker): State<Arc<TimeTracker>>,
    work_order_id: u64,
) -> Json<TimerStatusResponse> {
    Json(tracker.status(work_order_id))
}
