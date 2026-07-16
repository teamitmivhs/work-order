use axum::{routing::get, Router};
use sqlx::mysql::{MySqlConnectOptions, MySqlPoolOptions};
use std::{env, net::SocketAddr, time::Duration};
use tokio::net::TcpListener;

mod event_engine;
mod health;
mod models;
mod state;
mod web_api;

use health::health_check;
use state::AppState;
use web_api::events;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();

    let pool = connect_database().await;
    let state = AppState::new(pool);
    tokio::spawn(event_engine::run(state.clone()));

    let app = Router::new()
        .route("/health", get(health_check))
        .route("/events", get(events))
        .with_state(state);

    let port: u16 = env::var("PORT")
        .ok()
        .and_then(|value| value.parse().ok())
        .unwrap_or(9000);
    let address = SocketAddr::from(([0, 0, 0, 0], port));

    match TcpListener::bind(address).await {
        Ok(listener) => {
            tracing::info!("Event engine running on {}", address);
            if let Err(error) = axum::serve(listener, app).await {
                tracing::error!(%error, "Event engine server stopped");
                std::process::exit(1);
            }
        }
        Err(error) => {
            tracing::error!(%error, %address, "Failed to bind event engine");
            std::process::exit(1);
        }
    }
}

async fn connect_database() -> sqlx::MySqlPool {
    let host = env::var("DB_HOST").unwrap_or_else(|_| "127.0.0.1".into());
    let port = env::var("DB_PORT")
        .ok()
        .and_then(|value| value.parse().ok())
        .unwrap_or(3306);
    let user = env::var("DB_USER").expect("DB_USER must be set");
    let password = env::var("DB_PASSWORD").expect("DB_PASSWORD must be set");
    let database = env::var("DB_NAME").unwrap_or_else(|_| "dbwoit".into());
    let options = MySqlConnectOptions::new()
        .host(&host)
        .port(port)
        .username(&user)
        .password(&password)
        .database(&database);

    for attempt in 1..=30 {
        match MySqlPoolOptions::new()
            .max_connections(5)
            .acquire_timeout(Duration::from_secs(5))
            .connect_with(options.clone())
            .await
        {
            Ok(pool) => return pool,
            Err(error) if attempt < 30 => {
                tracing::warn!(attempt, %error, "Database connection failed; retrying");
                tokio::time::sleep(Duration::from_secs(2)).await;
            }
            Err(error) => {
                tracing::error!(%error, "Database connection failed");
                std::process::exit(1);
            }
        }
    }

    unreachable!()
}
