CREATE TABLE IF NOT EXISTS `weekly_shift_schedule` (
  `DayOfWeek` tinyint unsigned NOT NULL,
  `Position` tinyint unsigned NOT NULL,
  `MemberID` int NOT NULL,
  `UpdatedBy` int DEFAULT NULL,
  `UpdatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`DayOfWeek`, `Position`),
  UNIQUE KEY `uq_weekly_shift_member` (`DayOfWeek`, `MemberID`),
  KEY `idx_weekly_shift_member` (`MemberID`),
  CONSTRAINT `fk_weekly_shift_member` FOREIGN KEY (`MemberID`) REFERENCES `members` (`ID`) ON DELETE CASCADE,
  CONSTRAINT `fk_weekly_shift_updated_by` FOREIGN KEY (`UpdatedBy`) REFERENCES `members` (`ID`) ON DELETE SET NULL,
  CONSTRAINT `chk_weekly_shift_day` CHECK (`DayOfWeek` BETWEEN 1 AND 5),
  CONSTRAINT `chk_weekly_shift_position` CHECK (`Position` BETWEEN 1 AND 3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
