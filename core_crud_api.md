# Core CRUD + Pagination API Contract

Dokumen kontrak inti operasi CRUD dan pagination lintas modul aplikasi: pola standar 6 operasi, envelope list, lookup `/names`, dan tabel deviasi per domain.

- **Sumber kebenaran:** implementasi backend (semua controller `Api/V1/Dashboard/*`, `routes/api.php`, `ApiResponse.php`) diverifikasi terhadap dokumen alur master data internal.
- **Base path:** `/api/v1/dashboard`
- **Contoh utama:** domain **user** (`/users`) — dipilih karena paling lengkap (middleware `auth.jwt` + `permission:users,*`, filter kaya, relasi role/branch) sekaligus titik temu core 1–4 (Auth, RBAC, CRUD, Pagination). Penyimpangannya ditandai ❌.
- **Terkait:** `core_auth_api.md` (header `Authorization` & sesi `FULL_AUTH`), `core_role_permission_api.md` (semantik `permission:<resource>,<action>`).
- **Versi dokumen:** disusun 2026-10-21.

---

## 1. Ringkasan Kontrak

### 1.1 Aturan Umum

| Aspek | Ketentuan |
|---|---|
| Header | `Authorization: Bearer <access-token>` (sesi `FULL_AUTH`) + `Content-Type: application/json` untuk POST/PUT/PATCH |
| Format response | Sukses `{status:true, message, data}` / gagal `{status:false, message, data:""}` / gagal middleware `{error, code}` — lihat `core_auth_api.md` §1.1 |
| Path `{id}` | Wajib digit; non-digit → `400 {"message":"Invalid ID"}` (di beberapa route non-digit tidak match → 404) |
| ID di response | Selalu `id` (integer), timestamp `createdAt`/`updatedAt` format `Y-m-d H:i:s` |
| Soft delete | Trait `SoftDeletes` pada model → row terhapus tidak tampil di list/get (Category, Branch, Product, Supplier, User) |

### 1.2 Pola 6 Operasi (standar semua domain)

| # | Operasi | Method & Path | Status | `message` standar | `data` |
|---|---|---|---|---|---|
| 1 | Create | `POST /<resource>` | **201** | `<Resource> created successfully` | entity |
| 2 | List (paginated) | `GET /<resource>` | 200 | `Successfully retrieved all <resources>` | `{contents,page,size,totalData}` |
| 3 | Detail | `GET /<resource>/{id}` | 200 | `Successfully retrieved <resource>` | entity |
| 4 | Update | `PUT /<resource>/{id}` | 200 | `<Resource> updated successfully` | entity |
| 5 | Delete | `DELETE /<resource>/{id}` | 200 | `<Resource> deleted successfully` | `null` |
| 6 | Lookup | `GET /<resource>/names` | 200 | `Successfully retrieved all <resources> name` | `{contents,page,size,totalData}` (item `{id,name}`) |

**Error standar:**

| HTTP | `message` | Syarat |
|---|---|---|
| 400 | `Invalid ID` | `{id}` bukan digit |
| 400 | `<field> is required` | Field wajib kosong |
| 404 | `<resource> tidak ditemukan` | ID valid tapi data tidak ada |
| 403 | `permission denied` (`{error,code}`) | Role tanpa permission `permission:<resource>,<action>` |
| 500 | `internal server error` | Exception tak tertangani |

---

## 2. Contoh Utama: Domain User (`/users`)

Middleware semua operasi: `auth.jwt` + `permission:users,<action>`.

### 2.1 Create — `POST /api/v1/dashboard/users`

**Request**
```json
{
  "name": "Budi Santoso",
  "email": "budi@example.com",
  "roleId": 4,
  "branchId": 1
}
```
> `branchId` boleh `null`. Role `SUPERADMIN` tidak boleh dibuat; actor `MANAGER` wajib lewat flow user-activation-request.

**Sukses — 201 Created**
```json
{
  "status": true,
  "message": "User admin created successfully",
  "data": {
    "id": 12,
    "name": "Budi Santoso",
    "email": "budi@example.com",
    "roleId": 4,
    "roleName": "OWNER",
    "branchId": 1,
    "branchName": "Cabang Pusat",
    "status": "ACTIVE",
    "activationStatus": "REGISTERED",
    "emailVerifiedAt": null,
    "createdAt": "2026-04-22 12:00:00",
    "updatedAt": "2026-04-22 12:00:00",
    "userDetail": null
  }
}
```
Efek samping: password default di-set, activation token (30 menit) dibuat & dikirim ke email — kelanjutan flow di `core_auth_api.md` §2.6.

**Gagal — urutan validasi: duplikat → actor → role → guard → branch**

| HTTP | `message` | Catatan |
|---|---|---|
| 500 | `email already exists` | ❌ **Jangan ditiru** — duplikasi data = error klien, seharusnya `400` |
| 403 | `MANAGER must use user activation request flow` | Actor MANAGER tidak boleh buat user langsung |
| 404 | `role not found` | ❌ Pesan Inggris — pola standar `tidak ditemukan` |
| 403 | `cannot create SUPERADMIN user` | Guard role |
| 404 | `branch not found` | ❌ Pesan Inggris |
| 401/403 | error middleware | `AUTH_*` / `AUTHZ_*` (`core_auth_api.md` §3) |

> ❌ **Penyimpangan:** create user **tidak memvalidasi field wajib** (`name`/`email` kosong diterima) — modul baru wajib mengisi `<field> is required`.

---

### 2.2 List — `GET /api/v1/dashboard/users`

**Query params:** `search` (name/email), `status`, `roleId`, `page` (1), `size` (10).
> List **mengecualikan** `role_id = 1` (SUPERADMIN).

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Successfully retrieved all users",
  "data": {
    "contents": [
      {
        "id": 12,
        "name": "Budi Santoso",
        "email": "budi@example.com",
        "roleId": 4,
        "roleName": "OWNER",
        "branchId": 1,
        "branchName": "Cabang Pusat",
        "status": "ACTIVE",
        "activationStatus": "REGISTERED",
        "emailVerifiedAt": "2026-04-22 12:05:00",
        "createdAt": "2026-04-22 12:00:00",
        "updatedAt": "2026-04-22 12:00:00",
        "userDetail": {
          "fullName": "Budi Santoso",
          "nip": "EMP001",
          "joinDate": "2026-01-02",
          "employeeStatus": "FULL_TIME",
          "status": "ACTIVE"
        }
      }
    ],
    "page": 1,
    "size": 10,
    "totalData": 1
  }
}
```
Detail mekanisme pagination → bagian 4.

**Gagal:** hanya error middleware (query params tidak pernah error).

---

### 2.3 Detail — `GET /api/v1/dashboard/users/{id}`

**Sukses — 200 OK** — `data` berbentuk entity persis seperti entri `contents` di 2.2, `message`: `Successfully retrieved user`.

**Gagal**

| HTTP | `message` |
|---|---|
| 400 | `Invalid ID` |
| 404 | `user not found` — ❌ pola standar: `user tidak ditemukan` |

---

### 2.4 Update — `PUT /api/v1/dashboard/users/{id}`

**Request**
```json
{
  "name": "Budi Santoso",
  "email": "budi@example.com",
  "roleId": 4,
  "branchId": 1,
  "status": "ACTIVE"
}
```

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "User admin updated successfully",
  "data": { "...": "entity lengkap seperti 2.1" }
}
```

**Gagal — urutan: ID → duplikat → actor → role → guard → branch**

| HTTP | `message` | Catatan |
|---|---|---|
| 400 | `Invalid ID` | — |
| 404 | `user not found` | ❌ Inggris |
| 500 | `email already exists` | ❌ Jangan ditiru → seharusnya 400 |
| 403 | `MANAGER must use user activation request flow` | Actor MANAGER |
| 404 | `role not found` | ❌ Inggris |
| 403 | `cannot assign SUPERADMIN role` | Guard role tujuan |
| 404 | `branch not found` | ❌ Inggris |

---

### 2.5 Delete — `DELETE /api/v1/dashboard/users/{id}`

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "User admin deleted successfully",
  "data": null
}
```
Efek samping: token sesi dibersihkan (logout paksa), lalu **soft delete** (trait `SoftDeletes`).

**Gagal:** `400 Invalid ID`, `404 user not found` (❌ Inggris), error middleware.

---

### 2.6 Operasi Ekstra (di luar 6 pola — boleh, selama konsisten)

| Endpoint | Method | Sukses | Keterangan |
|---|---|---|---|
| `/users/{id}/deactivate` | PATCH | 200 `User admin deactivated successfully` + entity | Set `status=INACTIVE` + bersihkan semua token sesi |
| `/users/{id}/resend-activation` | POST | 200 `Activation code resent successfully`, `data:null` | Kirim ulang activation code — detail di `core_auth_api.md` §2.6 |

---

## 3. Lookup `/names` — Contoh Supplier & Category

Endpoint ringan untuk dropdown/autocomplete (payload hanya `id` + `name`), **user tidak memiliki endpoint ini**.

**Contoh:** `GET /api/v1/dashboard/suppliers/names?search=PT&page=1&size=10`
(middleware `permission:suppliers,read`)

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Successfully retrieved all suppliers name",
  "data": {
    "contents": [ { "id": 3, "name": "PT Maju Jaya" } ],
    "page": 1,
    "size": 10,
    "totalData": 1
  }
}
```

Category menyediakan pola sama dengan `message` `Successfully retrieved all categories name`.
**Catatan:** endpoint `/names` juga **terpaginasi** (bukan dump seluruh data).

---

## 4. Kontrak Pagination

### 4.1 Envelope (seragam di 16+ controller)

```json
{
  "contents": [ ... ],
  "page": 1,
  "size": 10,
  "totalData": 42
}
```

### 4.2 Query Params

| Param | Tipe | Default | Keterangan |
|---|---|---|---|
| `page` | integer | `1` | Nilai `0` di-reset ke `1`; nilai negatif **tidak divalidasi** (deviasi minor) |
| `size` | integer | `10` | Nilai `0` di-reset ke `10`; **tidak ada batas maksimum** di backend |
| `search` | string | `""` | `LIKE %…%`, kolom beda per domain (user: `name`,`email`; supplier: 7 kolom; product: `name`,`sku`,`barcode`) |
| filter lain | bervariasi | non-aktif | Contoh: `status` (category/branch/user), `roleId` (user), `moduleId`+`roleId` (role-permission) |

### 4.3 Rumus & Perilaku

```text
totalData  = COUNT(query dengan filter yang sama, tanpa limit)
offset     = (page - 1) * size
contents   = query → offset(offset) → limit(size)
jumlah halaman (tidak dikirim backend) = CEIL(totalData / size)
```

**Ketidakkonsistenan yang harus diketahui FE (tidak ada standar baku):**
1. **Urutan bawaan:** supplier/branch/user → `ORDER BY created_at DESC` (user juga `id DESC`); **category tanpa `ORDER BY`** (urutan tidak terjamin); role-permission list → `created_at DESC, role_id DESC, module_id DESC`.
2. `totalData` menghitung **baris**, bukan halaman; selalu muncul meski `contents` kosong (`totalData: 0`).
3. List role-permission memakai subquery group (`roleId+moduleId`), bukan row mentah — `totalData` = jumlah grup.

---

## 5. Tabel Deviasi per Domain

| Domain | Route terproteksi | Soft delete | `/names` | Validasi field wajib | Cek unik | Deviasi utama ❌ |
|---|---|---|---|---|---|---|
| **category** | ❌ **tanpa `auth.jwt`/`permission`** (`routes/api.php:74`) | ✅ | ✅ | ✅ `category name is required` (400) | — | Route terbuka; tanpa `ORDER BY` |
| **branch** | ✅ `permission:branch,*` | ✅ | ✅ | ❌ tak divalidasi | — | Pesan `Failed to create/update/delete branch` (500) gaya try/catch; tak ada validasi wajib |
| **product** | ✅ `role.strict:SUPERADMIN` + `permission:products,*` | ✅ | ✅ | ✅ `price is required` (400) | ✅ `sku already exists` (400) | — |
| **supplier** | ✅ `permission:suppliers,*` | ✅ | ✅ | ✅ `supplier name is required` + max length | ✅ `supplier name already exists` (400) | **Kandidat paling patuh pola** |
| **user** | ✅ `permission:users,*` | ✅ (sesi dipurge dulu) | ❌ tidak ada | ❌ tidak divalidasi | ✅ tapi `email already exists` → **500** | Pesan `user not found`/`role not found` Inggris; create tanpa validasi wajib; list menyembunyikan `role_id=1` |
| **role-permission** | ✅ `role.strict:SUPERADMIN` | n/a (tabel pivot) | n/a | ✅ | n/a | Butuh **step-up token** — kontrak khusus di `core_role_permission_api.md` |

**Prinsip untuk modul baru:** ikuti baris **supplier** (lengkap & bersih); jangan ikuti ❌ di atas.

---

## 6. Konvensi Menambah Modul CRUD Baru

1. Daftarkan route di `routes/api.php` dengan `Route::middleware(['auth.jwt'])->prefix('<plural>')` + `permission:<plural>,<action>` per operasi.
2. Implementasikan 6 operasi (§1.2) dengan `message` persis dari tabel standar.
3. Gunakan envelope pagination §4.1 + `search`/`page`/`size` + filter domain.
4. Model: trait `SoftDeletes`; response entity via `toResponse()` (camelCase, `createdAt`/`updatedAt`).
5. Error: `400 Invalid ID` / `<field> is required`, `404 <resource> tidak ditemukan` (**Bahasa Indonesia, konsisten**), validasi duplikat dengan `400`, bukan 500.
6. Sediakan `GET /<plural>/names` untuk dropdown.
7. Setiap operasi butuh permission `permission:<plural>,{create,read,update,delete}` — terima kasih sudah konsisten dengan core role-permission.
