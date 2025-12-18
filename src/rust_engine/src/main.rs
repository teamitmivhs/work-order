use axum::{
    routing::{get, post},
    Router,
};
use std::net::SocketAddr;
use tracing::{info};
use tracing_subscriber;

mod web_api;
mod time_tracker;
mod models;

#[tokio::main]
async fn main() {
    //Logger
    tracing_subscriber::fmt()
        .with_target(false)
        .compact()
        .init();

    info!("Starting Rust Time tracker engine...");


    let app = Router::new()
        .route("/health", get(web_api::health_check))
        .route("/api/time/start", post(time_tracker::start_time))
        .route("/api/time/stop", post(time_tracker::stop_time))
        .route("/api/time/status", get(time_tracker::get_timer_status));
    
    let addr = SocketAddr::from(([0, 0, 0, 0], 9000));

    info!("Time engine listening on {}", addr);

    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .expect("Failed to start server");
}
