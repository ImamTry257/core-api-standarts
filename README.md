# Core Architecture Knowledge Base

Knowledge base arsitektur inti ("5 core") untuk aplikasi API — dirancang **global (lintas project, lintas sesi, lintas AI agent)**: prinsip universal dipisah dari binding per aplikasi dan dari status yang sering berubah.

---

## 1) Mulai Dari Sini

```
README.md              ← kamu di sini: peta baca & aturan
   │
   ├─ core_principles.md     PRINSIP universal (stabil) — definisi 5 core, aturan baku
   │
   ├─ core_status.md         STATUS volatil — ✅/⬜, tanggal verifikasi, changelog
   │
   ├─ examples/              PROFILE per-project (contoh, bukan bagian core base)
   │      └─ core_index.md   Profile aplikasi StockFlow — salin → core_index_<project>.md (struktur di §4)
   │
   └─ Kontrak detail (per core, mengikuti kode aplikasi):
         core_auth_api.md             Core 1: Auth (9 endpoint, sesi, token, OTP)
         core_role_permission_api.md  Core 2: Role-Permission + step-up challenge
         core_crud_api.md             Core 3+4: Pola 6 operasi CRUD + pagination
         core_stepup_security_api.md  Core 5: Security policy, request action, role seal
```

**Urutan baca yang disarankan (untuk memahami keseluruhan):**
1. `core_principles.md` — konsep dua lapisan & prinsip kelima core.
2. `examples/core_index.md` — bagaimana prinsip itu terikat di aplikasi StockFlow.
3. Kontrak `core_*_api.md` — detail endpoint, request, response sukses & gagal.
4. `core_status.md` — apa yang sudah diverifikasi & riwayat keputusan.

**Untuk pekerjaan implementasi (tambah endpoint/modul):** cukup `core_principles.md` §4–§7 + kontrak core terkait.

---

## 2) Aturan Otoritas (saat dokumen saling bertentangan)

```
1. KODE SUMBER (repo aplikasi)      ← kebenaran tertinggi
2. core_*_api.md                    ← kontrak yang difokus dari kode
3. blueprint / reference (repo)     ← alur & kontrak turunan
4. core_status.md                   ← paling mudah usang
```

Jika kode berubah: **perbarui kontrak** (bukan mengabaikan kode), lalu catat di `core_status.md`.

---

## 3) Aturan Menjaga Agar Tetap Global & Tidak Usang

| Sifat file | Boleh berubah | Jangan dimasukkan |
|---|---|---|
| **Stabil** (`README.md`, `core_principles.md`) | Hanya saat arsitektur berubah — dengan baris changelog baru di `core_status.md` | Status pengerjaan, tanggal tugas, route spesifik |
| **Per-project** (`examples/core_index.md` / `examples/core_index_<project>.md`) | Saat modul/route project berubah | Prinsip umum, status lintas project |
| **Kontrak** (`core_*_api.md`) | Saat endpoint berubah (sinkron dengan kode) | Prinsip, changelog |
| **Volatil** (`core_status.md`) | Bebas — memang tempatnya | Detail kontrak (rujuk ke `core_*_api.md`) |

**Maintenance wajib:** ikuti checklist di `core_status.md` §4 (kode berubah → verifikasi ulang kontrak + perbarui tanggal).

---

## 4) Cara Pakai untuk Project Baru (lintas project)

1. Baca `core_principles.md` — prinsipnya berlaku umum, jangan ditulis ulang.
2. Salin `examples/core_index.md` menjadi **`examples/core_index_<nama-project>.md`** — isi §1 diagram lapisan (sesuai modul project baru), §2 pemetaan modul → core & doc, §3 endpoint core terkait (jika ada), §4 struktur knowledge.
3. Buat kontrak `core_<topik>_api.md` **per project** dengan memfokus dari kodenya sendiri (ikuti pola prinsip §2–§6; contoh terbaik ada di `core_auth_api.md` & `core_crud_api.md`).
4. Tambahkan baris project ke riwayat di `core_status.md`.
5. Jangan menimpa profile project yang sudah ada — setiap project satu file index.

**Standar implementasi yang wajib dipenuhi project baru** (ringkas dari prinsip):
- Sesi auth bertahap + endpoint terproteksi hanya `FULL_AUTH`.
- RBAC `<resource>,<action>` + step-up untuk aksi berisiko.
- Intent log **request action** + **security policy** per `(module, action)` sebelum aksi sensitif, lalu consume (`core_stepup_security_api.md`).
- 6 operasi CRUD (create 201, list/detail/update/delete 200) + lookup `/names`.
- Envelope `{contents,page,size,totalData}` + query `search,page,size`.
- Error: 400 input/`Invalid ID`/duplikat, 404 `<resource> tidak ditemukan`, 403 permission, 500 hanya gangguan server.

---

## 5) Ikhtisar 5 Core

| # | Core | Ruang lingkup | Kontrak (StockFlow) |
|---|---|---|---|
| 1 | Auth | Sesi `ANONYMOUS→PARTIAL→FULL_AUTH`, token, OTP, aktivasi, reset password | `core_auth_api.md` |
| 2 | Role-Permission | Permission `<resource>.<action>`, role, bypass superadmin, step-up challenge | `core_role_permission_api.md` |
| 3 | CRUD | 6 operasi standar, kode status & pesan baku, soft delete, validasi | `core_crud_api.md` |
| 4 | Pagination | Envelope `{contents,page,size,totalData}`, query standar (gabung core 3) | `core_crud_api.md` §4 |
| 5 | Step-up Security | Security policy `(module, action → risiko/validasi)`, request action (intent log + consume), role seal | `core_stepup_security_api.md` |

Status & tanggal verifikasi tiap core: `core_status.md` §1 — konsep & prinsipnya: `core_principles.md` (§2–§6).
