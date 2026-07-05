DROP PROCEDURE IF EXISTS add_member_lifecycle_columns;

DELIMITER //

CREATE PROCEDURE add_member_lifecycle_columns()
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'members' AND COLUMN_NAME = 'Division'
  ) THEN
    ALTER TABLE members ADD COLUMN Division varchar(50) DEFAULT NULL AFTER Role;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'members' AND COLUMN_NAME = 'AccountStatus'
  ) THEN
    ALTER TABLE members ADD COLUMN AccountStatus enum('pending','active','rejected','disabled') NOT NULL DEFAULT 'active' AFTER Avatar;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'members' AND COLUMN_NAME = 'MembershipStatus'
  ) THEN
    ALTER TABLE members ADD COLUMN MembershipStatus enum('active','alumni','inactive') NOT NULL DEFAULT 'active' AFTER AccountStatus;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'members' AND COLUMN_NAME = 'BatchYear'
  ) THEN
    ALTER TABLE members ADD COLUMN BatchYear varchar(20) DEFAULT NULL AFTER MembershipStatus;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'members' AND COLUMN_NAME = 'GraduationYear'
  ) THEN
    ALTER TABLE members ADD COLUMN GraduationYear int DEFAULT NULL AFTER BatchYear;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'members' AND COLUMN_NAME = 'CanHandleWorkOrder'
  ) THEN
    ALTER TABLE members ADD COLUMN CanHandleWorkOrder tinyint(1) NOT NULL DEFAULT 1 AFTER GraduationYear;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'members' AND COLUMN_NAME = 'RegisteredAt'
  ) THEN
    ALTER TABLE members ADD COLUMN RegisteredAt datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER CanHandleWorkOrder;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'members' AND COLUMN_NAME = 'ApprovedAt'
  ) THEN
    ALTER TABLE members ADD COLUMN ApprovedAt datetime DEFAULT NULL AFTER RegisteredAt;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'members' AND COLUMN_NAME = 'ApprovedBy'
  ) THEN
    ALTER TABLE members ADD COLUMN ApprovedBy int DEFAULT NULL AFTER ApprovedAt;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'orders' AND COLUMN_NAME = 'StartedAt'
  ) THEN
    ALTER TABLE orders ADD COLUMN StartedAt datetime DEFAULT NULL AFTER TimeSort;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'orders' AND COLUMN_NAME = 'Rating'
  ) THEN
    ALTER TABLE orders ADD COLUMN Rating tinyint DEFAULT NULL AFTER Notes;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'orders' AND COLUMN_NAME = 'AdminNotes'
  ) THEN
    ALTER TABLE orders ADD COLUMN AdminNotes text DEFAULT NULL AFTER Notes;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'orders' AND COLUMN_NAME = 'NotesQuality'
  ) THEN
    ALTER TABLE orders ADD COLUMN NotesQuality tinyint DEFAULT NULL AFTER Rating;
  END IF;

  IF NOT EXISTS (
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'orders' AND COLUMN_NAME = 'DocumentationPhoto'
  ) THEN
    ALTER TABLE orders ADD COLUMN DocumentationPhoto varchar(255) DEFAULT NULL AFTER NotesQuality;
  END IF;
END//

DELIMITER ;

CALL add_member_lifecycle_columns();

DROP PROCEDURE IF EXISTS add_member_lifecycle_columns;

UPDATE members
SET Division = CASE
    WHEN LOWER(TRIM(Role)) = 'soundman' THEN 'Soundman'
    WHEN LOWER(TRIM(Role)) = 'programmer' THEN 'Programmer'
    WHEN LOWER(TRIM(Role)) = 'maintenance' THEN 'Maintenance'
    WHEN LOWER(TRIM(Role)) IN ('data analyst', 'data-analyst', 'data_analyst') THEN 'Data Analyst'
    ELSE Division
  END,
  Role = CASE
    WHEN LOWER(TRIM(Role)) IN ('soundman', 'programmer', 'maintenance', 'data analyst', 'data-analyst', 'data_analyst') THEN 'Operator'
    ELSE Role
  END
WHERE Division IS NULL OR Division = '';

UPDATE members
SET Role = 'Operator'
WHERE LOWER(TRIM(Role)) IN ('staff', 'technician');

UPDATE members
SET CanHandleWorkOrder = 0
WHERE Role = 'Guest';
