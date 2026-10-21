# Core Step-up Security API Contract (Core 5)

Dokumen kontrak inti **Step-up Security & Request Action**: Security Policy CRUD (aturan risiko per aksi), Request Action CRUD (intent log sebelum aksi sensitif), Role Seal (bukti role aktif ber-TTL), serta titik temu dengan step-up challenge (Core 2).

- **Sumber kebenaran:** implementasi backend (`routes/api.php:29-60`, `SecurityPolicyController`, `RequestActionController`, `AuthController::sealRole/verifyRoleSeal`, `AuthService`, model `SecurityPolicy`/`RequestAction`, migration `2026_05_17_090100/090300`) — diverifikasi terhadap dokumen CRUD request-action internal.
- **Base path:** `/api/v1/dashboard`
- **Terkait:** `core_role_permission_api.md` §4 (step-up challenge & efek `CHALLENGE_VERIFIED`), §5 (error middleware `role.strict`); `core_auth_api.md` (sesi `FULL_AUTH`); `core_crud_api.md` (pola 6 operasi & envelope — dua domain ini **sebagian mematuhi**, lihat §7).
- **Versi dokumen:** disusun 2026-10-21.

---

## 1. Ringkasan Kontrak

### 1.1 Proteksi Route

| Endpoint prefix | Middleware | Catatan |
|---|---|---|
| `/role-permissions/*` (challenges, security-policies, request-actions, meta) | `auth.jwt` + **`role.strict:SUPERADMIN`** | Tanpa `permission:` — role HARUS `SUPERADMIN` persis; selain itu `403 {error:"permission denied", code:"AUTHZ_ROLE_DENIED"}` (lihat `core_role_permission_api.md` §5) |
| `/auth/role/seal`, `/auth/role/verify` | `auth.jwt` saja | Terbuka untuk semua role terautentikasi |

**Tidak ada step-up challenge** pada mutasi security-policy & request-action — perlindungannya cukup `role.strict`. Step-up challenge (`X-Step-Up-Context: role-permission`, password→OTP→`verifiedToken`) **wajib hanya untuk mutasi role-permission utama** — `core_role_permission_api.md` §3–§4.

### 1.2 Format Response

Sama dengan core lain (`core_auth_api.md` §1.1): sukses `{status:true, message, data}`; gagal bisnis `{status:false, message, data:""}`; gagal middleware `{error, code}`. **Semua error di dokumen ini (termasuk 400/401/403/404/409/422) berbentuk error bisnis** — dihasilkan `ApiResponse::error()`, bukan `middlewareError()`.

### 1.3 Model Data

**SecurityPolicy** (entity + nested targets):
```ts
{ id: number; name: string; description: string | null; isActive: boolean;
  createdAt: string | null; updatedAt: string | null;          // "Y-m-d H:i:s"
  targets: { id: number; module: string; action: string;
    riskLevel: 'LOW'|'MEDIUM'|'HIGH'|'CRITICAL';
    validationType: 'CONFIRM'|'NOTE'|'PASSWORD'|'OTP';
    requireNote: boolean; isActive: boolean; }[] }
```

**RequestAction:**
```ts
{ id: number; userId: number; userName: string; module: string; action: string;
  entityType: string | null; entityId: string | null;
  payload: Record<string, unknown> | null; payload_hash?: string;   // hash hanya di DB
  riskLevel: 'LOW'|'MEDIUM'|'HIGH'|'CRITICAL';
  validationType: 'CONFIRM'|'NOTE'|'PASSWORD'|'OTP';
  status: 'PENDING'|'CHALLENGE_REQUIRED'|'CHALLENGE_VERIFIED'|'EXECUTED'
        | 'REJECTED'|'EXPIRED'|'CANCELLED';                        // lihat §3.6 — hanya 4 pertama diproduksi
  securityChallengeId: number | null; expiresAt: string | null; consumedAt: string | null;
  createdAt: string | null; updatedAt: string | null;
  challenges: { challengeId: string; status: string | null; verifiedAt: string | null }[] }  // hanya di detail
```

### 1.4 Lifecycle Request Action

```
1. POST /security-policies          daftarkan target (module+action → riskLevel+validationType)
2. POST /request-actions            buat intent — status awal dari validationType:
                                      CONFIRM → PENDING, selain itu → CHALLENGE_REQUIRED
                                      expires_at = +30 menit, payload_hash = sha256(payload)
3. (jika perlu challenge) POST /security/challenges {requestActionId}
        → verify OTP (Core 2 §4.2) → status request action → CHALLENGE_VERIFIED
4. POST /request-actions/{id}/consume   → status EXECUTED, consumed_at terisi
5. DELETE /request-actions/{id}      hapus intent (hard delete, tanpa guard status)
```

---

## 2. Security Policy CRUD (`/role-permissions/security-policies`)

### 2.1 List — `GET /api/v1/dashboard/role-permissions/security-policies`

| Query | Default | Keterangan |
|---|---|---|
| `search` | `''` | LIKE pada `name` / `description` |
| `page` | 1 (min 1) | `max(1, int)` |
| `size` | 10 (min 1) | `max(1, int)`, **tanpa batas atas** |

Urut `created_at DESC`. Sukses **200** `Successfully retrieved security policies` → envelope `{contents, page, size, totalData}` (lihat `core_crud_api.md` §4.1).

### 2.2 Detail — `GET /api/v1/dashboard/role-permissions/security-policies/{id}`

| HTTP | `message` | Kondisi |
|---|---|---|
| 200 | `Successfully retrieved security policy` | ketemu (data = entity dgn `targets`) |
| 400 | `Invalid ID` | `{id}` bukan digit |
| 404 | **`security policy tidak ditemukan`** | tak ada (bahasa Indonesia — lihat §7.1) |

### 2.3 Create — `POST /api/v1/dashboard/role-permissions/security-policies`

**Request**
```json
{
  "name": "High Risk Mutation",
  "description": "Aksi berisiko tinggi",
  "is_active": true,
  "targets": [
    { "module": "role-permissions", "action": "create",
      "riskLevel": "HIGH", "validationType": "OTP",
      "requireNote": false, "isActive": true }
  ]
}
```

| Field | Wajib | Default / Validasi |
|---|---|---|
| `name` | Ya | kosong → `400 name is required` |
| `description` | Tidak | kosong → `null` |
| `is_active` | Tidak | default `true` |
| `targets` | Ya | array non-kosong → selain itu `400 targets is required` |
| `targets[].module`, `targets[].action` | Ya | kosong → entri **dilewati**; bila semua terskip → `400 targets is required` |
| `targets[].riskLevel` | Tidak | default `LOW`; di luar `LOW/MEDIUM/HIGH/CRITICAL` → fallback `LOW` |
| `targets[].validationType` | Tidak | default `CONFIRM`; di luar `CONFIRM/NOTE/PASSWORD/OTP` → fallback `CONFIRM` |
| `targets[].requireNote` / `isActive` | Tidak | default `false` / `true` |

Sukses **201** `Security policy created successfully` → entity (termasuk `targets` dengan `id` baru).

### 2.4 Update — `PUT /api/v1/dashboard/role-permissions/security-policies/{id}`

Validasi body sama persis §2.3. **Semua target diganti**: `syncTargets` = hard-delete seluruh target lama + insert ulang (id target berubah). 400 `Invalid ID` / 404 `security policy tidak ditemukan`. Sukses **200** `Security policy updated successfully` → entity.

### 2.5 Delete — `DELETE /api/v1/dashboard/role-permissions/security-policies/{id}`

400 `Invalid ID` / 404 `security policy tidak ditemukan`. Sukses **200** `Security policy deleted successfully`, `data: null`. Hard delete — **tanpa** `SoftDeletes`; target ikut hilang (FK cascade/relasi).

---

## 3. Request Action CRUD (`/role-permissions/request-actions`)

### 3.1 List — `GET /api/v1/dashboard/role-permissions/request-actions`

| Query | Default | Keterangan |
|---|---|---|
| `search` | `''` | LIKE pada `module` / `action` / `entity_id` |
| `status` | `''` | filter eksak `status` (nilai tak dikenal → hasil kosong, tanpa 400) |
| `page` / `size` | 1 / 10 (min 1) | tanpa batas atas `size` |

Urut **`id DESC`**. Sukses **200** `Successfully retrieved request actions` → envelope (pesan **beda** dari dokumen FE — lihat §7.2). Relasi `user:id,name` → `userName`.

### 3.2 Detail — `GET /api/v1/dashboard/role-permissions/request-actions/{id}`

| HTTP | `message` |
|---|---|
| 200 | `Successfully retrieved request action` — data = entity **+ `challenges[]`** |
| 400 | `Invalid ID` |
| 404 | `request action not found` (bahasa Inggris — lihat §7.1) |

### 3.3 Create — `POST /api/v1/dashboard/role-permissions/request-actions`

**Request**
```json
{ "module": "role-permissions", "action": "create",
  "entity_type": "role", "entity_id": "3",
  "payload": { "name": "Admin Gudang" } }
```

| Field | Wajib | Keterangan |
|---|---|---|
| `module`, `action` | Ya | kosong → `400 module and action are required` |
| `payload` | Ya | harus JSON object → `400 payload must be object` |
| `entity_type`, `entity_id` | Tidak | `null` bila kosong |

**Langkah server:**
1. Cari `SecurityPolicyTarget` dengan `module`+`action` **dan `is_active=true`** — tak ada → `404 security policy target not found`.
2. Buat row: `user_id` dari token, `riskLevel`/`validationType` **disalin dari target**, `payload_hash = sha256(json_encode(payload, JSON_UNESCAPED_UNICODE|JSON_UNESCAPED_SLASHES))`, `status = PENDING` jika `validationType==CONFIRM` selain itu `CHALLENGE_REQUIRED`, `expires_at = now + 30 menit`.

Sukses **201** `Request action created successfully` → entity.

### 3.4 Update — `PUT /api/v1/dashboard/role-permissions/request-actions/{id}`

| HTTP | `message` | Kondisi |
|---|---|---|
| 400 | `Invalid ID` | id non-digit |
| 400 | `module and action are required` | `module`/`action` kosong — **wajib ada di body meski tidak disimpan** (lihat §7.5) |
| 400 | `payload must be object` | `payload` bukan object |
| 404 | `request action not found` | tak ada |
| **409** | `request action already finalized` | status `EXECUTED` / `CANCELLED` |

Yang berubah: `entity_type`, `entity_id`, `payload` (+ `payload_hash` baru dihitung). `module`/`action`/`riskLevel`/`validationType`/`status` **tidak** berubah. Sukses **200** `Request action updated successfully` → entity.

### 3.5 Consume — `POST /api/v1/dashboard/role-permissions/request-actions/{id}/consume`

| HTTP | `message` | Kondisi |
|---|---|---|
| 200 | `Request action marked as consumed` | sukses → `status=EXECUTED`, `consumed_at=now()` |
| 400 | `Invalid ID` | id non-digit |
| 404 | `request action not found` | tak ada |
| **409** | `request action already executed` | sudah `EXECUTED` |

**Tidak mengecek `expires_at`** — intent yang lewat 30 menit tetap bisa di-consume (lihat §7.4).

### 3.6 Delete — `DELETE /api/v1/dashboard/role-permissions/request-actions/{id}`

400 `Invalid ID` / 404 `request action not found`. Sukses **200** `Request action deleted successfully`, `data: null`. **Tanpa guard status** — record `EXECUTED` tetap bisa dihapus (banding §3.4 yang menolak 409). Hard delete.

### 3.7 Catatan Status (enum DB vs yang diproduksi kode)

| Status | Diproduksi oleh | Keterangan |
|---|---|---|
| `PENDING` | create (validationType `CONFIRM`) | menunggu konfirmasi lalu consume |
| `CHALLENGE_REQUIRED` | create (validationType `NOTE/PASSWORD/OTP`) | perlu challenge dulu |
| `CHALLENGE_VERIFIED` | verify challenge terkait `requestActionId` (`core_role_permission_api.md` §4.2) | siap di-consume |
| `EXECUTED` | consume | terminal untuk update (409), tapi **bisa di-delete** |
| `REJECTED`, `EXPIRED`, `CANCELLED` | **tidak pernah** | ada di enum migration, tak ada kode/scheduler yang mengeset (lihat §7.4) |

---

## 4. Role Seal (`/auth/role/seal`, `/auth/role/verify`)

Membuktikan "user ini memegang role X saat ini" **terikat ke access token yang dipakai** — untuk keputusan aman berbasis role (mis. konfirmasi akses SUPERADMIN) tanpa menyimpan state di server.

### 4.1 Seal — `POST /api/v1/dashboard/auth/role/seal`

Header: `Authorization: Bearer <access-token>` (middleware `auth.jwt`). **Tanpa body.**

Sukses **200** `Role seal created`:
```json
{ "status": true, "message": "Role seal created",
  "data": { "sealedRole": "<string ter-encrypt>", "expiresAt": "2026-10-21T12:00:00+07:00" } }
```

- **TTL** = `min(exp access token, now + 15 menit)` → seal maksimal 15 menit.
- Isi seal (di-encrypt `Crypt`, **tidak terbaca klien**): `{uid, role, bind, iat, exp}` dengan `bind = HMAC-SHA256(access token, ROLE_SEAL_KEY)`.

**Gagal** (bentuk error bisnis):

| HTTP | `message` |
|---|---|
| 401 | `user not found` (`user_id` token ≤ 0) |
| 401 | `invalid or expired token` (access token kosong / `exp` sudah lewat / `role` pada token kosong) |

### 4.2 Verify — `POST /api/v1/dashboard/auth/role/verify`

**Request:** `{ "sealedRole": "<dari 4.1>" }`

Sukses **200** `Role seal verified` → `data: { "isValid": true, "isSuperadmin": <role aktif == "SUPERADMIN"> }`.

**Gagal:**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | body bukan JSON valid (`parsePayload` null) |
| 400 | `sealedRole is required` | field kosong |
| 401 | `invalid or expired token` | access token kosong |
| **422** | `invalid sealed role` | decrypt gagal / field seal tak lengkap (bukan seal milik aplikasi ini) |
| **422** | `sealed role expired` | `exp` seal sudah lewat |
| 401 | `role seal mismatch` | `uid`/`role` berbeda dari token aktif, **atau** `bind` HMAC tidak cocok — artinya access token sudah berganti (pernah di-refresh) sejak seal dibuat |

### 4.3 Sifat Keamanan

1. Seal **bukan token pengganti** — selalu dibawa bersama `Authorization` aktif yang sama.
2. Karena `bind` memakai access token: **refresh token membatalkan semua seal lama** (`role seal mismatch`).
3. Server tidak menyimpan state seal (stateless, ter-encrypt di klien).
4. `verify` tidak membedakan role mana pun — semua role terautentikasi boleh; jawaban `isSuperadmin` yang menentukan perilaku FE.

---

## 5. Integrasi dengan Step-up Challenge (Core 2)

Endpoint challenge ada di grup yang sama dan **sudah terdokumentasi penuh** di `core_role_permission_api.md` §4 — dokumen ini hanya mencatat sambungannya:

| Sambungan | Detail |
|---|---|
| Create challenge menerima `requestActionId` | `core_role_permission_api.md` §4.1 — mengisi `RolePermissionSecurityChallenge.request_action_id` |
| Verify OTP berhasil → request action `CHALLENGE_VERIFIED` | efek samping di `RolePermissionSecurityController:107-109` (§4.2 core 2) |
| Urutan wajib untuk aksi berisiko `OTP` | create request action (`CHALLENGE_REQUIRED`) → challenge → verify → **consume** |
| `X-Step-Up-Context: role-permission` | hanya untuk endpoint challenge — mutasi request-action/security-policy **tidak** memakainya |

---

## 6. Error Middleware Grup Ini

Berlaku untuk **semua** endpoint §2–§3 (bentuk `{error, code}` — lihat `core_role_permission_api.md` §5): `401` `AUTH_MISSING_HEADER` / `AUTH_INVALID_HEADER_FORMAT` / `AUTH_INVALID_OR_EXPIRED_TOKEN` / `AUTH_USER_NOT_FOUND` / `AUTH_USER_NOT_VERIFIED`, dan `403` `permission denied` + `AUTHZ_ROLE_DENIED` (role ≠ `SUPERADMIN`).

---

## 7. Catatan Deviasi Doc ↔ Implementasi

1. **Bahasa 404 campur** dalam satu grup: `security policy tidak ditemukan` (Indo) vs `request action not found` (Inggris) — prinsip satu bahasa (`core_principles.md` §4) dengan deviasi serupa tercatat di `core_crud_api.md` §5 (domain user/branch).
2. **Dokumen CRUD request-action (doc FE) ≠ kode** pada pesan: doc `Request actions fetched successfully` / `Request action fetched successfully` / `Request action consumed successfully`; kode `Successfully retrieved request actions` / `Successfully retrieved request action` / `Request action marked as consumed`. Pesan create/update/delete **cocok**. Contoh timestamp di doc ISO-8601 (`2026-05-19T08:30:00Z`), kode mengembalikan `Y-m-d H:i:s`.
3. **Security policy tidak punya doc sama sekali** di repo aplikasi — kontrak §2 diekstrak murni dari kode (dokumen request-action hanya mencakup §3).
4. **Status `REJECTED`/`EXPIRED`/`CANCELLED` tidak pernah diproduksi**: tidak ada endpoint/scheduler yang mengesetnya; `consume` tidak mengecek `expires_at` (TTL 30 menit hanya informatif); `CANCELLED` hanya disebut di guard update.
5. **Update meminta `module`/`action` wajib tapi mengabaikan nilainya** — divalidasi (`400`) lalu tidak disimpan (FE menampilkan read-only).
6. **Delete request-action tanpa guard status** — `EXECUTED` bisa dihapus (audit trail hilang), sedangkan update menolaknya (409).
7. **Mutasi security-policy & request-action tanpa step-up** — hanya `role.strict:SUPERADMIN`; berbeda dengan mutasi role-permission utama yang butuh challenge. Desain: request-action memang "wadah" intent-nya sendiri.
8. **Tidak memakai pola core 3 penuh**: tanpa endpoint `/names` lookup, tanpa `SoftDeletes` (hard delete), urutan list beda (`id DESC` vs `created_at DESC`), `size` tanpa batas atas. Dicatat sebagai varian, bukan pelanggaran — dua domain ini adalah konfigurasi keamanan, bukan data master.
9. **Fallback validasi diam-diam**: `riskLevel`/`validationType` tak dikenal di-set `LOW`/`CONFIRM` (bukan 400); target dengan `module`/`action` kosong dilewati tanpa pesan.
