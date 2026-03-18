-- ============================================================
-- DBWOIT - PostgreSQL Schema (Fresh Install)
-- Migrasi dari MySQL ke PostgreSQL
-- ============================================================

-- Drop tabel child dulu (foreign key order)
DROP TABLE IF EXISTS safetychecklist CASCADE;
DROP TABLE IF EXISTS executors        CASCADE;
DROP TABLE IF EXISTS orders           CASCADE;
DROP TABLE IF EXISTS members          CASCADE;

-- ============================================================
-- TABEL members
-- ============================================================
-- Perbedaan dari MySQL:
--   - AUTO_INCREMENT  → SERIAL (atau GENERATED ALWAYS AS IDENTITY)
--   - ENUM            → VARCHAR + CHECK CONSTRAINT (lebih portable)
--   - Backtick        → tidak ada (PostgreSQL pakai double-quote jika perlu, kolom lowercase aman tanpa quote)
CREATE TABLE members (
    id       SERIAL        PRIMARY KEY,
    name     VARCHAR(255)  NOT NULL,
    password VARCHAR(255)  NOT NULL DEFAULT '',
    role     VARCHAR(50)   DEFAULT NULL,
    status   VARCHAR(20)   NOT NULL DEFAULT 'offduty'
                           CHECK (status IN ('onjob','standby','support','nextshift','offduty')),
    avatar   VARCHAR(255)  NOT NULL DEFAULT 'no avatar',
    CONSTRAINT uq_members_name UNIQUE (name)
);

-- Seed data members (trailing space sudah dibersihkan)
INSERT INTO members (id, name, password, role, status, avatar) VALUES
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
(15, 'Faa''iz Rizqi Haryono',         '', 'programmer',   'offduty', 'faaiz.png'),
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
(36, 'Raissya Hanjani',               '', 'data analyst', 'offduty', 'raissya.png');

-- Sync sequence agar INSERT berikutnya tidak bentrok dengan ID yang sudah ada
SELECT setval('members_id_seq', (SELECT MAX(id) FROM members));

-- ============================================================
-- TABEL orders
-- ============================================================
-- Perbedaan dari MySQL:
--   - AUTO_INCREMENT → SERIAL
--   - TEXT           → TEXT (sama)
--   - DATETIME       → TIMESTAMP (atau tetap VARCHAR untuk display field)
--   - CompletedAt    → VARCHAR(20) karena Go menyimpan "HH:MM" bukan timestamp
--   - TimeSort       → TIME (sama)
CREATE TABLE orders (
    id            SERIAL        PRIMARY KEY,
    created_at    TIMESTAMP     DEFAULT NULL,
    updated_at    TIMESTAMP     DEFAULT NULL,
    deleted_at    TIMESTAMP     DEFAULT NULL,
    order_number  VARCHAR(50)   DEFAULT NULL,
    created_by    VARCHAR(100)  DEFAULT NULL,
    assigned_to   VARCHAR(100)  DEFAULT NULL,
    priority      VARCHAR(10)   DEFAULT NULL,
    time_display  VARCHAR(20)   DEFAULT NULL,
    time_sort     TIME          DEFAULT NULL,
    requester     VARCHAR(100)  DEFAULT NULL,
    location      VARCHAR(255)  DEFAULT NULL,
    device        VARCHAR(50)   DEFAULT NULL,
    problem       TEXT,
    working_hours VARCHAR(50)   DEFAULT NULL,
    status        VARCHAR(50)   DEFAULT NULL,
    completed_at  VARCHAR(20)   DEFAULT NULL,
    notes         TEXT          DEFAULT NULL
);

-- ============================================================
-- TABEL executors
-- ============================================================
-- Perbedaan dari MySQL:
--   - Nama kolom sudah menggunakan order_id / member_id (sesuai Go)
--   - ON DELETE CASCADE untuk kedua FK
CREATE TABLE executors (
    order_id  INT NOT NULL,
    member_id INT NOT NULL,
    PRIMARY KEY (order_id, member_id),
    CONSTRAINT fk_executors_order  FOREIGN KEY (order_id)  REFERENCES orders  (id) ON DELETE CASCADE,
    CONSTRAINT fk_executors_member FOREIGN KEY (member_id) REFERENCES members (id) ON DELETE CASCADE
);

-- ============================================================
-- TABEL safetychecklist
-- ============================================================
-- Perbedaan dari MySQL:
--   - Nama kolom sudah order_id (sesuai Go)
--   - VARCHAR(255) cukup untuk semua teks checklist
CREATE TABLE safetychecklist (
    order_id         INT          NOT NULL,
    safety_checklist VARCHAR(255) NOT NULL,
    PRIMARY KEY (order_id, safety_checklist),
    CONSTRAINT fk_safetychecklist_order FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);
