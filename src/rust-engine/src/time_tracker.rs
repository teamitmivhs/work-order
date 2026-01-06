use crate::state::{AppState, RunningTimer};
use std::time::{SystemTime, UNIX_EPOCH};

pub struct TimeTracker {
    pub state: AppState,
}

impl TimeTracker {
    pub fn new(state: AppState) -> Self {
        Self { state }
    }

    /// Get current unix timestamp in seconds
    fn now() -> i64 {
        SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default()
            .as_secs() as i64
    }

    /// Start a timer for a work order
    /// Returns error if timer already running for this order ID
    pub fn start(&self, work_order_id: u64, executor_id: u64) -> Result<i64, String> {
        if work_order_id == 0 {
            return Err("Work order ID cannot be 0".to_string());
        }

        let mut timers = self
            .state
            .timers
            .lock()
            .map_err(|e| format!("Failed to acquire lock: {}", e))?;

        if timers.contains_key(&work_order_id) {
            return Err(format!("Timer already running for work order {}", work_order_id));
        }

        let started_at = Self::now();

        timers.insert(
            work_order_id,
            RunningTimer {
                work_order_id,
                executor_id,
                started_at,
            },
        );

        tracing::info!(
            "Timer started for work order {} by executor {}",
            work_order_id,
            executor_id
        );

        Ok(started_at)
    }

    /// Stop a timer for a work order
    /// Returns (start_time, duration_in_seconds)
    pub fn stop(&self, work_order_id: u64) -> Result<(i64, i64), String> {
        if work_order_id == 0 {
            return Err("Work order ID cannot be 0".to_string());
        }

        let mut timers = self
            .state
            .timers
            .lock()
            .map_err(|e| format!("Failed to acquire lock: {}", e))?;

        let timer = timers
            .remove(&work_order_id)
            .ok_or_else(|| format!("Timer not found for work order {}", work_order_id))?;

        let stopped_at = Self::now();
        let duration = stopped_at - timer.started_at;

        tracing::info!(
            "Timer stopped for work order {}, duration: {} seconds",
            work_order_id,
            duration
        );

        Ok((timer.started_at, duration))
    }

    /// Get current status of a timer without stopping it
    /// Returns Option<(start_time, elapsed_seconds)>
    pub fn status(&self, work_order_id: u64) -> Result<Option<(i64, i64)>, String> {
        if work_order_id == 0 {
            return Err("Work order ID cannot be 0".to_string());
        }

        let timers = self
            .state
            .timers
            .lock()
            .map_err(|e| format!("Failed to acquire lock: {}", e))?;

        Ok(timers.get(&work_order_id).map(|t| {
            let elapsed = Self::now() - t.started_at;
            (t.started_at, elapsed)
        }))
    }

    /// Get count of active timers
    pub fn active_count(&self) -> Result<usize, String> {
        let timers = self
            .state
            .timers
            .lock()
            .map_err(|e| format!("Failed to acquire lock: {}", e))?;

        Ok(timers.len())
    }
}
