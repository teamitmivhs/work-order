-- ============================================================
-- DBWOIT - COMPLETE SCHEMA (FRESH INSTALL)
-- Versi lengkap yang sudah disesuaikan dengan Go backend
-- Gunakan file ini jika setup database dari awal
-- ============================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================
-- TABEL members
-- ============================================================
-- FIX: Tambah kolom Password untuk auth (login/register)
-- FIX: Role diubah dari ENUM sempit ke VARCHAR(50) agar bisa
--      menerima 'Operator' dan 'Admin' dari user_controller.go
-- FIX: Division dipisah dari Role untuk divisi teknis staff
-- FIX: Trailing space dibersihkan dari data rows 26-34
DROP TABLE IF EXISTS `members`;
CREATE TABLE `members` (
  `ID`       int          NOT NULL AUTO_INCREMENT,
  `Name`     varchar(255) NOT NULL,
  `Password` varchar(255) NOT NULL DEFAULT '',
  `Role`     varchar(50)  DEFAULT NULL,
  `Division` varchar(50)  DEFAULT NULL,
  `Status`   enum('onjob','standby','nextshift','offduty') NOT NULL DEFAULT 'offduty',
  `Avatar`   varchar(255) NOT NULL DEFAULT 'no avatar',
  `AccountStatus` enum('pending','active','rejected','disabled') NOT NULL DEFAULT 'active',
  `MembershipStatus` enum('active','alumni','inactive') NOT NULL DEFAULT 'active',
  `BatchYear` varchar(20) DEFAULT NULL,
  `GraduationYear` int DEFAULT NULL,
  `CanHandleWorkOrder` tinyint(1) NOT NULL DEFAULT 1,
  `RegisteredAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ApprovedAt` datetime DEFAULT NULL,
  `ApprovedBy` int DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `uq_members_name` (`Name`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `members` (`ID`,`Name`,`Password`,`Role`,`Division`,`Status`,`Avatar`,`AccountStatus`,`MembershipStatus`,`CanHandleWorkOrder`,`ApprovedAt`) VALUES
(1,  'Demo Programmer 01',        '', 'Operator', 'Programmer',   'offduty', 'aldi.png', 'active', 'active', 1, NOW()),
(2,  'Demo Data Analyst 01',      '', 'Operator', 'Data Analyst', 'offduty', 'raditya.png', 'active', 'active', 1, NOW()),
(3,  'Demo Maintenance 01',       '', 'Operator', 'Maintenance',  'offduty', 'azzam.png', 'active', 'active', 1, NOW()),
(4,  'Demo Programmer 02',        '', 'Operator', 'Programmer',   'offduty', 'ghani.png', 'active', 'active', 1, NOW()),
(5,  'Demo Maintenance 02',       '', 'Operator', 'Maintenance',  'offduty', 'gilang.png', 'active', 'active', 1, NOW()),
(6,  'Demo Programmer 03',        '', 'Operator', 'Programmer',   'offduty', 'charis.png', 'active', 'active', 1, NOW()),
(7,  'Demo Soundman 01',          '', 'Operator', 'Soundman',     'offduty', 'akbar.png', 'active', 'active', 1, NOW()),
(8,  'Demo Maintenance 03',       '', 'Operator', 'Maintenance',  'offduty', 'ghoni.png', 'active', 'active', 1, NOW()),
(9,  'Demo Programmer 04',        '', 'Operator', 'Programmer',   'offduty', 'ridwan.png', 'active', 'active', 1, NOW()),
(10, 'Demo Programmer 05',        '', 'Operator', 'Programmer',   'offduty', 'alvaro.png', 'active', 'active', 1, NOW()),
(11, 'Demo Programmer 06',        '', 'Operator', 'Programmer',   'offduty', 'dzaky.png', 'active', 'active', 1, NOW()),
(12, 'Demo Programmer 07',        '', 'Operator', 'Programmer',   'offduty', 'maynaldi.png', 'active', 'active', 1, NOW()),
(13, 'Demo Programmer 08',        '', 'Operator', 'Programmer',   'offduty', 'riva.png', 'active', 'active', 1, NOW()),
(14, 'Demo Programmer 09',        '', 'Operator', 'Programmer',   'offduty', 'farel.png', 'active', 'active', 1, NOW()),
(15, 'Demo Programmer 10',        '', 'Operator', 'Programmer',   'offduty', 'faaiz.png', 'active', 'active', 1, NOW()),
(16, 'Demo Maintenance 04',       '', 'Operator', 'Maintenance',  'offduty', 'royan.png', 'active', 'active', 1, NOW()),
(17, 'Demo Maintenance 05',       '', 'Operator', 'Maintenance',  'offduty', 'reyhansyah.png', 'active', 'active', 1, NOW()),
(18, 'Demo Maintenance 06',       '', 'Operator', 'Maintenance',  'offduty', 'naufal.png', 'active', 'active', 1, NOW()),
(19, 'Demo Data Analyst 02',      '', 'Operator', 'Data Analyst', 'offduty', 'rizki.png', 'active', 'active', 1, NOW()),
(20, 'Demo Data Analyst 03',      '', 'Operator', 'Data Analyst', 'offduty', 'althaf.png', 'active', 'active', 1, NOW()),
(21, 'Demo Data Analyst 04',      '', 'Operator', 'Data Analyst', 'offduty', 'dzakiyya.png', 'active', 'active', 1, NOW()),
(22, 'Demo Data Analyst 05',      '', 'Operator', 'Data Analyst', 'offduty', 'desvita.png', 'active', 'active', 1, NOW()),
(23, 'Demo Data Analyst 06',      '', 'Operator', 'Data Analyst', 'offduty', 'qiara.png', 'active', 'active', 1, NOW()),
(24, 'Demo Soundman 02',          '', 'Operator', 'Soundman',     'offduty', 'thoriq.png', 'active', 'active', 1, NOW()),
(25, 'Demo Data Analyst 07',      '', 'Operator', 'Data Analyst', 'offduty', 'purnomo.png', 'active', 'active', 1, NOW()),
(26, 'Demo Programmer 11',        '', 'Operator', 'Programmer',   'offduty', 'chelsea.png', 'active', 'active', 1, NOW()),
(27, 'Demo Maintenance 07',       '', 'Operator', 'Maintenance',  'offduty', 'nabil.png', 'active', 'active', 1, NOW()),
(28, 'Demo Data Analyst 08',      '', 'Operator', 'Data Analyst', 'offduty', 'wildan.png', 'active', 'active', 1, NOW()),
(29, 'Demo Programmer 12',        '', 'Operator', 'Programmer',   'offduty', 'gian.png', 'active', 'active', 1, NOW()),
(30, 'Demo Programmer 13',        '', 'Operator', 'Programmer',   'offduty', 'willy.png', 'active', 'active', 1, NOW()),
(31, 'Demo Maintenance 08',       '', 'Operator', 'Maintenance',  'offduty', 'pramadani.png', 'active', 'active', 1, NOW()),
(32, 'Demo Programmer 14',        '', 'Operator', 'Programmer',   'offduty', 'aira.png', 'active', 'active', 1, NOW()),
(33, 'Demo Programmer 15',        '', 'Operator', 'Programmer',   'offduty', 'azka.png', 'active', 'active', 1, NOW()),
(34, 'Demo Programmer 16',        '', 'Operator', 'Programmer',   'offduty', 'ihsan.png', 'active', 'active', 1, NOW()),
(35, 'Demo Maintenance 09',       '', 'Operator', 'Maintenance',  'offduty', 'dhimas.png', 'active', 'active', 1, NOW()),
(36, 'Demo Data Analyst 09',      '', 'Operator', 'Data Analyst', 'offduty', 'raissya.png', 'active', 'active', 1, NOW()),
(37, 'guest', '$2b$10$kZNK9YSmZhYDt83CQBzdg.OqF0S17Hge7O3BpURFRXVFUSZbQYSWy', 'Guest', NULL, 'offduty', 'default-avatar.png', 'active', 'active', 0, NOW());

-- ============================================================
-- TABEL orders
-- ============================================================
-- FIX: CompletedAt diubah dari DATETIME ke VARCHAR(20)
--      Go menyimpan string "HH:MM" bukan timestamp penuh
-- FIX: Tambah kolom Notes untuk fitur catatan evaluasi
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `ID`           int          NOT NULL AUTO_INCREMENT,
  `CreatedAt`    datetime     DEFAULT NULL,
  `UpdatedAt`    datetime     DEFAULT NULL,
  `DeletedAt`    datetime     DEFAULT NULL,
  `OrderNumber`  varchar(50)  DEFAULT NULL,
  `TrackingCode` varchar(20)  DEFAULT NULL,
  `CreatedBy`    varchar(100) DEFAULT NULL,
  `AssignedTo`   varchar(100) DEFAULT NULL,
  `Priority`     varchar(10)  DEFAULT NULL,
  `TimeDisplay`  varchar(20)  DEFAULT NULL,
  `TimeSort`     time         DEFAULT NULL,
  `StartedAt`    datetime     DEFAULT NULL,
  `Requester`    varchar(100) DEFAULT NULL,
  `Location`     varchar(255) DEFAULT NULL,
  `Device`       varchar(50)  DEFAULT NULL,
  `Problem`      text,
  `WorkingHours` varchar(50)  DEFAULT NULL,
  `Status`       varchar(50)  DEFAULT NULL,
  `CompletedAt`  varchar(20)  DEFAULT NULL,
  `Notes`        text         DEFAULT NULL,
  `AdminNotes`   text         DEFAULT NULL,
  `Rating`       tinyint      DEFAULT NULL,
  `NotesQuality` tinyint      DEFAULT NULL,
  `DocumentationPhoto` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================================
-- TABEL executors
-- ============================================================
-- FIX: Nama kolom diubah dari 'ID'/'Executors' menjadi 'order_id'/'member_id'
--      agar sesuai dengan semua query di workorder_repository.go
-- FIX: Tambah ON DELETE CASCADE agar executor otomatis terhapus
--      saat order dihapus (tidak perlu hapus manual di Go lagi)
-- FIX: Tambah foreign key ke members untuk integritas referensial
DROP TABLE IF EXISTS `executors`;
CREATE TABLE `executors` (
  `order_id`  int NOT NULL,
  `member_id` int NOT NULL,
  PRIMARY KEY (`order_id`, `member_id`),
  CONSTRAINT `fk_executors_order`  FOREIGN KEY (`order_id`)  REFERENCES `orders`  (`ID`) ON DELETE CASCADE,
  CONSTRAINT `fk_executors_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`ID`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================================
-- TABEL safetychecklist
-- ============================================================
-- FIX: Nama kolom diubah dari 'ID' menjadi 'order_id'
--      agar sesuai dengan semua query di workorder_repository.go
-- FIX: VARCHAR(50) → VARCHAR(255) karena teks checklist
--      seperti "Gunakan pelindung mata (goggles)" melebihi 50 karakter
-- FIX: Tambah ON DELETE CASCADE agar checklist otomatis terhapus
--      saat order dihapus
DROP TABLE IF EXISTS `safetychecklist`;
CREATE TABLE `safetychecklist` (
  `order_id`       int          NOT NULL,
  `SafetyChecklist` varchar(255) NOT NULL,
  PRIMARY KEY (`order_id`, `SafetyChecklist`),
  CONSTRAINT `fk_safetychecklist_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`ID`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================================
-- TABEL push_subscriptions
-- ============================================================
-- Menyimpan Web Push subscription per browser/perangkat user.
DROP TABLE IF EXISTS `push_subscriptions`;
CREATE TABLE `push_subscriptions` (
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

-- ============================================================
-- TABEL event_outbox
-- ============================================================
-- Menyimpan event bisnis secara transaksional sebelum dikirim oleh Rust.
DROP TABLE IF EXISTS `event_outbox`;
CREATE TABLE `event_outbox` (
  `ID`            bigint unsigned NOT NULL AUTO_INCREMENT,
  `EventType`     varchar(80)      NOT NULL,
  `AggregateType` varchar(40)      NOT NULL,
  `AggregateID`   bigint unsigned  NOT NULL,
  `CreatedAt`     datetime(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `ProcessedAt`   datetime(3)      DEFAULT NULL,
  PRIMARY KEY (`ID`),
  KEY `idx_event_outbox_pending` (`ProcessedAt`, `ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ============================================================
-- TABEL weekly_shift_schedule
-- ============================================================
-- Template jadwal mingguan Senin-Jumat, tepat 3 orang per hari.
DROP TABLE IF EXISTS `weekly_shift_schedule`;
CREATE TABLE `weekly_shift_schedule` (
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

-- ============================================================
-- TABEL shift_day_counter (legacy rollover)
-- ============================================================
-- Menyimpan tanggal terakhir rollover shift.
-- Saat hari berganti, semua member dengan status nextshift dipindah ke standby.
DROP TABLE IF EXISTS `shift_day_counter`;
CREATE TABLE `shift_day_counter` (
  `ID` tinyint NOT NULL,
  `LastDate` date NOT NULL,
  `LastDay` int NOT NULL,
  `LastMonth` int NOT NULL,
  `LastYear` int NOT NULL,
  `RolloverCount` int NOT NULL DEFAULT 0,
  `UpdatedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
