use axum::{
    routing::{get, post},
    Router,
};
use std::{env, net::SocketAddr, sync::Arc};
use tokio::net::TcpListener;

mod health;
mod models;
mod state;
mod time_tracker;
mod web_api;

use health::health_check;
use state::AppState;
use time_tracker::TimeTracker;
use web_api::{list_timers, start_timer, stop_timer, timer_status};

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let state = AppState::new();
    let tracker = Arc::new(TimeTracker::new(state));

    let app = Router::new()
        .route("/health", get(health_check))
        .route("/timer/start", post(start_timer))
        .route("/timer/stop", post(stop_timer))
        .route("/timer/:work_order_id", get(timer_status))
        // FIX: endpoint baru untuk list semua timer aktif
        .route("/timers", get(list_timers))
        .with_state(tracker);

    // FIX: port bisa dikonfigurasi via env variable PORT
    let port: u16 = env::var("PORT")
        .ok()
        .and_then(|p| p.parse().ok())
        .unwrap_or(9000);

    let addr = SocketAddr::from(([0, 0, 0, 0], port));

    match TcpListener::bind(addr).await {
        Ok(listener) => {
            tracing::info!("Time Tracker running on {}", addr);
            println!("⏱ Time Tracker running on {}", addr);

            if let Err(e) = axum::serve(listener, app).await {
                eprintln!("Server error: {}", e);
                std::process::exit(1);
            }
        }
        Err(e) => {
            eprintln!("Failed to bind to {}: {}", addr, e);
            std::process::exit(1);
        }
    }
}