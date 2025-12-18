use axum::{
    routing::{get, post},
    Router,
};
use std::{net::SocketAddr, sync::Arc};
use tokio::net::TcpListener;

mod models;
mod state;
mod time_tracker;
mod web_api;

use state::AppState;
use time_tracker::TimeTracker;
use web_api::{start_timer, stop_timer, timer_status};

#[tokio::main]
async fn main() {
    let state = AppState::new();
    let tracker = Arc::new(TimeTracker::new(state));

    let app = Router::new()
        .route("/timer/start", post(start_timer))
        .route("/timer/stop", post(stop_timer))
        .route("/timer/:work_order_id", get(timer_status))
        .with_state(tracker);

    let addr = SocketAddr::from(([0, 0, 0, 0], 9000));
    let listener = TcpListener::bind(addr).await.unwrap();

    println!("⏱ Time Tracker running on {}", addr);

    axum::serve(listener, app).await.unwrap();
}
