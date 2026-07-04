# Deployment Guide

Dokumen ini adalah panduan deployment production untuk project Work Order ITMIVHS dengan target:

- low maintenance
- data aman saat server restart atau container rebuild
- storage tetap terkendali untuk audit jangka panjang
- mudah dipahami orang berikutnya walaupun bukan developer utama

Panduan ini diasumsikan untuk server Ubuntu 24.04 dengan Podman, resource terbatas, dan disk data besar tersedia di `/mnt/disk3`.

## 1. Prinsip Production

Untuk project ini, prioritas production bukan fitur baru. Prioritasnya adalah:

1. data order dan foto tidak hilang saat restart
2. service bisa hidup lagi dengan satu prosedur sederhana
3. backup berjalan otomatis
4. storage tidak penuh diam-diam
5. orang lain tetap bisa mengoperasikan server tanpa harus paham seluruh codebase

Karena itu, deployment production harus memakai:

- persistent storage host
- backup otomatis
- ukuran upload yang dikendalikan
- satu entry point operasional yang konsisten

## 2. Arsitektur yang Dipakai

Komponen aplikasi:

- `nginx`: menyajikan frontend statis dan reverse proxy `/api`
- `backend`: service Go
- `db`: MySQL
- `time-tracker`: service Rust internal

Jangan tambahkan service baru kecuali benar-benar perlu. Semakin sedikit moving parts, semakin rendah biaya perawatan.

## 3. Direktori Production yang Disarankan

Jangan simpan data penting di filesystem container. Semua data audit harus dipasang ke host path yang jelas.

Gunakan struktur ini:

```text
/opt/work-order/                    -> source code repo
/mnt/disk3/work-order/mysql-data/   -> data MySQL
/mnt/disk3/work-order/static/       -> avatar + foto work order
/mnt/disk3/work-order/backups/      -> hasil backup
/mnt/disk3/work-order/logs/         -> log backup atau log operasional tambahan
```

Buat direktori:

```bash
sudo mkdir -p /opt/work-order
sudo mkdir -p /mnt/disk3/work-order/mysql-data
sudo mkdir -p /mnt/disk3/work-order/static
sudo mkdir -p /mnt/disk3/work-order/backups
sudo mkdir -p /mnt/disk3/work-order/logs
sudo chown -R $USER:$USER /opt/work-order
sudo chown -R $USER:$USER /mnt/disk3/work-order
```

## 4. Clone Repo

```bash
cd /opt
git clone https://github.com/teamitmivhs/work-order.git
cd /opt/work-order/src
```

## 5. File `.env` Production

Buat `src/.env`:

```env
DB_USER=workorder_prod
DB_PASSWORD=ganti_dengan_password_panjang
DB_NAME=dbwoit
MYSQL_ROOT_PASSWORD=ganti_dengan_root_password_panjang
JWT_SECRET=ganti_dengan_secret_random_minimal_32_karakter
INTERNAL_API_KEY=ganti_dengan_secret_internal_random
PUBLIC_UPLOAD_DIR=/static/public
```

Generate secret:

```bash
openssl rand -base64 48
```

Gunakan hasil random untuk:

- `JWT_SECRET`
- `INTERNAL_API_KEY`
- password database

## 6. Compose Production yang Disarankan

File `docker-compose.persistent.yml` di repo saat ini masih memakai path `/docker/...`, dan di server kamu itu sebelumnya gagal karena permission.

Untuk production, ubah mount persistence ke:

```yaml
db:
  volumes:
    - /mnt/disk3/work-order/mysql-data:/var/lib/mysql
    - ./db:/docker-entrypoint-initdb.d:ro,Z

backend:
  volumes:
    - /mnt/disk3/work-order/static:/static:z

nginx:
  volumes:
    - /mnt/disk3/work-order/static:/usr/share/nginx/html/static:ro,z
```

Catatan:

- backend menulis file upload ke `/static/public`
- nginx membaca file yang sama dari mount host yang sama
- dengan model ini, foto tetap aman walaupun container backend/nginx di-rebuild
- `db/complete_db.sql` berisi user dummy untuk simulasi lokal. Untuk production, jangan import seed user dummy tersebut. Pakai migration/schema saja, lalu buat user asli lewat register + approval atau import seed private di server.

## 7. Menjalankan Production

Di `src/`:

```bash
podman compose -f docker-compose.persistent.yml up -d --build
```

Cek status:

```bash
podman ps
curl -I http://localhost:4323
curl -I http://localhost:4323/login
curl -I http://localhost:4323/index
```

Kalau ingin expose langsung ke jaringan lokal:

```text
http://IP-SERVER:4323
```

Kalau nanti ingin final production yang lebih rapi, pindahkan nginx ke port `80` dan `443`.

## 8. URL Convention

Seluruh navigasi browser harus memakai clean URL:

- `/index`
- `/login`
- `/register`
- `/guest`
- `/summary`
- `/kaizen`
- `/techguide`
- `/shift`
- `/settings`
- `/staff`

Redirect `.html` tetap boleh hidup di nginx untuk kompatibilitas, tetapi link internal baru jangan kembali memakai `.html`.

## 9. Restart dan Update Harian

Operasional minimal:

Start ulang semua service:

```bash
cd /opt/work-order/src
podman compose -f docker-compose.persistent.yml up -d
```

Build ulang setelah update code:

```bash
cd /opt/work-order
git pull
cd src
podman compose -f docker-compose.persistent.yml up -d --build
```

Restart nginx saja:

```bash
podman restart work-order-nginx
```

Restart backend saja:

```bash
podman compose -f docker-compose.persistent.yml up -d --build backend
```

## 10. Backup yang Wajib

Backup minimum:

- dump MySQL harian
- arsip folder upload harian atau mingguan
- rotasi backup otomatis

Target retensi:

- backup harian: simpan 7 hari
- backup mingguan: simpan 4 minggu
- backup bulanan: simpan 6 bulan bila storage cukup

### 10.1 Backup Database

Contoh manual:

```bash
mkdir -p /mnt/disk3/work-order/backups/db
podman exec work-order-db \
  mysqldump -uworkorder_prod -p'ganti_dengan_password_panjang' dbwoit \
  > /mnt/disk3/work-order/backups/db/dbwoit-$(date +%F-%H%M).sql
```

### 10.2 Backup Upload Foto

Contoh manual:

```bash
mkdir -p /mnt/disk3/work-order/backups/files
tar -czf /mnt/disk3/work-order/backups/files/static-$(date +%F-%H%M).tar.gz \
  -C /mnt/disk3/work-order static
```

## 11. Script Backup Otomatis

Simpan sebagai `/opt/work-order/scripts/backup-work-order.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

BACKUP_ROOT="/mnt/disk3/work-order/backups"
DB_DIR="$BACKUP_ROOT/db"
FILES_DIR="$BACKUP_ROOT/files"
DATE_STAMP="$(date +%F-%H%M)"

mkdir -p "$DB_DIR" "$FILES_DIR"

podman exec work-order-db \
  mysqldump -uworkorder_prod -p'ganti_dengan_password_panjang' dbwoit \
  > "$DB_DIR/dbwoit-$DATE_STAMP.sql"

tar -czf "$FILES_DIR/static-$DATE_STAMP.tar.gz" \
  -C /mnt/disk3/work-order static

find "$DB_DIR" -type f -mtime +7 -delete
find "$FILES_DIR" -type f -mtime +14 -delete
```

Beri permission:

```bash
chmod +x /opt/work-order/scripts/backup-work-order.sh
```

Jalankan via cron:

```bash
crontab -e
```

Tambahkan:

```cron
30 1 * * * /opt/work-order/scripts/backup-work-order.sh >> /mnt/disk3/work-order/logs/backup.log 2>&1
```

## 12. Uji Restore

Backup tanpa uji restore belum cukup.

Minimal sebulan sekali:

1. restore dump database ke database test
2. ekstrak backup folder foto ke lokasi test
3. pastikan order dan dokumentasi bisa terbaca

Contoh restore database test:

```bash
podman exec -i work-order-db \
  mysql -uworkorder_prod -p'ganti_dengan_password_panjang' dbwoit \
  < /mnt/disk3/work-order/backups/db/nama-file.sql
```

## 13. Storage Management

Karena sistem menyimpan foto audit jangka panjang, storage policy harus jelas.

Aturan yang disarankan:

- avatar boleh diganti dan tidak perlu histori
- foto dokumentasi work order dianggap arsip
- ukuran upload harus dibatasi
- foto harus dikompres sebelum atau saat upload
- jangan simpan file mentah berukuran kamera asli

Target praktis:

- foto bukti maksimal sekitar `1280px`
- format JPEG atau WebP
- ukuran file final per foto sebaiknya di bawah `500KB` sampai `1.5MB`

Kalau ini tidak dilakukan, disk akan habis oleh foto, bukan oleh data order.

## 14. Healthcheck Operasional

Yang perlu dicek saat ada masalah:

```bash
podman ps
podman logs work-order-backend --tail 100
podman logs work-order-nginx --tail 100
podman logs work-order-db --tail 100
curl -I http://localhost:4323
curl -I http://localhost:4323/api/health
```

Kalau route health API berbeda, sesuaikan dengan endpoint backend yang aktif.

## 15. Recovery Plan Singkat

Kalau server restart:

```bash
cd /opt/work-order/src
podman compose -f docker-compose.persistent.yml up -d
```

Kalau backend rusak setelah update:

```bash
cd /opt/work-order
git log --oneline -n 10
git checkout <commit-sebelumnya-yang-stabil>
cd src
podman compose -f docker-compose.persistent.yml up -d --build
```

Jangan lakukan `git reset --hard` di server production tanpa alasan yang jelas.

## 16. Low-Maintenance Checklist

Project ini baru layak ditinggal kalau semua poin ini beres:

- compose production memakai host path persistent
- upload foto masuk ke storage host, bukan storage container sementara
- backup DB otomatis aktif
- backup file otomatis aktif
- restore test pernah dicoba
- ada satu dokumen runbook seperti file ini
- ada satu orang selain developer utama yang pernah mengikuti prosedur restart

## 17. Rekomendasi Lanjutan

Urutan kerja yang paling masuk akal:

1. ubah `docker-compose.persistent.yml` agar memakai `/mnt/disk3/work-order/...`
2. pastikan backend dan nginx benar-benar memakai mount static host yang sama
3. buat script backup otomatis di repo
4. tambahkan limit dan kompresi foto upload
5. uji restore dari backup

Kalau lima hal ini selesai, project akan jauh lebih aman untuk resource server yang pas-pasan dan jauh lebih rendah biaya perawatannya.
