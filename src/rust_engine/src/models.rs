use serde::{Deserialize, Serialize};

//request buat time tracking
//di call pas:
//- workorder di "take"
//- atau pas mulai pengerjaan 

#[derive(Deserialize, Debug)]
pub struct StartTimeRequest {
    pub work_order_id: u64,
    pub executor_id: u64,
}

#[derive(Serialize, Debug)]
pub struct StartTimeResponse {
    pub work_order_id: u64,
    pub started_at: i64,
}

/// Request buat stop time tracking
/// Di call pas:
/// - workorder di "complete"
#[derive(Deserialize, Debug)]
pub struct StopTimeRequest {
    pub work_order_id: u64,
}

///Response stop timer
#[derive(Serialize, Debug)]
pub struct StopTimeResponse {
    pub work_order_id: u64,
    pub stopped_at: i64,
    pub started_at: i64,
    pub duration_seconds: i64,
}

/// Status timer (buat polling)
#[derive(Serialize, Debug)]
pub struct TimerStatusResponse {
    pub work_order_id: u64,
    pub is_running: bool,
    pub started_at: Option<i64>,
    pub elapsed_seconds: Option<i64>,
}

///Error response
#[derive(Serialize, Debug)]
pub struct ErrorResponse {
    pub message: String,
}