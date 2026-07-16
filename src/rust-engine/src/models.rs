use serde::Serialize;

#[derive(Clone, Debug, PartialEq, Eq, Serialize)]
pub struct OutboxEvent {
    pub id: u64,
    pub event_type: String,
    pub aggregate_type: String,
    pub aggregate_id: u64,
}

#[cfg(test)]
mod tests {
    use super::OutboxEvent;

    #[test]
    fn serializes_browser_event_shape() {
        let event = OutboxEvent {
            id: 7,
            event_type: "work_order.completed".into(),
            aggregate_type: "work_order".into(),
            aggregate_id: 42,
        };

        let json = serde_json::to_value(event).unwrap();
        assert_eq!(json["id"], 7);
        assert_eq!(json["event_type"], "work_order.completed");
        assert_eq!(json["aggregate_id"], 42);
    }
}
