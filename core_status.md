# Core Status & Changelog

Status pengerjaan, tanggal verifikasi terhadap kode sumber, riwayat keputusan, dan aturan maintenance untuk knowledge base `core-architecture/`.

- **Sifat:** dokumen **volatil** — boleh sering berubah; prinsipnya ada di `core_principles.md`, profil aplikasi di `examples/core_index.md`.
- **Terakhir diperbarui:** 2026-10-21.

---

## 1) Status Core

| # | Core | File kontrak | Cakupan | Status | Terverifikasi ke kode |
|---|---|---|---|---|---|
| 1 | Auth | `core_auth_api.md` | 9 endpoint (`login`, `verify-code`, `resend-code`, activation×3, forgot/reset, `refresh-token`), state sesi `ANONYMOUS→PARTIAL→FULL_AUTH`, aturan password SHA-256, TTL token, error middleware | ✅ selesai | 2026-10-21 (`routes/api.php`, `AuthController`, `AuthService`) |
| 2 | Role-Permission | `core_role_permission_api.md` | 10 endpoint CRUD grouped `role+module`, 3 endpoint step-up challenge, rumus `payloadDigest`, guard `SUPERADMIN`, semantik middleware | ✅ selesai | 2026-10-21 (`RolePermissionController`, `RolePermissionSecurityService`) |
| 3 | CRUD | `core_crud_api.md` | Pola 6 operasi, kode status & pesan baku, soft delete, lookup `/names`, konvensi modul baru, tabel deviasi per domain | ✅ selesai | 2026-10-21 (`UserController` sebagai contoh utama + 5 domain pembanding) |
| 4 | Pagination | `core_crud_api.md` §4 | Envelope `{contents,page,size,totalData}`, query `search,page,size`, rumus offset/total | ✅ selesai (gabung ke core 3) | 2026-10-21 (16 controller memakai envelope sama) |
| 5 | Step-up Security & Request Action | `core_stepup_security_api.md` | Security Policy CRUD (5 endpoint), Request Action CRUD + consume (6 endpoint), Role Seal (2 endpoint), integrasi challenge Core 2 | ✅ selesai (disetujui jadi core ke-5) | 2026-10-21 (`SecurityPolicyController`, `RequestActionController`, `AuthController`, `AuthService`, migration `2026_05_17_*`) |

**Dokumen pendukung:**

| File | Status | Catatan |
|---|---|---|
| `README.md` | ✅ | Entry point & aturan otoritas |
| `core_principles.md` | ✅ | Prinsip universal lintas project |
| `examples/core_index.md` | ✅ | Profile aplikasi contoh (per-project, bukan core base) |
| `core_status.md` | ✅ | File ini |
| `structures/` | ✅ | Peta struktur per stack + stub — Laravel (6 stub PHP, `php -l` bersih); Go Fiber (skeleton `go build`/`go vet` hijau, `/health` jalan, 9 endpoint auth stub 501) |
| `flows/implementation-flow.md` | ✅ | Urutan implementasi 5 core dari nol (fase 0–4: fondasi → auth → CRUD+pagination → RBAC → step-up) |

---

## 2) Rincian Core 5 (mantan kandidat — disetujui & selesai)

| Bagian | Endpoint (prefix aplikasi) | Dokumentasi | Status |
|---|---|---|---|
| Step-up challenge role-permission | `/role-permissions/security/challenges*` | `core_role_permission_api.md` §4 | ✅ terdokumentasi |
| Request Action CRUD (intent log) | `/role-permissions/request-actions*` (+`/consume`) | **`core_stepup_security_api.md` §3** (doc FE lama, sudah diverifikasi — ada selisih pesan, lihat §7.2) | ✅ terdokumentasi |
| Security Policy CRUD | `/role-permissions/security-policies*` | **`core_stepup_security_api.md` §2** (doc repo: tidak ada — kontrak ekstrak murni dari kode) | ✅ terdokumentasi |
| Role seal (TTL ≤ 15 menit, HMAC binding) | `/auth/role/seal`, `/auth/role/verify` | **`core_stepup_security_api.md` §4** | ✅ terdokumentasi |

**Sisa yang diketahui (belum ditangani kode, tercatat di `core_stepup_security_api.md` §7.4):** status `REJECTED`/`EXPIRED`/`CANCELLED` tidak pernah diproduksi; tidak ada scheduler expiry; `consume` tidak mengecek `expires_at`.

---

## 3) Riwayat Keputusan

| Tanggal | Keputusan |
|---|---|
| 2026-10-21 | 4 core didefinisikan: auth, role-permission, CRUD, pagination |
| 2026-10-21 | Pagination digabung ke `core_crud_api.md` (bukan file terpisah) |
| 2026-10-21 | Contoh utama CRUD = domain **user** (hybrid), penyimpangan ditandai ❌; `/names` diilustrasi dengan supplier/category |
| 2026-10-21 | Kandidat core ke-5 (step-up security menyeluruh) dicatat, belum dibuat |
| 2026-10-21 | Knowledge di-refactor jadi **global-ready lintas project**: `README.md` + `core_principles.md` (stabil) / `core_index.md` (profile per project) / `core_status.md` (volatil) — konten lama dipindah, tidak dihapus |
| 2026-10-21 | **Core ke-5 disetujui & didokumentasikan** (`core_stepup_security_api.md`): menggabungkan security-policy CRUD, request-action CRUD + consume, dan role seal; challenge tetap di `core_role_permission_api.md` §4 (disambung, tidak diduplikasi) |
| 2026-10-21 | `core_index.md` **dipindah ke `examples/`** — root kini khusus core base universal (README, principles, 4 kontrak) + status; profile per-project tinggal di `examples/` |
| 2026-10-21 | **6 file (principles, status, 4 kontrak) dibersihkan dari identitas aplikasi** — 20 wording "stockflow" (judul, intro, email contoh, catatan ops) + 11 nama dokumen/repo (`fe-be-flow`, `flow_*.md`, `docs/api/*`) diganti generik; `examples/core_index.md` & README tetap menyebut aplikasi contoh |
| 2026-10-21 | Folder **`structures/`** dibuat: peta struktur + stub per stack (Laravel & Go Fiber); skeleton Go terverifikasi `go build`/`go vet` + `/health` berjalan; stub auth dilengkapi 9 endpoint — termasuk koreksi path `/refresh-token` (§2.9) dan rute resend-activation di group users (§2.6, bukan `/auth/activation/resend`) |
| 2026-10-21 | **`flows/implementation-flow.md`** dibuat: panduan urutan bootstrap 5 core (fase 0–4) — melengkapi `core_principles.md` §7 yang hanya menuntun menambah modul di sistem jalan; keputusan nama repo/entry → **Core API Playbook** (`core-api-playbook`) |

---

## 4) Checklist Maintenance

| Peristiwa | Tindakan wajib |
|---|---|
| Kode endpoint auth (login/OTP/activation/reset) berubah | Perbarui `core_auth_api.md` + tanggal verifikasi baris core 1 di atas |
| Kode role-permission / challenge / middleware berubah | Perbarui `core_role_permission_api.md` + baris core 2 |
| Pola CRUD/pesan/pagination berubah di ≥1 domain | Perbarui `core_crud_api.md` + baris core 3 & 4; evaluasi tabel deviasi |
| Kode security-policy / request-action / role seal berubah | Perbarui `core_stepup_security_api.md` + tanggal verifikasi baris core 5 di atas |
| Route/modul bisnis baru/tambah di aplikasi | Tambah baris di `examples/core_index.md` §2 |
| Prinsip berubah (jarang) | Edit `core_principles.md` + baris baru di Riwayat Keputusan |
| Project baru memakai knowledge ini | Buat `examples/core_index_<project>.md` baru — jangan menimpa profile yang ada |
| Struktur folder stack berubah (mis. tambah middleware/layer) | Perbarui `structures/<stack>/README.md` + stub terkait |
| Urutan/prasyarat fase implementasi berubah | Perbarui `flows/implementation-flow.md` |

**Otoritas saat konflik:** kode sumber > `core_*_api.md` > blueprint/reference > `core_status.md` (lihat `README.md`).
