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
