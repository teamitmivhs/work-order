use axum::{
    extract::State,
    response::sse::{Event, KeepAlive, Sse},
};
use std::{convert::Infallible, time::Duration};
use tokio_stream::{wrappers::BroadcastStream, Stream, StreamExt};

use crate::state::AppState;

pub async fn events(
    State(state): State<AppState>,
) -> Sse<impl Stream<Item = Result<Event, Infallible>>> {
    let sync = tokio_stream::once(Ok(Event::default().data("{\"event_type\":\"sync\"}")));
    let updates = BroadcastStream::new(state.events.subscribe()).filter_map(|message| {
        let message = message.ok()?;
        Event::default()
            .id(message.id.to_string())
            .json_data(message)
            .ok()
            .map(Ok)
    });

    Sse::new(sync.chain(updates)).keep_alive(KeepAlive::new().interval(Duration::from_secs(15)))
}
