use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::time::{SystemTime, UNIX_EPOCH};

#[derive(Debug, Clone)]
pub struct Timer {
    pub work_order_id: u64,
    pub executor_id: u64,
    pub started_at: i64,
    pub stopped_at: Option<i64>,
}

#[derive(Clone)]
pub struct TimeTracker {
    timers: Arc<Mutex<HashMap<u64, Timer>>>,
}

impl TimeTracker {
    pub fn new() -> Self {
        Self {
            timers: Arc::new(Mutex::new(HashMap::new())),
        }
    }

    fn now() -> i64 {
        SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .expect("time went backwards")
            .as_secs() as i64
    }

    /// Start timer (di call pas TAKE work order)
    pub fn start(
        &self,
        work_order_id: u64,
        executor_id: u64,
    ) -> Result<i64, String> {
        let mut timers = self.timers.lock().unwrap();

        if timers.contains_key(&work_order_id) {
            return Err("timer already running".into());
        }

        let started_at = Self::now();

        timers.insert(
            work_order_id,
            Timer {
                work_order_id,
                executor_id,
                started_at,
                stopped_at: None,
            },
        );

        Ok(started_at)
    }

    /// Stop timer (di call pas COMPLETE)
    pub fn stop(&self, work_order_id: u64) -> Result<(i64, i64), String> {
        let mut timers = self.timers.lock().unwrap();

        let timer = timers
            .remove(&work_order_id)
            .ok_or("timer not found")?;

        let stopped_at = Self::now();
        let elapsed = stopped_at - timer.started_at;

        Ok((timer.started_at, elapsed))
    }

    /// Status timer (polling frontend)
    pub fn status(&self, work_order_id: u64) -> Option<(i64, i64)> {
        let timers = self.timers.lock().unwrap();

        timers.get(&work_order_id).map(|timer| {
            let elapsed = Self::now() - timer.started_at;
            (timer.started_at, elapsed)
        })
    }
}
