use crate::state::{AppState, RunningTimer};
use std::time::{SystemTime, UNIX_EPOCH};

pub struct TimeTracker {
    pub state: AppState,
}

impl TimeTracker {
    pub fn new(state: AppState) -> Self {
        Self { state }
    }

    fn now() -> i64 {
        SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .expect("time went backwards")
            .as_secs() as i64
    }

    /// Start timer (dipanggil saat TAKE work order)
    pub fn start(
        &self,
        work_order_id: u64,
        executor_id: u64,
    ) -> Result<i64, String> {
        let mut timers = self.state.timers.lock().unwrap();

        if timers.contains_key(&workorder_id) {
            return Err("timer already running".into());
        }

        let started_at = Self::now();

        timers.insert(
            work_order_id
            RunningTimer {
                work_order_id,
                executor_id,
                started_at,
            },
        );

        Ok(started_at)
    }

    /// Stop timer (dipanggil saat COMPLETE)
    pub fn stop(&self, workorder_id: u64) -> Result<(i64, i64), String> {
        let mut timers = self.state.timers.lock().unwrap();

        let timer = timers
            .remove(&workorder_id)
            .ok_or("timer not found")?;

        let stopped_at = Self::now();
        let duration = stopped_at - timer.started_at;

        Ok((timer.started_at, duration))
    }

    /// Status timer (polling frontend)
    pub fn status(&self, workorder_id: u64) -> Option<(i64, i64)> {
        let timers = self.state.timers.lock().unwrap();

        timers.get(&workorder_id).map(|timer| {
            let elapsed = Self::now() - timer.started_at;
            (timer.started_at, elapsed)
        })
    }
}
