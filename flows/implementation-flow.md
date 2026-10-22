# Implementation Flow — Urutan Menerapkan 5 Core

Panduan **urutan bootstrap**: dari project kosong → kelima core terpasang penuh. Ini **alur kerja**, bukan kontrak (payload & pesan baku → `core_*_api.md`), bukan tutorial stack (letak folder → `structures/`), bukan prinsip (`core_principles.md`).

- **Kapan dipakai:** project baru, atau project lama yang direfactor agar memenuhi standar 5 core.
- **Kapan TIDAK:** project yang sudah berjalan dan hanya menambah modul bisnis → ikuti `core_principles.md` §7.
- **Aturan otoritas tetap berlaku** (`README.md` §2): kode sumber > kontrak > alur ini > status. Bila isi alur ini beda dari kode/kontrak — **alur inilah yang diperbaiki**.

---

## 1) Diagram Dependensi Antar-Fase

```
Fase 0         Fase 1            Fase 2                  Fase 3            Fase 4
Fondasi   ──→  Auth        ──→   CRUD + Pagination   ──→  RBAC        ──→   Step-up Security
helper,        9 endpoint,        6 operasi,               role/permission,    security-policy,
error handler  JWT, OTP, sesi     envelope, /names         challenge OTP       request-action, seal
```

| Fase | Prasyarat mutlak dari fase sebelum | Boleh paralel? |
|---|---|---|
| 0 | — | — |
| 1 | Fase 0 (helper response & error handler) | — |
| 2 | Fase 0 + **login & middleware JWT** dari Fase 1 | Ya — sisa Fase 1 (activation, forgot/reset, refresh) boleh menyusul |
| 3 | Fase 1 (**OTP** untuk challenge) + Fase 2 (pola list/read) | Tidak — permission justru dipasang balik ke route Fase 2 |
| 4 | Fase 1 (OTP) + Fase 3 (`role.strict:SUPERADMIN`, challenge yang dipinjam) | Tidak |

**Bukti dependensi dari kontrak:** Fase 4 meminjam challenge OTP milik Core 2 (`core_stepup_security_api.md` §1.4 langkah 3 → `core_role_permission_api.md` §4.2); proteksi Fase 4 memakai `role.strict:SUPERADMIN` yang baru ada setelah role dibuat di Fase 3 (`core_stepup_security_api.md` §1.1).

---

## 2) Fase 0 — Fondasi (tanpa endpoint bisnis)

**Prasyarat:** stack dipilih → ikuti `structures/<stack>/README.md` (tree + pemetaan).

**Langkah:**
1. Helper **3 bentuk response** → `core_auth_api.md` §1.1: sukses `{status, message, data}`; gagal bisnis `{status:false, message, data:""}`; gagal middleware `{error, code}` — jangan dicampur.
2. Helper **envelope list** `{contents, page, size, totalData}` → `core_crud_api.md` §4.1.
3. **Error handler global** (500/panic) memakai bentuk `{error, code}` → `core_auth_api.md` §3.
4. Konvensi route: prefix `/api/v1/...` (kontrak memakai `/api/v1/dashboard`) + default **`page=1`, `size=10`** (`core_crud_api.md` §4.2).

**Checklist verifikasi:**
- [ ] `GET /health` membalas envelope sukses `{status:true, ...}`.
- [ ] Paksa error 500 → balasan `{"error": ..., "code": ...}` (bukan `{status,...}`).
- [ ] List dummy mengembalikan persis 4 kunci envelope.

---

## 3) Fase 1 — Core 1: Auth (9 endpoint)

**Prasyarat:** tabel user (dengan status verified/activation), penyimpanan OTP, kanal kirim email/OTP (job/queue), helper Fase 0.
**Kontrak penuh:** `core_auth_api.md` §1–§3.

**Langkah berurutan:**
1. **Login** `POST /auth/login` (§2.1) → sesi `PARTIAL_AUTH`; aturan password/kode §1.3.
2. **Verify OTP** (§2.2) → sesi `FULL_AUTH`; **Resend** (§2.3).
3. **Activation** check / confirm / resend (§2.4–§2.6).
4. **Forgot** (§2.7) + **Reset password** (§2.8).
5. **Refresh token** (§2.9) — rotasi sekali pakai; TTL access `3600s`, refresh `604800s` (§1.4).
6. **Middleware proteksi** untuk semua endpoint terproteksi — daftar error `AUTH_*`/`AUTHZ_*` di `core_auth_api.md` §3; hanya token `FULL_AUTH` (`is_login_verified = true`, §1.5).

> Langkah 1–5 boleh sedikit bertukar urutan, tetapi **langkah 6 wajib selesai** sebelum Fase 2 memakai JWT.

**Checklist verifikasi:**
- [ ] Login sukses → 200 + data sesi §2.1; kredensial salah → pesan generik (tidak membocorkan akun ada/tidak).
- [ ] Tanpa header / format salah → 401 `AUTH_MISSING_HEADER` / `AUTH_INVALID_HEADER_FORMAT` (§3).
- [ ] Token kadaluwarsa → 401 `AUTH_INVALID_OR_EXPIRED_TOKEN`; belum verify OTP → `AUTH_USER_NOT_VERIFIED`.
- [ ] Endpoint tidak terdaftar → 404 (default, §3).

---

## 4) Fase 2 — Core 3+4: CRUD + Pagination

**Prasyarat:** helper Fase 0 + JWT dari Fase 1. Proteksi awal cukup `auth.jwt` — `permission:` dipasang di Fase 3.
**Kontrak penuh:** `core_crud_api.md` §1, §3, §4, §6.

**Langkah:**
1. Pilih **satu domain contoh** → bangun **6 operasi** (§1.2) dengan pesan & kode status persis dari tabel standar.
2. List: query `search,page,size` + envelope (§4.1), rumus offset/total (§4.3), `ORDER BY` eksplisit (§4.3 catatan 1).
3. Lookup **`GET /<plural>/names`** untuk tiap data master yang dirujuk FK (§3).
4. Validasi: `400 Invalid ID` / `<field> is required` / duplikat, `404 <resource> tidak ditemukan` (Bahasa Indonesia), soft delete untuk data master (§6 langkah 5).
5. Replikasi ke domain lain mengikuti **Konvensi Menambah Modul CRUD Baru** (§6, 7 langkah).

**Checklist verifikasi:**
- [ ] Create → **201**; list/detail/update/delete → 200; delete → `data: null` (§1.2).
- [ ] List balas envelope 4 kunci; `page=0` di-reset ke `1`, `size=0` ke `10` (§4.2).
- [ ] Detail id non-digit → 400 `Invalid ID`; id tak ada → 404.
- [ ] `/names` terpaginasi (memakai envelope yang sama, §3).

---

## 5) Fase 3 — Core 2: RBAC (role, permission, challenge)

**Prasyarat:** Fase 1 (JWT + OTP untuk challenge) + Fase 2 (pola read).
**Kontrak penuh:** `core_role_permission_api.md` §1–§5.

**Langkah:**
1. Model role + permission (+module) + seed role `SUPERADMIN` (§1.1).
2. **Middleware permission** `<resource>,<action>` — default deny, SUPERADMIN bypass (`core_principles.md` §3); kode error `AUTHZ_PERMISSION_DENIED`/`AUTHZ_ROLE_DENIED` (daftar di `core_auth_api.md` §3). Pasang ke seluruh route Fase 2.
3. **Endpoint read** role-permission (§2.1–§2.6) — tanpa step-up, semua `role.strict:SUPERADMIN` (§1.1).
4. **Step-up challenge** (§4.1–§4.3): header `X-Step-Up-Context: role-permission` → password → OTP → `verifiedToken` (alur dua lapis §1.3); hitung `payloadDigest` (§1.5 — **wajib identik** dengan payload mutasi); patuhi batas §1.4 (OTP 5 mnt, cooldown 60 dtk, token 10 mnt, 5× salah → 429).
5. **Mutasi** create/update/delete-group (§3.1–§3.3) dengan header `X-Role-Permission-Verified-Token` — kecuali **`DELETE /{id}` (hapus satu row) tanpa step-up** (§1.3).

**Checklist verifikasi:**
- [ ] Contoh `payloadDigest` di §1.5 direproduksi → hex sama persis.
- [ ] Mutasi tanpa `verifiedToken` → ditolak; dengan token terpakai ulang → ditolak (satu pemakaian, §1.4).
- [ ] Non-SUPERADMIN akses endpoint grup ini → 403 `AUTHZ_ROLE_DENIED` (§5).
- [ ] Create → **201**, endpoint lain → 200 (§6 catatan 7).
- [ ] `DELETE /{id}` tetap jalan tanpa challenge.

---

## 6) Fase 4 — Core 5: Step-up Security

**Prasyarat:** Fase 3 (`role.strict:SUPERADMIN`, challenge OTP) + Fase 1 (OTP).
**Kontrak penuh:** `core_stepup_security_api.md` §1–§5.

**Langkah:**
1. Migration model `SecurityPolicy (+Target)` & `RequestAction` — bentuk field di §1.3.
2. **Security-policy CRUD** (§2.1–§2.5) — proteksi `auth.jwt` + `role.strict:SUPERADMIN`, **tanpa step-up challenge** (§1.1; cukup role persis).
3. **Request-action** (§3): create intent — cari target aktif `module+action` dulu (tak ada → 404 `security policy target not found`, §3.3), status awal dari `validationType`, `expires_at = +30 menit` (§1.4) → **consume** (§3.5, `EXECUTED` + `consumed_at`). Challenge opsional memakai Core 2 §4.2 (§1.4 langkah 3).
4. **Role seal** (§4.1–§4.2): proteksi `auth.jwt` saja (§1.1); TTL `min(exp access token, now + 15 menit)` (§4.1); sifat keamanan §4.3 (stateless, terikat access token — refresh membatalkan seal lama).
5. **Integrasi** challenge ↔ request-action (§5).

**Checklist verifikasi:**
- [ ] `POST /request-actions` untuk `module+action` tanpa target aktif → 404 (§3.3).
- [ ] Consume saat sudah `EXECUTED` → 409 `request action already executed` (§3.5); update saat finalized → 409 `already finalized` (§3.4); **delete tanpa guard status** (§3.6).
- [ ] Mutasi security-policy/request-action TANPA challenge — cukup `role.strict` (§1.1).
- [ ] Seal non-SUPERADMIN tetap boleh (role apa pun, `auth.jwt` saja); seal kadaluwarsa → `role seal mismatch` (§4.3).
- [ ] Semua error grup ini berbentuk **error bisnis** `{status,...}` — bukan `{error,...}` (§1.2).

---

## 7) Pasca-Fase — Audit & Seremoni

1. **Audit proteksi:** semua route ber-proteksi `auth.jwt` + `permission:` sesuai prinsip (`core_principles.md` §7 langkah 2); route terbuka wajib by-design + terdokumentasi.
2. **Buat profile project** `examples/core_index_<nama-project>.md` (ikuti `README.md` §4) — jangan menimpa profile yang ada.
3. **Cocokkan kode ↔ kontrak:** beda → perbarui kontrak dulu, lalu baris status di `core_status.md` (checklist §4 file itu).
4. **Catat keputusan** baru di Riwayat Keputusan `core_status.md` §3.
5. **Stack baru?** tulis struktur + stub-nya ke `structures/` (panduan: `structures/README.md` §4).

---

**Referensi silang:** prinsip → `core_principles.md` · kontrak → `core_*_api.md` · struktur per stack → `structures/` · profil → `examples/` · status → `core_status.md`.
