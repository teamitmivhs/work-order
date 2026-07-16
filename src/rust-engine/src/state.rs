use sqlx::MySqlPool;
use std::sync::{atomic::AtomicU64, Arc};
use tokio::sync::broadcast;

use crate::models::OutboxEvent;

#[derive(Clone)]
pub struct AppState {
    pub pool: MySqlPool,
    pub events: broadcast::Sender<OutboxEvent>,
    pub processed_events: Arc<AtomicU64>,
}

impl AppState {
    pub fn new(pool: MySqlPool) -> Self {
        let (events, _) = broadcast::channel(256);
        Self {
            pool,
            events,
            processed_events: Arc::new(AtomicU64::new(0)),
        }
    }
}
