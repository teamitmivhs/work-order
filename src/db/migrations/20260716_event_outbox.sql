CREATE TABLE IF NOT EXISTS `event_outbox` (
  `ID`            bigint unsigned NOT NULL AUTO_INCREMENT,
  `EventType`     varchar(80)      NOT NULL,
  `AggregateType` varchar(40)      NOT NULL,
  `AggregateID`   bigint unsigned  NOT NULL,
  `CreatedAt`     datetime(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `ProcessedAt`   datetime(3)      DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `idx_event_outbox_pending` (`ProcessedAt`, `ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
