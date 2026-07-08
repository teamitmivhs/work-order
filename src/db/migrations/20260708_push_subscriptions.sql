CREATE TABLE IF NOT EXISTS `push_subscriptions` (
  `ID`           int          NOT NULL AUTO_INCREMENT,
  `UserID`       int          NOT NULL,
  `Endpoint`     text         NOT NULL,
  `EndpointHash` char(64)     NOT NULL,
  `P256DH`       varchar(255) NOT NULL,
  `Auth`         varchar(255) NOT NULL,
  `UserAgent`    varchar(255) DEFAULT NULL,
  `CreatedAt`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `UpdatedAt`    datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `uq_push_endpoint_hash` (`EndpointHash`),
  KEY `idx_push_user` (`UserID`),
  CONSTRAINT `fk_push_user` FOREIGN KEY (`UserID`) REFERENCES `members` (`ID`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
