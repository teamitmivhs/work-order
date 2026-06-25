# API Reference - IT Work Order System

## Base URL
```
http://localhost/api
```
> Semua request ke `/api/*` di-proxy oleh Nginx ke Go backend di port 8080.

---

## Response Format

Semua response menggunakan format konsisten:

**Success:**
```json
{
  "code": 200,
  "message": "Description",
  "data": { ... }
}
```

**Error:**
```json
{
  "code": 400,
  "message": "Error description",
  "details": "Optional additional context"
}
```

**Common HTTP Status Codes:**
| Code | Meaning |
|------|---------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request — input tidak valid |
| 401 | Unauthorized — token missing/expired |
| 403 | Forbidden — tidak punya akses |
| 404 | Not Found |
| 409 | Conflict — data sudah ada |
| 429 | Too Many Requests — rate limit |
| 500 | Internal Server Error |

---

## Authentication

### Token Usage
Semua protected endpoint wajib mengirim header:
```
Authorization: Bearer <token>
```

Token disimpan di `localStorage`:
```javascript
// Simpan setelah login
localStorage.setItem('userToken', data.data.token);

// Pakai di request
fetch('/api/workorders', {
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('userToken')}`,
    'Content-Type': 'application/json'
  }
});
```

**Token expiry:** 24 jam. Setelah expired, user harus login ulang.

---

## Auth Endpoints

### Register
`POST /api/register` — Public, rate limited (10 req/menit per IP)

**Request:**
```json
{
  "name": "username",
  "password": "Password123"
}
```

**Password requirements:**
- Minimum 8 karakter
- Harus ada huruf besar (A-Z)
- Harus ada huruf kecil (a-z)
- Harus ada angka (0-9)

**Response 201:**
```json
{
  "code": 201,
  "message": "Registration successful",
  "data": {
    "token": "eyJhbGci...",
    "member": {
      "id": 1,
      "name": "username",
      "role": "Operator",
      "status": "standby"
    }
  }
}
```

**Errors:** `400` invalid input / weak password, `409` username exists

---

### Login
`POST /api/login` — Public, rate limited (10 req/menit per IP)

**Request:**
```json
{
  "name": "username",
  "password": "Password123"
}
```

**Response 200:**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGci...",
    "member": {
      "id": 1,
      "name": "username",
      "role": "Operator",
      "status": "standby",
      "avatar": "avatar.png"
    }
  }
}
```

**Errors:** `401` invalid credentials

---

### Logout
`POST /api/logout` — 🔒 Protected

Menginstruksikan client untuk hapus token lokal. Token tidak di-blacklist di server (stateless JWT).

**Response 200:**
```json
{
  "code": 200,
  "message": "Logged out successfully. Please remove your token on the client side."
}
```

---

### Get Profile
`GET /api/profile` — 🔒 Protected

**Response 200:**
```json
{
  "id": 1,
  "name": "username",
  "role": "Operator",
  "status": "standby",
  "avatar": "avatar.png"
}
```

---

## Members

### Get All Members
`GET /api/members` — Public (tidak butuh token)

**Response 200:**
```json
[
  {
    "id": 1,
    "name": "Aldi Fadlurahman R",
    "role": "programmer",
    "status": "standby",
    "avatar": "aldi.png"
  }
]
```

> Selalu return array `[]`, tidak pernah `null`.
> Password tidak disertakan dalam response.

---

### Update Member Status
`PATCH /api/members/:id/status` — 🔒 Protected

**Request:**
```json
{
  "status": "onjob"
}
```

**Valid status values:** `standby`, `onjob`, `support`, `nextshift`, `offduty`

**Response 200:**
```json
{
  "code": 200,
  "message": "Member status updated successfully",
  "data": { "id": 1 }
}
```

**Errors:** `400` invalid status value

---

## Work Orders

### Get Work Orders
`GET /api/workorders` — Public

Response difilter berdasarkan role user (jika token disertakan):
- **Tanpa token / Admin** — semua work orders
- **Operator** — hanya orders yang di-assign ke user tersebut

**Response 200:**
```json
[
  {
    "id": 1,
    "priority": "high",
    "time": "14:30",
    "requester": "Budi",
    "location": "Gedung A - Lantai 2",
    "device": "Printer HP",
    "problem": "Paper jam di tray 2",
    "executors": [1, 2],
    "workingHours": 30,
    "status": "completed",
    "safetyChecklist": ["ga3", "ga4"],
    "completedAt": "15:05"
  }
]
```

> `workingHours` dalam menit (integer), bisa `null` jika belum selesai.
> `completedAt` berupa string display `"HH:MM"`, bisa kosong.

---

### Create Work Order
`POST /api/workorders` — 🔒 Protected

**Request:**
```json
{
  "priority": "high",
  "time_display": "14:30",
  "time_sort": "14:30:00",
  "requester": "Budi Santoso",
  "location": "Gedung A - Lantai 2, Ruang 201",
  "device": "Printer HP",
  "problem": "Paper jam di tray 2",
  "working_hours": "0 menit",
  "status": "pending",
  "executors": [],
  "safety_checklist": []
}
```

**Valid priority:** `low`, `medium`, `high`, `urgent`

**Response 201:**
```json
{
  "code": 201,
  "message": "Work order created successfully",
  "data": { "id": 42 }
}
```

**Errors:** `400` missing required fields / invalid priority

---

### Take Work Order
`POST /api/workorders/:id/take` — 🔒 Protected

Assign operator ke order dan mulai pengerjaan. Timer Rust otomatis dimulai setelah commit berhasil.

**Request:**
```json
{
  "status": "progress",
  "executors": [1, 2],
  "safety_checklist_items": ["ga3", "ga4", "ga5"]
}
```

**Validasi:**
- `status` harus `"progress"`
- `executors` minimal 1 item
- Semua executor yang dipilih harus berstatus `standby` (bukan `onjob`)
- Pengecekan status executor dilakukan dengan `SELECT FOR UPDATE` di dalam transaksi

**Response 200:**
```json
{
  "code": 200,
  "message": "Order taken successfully",
  "data": { "id": 42 }
}
```

**Errors:** `400` executor sudah onjob / status salah, `404` executor not found

---

### Complete Work Order
`PATCH /api/workorders/:id/complete` — 🔒 Protected

Tandai order selesai. Timer Rust otomatis dihentikan dan `working_hours` disimpan ke database.

**Request:**
```json
{
  "status": "completed",
  "completed_at_display": "15:05"
}
```

**Validasi:**
- Hanya member yang di-assign (`executors`) yang bisa complete
- Safety checklist harus sudah diisi (count > 0)
- `status` harus `"completed"`

**Response 200:**
```json
{
  "code": 200,
  "message": "Order completed successfully",
  "data": { "id": 42 }
}
```

**Errors:** `400` checklist belum diisi, `403` bukan executor order ini

---

### Update Work Order (Executors)
`PATCH /api/workorders/:id` — 🔒 Protected

Perbarui daftar executor untuk order yang masih pending. Dipakai untuk menambah/menghapus worker.

**Request:**
```json
{
  "executors": [1, 2, 3],
  "status": "pending"
}
```

> `status` opsional. Member yang dihapus dari `executors` otomatis di-reset ke `standby`.

**Response 200:**
```json
{
  "code": 200,
  "message": "Order updated successfully",
  "data": { "id": 42 }
}
```

---

### Delete Work Order
`DELETE /api/workorders/:id` — 🔒 Protected (Admin only)

Menghapus order beserta semua executor dan checklist terkait. Member yang sedang `onjob` untuk order ini otomatis di-reset ke `standby`.

**Response 200:**
```json
{
  "code": 200,
  "message": "Order deleted successfully",
  "data": { "id": 42 }
}
```

**Errors:** `403` bukan Admin

---

### Get Safety Checklist
`GET /api/workorders/:id/checklist` — 🔒 Protected

**Response 200:**
```json
{
  "checklist": ["ga3", "ga4", "ga5"]
}
```

---

### Update Safety Checklist
`PUT /api/workorders/:id/checklist` — 🔒 Protected

**Request:**
```json
{
  "checklist_items": ["ga3", "ga4", "ga5"]
}
```

**Response 200:**
```json
{
  "code": 200,
  "message": "Safety checklist updated successfully"
}
```

---

## Analytics

### Get Kaizen Metrics
`GET /api/kaizen` — 🔒 Protected

**Response 200:**
```json
{
  "totalKaizens": 50,
  "implementedKaizens": 35,
  "pendingKaizens": 10
}
```

> `totalKaizens` = completed + pending + progress (semua status dihitung).

---

## Rust Time Tracker (Internal)

Endpoint ini hanya bisa diakses dari Go backend (memerlukan header `X-Internal-Key`). Tidak bisa diakses langsung dari browser/frontend.

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `POST` | `/timer/start` | Mulai timer untuk order |
| `POST` | `/timer/stop` | Hentikan timer, return durasi |
| `GET` | `/timer/:id` | Status timer satu order |
| `GET` | `/timers` | List semua timer aktif |
| `GET` | `/health` | Health check (jumlah timer aktif) |

Timer start/stop dipanggil **otomatis** oleh Go backend saat `TakeOrder` dan `CompleteOrder` — tidak perlu dipanggil manual dari frontend.

---

## CORS

Backend dikonfigurasi dengan CORS:
- **Allowed Origins**: `*`
- **Allowed Methods**: `GET, POST, PUT, PATCH, DELETE, OPTIONS`
- **Allowed Headers**: `Origin, Content-Type, Authorization`

---

## Rate Limiting

| Endpoint | Limit |
|----------|-------|
| `POST /api/login` | 10 req/menit per IP |
| `POST /api/register` | 10 req/menit per IP |
| Semua endpoint lain | Tidak ada limit saat ini |

> Untuk production multi-instance, implementasi Redis-based rate limiter disarankan.

---

## Data Validation Rules

| Field | Rule |
|-------|------|
| `name` (username) | 3–50 karakter, required |
| `password` | Min 8 karakter, harus ada uppercase + lowercase + digit |
| `priority` | `low` / `medium` / `high` / `urgent` |
| Work order `status` | `pending` / `progress` / `completed` |
| Member `status` | `standby` / `onjob` / `support` / `nextshift` / `offduty` |
| `priority` (required) | Ya |
| `requester` (required) | Ya |
| `location` (required) | Ya |
| `device` (required) | Ya |
| `problem` (required) | Ya |