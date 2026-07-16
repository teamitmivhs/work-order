# Codex Server Handoff — Work Order ITMIVHS

Tanggal: 2026-06-30
Project: Work Order ITMIVHS
Target deploy: server `it@versatile`

Dokumen ini untuk Codex yang berjalan langsung di server agar bisa melanjutkan deployment dan hardening tanpa perlu menebak konteks.

---

## 1. Spek Server Target

Output `fastfetch` dari user:

- OS: Ubuntu 24.04.4 LTS x86_64
- Host: IBM System x3100 M5
- Kernel: Linux 6.8.0-124-generic
- CPU: Intel Xeon E3-1220 v3, 4 core, 3.50 GHz
- RAM: 7.60 GiB total, sekitar 2.68 GiB terpakai saat dicek
- Swap: 3.67 GiB
- Disk root `/`: 232.64 GiB total, 35.89 GiB terpakai
- Disk tambahan:
  - `/mnt/disk1`: 915.82 GiB
  - `/mnt/disk2`: 915.82 GiB
  - `/mnt/disk3`: 1.79 TiB, 62.49 GiB terpakai
- Local IP: `192.168.104.13/20`

Kesimpulan kapasitas:

- Aman untuk workload internal/sekolah kecil-menengah.
- CPU cukup untuk nginx + Go backend + MySQL + Rust event engine.
- RAM 7.6 GiB cukup, tetapi MySQL perlu dibatasi/tuning agar tidak mengambil RAM berlebihan.
- Risiko utama bukan spek CPU/RAM, melainkan persistence data, backup, HTTPS, secret production, dan path upload.

---

## 2. Service Project

Compose project menjalankan:

- `db`: MySQL 8.0
- `time-tracker`: Rust outbox consumer + SSE, port internal 9000
- `backend`: Go Gin backend, port internal 8080
- `nginx`: static frontend + reverse proxy `/api/`

File compose penting:

- `src/docker-compose.yml`
- `src/docker-compose.persistent.yml`
- `src/docker-compose.external-db.yml`

Rekomendasi deploy server:

- Pakai `src/docker-compose.persistent.yml` sebagai basis production karena MySQL data disimpan di host path.
- Jangan pakai destructive command seperti `docker compose down -v` kecuali user eksplisit meminta dan backup sudah ada.

---

## 3. Path Persistence yang Disarankan

Saat ini `docker-compose.persistent.yml` menggunakan:

```yaml
/docker/work-order/mysql-data:/var/lib/mysql
```

Untuk server ini lebih disarankan pakai disk besar, contoh:

```text
/mnt/disk3/work-order/mysql-data
```

Pertimbangkan juga membuat path upload persistent:

```text
/mnt/disk3/work-order/static-public
```

Lalu mount ke container backend/nginx agar upload avatar dan foto dokumentasi tidak hilang saat repo/container diganti.

Catatan kode:

- Backend sudah punya helper `PUBLIC_UPLOAD_DIR` di `src/backend/utils/uploads.go`.
- Jika deploy memakai path custom, set env:

```env
PUBLIC_UPLOAD_DIR=/static/public
```

dan pastikan volume host untuk upload dimount ke `/static/public` pada backend serta ke `/usr/share/nginx/html/static/public` pada nginx.

---

## 4. Environment Production

Pastikan `.env` production tidak memakai nilai default/lemah.

Wajib set kuat:

```env
DB_USER=...
DB_PASSWORD=...
DB_NAME=dbwoit
MYSQL_ROOT_PASSWORD=...
JWT_SECRET=...
```

Rekomendasi:

- `JWT_SECRET`: minimal 32 byte random.
- Jangan commit `.env` production.
- Pastikan `.gitignore` menjaga `.env` tetap tidak ikut commit.

---

## 5. HTTPS dan Domain

HTTPS sangat disarankan/wajib untuk production karena:

- login memakai JWT,
- ada upload foto,
- mobile browser notification butuh HTTPS atau localhost,
- cookie/storage/token lebih aman di HTTPS.

Pilihan:

- Reverse proxy di host dengan Caddy/Traefik/Nginx + Let's Encrypt jika server punya domain publik.
- Jika hanya LAN/internal, gunakan sertifikat internal/self-signed yang dipercaya device user, atau terminate HTTPS di reverse proxy existing.

Compose saat ini expose:

```yaml
ports:
  - "4323:80"
  - "8443:443"
```

Nginx config saat ini hanya listen port 80 di `src/nginx/nginx.conf`. Jika ingin 8443 aktif langsung dari container, perlu tambah konfigurasi TLS dan mount sertifikat. Alternatif lebih sederhana: expose HTTP internal lalu HTTPS ditangani reverse proxy host.

---

## 6. Backup Minimum

Siapkan backup harian minimal untuk:

1. MySQL database.
2. Folder upload static/public.
3. File `.env` production, disimpan aman dan tidak public.

Contoh objek yang harus dibackup:

```text
/mnt/disk3/work-order/mysql-data
/mnt/disk3/work-order/static-public
```

Lebih baik backup logical MySQL dengan `mysqldump`/`mysqlpump` selain backup raw folder.

Jangan lakukan upgrade schema atau deploy besar tanpa backup.

---

## 7. Tuning MySQL Ringan

Dengan RAM 7.6 GiB:

- Mulai dari `innodb_buffer_pool_size` sekitar `1G`.
- Bisa naik ke `2G` jika server dedicated dan memory masih lega.
- Jangan set terlalu besar karena server juga menjalankan Docker, Go backend, Rust service, nginx, dan OS.

Jika belum ada kebutuhan tuning, deploy dulu lalu pantau:

```bash
docker stats
free -h
df -h
```

---

## 8. Checklist Deploy untuk Codex Server

1. Cek Docker dan Compose tersedia:

```bash
docker --version
docker compose version
```

2. Tentukan lokasi project, misalnya:

```text
/opt/work-order
```

3. Siapkan folder persistent:

```text
/mnt/disk3/work-order/mysql-data
/mnt/disk3/work-order/static-public
```

4. Review dan sesuaikan `src/docker-compose.persistent.yml`:

- MySQL volume ke `/mnt/disk3/work-order/mysql-data`.
- Upload static public ke `/mnt/disk3/work-order/static-public`.
- Jangan expose port database ke publik.
- Pastikan backend dan nginx memakai volume upload yang sama.

5. Buat `.env` production di folder `src/`.

6. Build dan start:

```bash
cd src
docker compose -f docker-compose.persistent.yml up -d --build
```

7. Cek status:

```bash
docker compose -f docker-compose.persistent.yml ps
docker compose -f docker-compose.persistent.yml logs --tail=100
```

8. Cek endpoint:

```bash
curl -I http://localhost:4323/
curl -I http://localhost:4323/api/health
```

9. Setup HTTPS/reverse proxy.

10. Setup backup database dan upload.

---

## 9. Catatan dari Audit Lokal Terakhir

Beberapa perubahan penting yang sudah/ sedang ada di working tree lokal:

- Dark mode dashboard/Kaizen statistic cards diperbaiki agar tidak abu-abu terang.
- Endpoint ganti password asli ditambahkan:
  - `PATCH /api/profile/password`
  - route di `src/backend/routes/user_routes.go`
  - handler di `src/backend/controllers/user_controller.go`
  - repository update password di `src/backend/repository/member_repository.go`
- Timer backend mulai memakai `StartedAt DATETIME`:
  - migration startup di `src/backend/config/db.go`
  - schema dump di `src/db/complete_db.sql`
  - SQL migration di `src/db/migrations/20260625_member_lifecycle.sql`
  - query progress fallback ke `TimeSort` untuk data lama.

Validasi lokal terakhir:

```bash
env GOCACHE=/tmp/work-order-go-cache go test -C src/backend ./...
node --check src/static/assets/script.js
node --check src/static/assets/workorder-sw.js
git --no-pager diff --check
```

Semua lolos saat dokumen ini dibuat.

---

## 10. Hal yang Jangan Dilakukan Tanpa Izin User

- Jangan `docker compose down -v`.
- Jangan hapus volume MySQL.
- Jangan reset database.
- Jangan overwrite `.env` production jika sudah ada.
- Jangan expose MySQL ke internet.
- Jangan deploy tanpa backup jika sudah ada data production.
