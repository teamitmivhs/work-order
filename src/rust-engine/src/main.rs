use axum::{
    routing::{get, post},
    Router,
};
use std::{net::SocketAddr, sync::Arc};
use tokio::net::TcpListener;
use tracing_subscriber;

mod models;
mod state;
mod time_tracker;
mod web_api;

use state::AppState;
use time_tracker::TimeTracker;
use web_api::{start_timer, stop_timer, timer_status};

#[tokio::main]
async fn main() {
    // Initialize logging/tracing
    tracing_subscriber::fmt::init();

    let state = AppState::new();
    let tracker = Arc::new(TimeTracker::new(state));

    let app = Router::new()
        .route("/timer/start", post(start_timer))
        .route("/timer/stop", post(stop_timer))
        .route("/timer/:work_order_id", get(timer_status))
        .with_state(tracker);

    let addr = SocketAddr::from(([0, 0, 0, 0], 9000));
    
    match TcpListener::bind(addr).await {
        Ok(listener) => {
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
