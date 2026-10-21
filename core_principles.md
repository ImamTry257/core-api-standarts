# Core Principles — Kontrak Inti Lintas Project

Prinsip universal arsitektur "5 core" untuk aplikasi berbasis API. Dokumen ini **stabil dan tidak terikat satu project** — binding per aplikasi (route, middleware, teknologi) ada di `examples/core_index.md` (profile per project), detail kontrak ada di `core_*_api.md`.

- **Sifat:** prinsip / asas — hanya berubah bila arsitektur berubah, bukan karena kode berubah.
- **Status & changelog:** lihat `core_status.md`.
- **Entry point & aturan otoritas:** lihat `README.md`.
- **Contoh implementasi rujukan:** `examples/core_index.md`.

---

## 1) Konsep Dasar: Dua Lapisan

```
MODUL BISNIS   = domain & data sendiri (punya entitas, aturan bisnis, alur approval)
LAPISAN CORE   = kontrak lintas modul, TIDAK punya data domain sendiri
```

**Syarat sebuah hal disebut "core":**
1. Dipakai oleh **≥2 modul bisnis** (idealnya semua).
2. Tidak punya data domain sendiri (bukan entitas bisnis).
3. Kontraknya **wajib konsisten** — penyimpangan harus terdokumentasi sebagai deviasi.

**Keputusan arsitektur yang sudah diambil (sudah disepakati, jangan dibuka lagi kecuali ada alasan kuat):**

| # | Core | Ruang lingkup |
|---|---|---|
| 1 | **Auth** | Sesi, token, aktivasi akun, OTP, reset password |
| 2 | **Role-Permission (RBAC)** | Otorisasi endpoint per role/permission + step-up untuk aksi berisiko |
| 3 | **CRUD** | Pola operasi standar per entitas + pesan & kode status baku |
| 4 | **Pagination** | Envelope daftar + query standar (digabung ke core CRUD — keputusan 2026-10-21) |
| 5 | **Step-up Security & Request Action** | Intent log aksi sensitif, security policy per (module, action), role seal (disetujui sebagai core ke-5 — 2026-10-21) |

**Aturan konsumsi:** setiap modul bisnis baru WAJIB memenuhi kelima core.

---

## 2) Core 1 — Auth (Prinsip)

1. **Sesi bertahap** dengan level eksplisit: `ANONYMOUS → PARTIAL_AUTH (kredensial valid, OTP terkirim) → FULL_AUTH (OTP valid)`.
2. **Endpoint terproteksi hanya menerima `FULL_AUTH`** — token Partial dipakai seminimal mungkin (hanya untuk verifikasi OTP).
3. **Kode sekali pakai ber-expiry**: activation code, OTP, reset code — semua punya TTL jelas dan hangus setelah dipakai.
4. **Pesan kegagalan kredensial generik** (satu pesan untuk "user tidak ada" dan "password salah") agar tidak bocor informasi akun.
5. **Password disimpan ter-hash**; FE mengirim hash (SHA-256), backend menyimpan bcrypt.
6. **Refresh token dirotasi** setiap dipakai; akses token pendek, refresh panjang.
7. **Dua bentuk error**: error bisnis `{status,message,data}` vs error middleware `{error,code}` — FE wajib menangani keduanya.
8. **Audit trail** untuk event auth (login, verify OTP, reset).

> Detail kontrak: `core_auth_api.md`.

---

## 3) Core 2 — Role & Permission / RBAC (Prinsip)

1. **Izinkan per `<resource>,<action>`** — permission didefinisikan per kombinasi resource+action (mis. `<domain>.read`), role adalah kumpulan permission.
2. **Role super (setara SUPERADMIN) selalu bypass** pemeriksaan permission biasa, tetapi **dilindungi dari dimodifikasi** (guard: permission-nya sendiri tidak boleh diubah lewat API normal).
3. **Default tolak** — tanpa permission yang cocok → `403`.
4. **Step-up verification untuk aksi berisiko** (mutasi hak akses):
   - Bukti identitas ulang (password) → challenge → OTP → `verifiedToken` **ber-TTL pendek & satu pakai**.
   - Token diikat ke **payload digest** (SHA-256 dari payload request) — mengganti payload setelah verifikasi = gagal (`mismatch`).
   - Context header wajib agar challenge tidak bisa dipakai untuk konteks lain.
5. **Klasifikasi error otorisasi**: `401` (belum login/verifikasi/identitas salah) vs `403` (sudah login tapi tidak berhak).
6. **Role list & permission metadata** tersedia sebagai endpoint baca untuk dropdown/konfigurasi FE.

> Detail kontrak: `core_role_permission_api.md`.

---

## 4) Core 3 — CRUD (Prinsip)

**Pola 6 operasi standar untuk setiap entitas:**

| # | Operasi | Metode | Kode status | Bentuk sukses |
|---|---|---|---|---|
| 1 | Create | POST | **201** | entity |
| 2 | List | GET | 200 | envelope paginasi |
| 3 | Detail | GET `/{id}` | 200 | entity |
| 4 | Update | PUT `/{id}` | 200 | entity |
| 5 | Delete | DELETE `/{id}` | 200 | `null` |
| 6 | Lookup | GET `/names` | 200 | `{id,name}` terpaginasi |

**Aturan baku:**
1. `message` mengikuti pola `<Resource> created successfully` / `Successfully retrieved all <resources>` / `<Resource> updated successfully` / `<Resource> deleted successfully` / `Successfully retrieved all <resources> name`.
2. Kode status konsisten: **400** = input salah (`Invalid ID`, `<field> is required`, data duplikat), **404** = `<resource> tidak ditemukan`, **403** = permission, **500** = hanya gangguan server sungguhan. **Duplikasi data bukan 500.**
3. Pesan error konsisten satu bahasa (Bahasa Indonesia pada implementasi rujukan).
4. **Soft delete** (hapus logis) untuk data master — data terhapus tidak muncul di list.
5. Validasi field wajib & relasi (FK valid) **sebelum** menyimpan.
6. Response entity: `id` + field camelCase + `createdAt`/`updatedAt`.
7. Operasi ekstra di luar 6 pola (approve/reject/deactivate dsb.) diperbolehkan **asal** memakai format response & aturan error yang sama.

> Detail kontrak + tabel deviasi: `core_crud_api.md`.

---

## 5) Core 4 — Pagination (Prinsip)

1. **Satu envelope** untuk semua daftar: `{contents, page, size, totalData}`.
2. **Query standar**: `search` (LIKE), `page` (default 1), `size` (default 10) + filter khusus domain.
3. `totalData` = jumlah baris hasil filter (bukan jumlah halaman); backend **tidak** menghitung `totalPages` — FE menghitung `ceil(totalData/size)`.
4. Endpoint lookup untuk dropdown **juga terpaginasi** (payload ringan `{id,name}`) — bukan dump seluruh data.
5. Server **wajib** mengurutkan secara deterministik (standar ideal: `created_at DESC`) — ketiadaan `ORDER BY` adalah deviasi, bukan norma.
6. Validasi nilai query (negatif/batas atas `size`) adalah tanggung jawab server; ketiadaan batas `size` maksimum dicatat sebagai risiko.

> Detail kontrak: `core_crud_api.md` §4.

---

## 6) Core 5 — Step-up Security & Request Action (Prinsip)

1. **Intent sebelum aksi:** aksi sensitif didahului **request action** — log intent berisi payload + `payload_hash` (SHA-256) sebelum eksekusi nyata; aksi diakhiri dengan **consume** (status `EXECUTED`).
2. **Risiko dikonfigurasi, bukan hardcoded:** setiap `(module, action)` punya **security policy target** dengan `riskLevel` (`LOW..CRITICAL`) dan `validationType` (`CONFIRM|NOTE|PASSWORD|OTP`) — inilah yang menentukan tingkat verifikasi.
3. **Lifecycle status eksplisit & terminal:** `PENDING`/`CHALLENGE_REQUIRED` (awal) → `CHALLENGE_VERIFIED` (setelah OTP Core 2) → `EXECUTED` (terminal untuk edit: tolak dengan 409).
4. **Intent ber-TTL** — `expires_at` jelas pada pembuatan; status kedaluwarsa harus ditangani eksplisit (bila tidak ada job pembersih, catat sebagai deviasi).
5. **Role seal stateless:** bukti role aktif yang ter-encrypt, ber-TTL ≤ 15 menit, **terikat ke access token** (HMAC binding) — refresh token membatalkannya; selalu dibawa bersama Bearer yang sama, bukan sebagai token pengganti.
6. **Step-up challenge untuk mutasi berisiko** tetap milik Core 2 (`core_role_permission_api.md` §4); Core 5 adalah wadah intent & aturan risikonya.
7. **Konfigurasi keamanan ≠ data master** — boleh menyimpang dari pola core 3 (tanpa `/names`, hard delete) **asal terdokumentasi** sebagai deviasi.

> Detail kontrak: `core_stepup_security_api.md`.

---

## 7) Konvensi Penamaan & Langkah Menambah Modul

**Penamaan file knowledge:**

| Pola | Fungsi | Sifat |
|---|---|---|
| `README.md` | Entry point & aturan otoritas | global, stabil |
| `core_principles.md` | Prinsip universal (file ini) | global, stabil |
| `examples/core_index_<project>.md` | Profile pemetaan project → core (di folder `examples/`) | per-project |
| `core_status.md` | Status, tanggal verifikasi, changelog | volatil |
| `core_<topik>_api.md` | Kontrak detail satu core | per-project, mengikuti kode |

**Langkah menambah modul bisnis baru:**
1. Baca `README.md` → `core_principles.md` (file ini) → profile project → kontrak core terkait.
2. Daftarkan endpoint dengan proteksi auth + permission `<resource>,<action>` (jangan membuat route terbuka tanpa proteksi — kecuali memang publik by design dan didokumentasikan).
3. Implementasikan 6 operasi CRUD + envelope pagination sesuai prinsip §4–§5.
4. Validasi field wajib/relasi/duplikat dengan kode 400; 404 konsisten; soft delete untuk data master.
5. Tulis dokumen alur + kontrak detail di repo aplikasi, lalu tambahkan baris ke profile (`examples/core_index_<project>.md`).
6. Catat perubahan & tanggal verifikasi di `core_status.md`.

**Aturan maintenance:**
- **Otoritas: kode sumber > kontrak `core_*_api.md` > blueprint/reference > status.** Bila berbeda, kode yang menang dan kontrak wajib diperbarui.
- Saat kode salah satu core berubah → perbarui `core_*_api.md` terkait + baris verifikasi di `core_status.md`.
- Saat prinsip berubah → baris baru di changelog `core_status.md` (jangan edit diam-diam).
