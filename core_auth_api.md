# Core Auth API Contract

Dokumen kontrak API autentikasi: endpoint, request, response sukses, dan response gagal.

- **Sumber kebenaran:** implementasi backend (`routes/api.php`, `AuthController`, `AuthService`, `ApiResponse`, middleware) diverifikasi terhadap dokumen alur internal (bagian 7.1–7.8).
- **Base path:** `/api/v1/dashboard/auth`
- **Versi dokumen:** disusun 2026-10-21.

---

## 1. Ringkasan Kontrak

### 1.1 Format Response

**Sukses (semua endpoint, HTTP 200):**
```json
{
  "status": true,
  "message": "<pesan sukses>",
  "data": <objek | null | "">
}
```

**Gagal — error bisnis (dilempar `AppError` dari service, di-catch controller):**
```json
{
  "status": false,
  "message": "<pesan error persis dari service>",
  "data": ""
}
```

**Gagal — error middleware (auth/otorisasi, bentuk BERBEDA):**
```json
{
  "error": "<pesan error>",
  "code": "<ERROR_CODE>"
}
```

**Gagal — body bukan JSON valid:** HTTP 400 `{"status": false, "message": "invalid input", "data": ""}`.

### 1.2 Header

| Endpoint | Header wajib |
|---|---|
| Semua endpoint auth di bawah (publik) | `Content-Type: application/json` |
| Endpoint terproteksi / 7.6 versi riil | `Authorization: Bearer <access-token>` |

### 1.3 Aturan Password & Kode

- Field `password` / `confirmationPassword` wajib dikirim FE dalam format **SHA-256 hex 64 karakter**.
- Jika nilai bukan hex 64 karakter, backend otomatis di-hash SHA-256 (`AuthService::normalizeSha256Password`), lalu di-bcrypt untuk penyimpanan.
- **OTP:** 6 digit angka (`000000`–`999999`), berlaku **5 menit**, satu aktif per user (regenerasi menimpa yang lama).
- **Activation code & Reset code:** JWT satu-pemakaian, berlaku **30 menit**.
- Password dan confirmationPassword **wajib identik** (activation & reset).

### 1.4 Token

| Token | TTL | Keterangan |
|---|---|---|
| Access token (`jwt_expiration_seconds`) | 3600 detik | Hasil `login` (PARTIAL) dan `verify-code` (final) |
| Refresh token (`refresh_token_expiration_seconds`) | 604800 detik (7 hari) | Hasil `verify-code`; dirotasi oleh `refresh-token` |
| OTP | 5 menit | — |
| Activation code | 30 menit | — |
| Reset code | 30 menit | — |

### 1.5 State Sesi

`ANONYMOUS -> PARTIAL_AUTH (login sukses) -> FULL_AUTH (verify OTP sukses)`

Endpoint terproteksi hanya menerima token dari sesi `FULL_AUTH` (`is_login_verified = true`).

---

## 2. Endpoint

### 2.1 Login — `POST /api/v1/dashboard/auth/login`

**Request**
```json
{
  "email": "admin@example.com",
  "password": "<sha256-64hex>"
}
```

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Login success",
  "data": {
    "authLevel": "PARTIAL",
    "token": "<jwt-partial-token>",
    "expiresAt": "2026-04-22T12:00:00Z"
  }
}
```
Efek samping: OTP 6 digit dibuat (expiry 5 menit) dan dikirim ke email; `is_login_verified = false`.

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 401 | `email or password is wrong` | Email tidak terdaftar, password salah, atau akun belum aktivasi (password kosong) |
| 403 | `email or password is wrong` | Role `CHECKER_GUDANG` / `MAKER_GUDANG` — wajib login via varian ops (login NIP), bukan dashboard |
| 500 | `internal server error` | Exception tak tertangani |

> Pesan 401/403 sengaja sama (generik) agar tidak membocorkan keberadaan akun.

---

### 2.2 Verify OTP — `POST /api/v1/dashboard/auth/verify-code`

**Request**
```json
{
  "email": "admin@example.com",
  "otpCode": "123456",
  "accessType": "DASHBOARD"
}
```
> Backend hanya membaca `email` dan `otpCode`; `accessType` diabaikan (kompatibilitas kontrak blueprint).

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Verify otp success",
  "data": {
    "branchId": "1",
    "name": "Admin",
    "email": "admin@example.com",
    "token": "<full-access-token>",
    "refreshToken": "<refresh-token>",
    "expiresAt": "2026-04-22T12:00:00Z",
    "abilities": [
      { "resource": "dashboard", "actions": ["view"] }
    ]
  }
}
```
- `abilities` dikelompokkan dari permission role: `[ { "resource": "<module>", "actions": ["<aksi>", ...] } ]`.
- `branchId` bisa `""` (string kosong) jika user tanpa branch.
- Efek samping: OTP dihapus, `is_login_verified = true`, access+refresh token disimpan ke user.

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 401 | `user not found` | Email tidak terdaftar |
| 401 | `OTP not requested` | Tidak ada OTP aktif (belum login / sudah dipakai) |
| 401 | `OTP expired` | OTP lewat 5 menit |
| 401 | `invalid OTP code` | Kode salah |
| 500 | `internal server error` | Exception tak tertangani |

---

### 2.3 Resend OTP — `POST /api/v1/dashboard/auth/resend-code`

**Request**
```json
{
  "email": "admin@example.com"
}
```

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Resend otp success",
  "data": null
}
```
Efek samping: OTP baru (expiry baru 5 menit) menimpa OTP lama, dikirim ulang via email.

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 401 | `user not found` | Email tidak terdaftar |
| 500 | `internal server error` | Exception tak tertangani |

---

### 2.4 Activation Check — `POST /api/v1/dashboard/auth/activation/check`

**Request**
```json
{
  "activationCode": "<activation-token-jwt>"
}
```
> `activationCode` berupa JWT dengan claim `email`.

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Activation code is valid",
  "data": { "isSuccess": true }
}
```
Efek samping: `activation_status` user menjadi `ON_ACTIVATION`.

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 400 | `activationCode is required` | Field kosong |
| 401 | `activation code is invalid` | JWT tidak bisa diparse / claim `email` kosong / belum pernah generate token / token tidak cocok dengan data user |
| 401 | `user not found` | Email pada token tidak terdaftar |
| 401 | `activation code has been used` | Akun sudah aktivasi (`email_verified_at` terisi) |
| 401 | `activation code is expired` | Lewat 30 menit |
| 500 | `internal server error` | Exception tak tertangani |

---

### 2.5 Activation Confirm — `POST /api/v1/dashboard/auth/activation/confirm`

**Request**
```json
{
  "activationCode": "<activation-token-jwt>",
  "password": "<sha256-64hex>",
  "confirmationPassword": "<sha256-64hex>"
}
```

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Activation confirmed successfully",
  "data": { "isSuccess": true }
}
```
Efek samping: password di-hash & disimpan, `email_verified_at = now`, `activation_status = ACTIVATED`, token aktivasi di-null-kan (one-time use).

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 400 | `activationCode, password, and confirmationPassword are required` | Salah satu field kosong |
| 400 | `password and confirmationPassword do not match` | Password tidak sama |
| 401 | `activation code is invalid` / `user not found` / `activation code has been used` / `activation code is expired` | Sama seperti 2.4 |
| 500 | `internal server error` | Exception tak tertangani |

---

### 2.6 Activation Resend

**Blueprint:** `POST /api/v1/dashboard/auth/activation/resend` — **TIDAK ADA di route backend** (`routes/api.php:37`: *"resend activation must be accessed from dashboard after login"*).

**Endpoint riil:** `POST /api/v1/dashboard/users/{id}/resend-activation`
- `id` = ID user (digit)
- Header: `Authorization: Bearer <access-token>` (FULL_AUTH)
- Permission: `users,resend_activation`

**Request** (tanpa body; hanya path param `id`)
```
POST /api/v1/dashboard/users/12/resend-activation
Authorization: Bearer <full-access-token>
```

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Activation code resent successfully",
  "data": null
}
```
Efek samping: activation JWT baru (30 menit) dibuat, `activation_status = REQUEST_ACTIVATION`, email dikirim.

**Gagal**

| HTTP | `message` / `error` | Penyebab |
|---|---|---|
| 400 | `Invalid ID` | `id` bukan digit |
| 400 | `invalid input` | Body bukan JSON valid |
| 401 | `account already activated` | User sudah aktivasi |
| 401 | (middleware) `missing authorization header` dkk. | Lihat bagian 3 |
| 403 | `permission denied` | Role tanpa permission `users,resend_activation` |
| 404 | `user not found` | ID tidak ada |
| 500 | `internal server error` | Exception tak tertangani |

> Catatan: doc referensi lama menyebut path `auth/activation/resend` dengan `data: {isSuccess:true}` — kedua hal itu tidak berlaku pada implementasi saat ini.

---

### 2.7 Forgot Password — `POST /api/v1/dashboard/auth/forgot/password`

**Request**
```json
{
  "email": "user@example.com",
  "accessType": "DASHBOARD"
}
```
> Backend hanya membaca `email`; claim `accessType` pada reset token selalu `DASHBOARD` (hardcoded di `JwtService::generateActionToken`).

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Reset password code sent successfully",
  "data": null
}
```
Efek samping: reset JWT (30 menit) disimpan ke user (`forgot_password_token`, `forgot_password_token_expired_at`) dan dikirim via email.

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 400 | `email is required` | Email kosong |
| 401 | `user not found` | Email tidak terdaftar |
| 500 | `internal server error` | Exception tak tertangani |

---

### 2.8 Reset Password — `POST /api/v1/dashboard/auth/reset/password`

**Request**
```json
{
  "resetCode": "<reset-token-jwt>",
  "password": "<sha256-64hex>",
  "confirmationPassword": "<sha256-64hex>"
}
```

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Password reset successfully",
  "data": { "isSuccess": true }
}
```
Efek samping: password baru di-hash & disimpan, reset token di-null-kan (one-time use).

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 400 | `resetCode, password, and confirmationPassword are required` | Salah satu field kosong |
| 400 | `password and confirmationPassword do not match` | Password tidak sama |
| 401 | `reset code is invalid` | JWT tidak bisa diparse / claim `action != reset_password` atau `accessType != DASHBOARD` / user tidak ada / belum pernah request / token tidak cocok dengan data user |
| 401 | `user not found` | `userUuid` pada token tidak ditemukan |
| 401 | `reset code is expired` | Lewat 30 menit |
| 500 | `internal server error` | Exception tak tertangani |

---

### 2.9 Refresh Token — `POST /api/v1/dashboard/auth/refresh-token`

**Request**
```json
{
  "refreshToken": "<refresh-token>"
}
```

**Sukses — 200 OK**
```json
{
  "status": true,
  "message": "Refresh token success",
  "data": {
    "token": "<new-access-token>",
    "refreshToken": "<new-refresh-token>",
    "expiresAt": "2026-04-22T12:00:00Z"
  }
}
```
Efek samping: **rotasi** — access & refresh token lama diganti total di record user (refresh token lama langsung tidak berlaku).

**Gagal**

| HTTP | `message` | Penyebab |
|---|---|---|
| 400 | `invalid input` | Body bukan JSON valid |
| 400 | `refreshToken is required` | Field kosong |
| 401 | `refresh token is invalid or expired` | JWT tidak valid / claim tidak lengkap / user tidak ada / email tidak cocok / refresh token tidak cocok dengan record / sudah lewat 7 hari |
| 401 | `user is not verified, please verify OTP first` | Sesi belum `FULL_AUTH` |
| 500 | `internal server error` | Exception tak tertangani |

---

## 3. Error Middleware — Endpoint Terproteksi

Berlaku untuk semua endpoint dengan `auth.jwt` / `permission` / `role` (termasuk 2.6). Bentuk body berbeda dari error bisnis: `{"error": "...", "code": "..."}`.

| HTTP | `error` | `code` | Syarat |
|---|---|---|---|
| 401 | `missing authorization header` | `AUTH_MISSING_HEADER` | Tanpa header `Authorization` |
| 401 | `invalid format header` | `AUTH_INVALID_HEADER_FORMAT` | Bukan skema `Bearer <token>` |
| 401 | `invalid or expired token` | `AUTH_INVALID_OR_EXPIRED_TOKEN` | Signature/klaim/token tersimpan tidak cocok atau kedaluwarsa |
| 401 | `user not found` | `AUTH_USER_NOT_FOUND` | `sub` token tidak ada di database |
| 401 | `user is not verified, please verify OTP first` | `AUTH_USER_NOT_VERIFIED` | `is_login_verified = false` (belum verify OTP) |
| 403 | `permission denied` | `AUTHZ_PERMISSION_DENIED` | Permission (`permission:`) tidak dimiliki role |
| 403 | `permission denied` | `AUTHZ_ROLE_DENIED` | Role (`role:` / `role.strict:`) tidak sesuai |

**Umum:** endpoint yang tidak terdaftar → **404 Not Found** (default Laravel).

---

## 4. Catatan Deviasi Blueprint ↔ Implementasi

1. **7.6** `POST /auth/activation/resend` tidak ada di route — gunakan `POST /users/{id}/resend-activation` (terproteksi), `data` bernilai `null` (bukan `isSuccess`).
2. **7.2 & 7.7** field `accessType` pada request diabaikan backend.
3. **7.7** mengembalikan `401 user not found` (membocorkan keberadaan akun) — blueprint §6 meminta pesan generik.
4. **Password** plaintext tetap diterima (otomatis di-hash), berbeda dari catatan dokumen alur internal yang menyebut plaintext ditolak 400.
5. **7.2** response memuat `branchId`, `name`, `email` tambahan di luar blueprint.
6. **Error middleware** memakai bentuk `{error, code}` — bukan `{status, message, data}`; FE perlu menangani dua bentuk error.
7. **Refresh token dirotasi** setiap `refresh-token` dipanggil — token lama langsung hangus.
