use std::{
    collections::HashMap,
    sync::{Arc, Mutex},
};

#[derive(Debug, Clone)]
pub struct RunningTimer {
    pub work_order_id: u64,
    pub executor_id: u64,
    pub started_at: i64,
}

#[derive(Clone)]
pub struct AppState {
    pub timers: Arc<Mutex<HashMap<u64, RunningTimer>>>,
}

impl AppState {
    pub fn new() -> Self {
        Self {
            timers: Arc::new(Mutex::new(HashMap::new())),
        }
    }
}
