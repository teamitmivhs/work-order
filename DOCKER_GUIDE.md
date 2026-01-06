# IT Work Order System - Docker Deployment Guide

## Quick Start dengan Docker

Panduan lengkap untuk menjalankan IT Work Order System menggunakan Docker Compose.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 1.29+
- Git
- Minimal 2GB RAM tersedia
- Port 80, 3306, 8080 tidak digunakan

## Setup Instructions

### 1. Clone Repository

```bash
git clone https://github.com/teamitmivhs/work-order.git
cd work-order
```

### 2. Navigate to Source Directory

```bash
cd src
```

### 3. Start Docker Compose

```bash
# Build and start all services
docker compose up -d

# View logs untuk verify semua services
docker compose logs -f

# Atau lihat logs per service
docker compose logs -f backend
docker compose logs -f db
docker compose logs -f time-tracker
docker compose logs -f nginx
```

### 4. Verify All Services Running

```bash
# Check status
docker compose ps

# Expected output:
# NAME                           STATUS
# work-order-db                  Up (healthy)
# work-order-time-tracker        Up (healthy)
# work-order-backend             Up
# work-order-nginx               Up
```

### 5. Access Application

Buka browser dan navigasi ke:

- **Main Dashboard**: http://localhost
- **Login Page**: http://localhost/login.html
- **Register Page**: http://localhost/register.html
- **API Backend**: http://localhost/api
- **Time Tracker**: http://localhost:9000 (internal only)

## Login Credentials (Default)

### Test Admin Account
- **Username**: adminit2025
- **Password**: databaseit2045

### Guest Login
- Klik tombol "Login as Guest" di login page
- Dashboard akan membuka dengan popup "Create Orders"
- Guest bisa membuat orders tetapi tidak bisa manage/delete orders (hanya untuk admin)
- Klik tombol "Exit Guest" untuk keluar dari guest mode

## Database Initialization

Database MySQL akan otomatis diinisialisasi saat container pertama kali startup dengan SQL files di folder `db/`:

- `dbwoit_members.sql` - Tabel members/users
- `dbwoit_orders.sql` - Tabel work orders
- `dbwoit_executors.sql` - Tabel executor assignments
- `dbwoit_safetychecklist.sql` - Tabel safety checklist

## Stopping Services

```bash
# Stop all services (keep data)
docker compose down

# Stop dan remove volumes (reset database)
docker compose down -v
```

## Troubleshooting

### Service tidak start

```bash
# Check logs
docker compose logs

# Restart specific service
docker compose restart backend
docker compose restart db
```

### Database connection error

```bash
# Ensure MySQL is healthy
docker compose ps db

# Check MySQL logs
docker compose logs db

# Restart database
docker compose restart db
docker compose restart backend
```

### Port already in use

Edit `src/docker-compose.yml` dan ubah port mappings:

```yaml
nginx:
  ports:
    - "8000:80"  # Changed from 80:80
```

### Clear everything and restart fresh

```bash
cd src
docker compose down -v
docker compose up -d
```

## Service Architecture

```
┌─────────────────────────────────────┐
│         Nginx (Port 80)             │
│    (Frontend + API Proxy)           │
└──────────────────┬──────────────────┘
                   │
        ┌──────────┼──────────┐
        │          │          │
        ▼          ▼          ▼
    ┌─────┐  ┌──────────┐  ┌──────────┐
    │ Go  │  │   Rust   │  │  MySQL   │
    │ API │  │   Time   │  │   8.0    │
    │ 8080│  │  Tracker │  │   3306   │
    │     │  │   9000   │  │          │
    └─────┘  └──────────┘  └──────────┘
```

## Performance Notes

- First startup: ~1-2 menit (building images + init database)
- Subsequent startups: ~30-45 detik
- Memory usage: ~1.5GB saat running
- Disk space required: ~2GB untuk images + volumes

## Production Deployment

Untuk production deployment:

1. Gunakan external database (jangan container)
2. Update `.env` dengan credentials production
3. Implement proper SSL/TLS (nginx config)
4. Setup monitoring dan logging
5. Implement backup strategy untuk database
6. Use strong JWT_SECRET di environment

## Support

Untuk masalah atau pertanyaan, buka issue di GitHub atau hubungi team IT MIVHS.
