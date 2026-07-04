# Repository Guidelines

## Project Structure & Module Organization
- `src/` contains the runnable app. The frontend is plain HTML pages such as `index.html`, `login.html`, `guest.html`, `summary.html`, `kaizen.html`, `shift.html`, `settings.html`, and `staff.html`.
- `src/static/assets/` holds shared CSS, JavaScript, and the service worker.
- `src/static/public/` stores uploaded or bundled images used by the UI.
- `src/backend/` is the Go API, organized by `controllers/`, `routes/`, `repository/`, `services/`, `middleware/`, `models/`, and `utils/`.
- `src/rust-engine/` is the internal time-tracker service.
- `src/db/` contains schema and migration SQL.
- Deployment docs live at the repo root in `deployment.md`, `TODO.md`, and handoff files.

## Build, Test, and Development Commands
- `cd src && podman compose up -d --build` starts the full stack locally.
- `cd src && podman compose -f docker-compose.persistent.yml up -d --build` runs the persistent deployment variant.
- `cd src/backend && go test ./...` runs Go backend tests.
- `cd src/rust-engine && cargo test` runs Rust tests.
- `podman logs work-order-backend --tail 100` and `podman logs work-order-db --tail 100` are the first checks for runtime issues.

## Coding Style & Naming Conventions
- Keep HTML, CSS, and JavaScript ASCII-only unless the file already uses localized text.
- Match existing style: 2-space indentation in HTML/JS, Go `gofmt`, Rust `cargo fmt`.
- Use lower-case, dash-separated filenames for HTML pages and snake_case for Go package files.
- Prefer small, direct edits that follow the current patterns in the repo instead of introducing new frameworks or abstractions.

## Testing Guidelines
- No large automated frontend test suite exists yet, so verify UI changes manually in desktop and mobile views.
- For backend changes, run `go test ./...` and, when relevant, a local `podman compose up -d --build` smoke check.
- For Rust changes, run `cargo test` and confirm the health endpoint still responds.

## Commit & Pull Request Guidelines
- Use commit messages in the format `[type] message`, for example `[fix] clean login routes` or `[docs] add deployment guide`.
- Keep commits focused to one file or one logical change when possible.
- PRs should include a short summary, the user-facing impact, and screenshots for UI changes.

## Security & Configuration Tips
- Do not commit secrets or `.env` values.
- Treat MySQL data and uploaded files as persistent production data; keep them outside container-only storage.
- Prefer clean URLs such as `/index`, `/login`, and `/summary` instead of `.html` paths in navigation.
