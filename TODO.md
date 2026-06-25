# TODO List — IT Work Order System

## ✅ Completed

### Backend
- [x] Fix frontend JavaScript errors (scoping, undefined variables, duplicate listeners)
- [x] Implement user registration and login (`POST /api/register`, `POST /api/login`)
- [x] Secure API routes with JWT authentication middleware
- [x] Connect login/register pages to backend API
- [x] JWT session management (24-hour expiry)
- [x] Role-based access (Admin sees all orders, Operator sees assigned only)
- [x] Rate limiting on `/login` and `/register` (10 req/min per IP)
- [x] Add `POST /api/logout` endpoint
- [x] Add `PATCH /api/members/:id/status` endpoint
- [x] Add `PATCH /api/workorders/:id` endpoint for updating executors
- [x] Fix `TakeOrder` — remove reversed assignment check logic
- [x] Fix `CompleteOrder` — rows.Close() before tx.Exec (deadlock prevention)
- [x] Fix `DeleteOrder` — reset member status to standby when order deleted
- [x] Fix `UpdateMemberStatus` — now calls API instead of local-only update
- [x] Fix `GetKaizenMetrics` — count `progress` orders in total
- [x] Fix empty array responses (return `[]` not `null`)
- [x] Consistent response format `{ code, message, data }`
- [x] Fix `sql.NullString` → `*string` for JSON compatibility
- [x] Fix `Requester any` → `string` in WorkOrder struct
- [x] Add `Notes` field to work order for evaluation notes
- [x] Integrate `StartTimer` call in `TakeOrder`
- [x] Integrate `StopTimer` call in `CompleteOrder`, save duration to `working_hours`
- [x] Remove hardcoded credentials from db.go (require env vars)
- [x] Fix JWT secret evaluated at runtime (sync.Once)
- [x] Add `UpdateWorkOrderRequest` struct for PATCH endpoint

### Rust Time Tracker
- [x] Fix `std::sync::Mutex` → `tokio::sync::Mutex` (non-blocking async)
- [x] Fix `now()` — return `Result<i64>` instead of silent zero on clock error
- [x] Fix `stop()` — return actual `stopped_at` timestamp (not reconstructed)
- [x] Fix negative duration guard `.max(0)`
- [x] Add `executor_id` validation (not just `work_order_id`)
- [x] Add `list_active()` for all active timers
- [x] Add `GET /timers` endpoint
- [x] Add `X-Internal-Key` authentication on all endpoints
- [x] Add `PORT` env variable support
- [x] Remove unused `chrono` and `uuid` dependencies
- [x] Fix `impl Default` for `AppState`
- [x] Fix `Dockerfile` — copy `Cargo.lock` for reproducible builds

### Database
- [x] Add `Password` column to `members` table
- [x] Fix `Role` column from narrow ENUM to `VARCHAR(50)`
- [x] Rename `executors.ID` → `order_id`, `executors.Executors` → `member_id`
- [x] Rename `safetychecklist.ID` → `order_id`
- [x] Fix `SafetyChecklist VARCHAR(50)` → `VARCHAR(255)`
- [x] Fix `CompletedAt DATETIME` → `VARCHAR(20)` for display string
- [x] Add `Notes TEXT` column to `orders`
- [x] Add `ON DELETE CASCADE` to child tables
- [x] Add FK from `executors` to `members`
- [x] Clean trailing spaces in member data rows 26-34
- [x] Add `UNIQUE` constraint on `members.name`

### Frontend
- [x] Fix welcome banner text color (dark text on dark background)
- [x] Fix footer outside grid (`col-span-3` invalid outside grid)
- [x] Fix dropdown CSS hover vs JS click conflict
- [x] Fix `viewbox` typo → `viewBox` in all SVGs
- [x] Fix mobile menu — proper vertical dropdown
- [x] Fix hardcoded user profile (now loaded from JWT + API)
- [x] Fix `.status-completed` continuous blinking animation
- [x] Fix `summary.html` and `kaizen.html` using `localStorage` → API
- [x] Fix `colspan="10"` → `12` in summary.html table
- [x] Fix Kabel LAN step 3 truncated text
- [x] Fix login modal close via backdrop (classList vs style.display)
- [x] Remove duplicate `DOMContentLoaded` login/register listeners from script.js
- [x] Add loading skeleton for stats cards
- [x] Add `<br>` spacers replaced with CSS margin
- [x] Add JWT token saved to `localStorage` after login
- [x] Fix `script.js` crash on login/register pages (null check)
- [x] Fix API response unwrap `{ code, message, data }` format
- [x] Add `Authorization` header to all protected fetch calls
- [x] Fix GSAP null check before calling
- [x] Fix `logoutBtn` targeting wrong button index

### Docker
- [x] Add `.env` file requirement (no more hardcoded credentials)
- [x] Add `INTERNAL_API_KEY` env var to backend and time-tracker
- [x] Fix healthcheck endpoint from `HEAD /api/members` to no healthcheck (backend)
- [x] Fix `time-tracker` depends_on → `service_healthy`
- [x] Fix `nginx` depends_on → `service_started` (no healthcheck on backend)
- [x] Fix MySQL healthcheck using env vars not hardcoded credentials
- [x] Add `:ro` read-only flag to nginx volume mounts
- [x] Fix Rust Dockerfile — `COPY Cargo.toml Cargo.lock` → optional without Cargo.lock

---

## 🔴 High Priority

- [ ] **Production SSL/TLS** — Setup Let's Encrypt via nginx untuk HTTPS
- [ ] **GIN_MODE=release** — Set di docker-compose environment backend
- [ ] **Cargo.lock di Git** — Setelah build, simpan `Cargo.lock` untuk reproducible builds
  ```bash
  docker cp work-order-time-tracker:/app/Cargo.lock ./src/rust-engine/Cargo.lock
  ```

---

## 🟡 Medium Priority

- [ ] **Backend unit tests** — Controller, repository, dan model tests
- [ ] **Frontend tests** — Test untuk fungsi utama di `script.js`
- [ ] **Redis rate limiting** — Untuk deployment multi-instance
- [ ] **Timer persistence** — Rust timer hilang saat service restart; perlu recovery mechanism
- [ ] **`notes` endpoint** — Dedicated `PATCH /api/workorders/:id/notes` untuk simpan catatan evaluasi (saat ini pakai endpoint umum)
- [ ] **Tailwind build process** — Ganti CDN dengan PostCSS/Tailwind CLI untuk production

---

## 🟢 Future Enhancements

- [ ] Email / Telegram notifications saat order baru dibuat
- [ ] Export laporan ke PDF / Excel
- [ ] Mobile application (iOS / Android)
- [ ] Real-time notifications via WebSocket
- [ ] Advanced role-based permissions (selain Admin/Operator)
- [ ] Integration dengan tools IT eksternal
- [ ] Dashboard statistik lebih detail (grafik, tren)