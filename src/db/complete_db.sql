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
  `Status`   enum('onjob','standby','support','nextshift','offduty') NOT NULL DEFAULT 'offduty',
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
(1,  'Aldi Fadlurahman R',           '', 'Operator', 'Programmer',   'offduty', 'aldi.png', 'active', 'active', 1, NOW()),
(2,  'Raditya Ihsan Athallah',        '', 'Operator', 'Data Analyst', 'offduty', 'raditya.png', 'active', 'active', 1, NOW()),
(3,  'Azzam Alfarizhi',               '', 'Operator', 'Maintenance',  'offduty', 'azzam.png', 'active', 'active', 1, NOW()),
(4,  'Ghani Ilham Firdaus',           '', 'Operator', 'Programmer',   'offduty', 'ghani.png', 'active', 'active', 1, NOW()),
(5,  'Gilang Yoga Pangestu',          '', 'Operator', 'Maintenance',  'offduty', 'gilang.png', 'active', 'active', 1, NOW()),
(6,  'Charis Nur Noveli Alfaridzi',   '', 'Operator', 'Programmer',   'offduty', 'charis.png', 'active', 'active', 1, NOW()),
(7,  'M Akbar Putra P',              '', 'Operator', 'Soundman',     'offduty', 'akbar.png', 'active', 'active', 1, NOW()),
(8,  'Fahri Abdul Ghoni',             '', 'Operator', 'Maintenance',  'offduty', 'ghoni.png', 'active', 'active', 1, NOW()),
(9,  'Ridwan Bagoes Setiawan',        '', 'Operator', 'Programmer',   'offduty', 'ridwan.png', 'active', 'active', 1, NOW()),
(10, 'Moehammad Alvaro',              '', 'Operator', 'Programmer',   'offduty', 'alvaro.png', 'active', 'active', 1, NOW()),
(11, 'Dzaky Alvaro',                  '', 'Operator', 'Programmer',   'offduty', 'dzaky.png', 'active', 'active', 1, NOW()),
(12, 'Maynaldi Freza A',              '', 'Operator', 'Programmer',   'offduty', 'maynaldi.png', 'active', 'active', 1, NOW()),
(13, 'Muhammad Riva Nugraha',         '', 'Operator', 'Programmer',   'offduty', 'riva.png', 'active', 'active', 1, NOW()),
(14, 'Muhammad Farel Sustisna',       '', 'Operator', 'Programmer',   'offduty', 'farel.png', 'active', 'active', 1, NOW()),
(15, "Faa'iz Rizqi Haryono",          '', 'Operator', 'Programmer',   'offduty', 'faaiz.png', 'active', 'active', 1, NOW()),
(16, 'Royan Fadlan Musaminah',        '', 'Operator', 'Maintenance',  'offduty', 'royan.png', 'active', 'active', 1, NOW()),
(17, 'Muhammad Reyhansyah Hidayat',   '', 'Operator', 'Maintenance',  'offduty', 'reyhansyah.png', 'active', 'active', 1, NOW()),
(18, 'Naufal Abdilah Saputra',        '', 'Operator', 'Maintenance',  'offduty', 'naufal.png', 'active', 'active', 1, NOW()),
(19, 'Rizki Nuraulia',                '', 'Operator', 'Data Analyst', 'offduty', 'rizki.png', 'active', 'active', 1, NOW()),
(20, 'Khalishah Althaf',              '', 'Operator', 'Data Analyst', 'offduty', 'althaf.png', 'active', 'active', 1, NOW()),
(21, 'Dzakiyya Najdatul Rameyza',     '', 'Operator', 'Data Analyst', 'offduty', 'dzakiyya.png', 'active', 'active', 1, NOW()),
(22, 'Desvita Aurellia',              '', 'Operator', 'Data Analyst', 'offduty', 'desvita.png', 'active', 'active', 1, NOW()),
(23, 'Qiara Latifah Kaltsum',         '', 'Operator', 'Data Analyst', 'offduty', 'qiara.png', 'active', 'active', 1, NOW()),
(24, 'Dzaki Mathoriq',                '', 'Operator', 'Soundman',     'offduty', 'thoriq.png', 'active', 'active', 1, NOW()),
(25, 'Alvaro Purnomo',                '', 'Operator', 'Data Analyst', 'offduty', 'purnomo.png', 'active', 'active', 1, NOW()),
(26, 'Chelsea Aurelia',               '', 'Operator', 'Programmer',   'offduty', 'chelsea.png', 'active', 'active', 1, NOW()),
(27, 'Nabil Hilmy Zaenal',            '', 'Operator', 'Maintenance',  'offduty', 'nabil.png', 'active', 'active', 1, NOW()),
(28, 'Wildan Bait Maki',              '', 'Operator', 'Data Analyst', 'offduty', 'wildan.png', 'active', 'active', 1, NOW()),
(29, 'Gian Alvarezi Savatino Putra',  '', 'Operator', 'Programmer',   'offduty', 'gian.png', 'active', 'active', 1, NOW()),
(30, 'Jonathan Willy',                '', 'Operator', 'Programmer',   'offduty', 'willy.png', 'active', 'active', 1, NOW()),
(31, 'Pramadani Bintang Jasuma',      '', 'Operator', 'Maintenance',  'offduty', 'pramadani.png', 'active', 'active', 1, NOW()),
(32, 'Aira Nur Sabariyah Putri',      '', 'Operator', 'Programmer',   'offduty', 'aira.png', 'active', 'active', 1, NOW()),
(33, 'Azka Fakhri Alfito',            '', 'Operator', 'Programmer',   'offduty', 'azka.png', 'active', 'active', 1, NOW()),
(34, 'Ihsan Bintang Ghifari',         '', 'Operator', 'Programmer',   'offduty', 'ihsan.png', 'active', 'active', 1, NOW()),
(35, 'M. Dhimas Alfachry',            '', 'Operator', 'Maintenance',  'offduty', 'dhimas.png', 'active', 'active', 1, NOW()),
(36, 'Raissya Hanjani',               '', 'Operator', 'Data Analyst', 'offduty', 'raissya.png', 'active', 'active', 1, NOW()),
-- Akun guest permanen — password: guest123
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

SET FOREIGN_KEY_CHECKS = 1;
