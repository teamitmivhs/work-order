# TODO List — Work Order ITMIVHS

Status dokumen ini mengikuti kondisi project saat ini.

Estimasi readiness:

- core functionality: hampir selesai
- UI/UX consistency: masih butuh finishing pass
- production readiness: belum final
- low-maintenance readiness: belum final

## Completed

### Core app
- [x] Login, register, logout, dan JWT auth sudah berjalan
- [x] Role dasar sudah aktif untuk Admin, Operator, dan guest flow
- [x] Dashboard work order utama sudah berjalan untuk pending, on progress, completed
- [x] Summary, Kaizen, TechGuide, Shift Piket, Staff, dan Settings page sudah tersedia
- [x] Guest request flow sudah bisa membuat order dan cek status dengan tracking code
- [x] Guest notes sudah tersambung ke summary/kaizen data
- [x] Take order flow sudah mendukung pelaksana utama dan operator bantuan
- [x] Executor preview popup sudah tersedia

### Shift and member status
- [x] Halaman `shift.html` dipisah sebagai tempat pengelolaan status operator
- [x] Dashboard `Member Status` popup sekarang read-only
- [x] Kategori status aktif: `standby`, `onjob`, `nextshift`, `offduty`
- [x] Support state sudah dihapus dari flow utama
- [x] Shift day counter sudah ada
- [x] Saat hari berganti, member dengan status `nextshift` otomatis dipindah ke `standby`

### Backend and data
- [x] Backend Go, MySQL, nginx, dan Rust time-tracker sudah terhubung
- [x] Response API sudah konsisten ke format `{ code, message, data }`
- [x] Timer start/stop sudah terhubung ke workflow work order
- [x] Shift day counter repository dan migration sudah ditambahkan
- [x] Endpoint shift counter sudah tersedia
- [x] Patch status member dan patch work order sudah tersedia

### UI improvements already done
- [x] Banyak bug mobile layout sudah dibersihkan di dashboard, guest, kaizen, summary, settings, shift
- [x] Banyak bug dark mode/light mode sudah dibersihkan
- [x] Popup crop avatar di settings sudah diperbaiki untuk dark mode
- [x] Badge status member popup sudah diperbaiki agar lebih rapi
- [x] Status bar dashboard sekarang sudah punya `Off Duty`

### Deployment groundwork
- [x] `deployment.md` sudah dibuat sebagai draft panduan production low-maintenance
- [x] Clean URL routing sudah diterapkan untuk browser navigation dan service worker
- [x] `AGENTS.md` sudah ditambahkan sebagai contributor guide
- [x] Compose internal, persistent, dan external-db sudah tersedia

## High Priority

### Low-maintenance production
- [ ] Ubah `docker-compose.persistent.yml` agar memakai host path yang benar untuk server production
  Gunakan path seperti `/mnt/disk3/work-order/mysql-data` dan `/mnt/disk3/work-order/static`, bukan `/docker/...`.
- [ ] Pastikan upload foto benar-benar masuk ke storage host persistent
  Avatar dan foto dokumentasi audit tidak boleh bergantung pada lifecycle container.
- [ ] Buat script backup otomatis di repo
  Minimal backup MySQL dan folder upload.
- [ ] Uji restore dari backup
  Backup tanpa restore test belum cukup aman.
- [ ] Tambahkan runbook operasional singkat
  Fokus pada start, restart, update, rollback, dan recovery.

### Data and storage
- [ ] Tambahkan kebijakan kompresi/resize foto upload
  Ini penting untuk storage jangka panjang dan audit trail.
- [ ] Batasi ukuran dan format final file upload
  Target praktis: foto final tidak terlalu besar, tetap cukup jelas untuk audit.
- [ ] Audit lokasi penyimpanan avatar dan dokumentasi agar tidak ada file yang masih tersimpan di tempat sementara

### UI consistency pass
- [ ] Final pass light mode vs dark mode di semua popup utama
- [ ] Final pass mobile layout di halaman dengan data panjang atau avatar rusak
- [ ] Final pass desktop polish di dashboard popup dan status sections
- [ ] Final audit clean URL di semua navigasi yang tersisa dari page statis atau redirect helper

## Medium Priority

### Security and ops
- [ ] Set `GIN_MODE=release` untuk backend production
- [ ] Finalisasi `Cargo.lock` Rust engine untuk reproducible build
- [ ] Tambahkan healthcheck operasional yang benar-benar dipakai saat deploy
- [ ] Review limit upload dan validasi file dari sisi backend

### Reliability
- [ ] Review semua endpoint yang menyentuh status member agar konsisten dengan flow `shift.html`
- [ ] Review work order transitions: pending -> progress -> completed
- [ ] Review guest tracking flow untuk edge case order lama atau data foto kosong
- [ ] Audit fallback avatar/image handling di semua halaman

### Testing
- [ ] Tambahkan smoke test backend untuk endpoint penting
- [ ] Tambahkan minimal regression checklist manual sebelum deploy
- [ ] Tambahkan test atau script validasi untuk workflow guest, shift, dan complete order

## Future Enhancements

- [ ] Notifikasi WhatsApp / Telegram / email saat order selesai
- [ ] Export laporan audit ke PDF atau Excel
- [ ] Dashboard statistik yang lebih kaya untuk admin
- [ ] Permission matrix yang lebih detail per role/divisi
- [ ] Archive policy untuk data audit lama
- [ ] Monitoring sederhana untuk disk usage dan backup freshness

## Definition of Final

Project ini baru layak disebut final untuk internal production kalau poin berikut sudah terpenuhi:

- [ ] data tetap aman setelah server restart atau container rebuild
- [ ] backup otomatis aktif dan restore pernah diuji
- [ ] upload foto sudah dikendalikan agar storage tidak meledak
- [ ] UI sudah konsisten di mobile dan desktop, light dan dark
- [ ] deployment steps sudah terdokumentasi dan bisa dijalankan orang lain
- [ ] tidak ada workflow utama yang masih bergantung pada knowledge developer utama
