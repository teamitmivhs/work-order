# Docker Setup Guide for Production

## Problem
Backend dalam docker tidak bisa persist data ke MySQL yang berbeda dengan container, sehingga setiap kali container restart, data hilang.

## Solution: 3 Pilihan

### Option A: MySQL External + Docker Backend (RECOMMENDED)
MySQL berjalan di server host OS, Backend berjalan di Docker

**Kelebihan**:
- Data MySQL persisten di server
- Mudah backup dan maintenance
- Backend dalam Docker tetap terpisah

**Kekurangan**:
- Perlu setup MySQL manual di server

---

### Option B: MySQL Docker dengan Named Volume (GOOD)
MySQL berjalan dalam Docker container dengan persistent volume

**Kelebihan**:
- Semua dalam Docker, mudah deploy
- Data persisten di named volume
- Otomatis backup bersama container

**Kekurangan**:
- Volume perlu di-manage manual

---

### Option C: MySQL Docker dengan Bind Mount (OK)
MySQL berjalan dalam Docker, data di-mount ke directory host

**Kelebihan**:
- Data terlihat di filesystem host
- Mudah backup

**Kekurangan**:
- Path dependency, tidak portable

---

## Setup Steps

### Prerequisites
- Docker Desktop installed dan running
- PowerShell (Windows) atau Terminal (Linux/Mac)

---

## OPTION A: External MySQL + Docker Backend

### 1. Setup MySQL (Windows)

**Install MySQL Community Server**:
```powershell
# Option 1: Using Chocolatey (if installed)
choco install mysql

# Option 2: Download from https://dev.mysql.com/downloads/mysql/
```

**Start MySQL Service**:
```powershell
# If using Windows Service
net start MySQL80

# Or using mysql-server command
mysql-server start
```

**Create Database & User**:
```sql
mysql -u root -p

-- Create database
CREATE DATABASE dbwoit CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create user
CREATE USER 'adminit2025'@'localhost' IDENTIFIED BY 'databaseit2045';

-- Grant privileges
GRANT ALL PRIVILEGES ON dbwoit.* TO 'adminit2025'@'localhost';
FLUSH PRIVILEGES;

-- Import initial data (optional)
USE dbwoit;
SOURCE /path/to/src/db/dbwoit_members.sql;
SOURCE /path/to/src/db/dbwoit_orders.sql;
SOURCE /path/to/src/db/dbwoit_executors.sql;
SOURCE /path/to/src/db/dbwoit_safetychecklist.sql;
```

**Verify Connection**:
```powershell
mysql -u adminit2025 -p databaseit2045 -h 127.0.0.1 -e "SELECT VERSION();"
```

---

### 2. Update docker-compose.yml untuk External MySQL

**File**: `src/docker-compose.yml`

```yaml
version: '3.8'

services:
  # REMOVE MySQL service from docker-compose
  # It will connect to external MySQL

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: work-order-backend
    restart: unless-stopped
    environment:
      # Connect ke MySQL di host machine
      DB_HOST: host.docker.internal  # Special hostname untuk access host dari container
      DB_PORT: 3306
      DB_USER: adminit2025
      DB_PASSWORD: databaseit2045
      DB_NAME: dbwoit
      JWT_SECRET: "your-secret-key-min-32-chars"
    ports:
      - "8080:8080"
    networks:
      - workorder-net

  nginx:
    image: nginx:1.25-alpine
    container_name: work-order-nginx
    restart: unless-stopped
    depends_on:
      - backend
    ports:
      - "80:80"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/conf.d/default.conf
      - ./index.html:/usr/share/nginx/html/index.html
      - ./login.html:/usr/share/nginx/html/login.html
      - ./summary.html:/usr/share/nginx/html/summary.html
      - ./kaizen.html:/usr/share/nginx/html/kaizen.html
      - ./techguide.html:/usr/share/nginx/html/techguide.html
      - ./static:/usr/share/nginx/html/static
    networks:
      - workorder-net

networks:
  workorder-net:
    driver: bridge
```

---

### 3. Run Docker

```powershell
cd c:\Users\paron\work-order\src

# Rebuild backend image
docker build -t work-order-backend ./backend

# Start services
docker compose up -d

# Check logs
docker compose logs -f backend

# Check status
docker compose ps
```

---

## OPTION B: MySQL Docker dengan Named Volume (RECOMMENDED)

### 1. Create Named Volume

```powershell
docker volume create work-order-mysql-data
```

### 2. Update docker-compose.yml

**File**: `src/docker-compose.yml`

```yaml
version: '3.8'

services:
  db:
    image: mysql:8.0
    container_name: work-order-db
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: dbwoit
      MYSQL_USER: adminit2025
      MYSQL_PASSWORD: databaseit2045
    
    volumes:
      # Use named volume for persistent storage
      - work-order-mysql-data:/var/lib/mysql
      # Initialize database on first run
      - ./db:/docker-entrypoint-initdb.d
    ports:
      - "3306:3306"  # Expose port untuk debugging
    networks:
      - workorder-net
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "adminit2025", "-padatabaseit2045"]
      timeout: 20s
      retries: 10

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: work-order-backend
    restart: unless-stopped
    depends_on:
      db:
        condition: service_healthy
    environment:
      # Connect ke MySQL container dalam docker network
      DB_HOST: db
      DB_PORT: 3306
      DB_USER: adminit2025
      DB_PASSWORD: databaseit2045
      DB_NAME: dbwoit
      JWT_SECRET: "your-secret-key-min-32-chars"
    ports:
      - "8080:8080"
    networks:
      - workorder-net

  nginx:
    image: nginx:1.25-alpine
    container_name: work-order-nginx
    restart: unless-stopped
    depends_on:
      - backend
    ports:
      - "80:80"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/conf.d/default.conf
      - ./index.html:/usr/share/nginx/html/index.html
      - ./login.html:/usr/share/nginx/html/login.html
      - ./summary.html:/usr/share/nginx/html/summary.html
      - ./kaizen.html:/usr/share/nginx/html/kaizen.html
      - ./techguide.html:/usr/share/nginx/html/techguide.html
      - ./static:/usr/share/nginx/html/static
    networks:
      - workorder-net

# Named volume untuk persistent storage
volumes:
  work-order-mysql-data:
    external: true

networks:
  workorder-net:
    driver: bridge
```

---

### 2. Run Docker

```powershell
cd c:\Users\paron\work-order\src

# Build images
docker compose build

# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Check MySQL is running
docker exec work-order-db mysql -u adminit2025 -pdatabaseit2045 -e "SELECT VERSION();"

# Check backend is running
docker compose logs backend
```

---

## OPTION C: MySQL Docker dengan Bind Mount

```yaml
services:
  db:
    image: mysql:8.0
    container_name: work-order-db
    volumes:
      # Bind mount ke directory di host
      - C:/work-order-mysql-data:/var/lib/mysql  # Windows path
      # Or for Linux/Mac: - /opt/work-order-mysql-data:/var/lib/mysql
      - ./db:/docker-entrypoint-initdb.d
    # ... rest of config
```

---

## Common Commands

```powershell
# Start containers
docker compose up -d

# Stop containers
docker compose down

# View logs
docker compose logs -f backend
docker compose logs -f db

# Connect to MySQL in container
docker exec -it work-order-db mysql -u adminit2025 -pdatabaseit2045 -D dbwoit

# Execute SQL
docker exec work-order-db mysql -u adminit2025 -pdatabaseit2045 -D dbwoit -e "SELECT COUNT(*) FROM orders;"

# Check volume
docker volume ls
docker volume inspect work-order-mysql-data

# Remove all (WARNING: DELETE DATA if using bind mount)
docker compose down -v

# Rebuild without cache
docker compose build --no-cache
docker compose up -d
```

---

## Troubleshooting

### MySQL Connection Failed
```powershell
# Check if container is running
docker ps | findstr mysql

# Check logs
docker compose logs db

# Verify credentials
docker exec work-order-db mysql -u adminit2025 -pdatabaseit2045 -e "SELECT 1;"
```

### Port Already in Use
```powershell
# Find process using port
netstat -ano | findstr :3306

# Kill process
taskkill /PID <PID> /F

# Or change port in docker-compose.yml
# ports:
#   - "3307:3306"  # Map to different port
```

### Data Not Persisting
```powershell
# Check volume
docker volume inspect work-order-mysql-data

# Check if volume is mounted correctly
docker inspect work-order-db | findstr -A 10 Mounts

# Re-create volume
docker volume rm work-order-mysql-data
docker volume create work-order-mysql-data
docker compose up -d
```

### Backend Can't Connect to MySQL
```powershell
# Test from within backend container
docker exec work-order-backend ping db

# Or use DNS
docker exec work-order-backend nslookup db

# Check DNS in container
docker exec work-order-backend cat /etc/resolv.conf
```

---

## Production Checklist

- [ ] Change `JWT_SECRET` environment variable
- [ ] Use strong database password (not default)
- [ ] Enable database backups (scheduled)
- [ ] Set resource limits in docker-compose
- [ ] Use health checks
- [ ] Setup logging aggregation
- [ ] Configure restart policies
- [ ] Use named volumes for data
- [ ] Setup monitoring

---

## Backup & Restore

### Backup MySQL Data

```powershell
# Backup
docker exec work-order-db mysqldump -u adminit2025 -pdatabaseit2045 dbwoit > backup.sql

# Backup entire volume
docker run --rm -v work-order-mysql-data:/data -v ${PWD}:/backup ubuntu tar czf /backup/mysql-backup.tar.gz -C /data .
```

### Restore MySQL Data

```powershell
# Restore from SQL file
docker exec -i work-order-db mysql -u adminit2025 -pdatabaseit2045 dbwoit < backup.sql

# Restore from volume backup
docker run --rm -v work-order-mysql-data:/data -v ${PWD}:/backup ubuntu tar xzf /backup/mysql-backup.tar.gz -C /data
```

