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
-- FIX: Trailing space dibersihkan dari data rows 26-34
DROP TABLE IF EXISTS `members`;
CREATE TABLE `members` (
  `ID`       int          NOT NULL AUTO_INCREMENT,
  `Name`     varchar(255) NOT NULL,
  `Password` varchar(255) NOT NULL DEFAULT '',
  `Role`     varchar(50)  DEFAULT NULL,
  `Status`   enum('onjob','standby','support','nextshift','offduty') NOT NULL DEFAULT 'offduty',
  `Avatar`   varchar(255) NOT NULL DEFAULT 'no avatar',
  PRIMARY KEY (`ID`),
  UNIQUE KEY `uq_members_name` (`Name`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `members` (`ID`,`Name`,`Password`,`Role`,`Status`,`Avatar`) VALUES
(1,  'Aldi Fadlurahman R',           '', 'programmer',   'offduty', 'aldi.png'),
(2,  'Raditya Ihsan Athallah',        '', 'data analyst', 'offduty', 'raditya.png'),
(3,  'Azzam Alfarizhi',               '', 'maintenance',  'offduty', 'azzam.png'),
(4,  'Ghani Ilham Firdaus',           '', 'programmer',   'offduty', 'ghani.png'),
(5,  'Gilang Yoga Pangestu',          '', 'maintenance',  'offduty', 'gilang.png'),
(6,  'Charis Nur Noveli Alfaridzi',   '', 'programmer',   'offduty', 'charis.png'),
(7,  'M Akbar Putra P',              '', 'soundman',     'offduty', 'akbar.png'),
(8,  'Fahri Abdul Ghoni',             '', 'maintenance',  'offduty', 'ghoni.png'),
(9,  'Ridwan Bagoes Setiawan',        '', 'programmer',   'offduty', 'ridwan.png'),
(10, 'Moehammad Alvaro',              '', 'programmer',   'offduty', 'alvaro.png'),
(11, 'Dzaky Alvaro',                  '', 'programmer',   'offduty', 'dzaky.png'),
(12, 'Maynaldi Freza A',              '', 'programmer',   'offduty', 'maynaldi.png'),
(13, 'Muhammad Riva Nugraha',         '', 'programmer',   'offduty', 'riva.png'),
(14, 'Muhammad Farel Sustisna',       '', 'programmer',   'offduty', 'farel.png'),
(15, "Faa'iz Rizqi Haryono",          '', 'programmer',   'offduty', 'faaiz.png'),
(16, 'Royan Fadlan Musaminah',        '', 'maintenance',  'offduty', 'royan.png'),
(17, 'Muhammad Reyhansyah Hidayat',   '', 'maintenance',  'offduty', 'reyhansyah.png'),
(18, 'Naufal Abdilah Saputra',        '', 'maintenance',  'offduty', 'naufal.png'),
(19, 'Rizki Nuraulia',                '', 'data analyst', 'offduty', 'rizki.png'),
(20, 'Khalishah Althaf',              '', 'data analyst', 'offduty', 'althaf.png'),
(21, 'Dzakiyya Najdatul Rameyza',     '', 'data analyst', 'offduty', 'dzakiyya.png'),
(22, 'Desvita Aurellia',              '', 'data analyst', 'offduty', 'desvita.png'),
(23, 'Qiara Latifah Kaltsum',         '', 'data analyst', 'offduty', 'qiara.png'),
(24, 'Dzaki Mathoriq',                '', 'soundman',     'offduty', 'thoriq.png'),
(25, 'Alvaro Purnomo',                '', 'data analyst', 'offduty', 'purnomo.png'),
(26, 'Chelsea Aurelia',               '', 'programmer',   'offduty', 'chelsea.png'),
(27, 'Nabil Hilmy Zaenal',            '', 'maintenance',  'offduty', 'nabil.png'),
(28, 'Wildan Bait Maki',              '', 'data analyst', 'offduty', 'wildan.png'),
(29, 'Gian Alvarezi Savatino Putra',  '', 'programmer',   'offduty', 'gian.png'),
(30, 'Jonathan Willy',                '', 'programmer',   'offduty', 'willy.png'),
(31, 'Pramadani Bintang Jasuma',      '', 'maintenance',  'offduty', 'pramadani.png'),
(32, 'Aira Nur Sabariyah Putri',      '', 'programmer',   'offduty', 'aira.png'),
(33, 'Azka Fakhri Alfito',            '', 'programmer',   'offduty', 'azka.png'),
(34, 'Ihsan Bintang Ghifari',         '', 'programmer',   'offduty', 'ihsan.png'),
(35, 'M. Dhimas Alfachry',            '', 'maintenance',  'offduty', 'dhimas.png'),
(36, 'Raissya Hanjani',               '', 'data analyst', 'offduty', 'raissya.png'),
-- Akun guest permanen — password: guest123
(37, 'guest', '$2b$10$kZNK9YSmZhYDt83CQBzdg.OqF0S17Hge7O3BpURFRXVFUSZbQYSWy', 'Guest', 'offduty', 'default-avatar.png');

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
  `CreatedBy`    varchar(100) DEFAULT NULL,
  `AssignedTo`   varchar(100) DEFAULT NULL,
  `Priority`     varchar(10)  DEFAULT NULL,
  `TimeDisplay`  varchar(20)  DEFAULT NULL,
  `TimeSort`     time         DEFAULT NULL,
  `Requester`    varchar(100) DEFAULT NULL,
  `Location`     varchar(255) DEFAULT NULL,
  `Device`       varchar(50)  DEFAULT NULL,
  `Problem`      text,
  `WorkingHours` varchar(50)  DEFAULT NULL,
  `Status`       varchar(50)  DEFAULT NULL,
  `CompletedAt`  varchar(20)  DEFAULT NULL,
  `Notes`        text         DEFAULT NULL,
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