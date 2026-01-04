# Quick Start - Docker Setup

## Prerequisites
- Docker Desktop installed & running
- PowerShell / Command Line
- Project folder: `c:\Users\paron\work-order`

---

## Step 1: Build Backend Image

```powershell
cd c:\Users\paron\work-order\src\backend

# Ensure go build succeeded
go build -o work-order-backend.exe .

# If successful, navigate back
cd ..
```

---

## Step 2: Create Named Volume (First Time Only)

```powershell
docker volume create work-order-mysql-data
```

---

## Step 3: Start Docker Compose

```powershell
cd c:\Users\paron\work-order\src

# Start all services (build + run)
docker compose up -d

# View logs
docker compose logs -f

# Or specific service
docker compose logs -f backend
docker compose logs -f db
```

---

## Step 4: Verify Services

```powershell
# Check all containers running
docker compose ps

# Check if MySQL is healthy
docker exec work-order-db mysql -u adminit2025 -pdatabaseit2045 -e "SELECT VERSION();"

# Check backend logs
docker compose logs backend | tail -20
```

---

## Step 5: Access Application

Open browser:
- **Frontend**: http://localhost
- **Backend API**: http://localhost:8080/api
- **Nginx**: http://localhost

---

## Common Issues & Solutions

### 1. "Docker is not running"
```powershell
# Start Docker Desktop from Windows Start Menu
# Or check if Docker service is running
docker ps
```

### 2. "Port 80 already in use"
```powershell
# Find process using port 80
netstat -ano | findstr :80

# Kill it or change port in docker-compose.yml
```

### 3. "Backend can't connect to MySQL"
```powershell
# Check MySQL container
docker ps | findstr mysql

# View MySQL logs
docker compose logs db

# Test connection
docker exec work-order-db mysql -u adminit2025 -pdatabaseit2045 -e "SELECT 1;"
```

### 4. "No such file or directory: db/"
```powershell
# Ensure you're in src directory
cd c:\Users\paron\work-order\src

# Check db folder exists
dir db

# If not, files should be in db/ folder for database initialization
```

---

## Stop Services

```powershell
# Stop containers (keep volumes)
docker compose down

# Stop and remove volumes (WARNING: DELETES DATA)
docker compose down -v
```

---

## View Data

```powershell
# Connect to MySQL in container
docker exec -it work-order-db mysql -u adminit2025 -pdatabaseit2045 -D dbwoit

# List tables
SHOW TABLES;

# View orders
SELECT * FROM orders;

# View members
SELECT * FROM members;

# Exit
EXIT;
```

---

## Update Code & Rebuild

```powershell
cd c:\Users\paron\work-order\src

# Rebuild backend image (after code changes)
docker compose build --no-cache

# Restart services
docker compose up -d

# View new logs
docker compose logs -f backend
```

---

## Production Notes

⚠️ **Before deploying to production**:

1. **Change JWT_SECRET**
   - Edit `src/docker-compose.yml`
   - Change `JWT_SECRET` value to something secure
   
2. **Change Database Password**
   - Edit `src/docker-compose.yml`
   - Change `MYSQL_PASSWORD` to something strong

3. **Enable Backups**
   - Set up regular database backups
   - Use named volumes or bind mounts

4. **Monitor Logs**
   - Set up log aggregation
   - Monitor error rates

5. **Set Resource Limits**
   - Add resource constraints in docker-compose.yml
   - Prevent container from consuming all CPU/memory

---

## Next Steps

✅ Backend is running with:
- JWT Authentication
- Database persistence
- Role-based access control
- Safety checklist system
- Kaizen metrics

📝 TODO:
- [ ] Update frontend to send JWT tokens
- [ ] Test authentication flows
- [ ] Implement WebSocket for real-time updates
- [ ] Integrate Rust time tracker service
- [ ] Setup email notifications

---

## Documentation

- **API Reference**: See `API_REFERENCE.md`
- **Improvements**: See `IMPROVEMENTS.md`
- **Docker Guide**: See `DOCKER_SETUP.md`
- **Deployment**: See `BACKEND_REFACTORING_SUMMARY.md`

