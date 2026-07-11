# Work Order Progress Check

Update: 2026-07-11 12:49 WIB

## Progress

90% selesai.

Fitur utama sudah selesai, sudah dites lokal, dan sudah di-commit per file. Sisa utama adalah push ke GitHub dan manual QA singkat di browser untuk flow reject dari UI.

## Selesai

- Notification lifecycle work order sudah ditambahkan untuk event penting seperti guest create order, operator join, admin approve/start, add worker, dan complete.
- Push notification sudah ditargetkan ke role/user yang relevan, bukan hanya broadcast asal.
- Filter work order sudah ditambah untuk requester dan tanggal.
- Backend sudah mengembalikan `createdAt` dan migrasi kolom `orders.CreatedAt` sudah ditambahkan.
- Flow reject work order sudah ditambahkan sebagai status `rejected`, bukan hard delete, supaya guest tetap bisa tracking.
- Sebelum reject muncul popup konfirmasi.
- Setelah admin pilih `Ya`, muncul popup catatan wajib untuk requester.
- Catatan reject disimpan ke `AdminNotes`.
- Guest tracking menampilkan status `Rejected` dan catatan reject.
- Dashboard punya tab/count `Rejected`.

## Validasi

- `git diff --check`: pass
- `node --check src/static/assets/script.js`: pass
- Inline script check `src/guest.html`: pass
- `cd src/backend && GOCACHE=/tmp/go-build go test ./...`: pass
- `cd src && podman compose up -d --build`: pass
- `curl -I http://127.0.0.1:4323/guest`: HTTP 200
- Backend route `PATCH /api/workorders/:id/reject` sudah terdaftar di log Gin.

## Commit Terakhir

- `2b09fbd [feat] show rejected work order reason`
- `a2a5787 [feat] reject work orders with requester note`
- `8297648 [style] add rejected work order badges`
- `7f1cda6 [feat] add rejected work order tab`
- `f2dc09c [feat] route work order reject`
- `75c6307 [feat] add work order reject handler`
- `69197c8 [feat] reject pending work orders`
- `f62c4ca [feat] add reject work order request`
- `c954a31 [feat] filter work orders by requester and date`
- `c971bd0 [feat] add work order filters`
- `cdc04bf [feat] notify work order lifecycle`
- `409a91a [feat] target work order push recipients`
- `9b9a3e4 [feat] migrate work order created date`
- `0ea77ea [feat] return work order created date`
- `201a968 [feat] expose work order created date model`

## Belum Selesai

- Push ke GitHub belum berhasil. Percobaan pertama gagal karena sandbox tidak bisa resolve `github.com`; percobaan eskalasi network dibatalkan oleh user.
- Manual QA browser untuk reject flow belum dilakukan setelah commit: klik reject, konfirmasi, isi catatan, lalu cek tracking guest.
- File untracked `src/static/public/avatar_38_1783413710.jpg` masih ada dan sengaja tidak disentuh.

## Catatan Teknis

- Hard delete lama tidak diubah. Flow baru memakai reject agar requester tetap bisa melihat status dan alasan.
- Reject hanya boleh untuk admin dan hanya untuk work order dengan status `pending`.
- Reason reject divalidasi backend: wajib, minimal 3 karakter, maksimal 1000 karakter.
