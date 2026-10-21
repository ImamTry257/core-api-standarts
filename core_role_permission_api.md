# Core Role Permission API Contract

Dokumen kontrak API manajemen role-permission dan step-up security (challenge): endpoint, request, response sukses, dan response gagal.

- **Sumber kebenaran:** implementasi backend (`routes/api.php`, `RolePermissionController`, `RolePermissionSecurityController`, `RolePermissionSecurityService`) diverifikasi terhadap dokumen referensi API internal.
- **Base path:** `/api/v1/dashboard/role-permissions`
- **Scope:** CRUD role-permission (10 endpoint) + step-up challenge (3 endpoint). Otorisasi middleware & `roles/names` di luar scope (lihat `core_auth_api.md` §3 untuk bentuk error middleware).
- **Versi dokumen:** disusun 2026-10-21.

---

## 1. Ringkasan Kontrak

### 1.1 Autentikasi & Role

| Aspek | Ketentuan |
|---|---|
| Header | `Authorization: Bearer <access-token>` (sesi `FULL_AUTH`) |
| Role | Seluruh endpoint: **`SUPERADMIN`** (middleware `auth.jwt` + `role.strict:SUPERADMIN`) → non-SUPERADMIN `403 {error:"permission denied", code:"AUTHZ_ROLE_DENIED"}` |
| Content-Type | `application/json` (untuk POST/PUT) |

### 1.2 Format Response

Sama dengan kontrak auth (`core_auth_api.md` §1.1):
- Sukses: `{"status": true, "message": "...", "data": ...}`
- Gagal: `{"status": false, "message": "...", "data": ""}`
- Gagal middleware: `{"error": "...", "code": "..."}`

### 1.3 Dua Lapis Step-up Security

Mutasi **create / update / delete-group** wajib melewati step-up verification:

```
1. POST /security/challenges          → password + intent + payloadDigest
                                         (header X-Step-Up-Context: role-permission)
2. POST /security/challenges/{id}/verify → otpCode → verifiedToken (berlaku 10 menit)
3. Kirim mutasi dengan header         X-Role-Permission-Verified-Token: <verifiedToken>
   → backend cek ulang payloadDigest request vs digest di token
```

`DELETE /{id}` (hapus satu row) **tidak** butuh step-up.

### 1.4 Batas Waktu & Limit

| Item | Nilai |
|---|---|
| OTP challenge | 5 menit |
| Cooldown resend OTP | 60 detik |
| Masa berlaku `verifiedToken` | 10 menit |
| Max percobaan OTP salah | 5 → `429 too many invalid otp attempts` |
| `verifiedToken` | satu pemakaian (di-`consumed` saat mutasi sukses) |

### 1.5 Rumus `payloadDigest` (wajib identik)

`payloadDigest` = SHA-256 hex dari JSON **tanpa spasi**, dengan array ID sudah **di-unique, difilter >0, diurutkan naik**:

| intent | JSON |
|---|---|
| CREATE | `{"roleId":<n>,"moduleIds":[...],"permissionIds":[...]}` |
| UPDATE | `{"id":<anchorId>,"roleId":<n>,"moduleIds":[...],"permissionIds":[...]}` |
| DELETE | `{"roleId":<n>,"moduleId":<n>}` |

**Contoh perhitungan (siap pakai):**

```text
CREATE  roleId=6, moduleIds=[3], permissionIds=[19,20]
        {"roleId":6,"moduleIds":[3],"permissionIds":[19,20]}
        → 1e114392486f9e6fd7f512f2d9901f2f7c438a618c1fa5b30a22d51989ef2365

UPDATE  id=100, roleId=6, moduleIds=[3], permissionIds=[19,22]
        {"id":100,"roleId":6,"moduleIds":[3],"permissionIds":[19,22]}
        → 73ab9fc0038a0b11de43beadc52c971a9ceccb9f7f7cbf37931f0527dfd36465

DELETE  roleId=4, moduleId=3
        {"roleId":4,"moduleId":3}
        → cebb4f2784d5d34371bf88e2db59cd000d16fa67f87a360357c22134a6664a8
```

> Digest challenge harus dihitung dari **payload request mutasi yang sama** yang akan dikirim kemudian, karena backend membandingkan ulang dengan payload riil.

---

## 2. Endpoint Read (tanpa step-up)

### 2.1 List Grouped — `GET /api/v1/dashboard/role-permissions`

**Query params:** `search` (nama role/module), `roleId` (0=semua), `moduleId` (0=semua), `page` (1), `size` (10).

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Successfully retrieved role permissions",
  "data": {
    "contents": [
      {
        "roleId": 6,
        "roleName": "DUMMY_ROLE_PERMISSION",
        "moduleId": 3,
        "moduleName": "users",
        "operatorBy": 1,
        "operatorName": "Super Admin",
        "updatedAt": "2026-05-05 10:22:00"
      }
    ],
    "page": 1,
    "size": 10,
    "totalData": 1
  }
}
```
Satu entri = satu grup `roleId + moduleId` (bukan per permission).

**Gagal:** `401/403` error middleware (bagian 5). Tidak ada error body-query.

---

### 2.2 Detail Grup — `GET /api/v1/dashboard/role-permissions/{id}?roleId=&moduleId=`

**Query params wajib:** `roleId`, `moduleId`. Path `{id}` (angka) **diabaikan** — hanya untuk backward compatibility; wajib angka, selain itu route tidak match → `404`.

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Successfully retrieved role permission",
  "data": {
    "anchorId": 100,
    "roleId": 6,
    "roleName": "DUMMY_ROLE_PERMISSION",
    "moduleId": 3,
    "moduleName": "users",
    "permissions": [
      { "id": 19, "name": "users.read" },
      { "id": 20, "name": "users.update" }
    ],
    "operatorBy": 1,
    "operatorName": "Super Admin",
    "updatedAt": "2026-05-05 10:22:00"
  }
}
```
> `anchorId` = ID row pertama grup — dipakai sebagai `{id}` untuk `PUT` (2.8).

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `roleId and moduleId are required` | Salah satu query param ≤ 0/kosong |
| 404 | `role tidak ditemukan` | `roleId` tidak ada |
| 404 | `module tidak ditemukan` | `moduleId` tidak ada |

---

### 2.3 Meta Modules — `GET /api/v1/dashboard/role-permissions/meta/modules`

**Query:** `search`, `page` (1), `size` (10).

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Successfully retrieved module metadata",
  "data": {
    "contents": [ { "id": 3, "name": "users" } ],
    "page": 1, "size": 10, "totalData": 1
  }
}
```

**Gagal:** hanya error middleware.

---

### 2.4 Meta Permissions — `GET /api/v1/dashboard/role-permissions/meta/permissions`

**Query wajib:** `moduleId`. Opsional: `search`, `page`, `size`.

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Successfully retrieved permission metadata",
  "data": {
    "contents": [
      { "id": 19, "name": "users.read" },
      { "id": 20, "name": "users.update" }
    ],
    "page": 1, "size": 10, "totalData": 2
  }
}
```
Hanya permission dengan prefix `<moduleName>.` yang dikembalikan.

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `moduleId is required` | `moduleId` ≤ 0/kosong |
| 404 | `module tidak ditemukan` | Module tidak ada |

---

### 2.5 Modules + Permission — `GET /api/v1/dashboard/role-permissions/modules-permission`

**Query:** `search` (nama module).

**Sukses — 200 OK** (tanpa pagination)
```json
{
  "status": true,
  "message": "Successfully retrieved modules-permission metadata",
  "data": {
    "contents": [
      {
        "moduleId": 3,
        "moduleName": "users",
        "permissions": [
          { "id": 19, "name": "users.read" },
          { "id": 20, "name": "users.update" }
        ]
      }
    ]
  }
}
```
Module tanpa permission ikut difilter keluar.

**Gagal:** hanya error middleware.

---

### 2.6 Permission by Role — `GET /api/v1/dashboard/role-permissions/by-role?roleId=4`

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Successfully retrieved existing permissions for role",
  "data": {
    "roleId": 4,
    "roleName": "OWNER",
    "permissions": [
      {
        "permissionId": 1,
        "permissionName": "dashboard.view",
        "moduleId": 1,
        "moduleName": "dashboard"
      }
    ]
  }
}
```

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `roleId is required` | `roleId` ≤ 0/kosong |
| 404 | `role tidak ditemukan` | Role tidak ada |

---

## 3. Endpoint Mutasi (step-up wajib)

Header umum:
```
Authorization: Bearer <full-access-token>
X-Role-Permission-Verified-Token: <verifiedToken dari 4.2>
Content-Type: application/json
```

### 3.1 Create — `POST /api/v1/dashboard/role-permissions`

**Request**
```json
{
  "roleId": 6,
  "moduleIds": [3],
  "permissionIds": [19, 20]
}
```

**Sukses — 201 Created**
```json
{
  "status": true,
  "message": "Role permission created successfully",
  "data": {
    "createdCount": 2,
    "contents": [
      {
        "id": 100,
        "roleId": 6,
        "roleName": "DUMMY_ROLE_PERMISSION",
        "moduleId": 3,
        "moduleName": "users",
        "permissionId": 19,
        "permissionName": "users.read",
        "operatorBy": 1,
        "operatorName": "Super Admin",
        "createdAt": "2026-05-05 10:22:00",
        "updatedAt": "2026-05-05 10:22:00"
      }
    ]
  }
}
```
- Pair `(role, module, permission)` yang sudah ada dilewati (idempotent) → `createdCount` bisa < jumlah `permissionIds`.
- Module tiap permission ditentukan dari **prefix nama permission** terpanjang (mis. `users.read` → module `users`).

**Gagal — urutan validasi: body → ID → role → SUPERADMIN → module → permission-match → step-up**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 400 | `roleId, moduleIds, and permissionIds are required` | `roleId` ≤ 0, `moduleIds`/`permissionIds` bukan array, atau `permissionIds` kosong |
| 400 | `moduleIds must contain valid module IDs` | Semua ID di `moduleIds` tidak valid (≤0/non-angka) |
| 400 | `permissionIds must contain valid permission IDs` | Semua ID di `permissionIds` tidak valid |
| 404 | `role tidak ditemukan` | `roleId` tidak ada |
| 400 | `cannot modify SUPERADMIN role permissions` | Target role `SUPERADMIN` |
| 404 | `module tidak ditemukan` | Ada `moduleId` yang tidak ada |
| 404 | `permission tidak ditemukan` | Ada `permissionId` yang tidak ada |
| 400 | `permission does not match selected modules` | Prefix permission tidak sesuai `moduleIds` |
| 403 | `step-up verification required` | Header `X-Role-Permission-Verified-Token` kosong |
| 422 | `invalid verified token` | Token tidak bisa didekripsi / claim tidak lengkap / challenge tidak ada |
| 422 | `verified token expired` | Lewat 10 menit |
| 422 | `verified token already used` | Token sudah dipakai mutasi lain |
| 403 | `step-up verification mismatch` | User/intent/`payloadDigest` token ≠ payload request |
| 500 | `internal server error` | Exception tak tertangani |

---

### 3.2 Update — `PUT /api/v1/dashboard/role-permissions/{id}`

`{id}` = `anchorId` dari 2.2. Bersifat **sinkronisasi**: permission yang tidak ada di payload dihapus, yang baru ditambahkan.

**Request**
```json
{
  "roleId": 6,
  "moduleIds": [3],
  "permissionIds": [19, 22]
}
```

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Role permission updated successfully",
  "data": {
    "addedCount": 1,
    "removedCount": 1,
    "contents": [
      {
        "id": 101,
        "roleId": 6,
        "roleName": "DUMMY_ROLE_PERMISSION",
        "moduleId": 3,
        "moduleName": "users",
        "permissionId": 22,
        "permissionName": "users.delete",
        "operatorBy": 1,
        "operatorName": "Super Admin",
        "createdAt": "2026-05-05 10:30:00",
        "updatedAt": "2026-05-05 10:30:00"
      }
    ]
  }
}
```
> `permissionIds: []` = **cabut semua** permission role tersebut di module terpilih.
> `role_id` anchor dan keberadaan `module_id` di `moduleIds` harus cocok dengan anchor row.

**Gagal — urutan: path → body → anchor-match → role → SUPERADMIN → module → step-up → permission-match**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `Invalid ID` | `{id}` bukan digit |
| 404 | `role permission tidak ditemukan` | Row anchor tidak ada |
| 400 | `invalid input` | Body bukan JSON valid |
| 400 | `roleId, moduleIds, and permissionIds are required` | `roleId` ≤ 0 / `moduleIds`/`permissionIds` bukan array |
| 400 | `moduleIds must contain valid module IDs` | Semua ID module tidak valid |
| 400 | `id does not match roleId/moduleIds` | Anchor tidak milik `roleId` / `moduleIds` tidak memuat module anchor |
| 404 | `role tidak ditemukan` | `roleId` tidak ada |
| 400 | `cannot modify SUPERADMIN role permissions` | Target role `SUPERADMIN` |
| 404 | `module tidak ditemukan` | Ada `moduleId` yang tidak ada |
| 403 | `step-up verification required` / `422` / `403 mismatch` | Sama seperti 3.1 (digest UPDATE memuat `id`) |
| 404 | `permission tidak ditemukan` | Ada `permissionId` tidak ada (hanya dicek jika `permissionIds` tidak kosong) |
| 400 | `permission does not match selected modules` | Prefix permission tak sesuai module |
| 500 | `internal server error` | Exception tak tertangani |

---

### 3.3 Delete Group — `DELETE /api/v1/dashboard/role-permissions/group?roleId=4&moduleId=3`

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Role permission group deleted successfully",
  "data": { "deletedCount": 5 }
}
```

**Gagal — urutan: param → role/module → SUPERADMIN → step-up**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `roleId and moduleId are required` | Salah satu param ≤ 0/kosong |
| 404 | `role tidak ditemukan` / `module tidak ditemukan` | ID tidak ada |
| 400 | `cannot modify SUPERADMIN role permissions` | Target role `SUPERADMIN` |
| 403 / 422 / 403 | step-up errors | Sama seperti 3.1 (digest DELETE = `{"roleId":N,"moduleId":N}`) |

---

### 3.4 Delete Row — `DELETE /api/v1/dashboard/role-permissions/{id}`

**Tanpa step-up token.**

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Role permission deleted successfully",
  "data": null
}
```

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `Invalid ID` | `{id}` bukan digit (non-numeric juga tidak match route → 404) |
| 404 | `role permission tidak ditemukan` | Row tidak ada |

> Tidak ada guard `cannot modify SUPERADMIN` di endpoint ini — hanya row-by-row.

---

## 4. Step-up Challenge

Header umum untuk ketiganya:
```
Authorization: Bearer <full-access-token>
X-Step-Up-Context: role-permission        ← wajib, selain itu 400 invalid step-up context
Content-Type: application/json
```

### 4.1 Create Challenge — `POST /api/v1/dashboard/role-permissions/security/challenges`

**Request**
```json
{
  "password": "<sha256-64hex>",
  "intent": "CREATE",
  "payloadDigest": "1e114392486f9e6fd7f512f2d9901f2f7c438a618c1fa5b30a22d51989ef2365",
  "requestActionId": 0,
  "module": "users",
  "action": "create"
}
```
| Field | Wajib | Keterangan |
|---|---|---|
| `password` | Ya | Password SUPERADMIN (SHA-256 hex; plaintext otomatis di-hash) |
| `intent` | Ya | `CREATE` / `UPDATE` / `DELETE` |
| `payloadDigest` | Ya | SHA-256 hex 64 huruf kecil (rumus §1.5) |
| `requestActionId` | Tidak | ID intent log Request Action (lifecycle & status → `core_stepup_security_api.md` §3) |
| `module`, `action` | Tidak | Konteks audit |

**Sukses — 201 Created**
```json
{
  "status": true,
  "message": "Role permission security challenge created",
  "data": {
    "challengeId": "rpsec_<32-hex>",
    "requestActionId": null,
    "expiresAt": "2026-04-22T12:00:00+07:00",
    "resendAvailableAt": "2026-04-22T12:01:00+07:00"
  }
}
```
Efek samping: OTP 6 digit dikirim ke email, berlaku 5 menit.

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 400 | `invalid step-up context` | Header `X-Step-Up-Context` ≠ `role-permission` |
| 403 | `permission denied` | Role aktif bukan `SUPERADMIN` |
| 401 | `user not found` | User pada token tidak ada |
| 400 | `password is required` | `password` kosong |
| 401 | `password is wrong` | Password salah |
| 400 | `intent must be CREATE, UPDATE, or DELETE` | Intent tak dikenal |
| 400 | `payloadDigest must be valid sha256 hex` | Bukan 64-hex |
| 404 | `challenge not found` | (lihat 4.2/4.3) |
| 500 | `internal server error` | Exception tak tertangani |

---

### 4.2 Verify Challenge — `POST /api/v1/dashboard/role-permissions/security/challenges/{challengeId}/verify`

**Request**
```json
{ "otpCode": "123456" }
```

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Role permission security challenge verified",
  "data": {
    "verifiedToken": "<encrypted-token>",
    "verifiedUntil": "2026-04-22T12:10:00+07:00"
  }
}
```
- `verifiedToken` → dikirim sebagai header `X-Role-Permission-Verified-Token` pada 3.1–3.3.
- Efek samping: OTP dihapus, challenge `verified_at` terisi; jika terkait `requestActionId`, status Request Action → `CHALLENGE_VERIFIED`.

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 400 | `invalid step-up context` | Header context salah |
| 400 | `challengeId is required` | Path kosong |
| 400 | `otpCode is required` | `otpCode` kosong |
| 404 | `challenge not found` | Challenge bukan milik user aktif |
| 410 | `challenge expired` | OTP lewat 5 menit |
| 422 | `challenge already used` | Challenge sudah `consumed` |
| 401 | `invalid OTP code` | Salah (percobaan ke-1–4) |
| 429 | `too many invalid otp attempts` | Salah ≥ 5 kali |
| 500 | `internal server error` | Exception tak tertangani |

---

### 4.3 Resend Challenge OTP — `POST /api/v1/dashboard/role-permissions/security/challenges/{challengeId}/resend`

**Request:** tanpa body.

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Role permission security challenge OTP resent",
  "data": {
    "challengeId": "rpsec_<32-hex>",
    "expiresAt": "2026-04-22T12:05:00+07:00",
    "resendAvailableAt": "2026-04-22T12:06:00+07:00"
  }
}
```
Efek samping: OTP baru + expiry 5 menit baru, cooldown 60 detik.

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid step-up context` | Header context salah |
| 404 | `challenge not found` | Challenge bukan milik user |
| 410 | `challenge expired` | OTP lama sudah lewat 5 menit |
| 422 | `challenge already used` | Challenge sudah dipakai |
| 429 | `resend OTP is not available yet` | < 60 detik sejak kirim terakhir |
| 500 | `internal server error` | Exception tak tertangani |

---

## 5. Error Middleware (berlaku semua endpoint)

Bentuk: `{"error": "...", "code": "..."}` — **berbeda** dari error bisnis `{status, message, data}`.

| HTTP | `error` | `code` |
|---|---|---|
| 401 | `missing authorization header` | `AUTH_MISSING_HEADER` |
| 401 | `invalid format header` | `AUTH_INVALID_HEADER_FORMAT` |
| 401 | `invalid or expired token` | `AUTH_INVALID_OR_EXPIRED_TOKEN` |
| 401 | `user not found` | `AUTH_USER_NOT_FOUND` |
| 401 | `user is not verified, please verify OTP first` | `AUTH_USER_NOT_VERIFIED` |
| 403 | `permission denied` | `AUTHZ_ROLE_DENIED` (role ≠ SUPERADMIN) |

---

## 6. Catatan Deviasi Doc ↔ Implementasi

1. **`X-Step-Up-Context: role-permission`** wajib untuk endpoint challenge — tidak terdokumentasi di doc referensi.
2. Doc tidak mencatat beberapa **404** pada create: `role/module/permission tidak ditemukan`.
3. **Urutan validasi tidak seragam:** create memvalidasi permission-module sebelum step-up; update memvalidasi step-up sebelum permission-module; delete group memvalidasi SUPERADMIN target sebelum step-up.
4. **`DELETE /group`** punya guard `cannot modify SUPERADMIN role permissions` yang tidak ada di doc.
5. Doc menyebut `404 role permission tidak ditemukan` untuk update — sesuai kode (`RolePermissionController.php:446`).
6. `permissionIds: []` pada **update** sah = mencabut semua permission module tsb; pada **create** ditolak `400`.
7. Create mengembalikan **201**, endpoint lain 200.
