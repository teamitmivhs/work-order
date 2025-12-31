# 🚀 Deployment Guide - Work Order System

## Prerequisites

- Docker & Docker Compose installed
- MySQL 8.0+ (untuk external-db option)
- Linux/Unix server (atau WSL untuk Windows)

---

## 📋 Deployment Options

### **OPTION 1: MySQL External (RECOMMENDED untuk Production)**

Gunakan MySQL yang sudah berjalan di server (bukan dalam container).

#### Setup:

1. **Install MySQL di server (jika belum ada)**
   ```bash
   sudo apt-get update
   sudo apt-get install mysql-server
   ```

2. **Setup database**
   ```bash
   mysql -u root -p < src/db/dbwoit_members.sql
   mysql -u root -p < src/db/dbwoit_orders.sql
   mysql -u root -p < src/db/dbwoit_executors.sql
   mysql -u root -p < src/db/dbwoit_safetychecklist.sql
   ```

3. **Copy .env dari .env.example**
   ```bash
   cp src/.env.example src/.env
   ```

4. **Edit src/.env untuk point ke MySQL external**
   ```bash
   DB_HOST=localhost  # atau IP server MySQL
   DB_PORT=3306
   DB_USER=adminit2025
   DB_PASSWORD=databaseit2045
   DB_NAME=dbwoit
   JWT_SECRET=your-strong-secret-key-here
   ```

5. **Start Docker containers**
   ```bash
   cd src
   docker-compose -f docker-compose.external-db.yml up -d
   ```

6. **Verify services**
   ```bash
   docker-compose -f docker-compose.external-db.yml ps
   ```

**Keuntungan:**
✅ Database persist terlepas dari container state
✅ Mudah backup & manage database
✅ Mudah migrate ke server lain
✅ Performance lebih baik (tidak ada virtualisasi DB)

---

### **OPTION 2: MySQL dalam Docker dengan Persistent Storage**

MySQL berjalan dalam container tapi data disimpan di host volume.

#### Setup:

1. **Buat directory untuk persistent storage**
   ```bash
   sudo mkdir -p /docker/work-order/mysql-data
   sudo chown 999:999 /docker/work-order/mysql-data
   chmod 700 /docker/work-order/mysql-data
   ```

2. **Copy .env**
   ```bash
   cp src/.env.example src/.env
   ```

3. **Start containers dengan persistent DB**
   ```bash
   cd src
   docker-compose -f docker-compose.persistent.yml up -d
   ```

4. **Verify database initialized**
   ```bash
   docker exec work-order-db mysql -u adminit2025 -pdatabaseit2045 dbwoit -e "SHOW TABLES;"
   ```

**Keuntungan:**
✅ Self-contained (semua dalam Docker)
✅ Mudah deployment ke server manapun
✅ Data persist di `/docker/work-order/mysql-data`

**Catatan:** Pastikan directory `/docker/work-order/mysql-data` exist dan accessible!

---

## 🔧 Environment Variables

Buat file `.env` di `src/` directory:

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=adminit2025
DB_PASSWORD=databaseit2045
DB_NAME=dbwoit

# Security (CHANGE THIS!)
JWT_SECRET=your-super-secret-key-generate-random-string-here

# Services
TIME_TRACKER_URL=http://time-tracker:9000
```

---

## 📊 Services & Ports

```
┌─────────────────────────────────────────────┐
│  Nginx (Reverse Proxy)                      │
│  Port: 80 (HTTP), 443 (HTTPS)              │
└──────────────┬──────────────────────────────┘
               │
        ┌──────┴──────┐
        ▼              ▼
   Backend        Time Tracker
   Port: 8080     Port: 9000
   (Go API)       (Rust)
        │
        ▼
   MySQL Database
   Port: 3306
```

---

## 🚀 Commands

### Start Services
```bash
cd src

# Option 1: External MySQL
docker-compose -f docker-compose.external-db.yml up -d

# Option 2: MySQL dalam Docker
docker-compose -f docker-compose.persistent.yml up -d
```

### Stop Services
```bash
docker-compose down
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f nginx
docker-compose logs -f db
```

### Database Commands
```bash
# Access MySQL container
docker exec -it work-order-db mysql -u adminit2025 -pdatabaseit2045 dbwoit

# Backup database
docker exec work-order-db mysqldump -u adminit2025 -pdatabaseit2045 dbwoit > backup-$(date +%Y%m%d_%H%M%S).sql

# Restore database
docker exec -i work-order-db mysql -u adminit2025 -pdatabaseit2045 dbwoit < backup.sql
```

### Restart Services
```bash
docker-compose restart

# Specific service
docker-compose restart backend
```

---

## 🔍 Testing API Endpoints

### Register User
```bash
curl -X POST http://localhost/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "testuser",
    "password": "TestPassword123"
  }'
```

Response:
```json
{
  "message": "Registration successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "member": {
      "id": 1,
      "name": "testuser",
      "role": "Operator",
      "status": "standby"
    }
  }
}
```

### Login
```bash
curl -X POST http://localhost/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "name": "testuser",
    "password": "TestPassword123"
  }'
```

### Get Profile (with Auth)
```bash
TOKEN="your-jwt-token-from-login"
curl -X GET http://localhost/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

### Create Work Order
```bash
curl -X POST http://localhost/api/workorders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "priority": "high",
    "requester": "John Doe",
    "location": "Office Room 101",
    "device": "Laptop",
    "problem": "Screen tidak menyala"
  }'
```

### Get All Work Orders
```bash
curl -X GET http://localhost/api/workorders \
  -H "Authorization: Bearer $TOKEN"
```

### Get Kaizen Metrics
```bash
curl -X GET http://localhost/api/kaizen \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🔒 Security Tips

1. **Change JWT_SECRET** di `.env` file - jangan gunakan default!
   ```bash
   # Generate random secret
   openssl rand -base64 32
   ```

2. **Change database password** - jangan gunakan default `databaseit2045`

3. **Enable HTTPS** - setup SSL certificate menggunakan Let's Encrypt
   ```bash
   # Update nginx config untuk HTTPS
   ```

4. **Set strong MySQL password**
   ```bash
   ALTER USER 'adminit2025'@'localhost' IDENTIFIED BY 'new-strong-password';
   ```

5. **Backup database regularly**
   ```bash
   # Daily backup script
   0 2 * * * /home/user/backup-db.sh
   ```

---

## 🛠️ Troubleshooting

### Backend tidak bisa connect ke MySQL

**Symptom:** Backend continuously restarting atau error log shows "database connection failed"

**Solution:**
```bash
# Check if MySQL is running
docker ps | grep mysql

# Check backend logs
docker logs work-order-backend

# Verify MySQL is accessible
docker exec work-order-backend mysql -h db -u adminit2025 -pdatabaseit2045 -e "SELECT 1"

# Verify environment variables
docker inspect work-order-backend | grep -A 20 "Env"
```

### MySQL data hilang setelah restart

**Solution:**
- Pastikan menggunakan `docker-compose.persistent.yml`
- Verify volume mount: `docker inspect work-order-db | grep Mounts`
- Check folder permissions: `ls -la /docker/work-order/mysql-data`

### Cannot access frontend

**Solution:**
```bash
# Check nginx logs
docker logs work-order-nginx

# Verify nginx config
docker exec work-order-nginx cat /etc/nginx/conf.d/default.conf

# Check if ports are exposed
netstat -tuln | grep 80
```

---

## 📦 Backup & Recovery

### Backup Database
```bash
#!/bin/bash
# save as: backup-db.sh

BACKUP_DIR="/home/user/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/work-order-db-$DATE.sql"

mkdir -p $BACKUP_DIR

docker exec work-order-db mysqldump \
  -u adminit2025 \
  -pdatabaseit2045 \
  dbwoit > $BACKUP_FILE

# Keep only last 7 days of backups
find $BACKUP_DIR -name "work-order-db-*.sql" -mtime +7 -delete

echo "Backup completed: $BACKUP_FILE"
```

### Restore Database
```bash
docker exec -i work-order-db mysql \
  -u adminit2025 \
  -pdatabaseit2045 \
  dbwoit < /path/to/backup.sql
```

---

## 📈 Performance Optimization

1. **Enable MySQL Query Cache** (jika MySQL < 8.0)
   ```sql
   SET GLOBAL query_cache_type = ON;
   SET GLOBAL query_cache_size = 100MB;
   ```

2. **Add indexes** untuk frequently queried columns
   ```sql
   CREATE INDEX idx_orders_status ON orders(Status);
   CREATE INDEX idx_executors_member ON executors(Executors);
   ```

3. **Use connection pooling** di backend (optional future enhancement)

4. **Enable Redis caching** (optional future enhancement)

---

## ✅ Production Checklist

- [ ] Environment variables di `.env` sudah di-set dengan nilai production
- [ ] JWT_SECRET sudah di-change ke random string yang panjang
- [ ] Database password sudah di-change dari default
- [ ] HTTPS/SSL sudah di-setup
- [ ] Backup strategy sudah di-implement
- [ ] Monitoring & alerting sudah di-setup
- [ ] Log aggregation sudah di-setup
- [ ] Database indexes sudah di-create
- [ ] CORS settings sudah di-restrict (jangan "*")
- [ ] Rate limiting sudah di-implement

---

## 📞 Support

Jika ada issue:
1. Check logs: `docker-compose logs -f`
2. Verify configuration di `.env`
3. Ensure database initialized: `docker exec work-order-db mysql -u adminit2025 -pdatabaseit2045 dbwoit -e "SHOW TABLES;"`
4. Restart services: `docker-compose restart`
