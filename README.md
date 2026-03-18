# IT Work Order System – MIVHS

A work order and helpdesk system for **TEAM IT MIVHS**, designed to provide a simple, fast, and user-friendly interface for managing IT requests and devices at SMK MITRA INDUSTRI MM2100.

## 👥 About the Project

Built for internal use by the **IT MIVHS Team**. A complete **full-stack microservices-based work order management system** with:

- ✅ **Frontend UI** (HTML, TailwindCSS, Vanilla JS)
- ✅ **Backend API** (Go + Gin Framework)
- ✅ **Rust Time Tracker Service** (Axum Framework)
- ✅ **MySQL Database** with health checks
- ✅ **Nginx Reverse Proxy**
- ✅ **Multi-service Docker Compose orchestration**
- ✅ **JWT Authentication** with role-based access (Admin / Operator)
- ✅ **Real-time status tracking**
- ✅ **Team member management**
- ✅ **Safety checklist system**
- ✅ **Performance evaluation (Kaizen)**
- ✅ **Audit trail & history tracking**
- ✅ **Evaluation notes per work order**
- ✅ **Responsive design** (mobile-friendly)
- ✅ **Rate limiting** on auth endpoints
- ✅ **Internal API key** between Go backend and Rust service

---

## 🚀 Getting Started

### Prerequisites
- Docker Engine 20.10+
- Docker Compose v2+
- Port 80 and 443 available on host

### 1. Clone Repository
```bash
git clone https://github.com/teamitmivhs/work-order.git
cd work-order/src
```

### 2. Create Environment File
```bash
cp .env.example .env
```

Edit `.env` and fill in all required values:

```env
DB_HOST=db
DB_PORT=3306
DB_USER=adminit2025
DB_PASSWORD=your_secure_password
DB_NAME=dbwoit
MYSQL_ROOT_PASSWORD=your_root_password
JWT_SECRET=your_random_secret_min_32_chars
INTERNAL_API_KEY=your_random_hex_key
```

Generate secure keys:
```bash
# JWT_SECRET
openssl rand -base64 32

# INTERNAL_API_KEY
openssl rand -hex 32
```

### 3. Run with Docker Compose
```bash
docker compose up -d --build
```

### 4. Access the Application
| Page | URL |
|------|-----|
| Dashboard | http://localhost |
| Login | http://localhost/login.html |
| Register | http://localhost/register.html |
| Summary | http://localhost/summary.html |
| Kaizen | http://localhost/kaizen.html |
| TechGuide | http://localhost/techguide.html |

---

## 🏗️ Architecture

```
Browser
  ↓
┌─────────────────────────────┐
│     Nginx (Port 80/443)     │
│  Static files + API proxy   │
└──────────────┬──────────────┘
               │
               ▼
┌──────────────────────────────┐
│    Go Backend (Port 8080)    │
│   REST API + JWT Auth        │
└──────┬───────────────┬───────┘
       │               │
       ▼               ▼
┌────────────┐  ┌──────────────────┐
│ MySQL 8.0  │  │ Rust Time Tracker│
│ (Port 3306)│  │   (Port 9000)    │
└────────────┘  └──────────────────┘
```

---

## 📌 Features

### Core Functionality
- **Work Order Management** — Create, assign, track, and complete IT work orders
- **Team Member Management** — Monitor technician availability and status (Stand By, On Job, Support, Next Shift, Off Duty)
- **Safety Checklist System** — Location-specific safety checklists enforced before work begins
- **Automatic Time Tracking** — Working hours calculated by Rust service and saved to database
- **Evaluation Notes** — Add and view notes per completed work order
- **Performance Analytics (Kaizen)** — Completion rate, average time, and documentation quality metrics

### Security & Auth
- **JWT Authentication** — 24-hour token expiry
- **Role-based Access** — Admin sees all orders; Operator sees only assigned orders
- **Rate Limiting** — 10 requests/minute per IP on `/login` and `/register`
- **Internal API Key** — Shared secret between Go backend and Rust time tracker
- **Password Policy** — Minimum 8 chars, must include uppercase, lowercase, and digit

### Technical
- **Microservices Architecture** — Go, Rust, MySQL, Nginx as separate containers
- **Docker health checks** — MySQL and Rust service health-checked before backend starts
- **Retry logic** — Backend retries DB connection up to 30 times on startup
- **Empty array responses** — API always returns `[]` not `null` for list endpoints
- **Consistent response format** — All responses use `{ code, message, data }` structure
- **ON DELETE CASCADE** — Child records (executors, checklist) auto-deleted with parent order

---

## 🔄 Workflow

### 1. Create Work Order
```
User → Click "Create Orders" → Fill form (priority, requester, location, device, problem)
     → Submit → Status: "Pending"
```

### 2. Take Order
```
Operator → Click take button → Select standby operators → Review safety checklist
         → Confirm → Status: "On Progress" | Operator status: "On Job"
         → Rust time tracker starts automatically
```

### 3. Complete Order
```
Operator → Click done button → Status: "Completed" | Operator status: "Stand By"
         → Rust time tracker stops → Working hours saved to database
```

### 4. Review & Evaluate
```
Admin → Summary page → View history → Add evaluation notes
      → Kaizen page → View completion rate, avg time, notes quality
```

---

## 🐳 Docker Services

| Service | Image/Source | Port | Purpose |
|---------|-------------|------|---------|
| `db` | `mysql:8.0` | 3306 (internal) | Primary database |
| `time-tracker` | `./rust-engine` | 9000 (internal) | Work duration tracking |
| `backend` | `./backend` | 8080 (internal) | REST API server |
| `nginx` | `nginx:1.25-alpine` | 80, 443 | Reverse proxy + static files |

### Startup Order
```
db (healthy) ──┐
               ├──► backend ──► nginx
time-tracker ──┘
   (healthy)
```

### Available Compose Files
| File | Use Case |
|------|----------|
| `docker-compose.yml` | Default — MySQL in Docker (named volume) |
| `docker-compose.persistent.yml` | MySQL in Docker (host path `/docker/work-order/mysql-data`) |
| `docker-compose.external-db.yml` | MySQL external (set `DB_HOST` in `.env`) |

---

## 📊 Database Schema

| Table | Description |
|-------|-------------|
| `members` | Team member accounts, roles, and status |
| `orders` | Work order records with priority, status, timestamps |
| `executors` | Many-to-many: orders ↔ members assignment |
| `safetychecklist` | Safety checklist items per order |

**Key schema notes:**
- `executors` uses columns `order_id` / `member_id` (not `ID` / `Executors`)
- `safetychecklist` uses column `order_id` (not `ID`)
- `orders.completed_at` is `VARCHAR(20)` storing display string e.g. `"14:30"`
- `orders.notes` column added for evaluation notes
- `members.password` column required for authentication
- Both child tables have `ON DELETE CASCADE` to parent `orders`

---

## 🔐 Authentication

```
POST /api/register   → Create account (public)
POST /api/login      → Get JWT token (public, rate limited)
POST /api/logout     → Logout (protected)
GET  /api/profile    → Get current user (protected)
```

Token usage:
```javascript
fetch('/api/workorders', {
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('userToken')}`,
    'Content-Type': 'application/json'
  }
})
```

**Roles:**
- `Admin` — Full access, can delete orders, sees all work orders
- `Operator` — Can take/complete orders, sees only assigned orders

---

## 🔧 Development (without Docker)

### Backend
```bash
cd src/backend

# Set environment variables (Windows PowerShell)
$env:DB_HOST="localhost"; $env:DB_PORT="3306"
$env:DB_USER="adminit2025"; $env:DB_PASSWORD="your_password"
$env:DB_NAME="dbwoit"; $env:JWT_SECRET="dev-secret"
$env:INTERNAL_API_KEY="dev-key"

go mod tidy
go run main.go
```

### Rust Time Tracker
```bash
cd src/rust-engine

# Generate Cargo.lock (only needed once)
cargo generate-lockfile

cargo run
```

### After Docker build — save Cargo.lock
```bash
docker cp work-order-time-tracker:/app/Cargo.lock ./src/rust-engine/Cargo.lock
```

---

## 🔍 Monitoring

```bash
# Check all service status
docker compose ps

# View logs
docker compose logs -f backend
docker compose logs -f time-tracker
docker compose logs -f db

# Restart a service
docker compose restart backend

# Full reset (deletes database data)
docker compose down -v && docker compose up -d --build
```

---

## 🚀 Future Enhancements

- [ ] Email / Telegram notifications for new orders
- [ ] Advanced reporting & analytics export (PDF/Excel)
- [ ] Mobile application (iOS/Android)
- [ ] Real-time notifications via WebSocket
- [ ] Production SSL/TLS setup via Let's Encrypt
- [ ] Redis-based rate limiting (for multi-instance deployment)
- [ ] Timer state persistence (survive Rust service restart)

---

## 📝 License

MIT License — see [LICENSE](LICENSE) for details.

---

**Developed with ❤️ by IT MIVHS Team**  
*Empowering efficient IT work order management*