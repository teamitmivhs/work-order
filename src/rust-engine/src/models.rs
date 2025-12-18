use serde::{Deserialize, Serialize};

/// Request buat start time tracking
/// Dipanggil saat:
/// - work order di "TAKE"
/// - mulai pengerjaan
#[derive(Debug, Deserialize)]
pub struct StartTimeRequest {
    pub work_order_id: u64,
    pub executor_id: u64,
}

#[derive(Debug, Serialize)]
pub struct StartTimeResponse {
    pub work_order_id: u64,
    pub started_at: i64,
}

/// Request buat stop time tracking
/// Dipanggil saat:
/// - work order di "COMPLETE"
#[derive(Debug, Deserialize)]
pub struct StopTimeRequest {
    pub work_order_id: u64,
}

#[derive(Debug, Serialize)]
pub struct StopTimeResponse {
    pub work_order_id: u64,
    pub started_at: i64,
    pub stopped_at: i64,
    pub duration_seconds: i64,
}

/// Status timer (buat polling frontend)
#[derive(Debug, Serialize)]
pub struct TimerStatusResponse {
    pub work_order_id: u64,
    pub is_running: bool,
    pub started_at: Option<i64>,
    pub elapsed_seconds: Option<i64>,
}

/// Error response standar API
#[derive(Debug, Serialize)]
pub struct ErrorResponse {
    pub message: String,
}
