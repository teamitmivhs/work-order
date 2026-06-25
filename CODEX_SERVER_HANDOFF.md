# Codex Server Handoff

Dokumen ini dipakai untuk meneruskan perubahan dari Codex lokal/VS Code ke Codex yang berjalan di Ubuntu server setelah perubahan di-push ke Git.

## Server Context

- Production path: `/home/it/work-order`
- Runtime: Docker Compose
- Database: MySQL
- Main stack:
  - Frontend: static HTML, Tailwind CDN, vanilla JavaScript
  - Backend: Go Gin
  - Time tracker: Rust Axum
  - Proxy/static serving: Nginx
  - Storage: MySQL

## Local Codex Responsibilities

Setiap task perubahan harus menghasilkan:

1. Ringkasan perubahan yang jelas.
2. Daftar file yang berubah.
3. Hasil verifikasi lokal yang relevan.
4. Instruksi untuk Codex server setelah `git pull`.
5. Commit message per file dengan format:

```text
[type] message
```

Contoh tipe:

- `[fix]` untuk bug fix
- `[feat]` untuk fitur baru
- `[docs]` untuk dokumentasi
- `[refactor]` untuk perubahan struktur tanpa perubahan behavior
- `[test]` untuk test
- `[chore]` untuk maintenance

## Per-File Commit Message Format

Gunakan satu baris per file:

```text
path/to/file.ext -> [type] message
```

Contoh:

```text
src/backend/controllers/workorder_controller.go -> [fix] restrict work order deletion to admins
src/.env.example -> [docs] add safe environment template
DOCKER_GUIDE.md -> [docs] sync database init file name
```

## Handoff Template

Gunakan template ini setiap selesai mengubah kode:

```markdown
# Server Handoff

## Branch / Commit

- Branch:
- Commit:

## Summary

-

## Changed Files

-

## Verification Done Locally

-

## Server Steps After Git Push

1. `cd /home/it/work-order`
2. `git status --short`
3. `git pull`
4. Run verification commands:
   - `cd src/backend && go test ./...`
   - `cd ../rust-engine && cargo test`
5. If code changes affect Docker runtime, review diff first, then rebuild intentionally:
   - `cd /home/it/work-order/src`
   - `docker compose up -d --build`
6. Check containers:
   - `docker compose ps`
   - `docker compose logs --tail=100 backend`
   - `docker compose logs --tail=100 nginx`

## Production Safety Notes

- Do not print or commit `src/.env`.
- Do not run destructive database reset commands unless explicitly approved.
- Do not run `docker compose down -v` on production unless explicitly approved.
- Review `git diff` before rebuilding production containers.

## Per-File Commit Messages

- `path/to/file` -> `[type] message`
```

## Current Known Follow-Up Candidates

- Restrict `DELETE /api/workorders/:id` to Admin users.
- Add a safe `src/.env.example` because docs mention it.
- Sync Docker docs with the real SQL init file name: `complete_db.sql`.
- Re-check auth rate limiting claim versus actual backend route implementation.

## Latest Handoff: Staff Lifecycle Approval

## Summary

- Register no longer requires a pre-seeded member name.
- New registrations become `pending` staff requests after an active Admin exists.
- If no active Admin exists, the first registered account becomes the bootstrap Admin.
- Admin-only staff lifecycle endpoints were added for pending approval, reject, disable, alumni, and graduate batch.
- Dashboard member API now returns only active, non-guest staff who can handle work orders.
- `staff.html` provides the admin approval page with Pending, Active, Alumni, Disabled, and Rejected tabs.
- Fresh schema and non-destructive migration support member lifecycle fields.
- Docker Compose variants now mount `staff.html`.

## Changed Files

- `src/backend/config/db.go`
- `src/backend/controllers/user_controller.go`
- `src/backend/main.go`
- `src/backend/models/workorder_struct.go`
- `src/backend/repository/member_repository.go`
- `src/backend/routes/user_routes.go`
- `src/db/complete_db.sql`
- `src/db/migrations/20260625_member_lifecycle.sql`
- `src/docker-compose.yml`
- `src/docker-compose.persistent.yml`
- `src/docker-compose.external-db.yml`
- `src/register.html`
- `src/staff.html`

## Verification Done Locally

- `GOCACHE=/tmp/go-build-cache go test ./...` from `src/backend` passed.
- Local backend was rebuilt and started at `http://localhost:4323`.
- Local frontend dev image was rebuilt and started at `http://localhost:4324`.
- `GET /health` returned `{"status":"healthy"}`.
- `GET /staff.html` returned `200 OK`.
- `GET /api/workorders` through frontend returned `{"code":200,"data":[]}`.
- Register bootstrap Admin was verified.
- Register pending staff was verified.
- Admin pending list was verified.
- Admin approve endpoint was verified.
- Approved staff appeared in `GET /api/members`.
- Manual migration SQL was tested successfully against the local MySQL container.

## Server Steps After Git Push

1. `cd /home/it/work-order`
2. `git status --short`
3. `git pull`
4. `cd src/backend && go test ./...`
5. Review DB lifecycle migration:
   - Backend runs non-destructive column migration automatically on startup.
   - Optional manual reference: `src/db/migrations/20260625_member_lifecycle.sql`
6. `cd /home/it/work-order/src`
7. `docker compose up -d --build backend nginx`
8. Check:
   - `docker compose ps`
   - `docker compose logs --tail=100 backend`
   - `curl -I http://localhost:4323/staff.html`

## Production Safety Notes

- Do not reset the MySQL volume.
- Do not run `docker compose down -v`.
- Existing staff are preserved and default to active.
- Alumni/disabled staff are hidden from dashboard handler lists but remain available for work order history references.
- Local compose rebuild hit an unrelated Rust syntax error in `src/rust-engine/src/time_tracker.rs`; fix or exclude that local change before running full `docker compose up -d --build` if it exists on the server branch.

## Per-File Commit Messages

- `src/backend/config/db.go` -> `[feat] add automatic member lifecycle migration`
- `src/backend/controllers/user_controller.go` -> `[feat] add pending registration and admin staff lifecycle actions`
- `src/backend/main.go` -> `[feat] serve staff management page`
- `src/backend/models/workorder_struct.go` -> `[feat] add member lifecycle fields`
- `src/backend/repository/member_repository.go` -> `[feat] add member approval and graduation persistence`
- `src/backend/routes/user_routes.go` -> `[feat] add admin staff management routes`
- `src/db/complete_db.sql` -> `[feat] add member lifecycle columns to fresh schema`
- `src/db/migrations/20260625_member_lifecycle.sql` -> `[feat] add member lifecycle migration script`
- `src/docker-compose.yml` -> `[feat] mount staff management page in default compose`
- `src/docker-compose.persistent.yml` -> `[feat] mount staff management page in persistent compose`
- `src/docker-compose.external-db.yml` -> `[feat] mount staff management page in external database compose`
- `src/register.html` -> `[feat] submit registrations for admin approval`
- `src/staff.html` -> `[feat] add admin staff approval page`

## Latest Handoff: Admin Staff Navbar Link

## Summary

- Added a Staff Management navbar button beside the profile card.
- The navbar button and profile dropdown Staff Management link are hidden by default.
- Both links become visible only when the logged-in JWT role is `Admin`.

## Changed Files

- `src/index.html`

## Verification Done Locally

- Rebuilt and restarted local frontend dev container.
- Verified served `index.html` contains `adminStaffNavLink` and `adminStaffDropdownLink`.
- Verified `GET /staff.html` still returns `200 OK`.

## Server Steps After Git Push

1. `cd /home/it/work-order`
2. `git pull`
3. `cd src`
4. `docker compose up -d --build nginx`
5. Open dashboard as Admin and confirm the `Staff` button appears beside the profile area.

## Per-File Commit Messages

- `src/index.html` -> `[feat] show staff management link for admins`

## Latest Handoff: Staff Page Back Button

## Summary

- Added a Back button to the `staff.html` header.
- The button uses browser history when available and falls back to `index.html`.

## Changed Files

- `src/staff.html`

## Verification Done Locally

- Rebuilt and restarted local frontend dev container.
- Verified served `staff.html` contains `backBtn`.
- Verified `GET /staff.html` returns `200 OK`.

## Per-File Commit Messages

- `src/staff.html` -> `[feat] add back button to staff management page`

## Latest Handoff: Staff Page Mobile Overhaul

## Summary

- Overhauled `staff.html` mobile layout.
- Header actions become a compact mobile grid.
- Status tabs become horizontal pill controls on mobile.
- Desktop keeps the table layout.
- Mobile now renders staff as readable cards with role, batch, lifecycle, work-order handling flag, and full-width actions.
- Approval and graduate modals now use mobile-friendly button stacks and scroll bounds.

## Changed Files

- `src/staff.html`

## Verification Done Locally

- `node --check` passed for the inline staff page script.
- Rebuilt and restarted local frontend dev container.
- Verified `GET /staff.html` returns `200 OK`.
- Verified served `staff.html` contains `staffCards`, `mobile-card`, `header-actions`, and `tab-scroll`.

## Per-File Commit Messages

- `src/staff.html` -> `[feat] overhaul staff management mobile layout`

## Latest Handoff: Global UI Refresh

## Summary

- Added a shared modern UI layer in `src/static/assets/style.css` for all main pages.
- Refreshed page background, panels, cards, tables, inputs, modals, status chips, skeleton loading, hover states, and CSS3 motion.
- Kept the existing slate/blue/white operational palette and status colors.
- Refreshed `login.html` and `register.html` auth cards with matching minimal styling and CSS entry animation.
- Refreshed `settings.html` cards because its inline CSS overrides the shared stylesheet.
- Connected `staff.html` to the shared stylesheet and added staggered row/card animation.

## Changed Files

- `src/static/assets/style.css`
- `src/staff.html`
- `src/login.html`
- `src/register.html`
- `src/settings.html`

## Verification Done Locally

- `node --check` passed for the inline staff page script.
- Rebuilt and restarted local frontend dev container with Podman.
- Verified these pages return `200 OK`: `index.html`, `login.html`, `register.html`, `staff.html`, `settings.html`, `summary.html`, `kaizen.html`, and `techguide.html`.
- Verified served `style.css` contains the new `MODERN UI REFRESH` layer.
- Verified served `staff.html` contains shared stylesheet link, `content-panel`, `table-row`, and `mobile-card`.
- Note: regular sandbox curl was intermittent against the rootless port, but escalated curl and Nginx logs verified the frontend container is serving correctly.

## Server Steps After Git Push

1. `cd /home/it/work-order`
2. `git pull`
3. `cd src`
4. `docker compose up -d --build nginx`
5. Open `index.html`, `staff.html`, `login.html`, `register.html`, `settings.html`, `summary.html`, `kaizen.html`, and `techguide.html` for a quick visual smoke test.

## Per-File Commit Messages

- `src/static/assets/style.css` -> `[feat] add global modern UI refresh styles`
- `src/staff.html` -> `[feat] align staff management page with global UI refresh`
- `src/login.html` -> `[feat] refresh login page visual polish`
- `src/register.html` -> `[feat] refresh registration page visual polish`
- `src/settings.html` -> `[feat] refresh settings page surfaces`
- `CODEX_SERVER_HANDOFF.md` -> `[docs] document global UI refresh handoff`

## Latest Handoff: Split Staff Role And Division

## Summary

- Split staff authorization role from technical division.
- `Role` is now for access level only: `Admin` or `Operator`.
- Added `Division` for technical team grouping: `Soundman`, `Programmer`, `Maintenance`, `Data Analyst`.
- Existing legacy role values such as `programmer`, `soundman`, `maintenance`, and `data analyst` are migrated into `Division`, then `Role` becomes `Operator`.
- Staff approval now submits both `role` and `division`.
- Dashboard operator detail/search now displays division instead of authorization role.
- Settings profile now shows both role and division.

## Changed Files

- `src/backend/config/db.go`
- `src/backend/controllers/user_controller.go`
- `src/backend/models/workorder_struct.go`
- `src/backend/repository/member_repository.go`
- `src/db/complete_db.sql`
- `src/db/migrations/20260625_member_lifecycle.sql`
- `src/staff.html`
- `src/index.html`
- `src/settings.html`
- `src/static/assets/script.js`

## Verification Done Locally

- `GOCACHE=/tmp/go-build-cache go test ./...` passed in `src/backend`.
- `node --check` passed for `src/static/assets/script.js`.
- `node --check` passed for the inline `staff.html` script.
- Rebuilt and restarted the local backend dev container.
- Rebuilt and restarted the local frontend dev container.
- Verified MySQL backfill: legacy divisions now appear in `Division`, with `Role = Operator`.
- Verified approve endpoint accepts payload with `role` and `division`.
- Verified served `staff.html` contains the `Divisi` column/dropdown and sends `division`.
- Verified served `script.js` uses `member.division` in operator search/detail.

## Server Steps After Git Push

1. `cd /home/it/work-order`
2. `git pull`
3. `cd src/backend && go test ./...`
4. `cd /home/it/work-order/src`
5. `docker compose up -d --build backend nginx`
6. Check migration result in MySQL:
   - `SELECT ID, Name, Role, Division FROM members WHERE Role <> 'Guest' LIMIT 10;`
7. Open `staff.html` as Admin and confirm approval has separate `Role` and `Divisi` fields.

## Per-File Commit Messages

- `src/backend/config/db.go` -> `[feat] migrate member division separately from role`
- `src/backend/controllers/user_controller.go` -> `[feat] approve staff with role and division`
- `src/backend/models/workorder_struct.go` -> `[feat] add member division field`
- `src/backend/repository/member_repository.go` -> `[feat] persist member division data`
- `src/db/complete_db.sql` -> `[feat] split seeded member roles and divisions`
- `src/db/migrations/20260625_member_lifecycle.sql` -> `[feat] add member division migration`
- `src/staff.html` -> `[feat] separate staff role and division controls`
- `src/index.html` -> `[feat] show operator division in profile modal`
- `src/settings.html` -> `[feat] show profile division separately`
- `src/static/assets/script.js` -> `[feat] render member division in operator UI`
- `CODEX_SERVER_HANDOFF.md` -> `[docs] document staff role division split`

## Latest Handoff: Navbar Guest Flow And Staff Role Change

## Summary

- Improved dashboard navbar/profile visibility on desktop and mobile.
- Moved the Staff Management entry into the main navbar as an Admin-only modern pill.
- Restyled the mobile nav menu so it opens as a compact panel instead of a rough full-width row.
- Removed the dashed/shimmer/ripple animation from the Create Order buttons.
- Added `guest.html` as a dedicated guest-only work order submission page.
- Guest login now redirects to `guest.html`, and guests are redirected away from the internal dashboard.
- `GET /api/workorders` is now protected; Admin can list all, Operator can list assigned orders, Guest gets `403`.
- `summary.html` and `kaizen.html` now send JWT auth when fetching work orders and redirect guests to the guest page.
- Admin staff management supports changing active staff role between `Admin` and `Operator`.

## Changed Files

- `src/index.html`
- `src/login.html`
- `src/guest.html`
- `src/static/assets/script.js`
- `src/summary.html`
- `src/kaizen.html`
- `src/backend/main.go`
- `src/backend/routes/workorder_routes.go`
- `src/backend/controllers/workorder_controller.go`
- `src/backend/routes/user_routes.go`
- `src/backend/controllers/user_controller.go`
- `src/backend/repository/member_repository.go`
- `src/staff.html`
- `src/docker-compose.yml`
- `src/docker-compose.persistent.yml`
- `src/docker-compose.external-db.yml`

## Verification Done Locally

- `GOCACHE=/tmp/go-build-cache go test ./...` passed in `src/backend`.
- `node --check src/static/assets/script.js` passed.
- `node --check` passed for extracted inline scripts from `index.html`, `guest.html`, and `staff.html`.
- Rebuilt and restarted local backend dev container at `http://localhost:4323`.
- Rebuilt and restarted local frontend dev container at `http://localhost:4324`.
- Verified `GET /api/workorders` without token returns `401`.
- Verified guest token `GET /api/workorders` returns `403`.
- Verified Admin token `GET /api/workorders` returns `200`.
- Verified guest can create a work order through `POST /api/workorders`; the smoke-test order was deleted afterward.
- Verified Admin can promote an active Operator to `Admin` through `PATCH /api/admin/members/:id/role`.
- Verified Admin can demote that same staff account back to `Operator`; the smoke-test member was deleted afterward.
- Verified served `index.html` contains the new Admin-only staff nav link and guest redirect.
- Verified served `guest.html` contains the guest order form.
- Verified `GET /guest.html` and `GET /staff.html` return `200 OK`.

## Server Steps After Git Push

1. `cd /home/it/work-order`
2. `git status --short`
3. `git pull`
4. `cd src/backend && go test ./...`
5. `cd /home/it/work-order/src`
6. `docker compose up -d --build backend nginx`
7. Check:
   - `docker compose ps`
   - `docker compose logs --tail=100 backend`
   - `curl -I http://localhost/guest.html`
   - `curl -I http://localhost/staff.html`
8. Browser smoke test:
   - Login as Guest and confirm it opens `guest.html`.
   - Confirm Guest cannot open the dashboard order list.
   - Login as Admin and confirm Staff appears in the navbar.
   - Open `staff.html`, pick an active Operator, and confirm role can be changed to Admin.

## Production Safety Notes

- Do not reset the MySQL volume.
- Do not run `docker compose down -v`.
- Ensure the permanent `guest` account exists with role `Guest`, active status, and `CanHandleWorkOrder = 0`.
- `GET /api/workorders` now requires JWT; any external integration that reads work orders must send `Authorization: Bearer <token>`.
- `summary.html` and `kaizen.html` now require an internal logged-in user because they read work order data.

## Per-File Commit Messages

- `src/index.html` -> `[feat] improve dashboard navbar and isolate guest access`
- `src/login.html` -> `[feat] route guest login to guest order page`
- `src/guest.html` -> `[feat] add guest-only work order page`
- `src/static/assets/script.js` -> `[feat] block guest dashboard access and authenticate order fetch`
- `src/summary.html` -> `[fix] authenticate summary work order fetch`
- `src/kaizen.html` -> `[fix] authenticate kaizen work order fetch`
- `src/backend/main.go` -> `[feat] serve guest work order page`
- `src/backend/routes/workorder_routes.go` -> `[fix] protect work order list endpoint`
- `src/backend/controllers/workorder_controller.go` -> `[fix] block guest work order listing`
- `src/backend/routes/user_routes.go` -> `[feat] add admin member role route`
- `src/backend/controllers/user_controller.go` -> `[feat] add admin member role change handler`
- `src/backend/repository/member_repository.go` -> `[feat] persist admin member role changes`
- `src/staff.html` -> `[feat] add active staff role change action`
- `src/docker-compose.yml` -> `[feat] mount guest work order page`
- `src/docker-compose.persistent.yml` -> `[feat] mount guest work order page`
- `src/docker-compose.external-db.yml` -> `[feat] mount guest work order page`
- `CODEX_SERVER_HANDOFF.md` -> `[docs] document navbar guest flow and role change handoff`

## Latest Handoff: Rootless Podman Compose Fix

## Summary

- Fixed local `podman compose up -d` on Fedora/rootless Podman.
- Added SELinux `:Z` labels to bind-mounted DB init, static, Nginx config, and HTML files.
- Changed Nginx host TLS mapping from privileged `443:443` to rootless-safe `8443:443`.
- HTTP local development remains at `http://localhost:4323`.

## Changed Files

- `src/docker-compose.yml`
- `src/docker-compose.persistent.yml`
- `src/docker-compose.external-db.yml`
- `CODEX_SERVER_HANDOFF.md`

## Verification Done Locally

- `systemctl --user start podman.socket` started the user Podman API socket needed by the external compose provider.
- `podman compose up -d` succeeded from `src`.
- Final running containers:
  - `work-order-db` healthy
  - `work-order-time-tracker` healthy
  - `work-order-backend` running
  - `work-order-nginx` running on `4323:80` and `8443:443`
- `curl -I http://127.0.0.1:4323/` returned `200 OK`.
- `curl -i http://127.0.0.1:4323/api/workorders` returned `401 Missing authorization header`, confirming the rebuilt backend auth behavior is active.

## Local Run Notes

1. `cd /home/parothegreat/work-order/src`
2. `systemctl --user start podman.socket`
3. `podman compose up -d`
4. Open `http://localhost:4323`

If backend code changed but Rust time-tracker has unrelated local syntax errors, rebuild backend only:

1. `cd /home/parothegreat/work-order/src/backend`
2. `podman build -t src-backend -f Dockerfile .`
3. `cd /home/parothegreat/work-order/src`
4. `podman compose up -d --no-build --force-recreate backend nginx`

## Production Safety Notes

- The `8443:443` mapping is rootless-safe for local Podman.
- If production Docker must bind public HTTPS on port 443 directly, either adjust the production compose mapping back to `443:443` intentionally or terminate TLS at an external reverse proxy.
- Do not remove MySQL volumes while fixing local compose startup.

## Per-File Commit Messages

- `src/docker-compose.yml` -> `[fix] support rootless podman compose mounts and ports`
- `src/docker-compose.persistent.yml` -> `[fix] support rootless podman persistent compose`
- `src/docker-compose.external-db.yml` -> `[fix] support rootless podman external database compose`
- `CODEX_SERVER_HANDOFF.md` -> `[docs] document rootless podman compose fix`

## Latest Handoff: Register Division And Numeric Batch

## Summary

- Registration form now asks for `Divisi` through a selector, not free text.
- Available registration divisions: `Soundman`, `Programmer`, `Maintenance`, `Data Analyst`.
- Registration `Angkatan` now uses a numeric input such as `13` or `14`, not school year text.
- Backend register payload accepts and validates `division`.
- Backend validates `Angkatan` as a number from `1` to `99`.
- Admin approval and graduate batch forms now use numeric `Angkatan` inputs too.
- Pending staff now preserve the selected division and numeric angkatan before admin approval.

## Changed Files

- `src/register.html`
- `src/staff.html`
- `src/backend/controllers/user_controller.go`
- `src/backend/repository/member_repository.go`
- `CODEX_SERVER_HANDOFF.md`

## Verification Done Locally

- `GOCACHE=/tmp/go-build-cache go test ./...` passed in `src/backend`.
- Register inline script syntax check passed.
- Rebuilt backend image `src-backend`.
- Recreated backend and Nginx with `podman compose up -d --no-build --force-recreate backend nginx`.
- Smoke tested `POST /api/register` with `division = Data Analyst` and `batchYear = 14`.
- Verified MySQL stored the pending member with `Division = Data Analyst` and `BatchYear = 14`.
- Removed the smoke-test member from local MySQL afterward.
- Verified served `register.html` contains the division selector and numeric angkatan input.

## Server Steps After Git Push

1. `cd /home/it/work-order`
2. `git pull`
3. `cd src/backend && go test ./...`
4. `cd /home/it/work-order/src`
5. `docker compose up -d --build backend nginx`
6. Open `/register.html` and confirm:
   - `Angkatan` is numeric.
   - `Divisi` is a dropdown selector.
7. Register a test pending staff and confirm Admin approval page shows the selected division and angkatan.

## Per-File Commit Messages

- `src/register.html` -> `[feat] collect division and numeric batch during registration`
- `src/staff.html` -> `[feat] use numeric batch inputs in staff management`
- `src/backend/controllers/user_controller.go` -> `[feat] validate registration division and numeric batch`
- `src/backend/repository/member_repository.go` -> `[feat] persist registration division and batch`
- `CODEX_SERVER_HANDOFF.md` -> `[docs] document register division and numeric batch handoff`

## Latest Handoff: UI Logic Fix Batch

## Summary

- Register division dropdown now has proper right padding and a custom selector arrow.
- Dashboard navbar desktop layout now uses a three-column grid so the center nav sits more centrally and buttons have more breathing room.
- Dashboard profile dropdown now unwraps `/api/profile` responses correctly and formats user status labels.
- The profile dropdown Standby action now reads the wrapped `/api/status` response correctly and updates the visible status.
- Settings profile avatar now unwraps `/api/profile` correctly, treats `no avatar` and `default-avatar.png` as default, preserves the uploaded avatar URL, and shows `Foto profil berhasil di update`.
- Static uploads are shared between backend and Nginx in compose files with shared SELinux labels.
- Completed work order duration is stored as numeric minutes by backend and formatted in UI as:
  - `Less than 1 minute`
  - `54 minutes`
  - `1 hour 34 minutes`
- `summary.html` and `kaizen.html` now use the same working-time formatter.
- `summary.html` now includes initial rating/quality cards for Current Rating, Notes Quality, and Average Work Time.
- `staff.html` now opens on Active Staff, puts Active Staff first, removes the Back button, keeps Dashboard, and adds safer mobile tab spacing.

## Changed Files

- `src/register.html`
- `src/index.html`
- `src/settings.html`
- `src/summary.html`
- `src/kaizen.html`
- `src/staff.html`
- `src/backend/repository/workorder_repository.go`
- `src/docker-compose.yml`
- `src/docker-compose.persistent.yml`
- `src/docker-compose.external-db.yml`
- `CODEX_SERVER_HANDOFF.md`

## Verification Done Locally

- Inline scripts passed syntax checks for:
  - `index.html`
  - `settings.html`
  - `summary.html`
  - `kaizen.html`
  - `staff.html`
  - `register.html`
- `GOCACHE=/tmp/go-build-cache go test ./...` passed in `src/backend`.
- Rebuilt backend image `src-backend`.
- Recreated backend and Nginx with `podman compose up -d --no-build --force-recreate backend nginx`.
- Verified served `register.html` contains `select-field`, `Pilih divisi`, and numeric angkatan placeholder.
- Verified served `staff.html` has Active Staff as active first tab and no `backBtn`.
- Verified served `summary.html` contains Current Rating, Notes Quality, Working Time, and `Less than 1 minute` formatter.
- Created a temporary local admin, verified `/api/profile` wrapper shape, verified `POST /api/status` returns wrapped `data.member.status = standby`, then deleted the temporary member.

## Server Steps After Git Push

1. `cd /home/it/work-order`
2. `git pull`
3. `cd src/backend && go test ./...`
4. `cd /home/it/work-order/src`
5. `docker compose up -d --build backend nginx`
6. Browser smoke test:
   - Register page division dropdown spacing.
   - Dashboard navbar centering and spacing.
   - Profile dropdown Standby button.
   - Settings avatar upload and success feedback.
   - Summary/Kaizen completed working-time display.
   - Staff mobile tabs and default Active Staff tab.

## Production Safety Notes

- Do not reset MySQL volumes.
- Static avatar uploads require backend and Nginx to share the same `./static` host directory.
- If production Docker does not use SELinux/rootless Podman, `:z` labels are still acceptable on Podman; Docker Compose may ignore or reject SELinux flags depending on environment. Review production runtime before deploy if it is not Podman.
- The Summary rating formula is an initial product heuristic:
  - Notes Quality counts completed orders with notes of at least 20 characters.
  - Current Rating combines notes quality, average work time, and a base completion factor.

## Per-File Commit Messages

- `src/register.html` -> `[fix] improve division selector spacing`
- `src/index.html` -> `[fix] center dashboard navbar and sync profile status`
- `src/settings.html` -> `[fix] persist uploaded profile avatar feedback`
- `src/summary.html` -> `[feat] add summary rating and working time formatting`
- `src/kaizen.html` -> `[fix] format completed working time consistently`
- `src/staff.html` -> `[fix] improve staff mobile tabs and default active view`
- `src/backend/repository/workorder_repository.go` -> `[fix] store completed work duration as numeric minutes`
- `src/docker-compose.yml` -> `[fix] share static upload mount between backend and nginx`
- `src/docker-compose.persistent.yml` -> `[fix] share static upload mount in persistent compose`
- `src/docker-compose.external-db.yml` -> `[fix] share static upload mount in external db compose`
- `CODEX_SERVER_HANDOFF.md` -> `[docs] document ui logic fix batch`
