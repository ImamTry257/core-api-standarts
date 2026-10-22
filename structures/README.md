# Structures — Peta Struktur Folder per Stack

Folder berisi **penerapan 5 core ke struktur project nyata per teknologi**: tree folder, pemetaan folder → core, dan stub penunjuk. Sifat file di sini **per-stack** (bukan universal) — prinsipnya tetap di `core_principles.md`, kontrak endpoint di `core_*_api.md`.

- **Bukan tutorial framework** — hanya menjawab: *letakkan auth/RBAC/CRUD/pagination/step-up di folder mana.*
- **Stub = penunjuk berkomentar, bukan kode aktif** — setiap stub mencantumkan rujukan `§` kontrak yang harus dipenuhi saat diimplementasi.

---

## 1) Daftar Stack

| Stack | Folder | Sumber belajar (diadopsi) | Kontrak rujukan |
|---|---|---|---|
| **Laravel** | `laravel/` | `fe-be-stockflow` (struktur nyata: routes, Controllers `Api/V1/{Dashboard,Ops}`, 5 middleware, Services, Support) | 4 file `core_*_api.md` |
| **Go (Fiber)** | `golang-fiber/` | `learning-api-gateway` (layout `cmd/internal`, Fiber v2, router/middleware/handler/usecase/dto) + `learning-auth-services` (lapisan auth repository→usecase→models + jwt utils) | 4 file `core_*_api.md` |

---

## 2) Peta Umum Folder → Core (berlaku semua stack)

| Core | Yang harus ada | Pola penempatan |
|---|---|---|
| 1 — Auth | middleware autentikasi, handler login/OTP/activation/refresh, service auth, util token/OTP | middleware + handler/usecase khusus auth |
| 2 — RBAC | middleware permission/role, handler role-permission, service challenge | middleware otorisasi + domain role |
| 3 — CRUD | controller/handler per domain, helper response standar | controller per resource |
| 4 — Pagination | helper envelope list + penerapan query `search,page,size` | di helper response + controller list |
| 5 — Step-up | handler security-policy, request-action, role-seal | domain keamanan (biasanya bersama RBAC) |

---

## 3) Cara Memakai Folder Ini

1. Pilih stack → baca `README.md` di folder stack-nya (tree + tabel pemetaan lengkap ke tiap file kontrak).
2. Untuk bikin project baru: ikuti tree-nya, lalu isi stub satu per satu sambil membuka kontrak `§` yang tertulis di komentar stub.
3. Stub **tidak** dijalankan apa adanya — ia penanda "isi di sini sesuai kontrak".

## 4) Cara Menambah Stack Baru

1. Pilih project contoh yang sudah jalan (sudah menerapkan ke-5 core sebagian).
2. Susun `README.md`: tree nyata + tabel pemetaan folder → core 1–5 (file konkret, bukan abstrak).
3. Buat `stubs/` minimal: middleware auth, middleware permission, helper response, 1 controller CRUD, 1 handler step-up.
4. Catat sumber belajar & deviasi (mis. bila envelope list di proyek asli tidak sesuai `core_crud_api.md` §4.1 — tulis apa adanya, jangan disamarkan).
5. Daftarkan stack ke tabel §1 di halaman ini.
