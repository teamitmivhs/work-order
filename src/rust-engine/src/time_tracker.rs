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
    /// FIX: return Result<i64, String> agar clock error tidak diam-diam return 0
    fn now() -> Result<i64, String> {
        SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .map(|d| d.as_secs() as i64)
            .map_err(|e| format!("System clock error: {}", e))
    }

    /// Start a timer for a work order
    /// Returns started_at timestamp (unix seconds)
    /// Error jika timer sudah berjalan untuk order ini
    pub async fn start(&self, work_order_id: u64, executor_id: u64) -> Result<i64, String> {
        // FIX: validasi executor_id juga, bukan hanya work_order_id
        if work_order_id == 0 {
            return Err("Work order ID cannot be 0".to_string());
        }
        if executor_id == 0 {
            return Err("Executor ID cannot be 0".to_string());
        }

        // FIX: ambil timestamp SEBELUM lock agar durasi lock sesingkat mungkin
        let started_at = Self::now()?;

        // FIX: await karena sekarang pakai tokio::sync::Mutex
        let mut timers = self.state.timers.lock().await;

        if timers.contains_key(&work_order_id) {
            return Err(format!(
                "Timer already running for work order {}",
                work_order_id
            ));
        }

        timers.insert(
            work_order_id,
            RunningTimer {
                work_order_id,
                executor_id,
                started_at,
            },
        );

        tracing::info!(
            work_order_id,
            executor_id,
            started_at,
            "Timer started"
        );

        Ok(started_at)
    }

    /// Stop a timer for a work order
    /// FIX: return (started_at, stopped_at, duration) — tuple 3 elemen
    /// sehingga stopped_at bisa dikirim sebagai actual timestamp, bukan rekonstruksi
    pub async fn stop(&self, work_order_id: u64) -> Result<(i64, i64, i64), String> {
        if work_order_id == 0 {
            return Err("Work order ID cannot be 0".to_string());
        }

        // FIX: ambil timestamp sebelum lock
        let stopped_at = Self::now()?;

        let mut timers = self.state.timers.lock().await;

        let timer = timers
            .remove(&work_order_id)
            .ok_or_else(|| format!("Timer not found for work order {}", work_order_id))?;

        // FIX: guard durasi negatif — bisa terjadi akibat NTP sync / clock skew
        let duration = (stopped_at - timer.started_at).max(0);

        tracing::info!(
            work_order_id,
            duration_seconds = duration,
            "Timer stopped"
        );

        Ok((timer.started_at, stopped_at, duration))
    }

    /// Get current status of a timer tanpa menghentikannya
    /// Returns Option<(started_at, elapsed_seconds)>
    pub async fn status(&self, work_order_id: u64) -> Result<Option<(i64, i64)>, String> {
        if work_order_id == 0 {
            return Err("Work order ID cannot be 0".to_string());
        }

        let now = Self::now()?;
        let timers = self.state.timers.lock().await;

        Ok(timers.get(&work_order_id).map(|t| {
            // FIX: guard elapsed negatif
            let elapsed = (now - t.started_at).max(0);
            (t.started_at, elapsed)
        }))
    }

    /// Get semua timer yang sedang aktif beserta info lengkapnya
    /// FIX: fungsi baru untuk endpoint list — menggantikan active_count() yang tidak diekspos
    pub async fn list_active(&self) -> Result<Vec<(u64, u64, i64, i64)>, String> {
        let now = Self::now()?;
        let timers = self.state.timers.lock().await;

        let result = timers
            .values()
            .map(|t| {
                let elapsed = (now - t.started_at).max(0);
                (t.work_order_id, t.executor_id, t.started_at, elapsed)
            })
            .collect();

        Ok(result)
    }

    /// Get jumlah timer aktif (dipakai oleh health check)
    pub async fn active_count(&self) -> Result<usize, String> {
        let timers = self.state.timers.lock().await;
        Ok(timers.len())
    }
}