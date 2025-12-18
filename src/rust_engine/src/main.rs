use axum::{
    routing::{get, post},
    Router,
};
use std::sync::Arc;
use std::net::SocketAddr;

mod state;
mod models;
mod web_api;

use state::TimeTracker;
use web_api::{start_timer, stop_timer, timer_status};

#[tokio::main]
async fn main() {
    // Shared state (in-memory)
    let tracker = Arc::new(TimeTracker::new());

    // Router
    let app = Router::new()
        .route("/timer/start", post(start_timer))
        .route("/timer/stop", post(stop_timer))
        .route("/timer/:work_order_id", get(timer_status))
        .with_state(tracker);

    let addr = SocketAddr::from(([0, 0, 0, 0], 9000));
    println!("⏱ Time Tracker running on {}", addr);

    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
