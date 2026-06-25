# Repository Cleanup & Refactoring Summary

**Last Updated**: March 2026  
**Status**: ✅ COMPLETE

---

## Overview

Dokumen ini merangkum semua perubahan signifikan yang dilakukan pada sistem Work Order MIVHS, mencakup perbaikan bug, refactoring, dan peningkatan fitur di seluruh stack (Go backend, Rust engine, frontend, database, dan Docker).

---

## 1. Go Backend

### Bug Fixes
- **`TakeOrder` logika assignment terbalik** — Pengecekan `IsMemberAssigned` sebelum take dihapus (logika keliru). Sekarang cek dilakukan hanya di `CompleteOrder`.
- **`CompleteOrder` resource leak** — `rows.Close()` ditambah secara eksplisit sebelum `tx.Exec()` untuk mencegah deadlock transaksi.
- **`DeleteOrder` tidak reset status member** — Sekarang member yang `onjob` di-reset ke `standby` saat order dihapus.
- **`UpdateMemberStatus` tidak sync ke API** — Sekarang memanggil `PATCH /api/members/:id/status` ke database.
- **Race condition di `TakeOrder`** — Pengecekan status executor dipindah ke dalam transaksi dengan `SELECT FOR UPDATE`.
- **`sql.NullString` tidak kompatibel JSON** — Diganti `*string` pada field `CompletedAt`.
- **`Requester any`** — Diganti `string` karena selalu berupa string dari database.
- **Duplikat `exitGuestBtn` listener** — Digabung menjadi satu handler.
- **Empty slice vs null** — Semua fungsi repository diinisialisasi dengan `make([]T, 0)` bukan `var slice []T`.

### Endpoints Baru
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `POST` | `/api/logout` | Logout user |
| `PATCH` | `/api/members/:id/status` | Update status member |
| `PATCH` | `/api/workorders/:id` | Update executor list order |

### Perubahan Struktur
- **`WorkOrderRequest.CompletedAt`** — `sql.NullString` → `*string`
- **`WorkOrder.Requester`** — `any` → `string`
- **`UpdateWorkOrderRequest`** — struct baru untuk PATCH endpoint
- **Response format** — Semua response konsisten `{ code, message, data }`

### Keamanan
- **`db.go`** — `DB_USER` dan `DB_PASSWORD` wajib dari env (tidak ada fallback hardcoded)
- **`jwt.go`** — Secret dibaca dengan `sync.Once` saat pertama dibutuhkan, bukan saat package init
- **Rate limiting** — 10 req/menit per IP pada `/login` dan `/register`
- **`INTERNAL_API_KEY`** — Shared secret antara Go dan Rust

### Integrasi Timer
- `TakeOrder` sekarang memanggil `services.StartTimer()` setelah commit
- `CompleteOrder` sekarang memanggil `services.StopTimer()`, durasi disimpan ke `working_hours`
- `time_tracker.go` (Go client) diperbarui: return `(int64, error)`, kirim `X-Internal-Key` header, baca URL saat runtime bukan startup

---

## 2. Rust Time Tracker

### Bug Fixes
- **`std::sync::Mutex` blocking** — Diganti `tokio::sync::Mutex` untuk async context
- **`now()` silent zero** — Diganti return `Result<i64, String>` agar error clock terdeteksi
- **`stop()` `stopped_at` salah** — Sekarang return tuple `(started_at, stopped_at, duration)` dengan `stopped_at` aktual
- **Durasi negatif** — Ditambah `.max(0)` guard
- **`executor_id` tidak divalidasi** — Sekarang validasi sama dengan `work_order_id`
- **`chrono` dan `uuid`** — Dihapus dari `Cargo.toml` (tidak dipakai, menambah build time)

### Fitur Baru
- **`X-Internal-Key` authentication** — Semua endpoint diproteksi
- **`GET /timers`** — List semua timer aktif
- **`PORT` env variable** — Port bisa dikonfigurasi
- **`list_active()`** — Method baru di `TimeTracker`
- **`impl Default for AppState`** — Idiomatis Go/Rust
- **Health check** menampilkan jumlah timer aktif

---

## 3. Database Schema

### Perubahan Tabel `members`
| Perubahan | Sebelum | Sesudah |
|-----------|---------|---------|
| Kolom Password | Tidak ada | `VARCHAR(255) NOT NULL DEFAULT ''` |
| Kolom Role | `ENUM('programmer','maintenance',...)` | `VARCHAR(50)` |
| Constraint | Tidak ada | `UNIQUE (name)` |
| Data rows 26-34 | Ada trailing space | Sudah di-trim |

### Perubahan Tabel `executors`
| Perubahan | Sebelum | Sesudah |
|-----------|---------|---------|
| Nama kolom | `ID`, `Executors` | `order_id`, `member_id` |
| FK ke members | Tidak ada | `ON DELETE CASCADE` |

### Perubahan Tabel `safetychecklist`
| Perubahan | Sebelum | Sesudah |
|-----------|---------|---------|
| Nama kolom | `ID` | `order_id` |
| Panjang kolom | `VARCHAR(50)` | `VARCHAR(255)` |
| FK | Tidak ada cascade | `ON DELETE CASCADE` |

### Perubahan Tabel `orders`
| Perubahan | Sebelum | Sesudah |
|-----------|---------|---------|
| `CompletedAt` | `DATETIME` | `VARCHAR(20)` |
| `Notes` | Tidak ada | `TEXT DEFAULT NULL` |

---

## 4. Frontend

### Bug Fixes
- **Welcome banner tidak terbaca** — `text-[#333]` di background gelap → `text-white`
- **Footer layout rusak** — `div.col-span-3` di luar grid → `<footer>` semantic
- **Dropdown conflict** — CSS `group-hover:block` + JS `toggle` → JS only
- **`viewbox` typo** — Diperbaiki ke `viewBox` di semua SVG (8+ instance)
- **Mobile menu salah posisi** — Ditambah CSS `.mobile-menu-active` untuk dropdown vertikal
- **User profile hardcoded** — Sekarang dibaca dari JWT token + `GET /api/profile`
- **`.status-completed` blink terus** — Dihapus animasi, diganti warna statis
- **`summary.html` + `kaizen.html` pakai localStorage** — Dimigrasi ke API
- **`colspan="10"` salah** — Diperbaiki ke `12` sesuai jumlah kolom
- **Kabel LAN step 3 terpotong** — Konten dilengkapi
- **Modal backdrop tidak bisa tutup** — Konsistenkan `classList` vs `style.display`

### script.js Fixes
- **Crash di halaman login** — Ditambah null check untuk elemen navbar
- **Duplikat login/register listener** — Blok `DOMContentLoaded` kedua dihapus
- **API response format** — Ditambah `unwrapData()` untuk handle `{ code, message, data }`
- **Missing `Authorization` header** — Ditambah `authHeaders()` helper di semua protected fetch
- **GSAP null check** — Dibungkus `if (typeof gsap !== 'undefined')`
- **`logoutBtn` index salah** — `querySelector('button')` → `querySelectorAll('button')[1]`
- **Mobile menu resize** — Ditambah `window.addEventListener('resize', ...)` reset

---

## 5. Docker & Infrastructure

### docker-compose.yml
- **Credentials tidak hardcode** — Semua pakai `${VAR}` dari `.env`
- **`INTERNAL_API_KEY`** — Ditambah ke backend dan time-tracker
- **`time-tracker` depends_on** — `service_started` → `service_healthy`
- **`nginx` depends_on** — Dari `condition: service_healthy` (backend tidak punya healthcheck) → `- backend`
- **Volume mounts** — Ditambah `:ro` (read-only) untuk file HTML/CSS/JS di nginx
- **Backend healthcheck dihapus** — `wget --spider HEAD /health` selalu 404 karena Gin hanya daftarkan GET. Backend sudah punya retry logic ke DB sendiri.

### Dockerfile Rust
- **`COPY Cargo.toml Cargo.lock`** — Diubah opsional (tanpa Cargo.lock jika belum ada)
- **`HEALTHCHECK`** — Ditambah menggunakan `/health` endpoint
- **`touch src/main.rs`** — Ditambah agar cargo tidak skip recompile

### .env
- **`.env.example`** — Ditambah `INTERNAL_API_KEY`
- **`.env.production`** — Diperbarui sesuai variable baru

---

## Struktur Repository Saat Ini

```
work-order/
├── .gitignore
├── README.md
├── DOCKER_GUIDE.md
├── API_REFERENCE.md
├── TODO.md
├── CLEANUP_SUMMARY.md
├── nginx/
│   └── nginx.conf
└── src/
    ├── .env                    ← dibuat dari .env.example (tidak di-commit)
    ├── .env.example
    ├── .env.production
    ├── docker-compose.yml
    ├── docker-compose.external-db.yml
    ├── docker-compose.persistent.yml
    ├── backend/
    │   ├── Dockerfile
    │   ├── main.go
    │   ├── go.mod
    │   ├── go.sum
    │   ├── wait-for-db.sh      ← tidak dipakai (digantikan healthcheck compose)
    │   ├── config/
    │   ├── controllers/
    │   ├── middleware/
    │   ├── models/
    │   ├── repository/
    │   ├── routes/
    │   ├── services/
    │   └── utils/
    ├── rust-engine/
    │   ├── Cargo.toml
    │   ├── Cargo.lock          ← generate dengan: cargo generate-lockfile
    │   ├── Dockerfile
    │   └── src/
    ├── db/
    │   └── schema_mysql.sql    ← schema lengkap (ganti 4 file SQL lama)
    ├── nginx/
    │   └── nginx.conf
    ├── static/
    │   ├── assets/
    │   │   ├── script.js
    │   │   └── style.css
    │   └── public/             ← avatar images
    ├── index.html
    ├── login.html
    ├── register.html
    ├── summary.html
    ├── kaizen.html
    └── techguide.html
```

---

**Updated**: March 2026  
**Repository**: teamitmivhs/work-order