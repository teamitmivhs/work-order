# Docker Deployment Guide — IT Work Order System

## Prerequisites

- Docker Engine 20.10+
- Docker Compose v2+
- Minimal 2GB RAM tersedia
- Port 80 dan 443 tidak digunakan

---

## Quick Start

### 1. Navigate ke folder src
```bash
cd work-order/src
```

### 2. Buat file `.env`
```bash
cp .env.example .env
```

Edit `.env` dan isi semua nilai:
```env
DB_HOST=db
DB_PORT=3306
DB_USER=adminit2025
DB_PASSWORD=ganti_dengan_password_aman
DB_NAME=dbwoit
MYSQL_ROOT_PASSWORD=ganti_dengan_root_password
JWT_SECRET=ganti_random_string_min_32_karakter
INTERNAL_API_KEY=ganti_random_hex_32_karakter
```

Generate secure keys:
```bash
# Linux/Mac
openssl rand -base64 32   # untuk JWT_SECRET
openssl rand -hex 32      # untuk INTERNAL_API_KEY

# Windows PowerShell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))
-join ((1..32) | ForEach-Object { '{0:x2}' -f (Get-Random -Maximum 256) })
```

### 3. Build dan jalankan
```bash
docker compose up -d --build
```

### 4. Verifikasi semua service running
```bash
docker compose ps
```

Expected output:
```
NAME                           STATUS
work-order-db                  Up (healthy)
work-order-time-tracker        Up (healthy)
work-order-backend             Up
work-order-nginx               Up
```

### 5. Akses aplikasi
| Halaman | URL |
|---------|-----|
| Dashboard | http://localhost |
| Login | http://localhost/login.html |
| Register | http://localhost/register.html |
| Summary | http://localhost/summary.html |
| Kaizen | http://localhost/kaizen.html |
| TechGuide | http://localhost/techguide.html |

---

## Pilihan Docker Compose

Ada 3 file compose sesuai kebutuhan deployment:

### `docker-compose.yml` (Default)
MySQL berjalan di dalam Docker dengan **named volume**.

```bash
docker compose up -d --build
```

Cocok untuk: development lokal, server baru.

---

### `docker-compose.persistent.yml`
MySQL berjalan di dalam Docker, data disimpan di **path host** `/docker/work-order/mysql-data`.

```bash
docker compose -f docker-compose.persistent.yml up -d --build
```

Cocok untuk: server production dimana path backup sudah dikonfigurasi.

---

### `docker-compose.external-db.yml`
MySQL **tidak** dijalankan di Docker — menggunakan MySQL server eksternal.

```bash
# Set DB_HOST di .env ke IP/hostname MySQL eksternal
# Contoh: DB_HOST=192.168.1.100
docker compose -f docker-compose.external-db.yml up -d --build
```

Cocok untuk: server production yang sudah punya MySQL dedicated.

---

## Startup Order

Service dimulai dalam urutan berikut berdasarkan health check:

```
db (healthy) ──┐
               ├──► backend ──► nginx
time-tracker ──┘
   (healthy)
```

- `db` dan `time-tracker` harus **healthy** sebelum `backend` dimulai
- `backend` dimulai setelah container-nya running (tanpa health check)
- `nginx` dimulai setelah `backend` container running

> **Catatan:** Backend punya retry logic — akan mencoba koneksi ke database hingga 30x dengan jeda 2 detik, sehingga tetap aman meski MySQL butuh waktu ekstra untuk siap.

---

## Database Initialization

SQL files di folder `./db/` dijalankan otomatis saat container MySQL **pertama kali** dibuat:

```
src/db/
└── schema_mysql.sql    ← schema lengkap (orders, members, executors, safetychecklist)
```

> Jika container sudah pernah dibuat sebelumnya, SQL init tidak dijalankan ulang. Untuk reset database, gunakan `docker compose down -v`.

---

## Environment Variables

| Variable | Wajib | Deskripsi |
|----------|-------|-----------|
| `DB_USER` | ✅ | Username MySQL |
| `DB_PASSWORD` | ✅ | Password MySQL |
| `JWT_SECRET` | ✅ | Secret key untuk JWT token |
| `INTERNAL_API_KEY` | ✅ | Shared key antara Go backend dan Rust service |
| `MYSQL_ROOT_PASSWORD` | ✅ | Root password MySQL (hanya untuk compose dengan db container) |
| `DB_HOST` | ✅ | Hostname DB (`db` untuk in-Docker, IP untuk external) |
| `DB_PORT` | ❌ | Default `3306` |
| `DB_NAME` | ❌ | Default `dbwoit` |

---

## Perintah Umum

```bash
# Lihat status semua service
docker compose ps

# Lihat log semua service
docker compose logs -f

# Lihat log service tertentu
docker compose logs -f backend
docker compose logs -f db
docker compose logs -f time-tracker

# Restart service tertentu
docker compose restart backend

# Stop semua service (data tetap tersimpan)
docker compose down

# Stop dan hapus semua data (reset database)
docker compose down -v

# Rebuild dan restart setelah perubahan kode
docker compose up -d --build
```

---

## Troubleshooting

### Backend container unhealthy / gagal start

```bash
docker logs work-order-backend
```

Kemungkinan penyebab:
- `.env` belum dibuat atau variable kosong
- MySQL belum selesai init saat backend start (biasanya teratasi otomatis karena retry logic)

---

### Database connection error

```bash
# Cek status MySQL
docker compose ps db

# Cek log MySQL
docker compose logs db

# Restart MySQL lalu backend
docker compose restart db
docker compose restart backend
```

---

### Port 80 already in use

Edit `docker-compose.yml`, ubah port nginx:
```yaml
nginx:
  ports:
    - "8000:80"   # akses via http://localhost:8000
```

---

### Cargo.lock tidak ada (Rust build error)

Jika error `"/Cargo.lock": not found` saat build:

```bash
# Opsi 1: Generate dari lokal (jika Rust terinstall)
cd src/rust-engine
cargo generate-lockfile

# Opsi 2: Ambil dari container setelah build berhasil
docker cp work-order-time-tracker:/app/Cargo.lock ./src/rust-engine/Cargo.lock
```

---

### Reset dan mulai dari awal

```bash
cd src
docker compose down -v
docker compose up -d --build
```

---

## Service Architecture Detail

### Nginx
- Serve file HTML/CSS/JS statis dari volume mount
- Proxy `/api/*` ke Go backend di port 8080
- Port 80 dan 443 exposed ke host

### Go Backend
- REST API dengan Gin framework
- Koneksi ke MySQL dengan retry 30x
- Memanggil Rust time tracker via HTTP dengan `X-Internal-Key` header
- JWT authentication, rate limiting pada auth endpoints

### Rust Time Tracker
- Menyimpan timer aktif di in-memory HashMap
- Endpoint hanya bisa diakses dengan `X-Internal-Key` header
- Health check tersedia di `GET /health`
- Port 9000, tidak exposed ke host (internal only)

### MySQL
- Inisialisasi otomatis dari `./db/*.sql`
- Health check via `mysqladmin ping`
- Data persistent via named volume atau host path

---

## Production Checklist

- [ ] Ganti semua default password di `.env`
- [ ] Generate `JWT_SECRET` dan `INTERNAL_API_KEY` yang kuat
- [ ] Gunakan `docker-compose.external-db.yml` atau `docker-compose.persistent.yml`
- [ ] Setup SSL/TLS di nginx config
- [ ] Setup backup database berkala
- [ ] Batasi akses port 3306 dari luar
- [ ] Set `GIN_MODE=release` di environment backend