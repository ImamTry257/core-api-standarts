# Core Index — Profile Aplikasi StockFlow

Profile pemetaan **aplikasi StockFlow (dashboard `fe-be-flow`)** ke lapisan core. File ini **per-project** dan tinggal di `examples/` — prinsip universal ada di `core_principles.md`, status & changelog di `core_status.md`, entry point di `README.md`.

> **Untuk project lain:** jangan menimpa file ini — buat `examples/core_index_<nama-project>.md` baru (salin struktur §1–§4 di bawah).

- **Profile:** StockFlow Dashboard — base path `/api/v1/dashboard` (varian ops: `/api/v1/ops`)
- **Repo sumber:** `fe-be-flow/` (route di `routes/api.php`, kontrak di `app/Http/Controllers/Api/V1/*`)
- **Prinsip yang dipakai:** `core_principles.md` (5 core: auth, role-permission, CRUD, pagination, step-up security)
- **Status core:** `core_status.md`

---

## 1) Diagram Lapisan (StockFlow)

```
┌─────────────────────────── MODUL BISNIS (domain & data sendiri) ───────────────────────────┐
│ master data        │ inventory          │ POS penjualan │ reporting │ transfer │ owner     │
│ category, branch,  │ receiving, opname, │ pembayaran    │ modul G   │ antar    │ modul F   │
│ product, supplier, │ product request    │ QRIS/transfer │           │ cabang   │           │
│ user, role         │                    │               │           │          │           │
├────────────────────┼────────────────────┼───────────────┼───────────┼──────────┼───────────┤
│ request-action & security policy — intent log & aturan risiko (Core 5, lihat core_status.md §2) │
└────────────────────────────────┬──────────────────────────────────────────────────────────┘
                                 │  semua modul memakai kontrak yang sama
┌────────────────────────────────▼──────────────── LAPISAN CORE ─────────────────────────────┐
│ 1. Auth            sesi/token/OTP/aktivasi    → core_auth_api.md            ✅             │
│ 2. Role-Permission RBAC + step-up challenge   → core_role_permission_api.md ✅             │
│ 3. CRUD            pola 6 operasi per entitas → core_crud_api.md            ✅             │
│ 4. Pagination      envelope list + query      → core_crud_api.md §4 (gabung)✅             │
│ 5. Step-up Security intent log + policy + seal → core_stepup_security_api.md✅             │
│                                                                                           │
│ Peta ini (profile StockFlow)         → examples/core_index.md (file ini)                   │
│ Prinsip lintas project               → core_principles.md                                │
│ Status & changelog                   → core_status.md                                    │
└───────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2) Pemetaan Modul Bisnis → Core & Dokumen

| Modul bisnis | Route prefix | Core 1 Auth | Core 2 RBAC | Core 3 CRUD | Core 4 Pagination | Core 5 Step-up | Blueprint (alur) | Reference (kontrak detail) |
|---|---|---|---|---|---|---|---|---|
| **Master data** (category, branch, product, supplier, user, role) | `/categories`, `/branches`, `/products`, `/suppliers`, `/users`, `/roles` | ✅ Bearer | ✅ `permission:<domain>,*` (kecuali category terbuka ❌) | ✅ contoh utama `/users` | ✅ `search,page,size` + `/names` | — | `flow_master_data_dashboard.md` | `api_master_data`, `api_category`, `api_branch`, `api_product`, `api_supplier`, `api_user` `.md` |
| **Inventory** (receiving, opname, product request) | `/inventory/*`, `/stock-opnames`, `/product-requests` | ✅ Bearer | ✅ `permission:inventory.receiving,*`, `stock.opname.*` | ⚠️ pola approve/reject (bukan murni CRUD) | ✅ list paginated | — | `flow_barang_masuk_maker_gudang.md` | `api_barang_masuk`, `api_permintaan_produk_maker_checker_admin_gudang`, `api_auth_inventory` (bagian inventory) `.md` |
| **POS Penjualan** (pembayaran QRIS/transfer) | POS routes | ✅ Bearer | ✅ role-based | ⚠️ transaksi (non-CRUD) | ✅ list transaksi | — | `flow_pos_penjualan_payment_qris_transfer.md` | — |
| **Reporting** | reporting routes | ✅ Bearer | ✅ `permission:*` | ⚠️ read-only | ✅ paginated | — | `flow_modul_g_reporting.md` | — |
| **Transfer antar cabang** | transfer routes | ✅ Bearer | ✅ maker-checker | ⚠️ approval flow | ✅ list | — | `transfer_barang_antar_cabang.md` | — |
| **Owner dashboard (modul F)** | owner routes | ✅ Bearer | ✅ OWNER/SUPERADMIN | ⚠️ read + approve | ✅ | — | `flow_modul_f_dashboard_owner.md` | — |
| **Request-action & security policy** | `/role-permissions/request-actions`, `/security-policies`, `/auth/role/seal` | ✅ Bearer | ✅ `role.strict:SUPERADMIN` | ✅ CRUD + consume (varian) | ✅ envelope | ✅ kontrak lengkap | — | **`core_stepup_security_api.md`** (doc FE lama `request-action-crud.md`, diverifikasi) |

**Lokasi dokumen sumber:** `fe-be-flow/docs/api/blueprint/*` (alur) dan `fe-be-flow/docs/api/reference/*` (kontrak per domain).

**Istilah kunci lintas dokumen (binding StockFlow):**
- `Bearer` = token hasil `verify-code` (FULL_AUTH) — `core_auth_api.md`.
- `permission:<resource>,<action>` = middleware Laravel yang mencocokkan permission `<resource>.<action>` — `core_role_permission_api.md` §5.
- Envelope list & pesan `<X> created successfully` dsb. — `core_crud_api.md`.
- Step-up token untuk mutasi role-permission — `core_role_permission_api.md` §4; intent log & aturan risiko (request-action, security policy, role seal) — `core_stepup_security_api.md`.

**Varian di luar kelima core (dicatat, bukan bagian core):**
- **Ops auth** (`/v1/ops/*`, login via NIP, middleware `auth.ops`) — varian core 1 untuk role gudang (`MAKER_GUDANG`, `CHECKER_GUDANG`), terpisah dari dashboard.
- **Smoke test resume** — `api_smoke_test_resume_2026-04-29.md` (hasil uji end-to-end 2026-04-29).

---

## 3) Endpoint Core 5 (StockFlow)

Pemetaan endpoint (prinsip di `core_principles.md` §6, status & verifikasi di `core_status.md` §2):

| Bagian | Endpoint | Kontrak |
|---|---|---|
| Security Policy CRUD | `/role-permissions/security-policies*` (5 endpoint) | `core_stepup_security_api.md` §2 |
| Request Action CRUD + consume | `/role-permissions/request-actions*` (6 endpoint) | `core_stepup_security_api.md` §3 |
| Role seal | `POST /auth/role/seal`, `POST /auth/role/verify` | `core_stepup_security_api.md` §4 |
| Step-up challenge role-permission | `POST /role-permissions/security/challenges`, `/verify`, `/resend` | `core_role_permission_api.md` §4 (Core 2, disambung di core 5 §5) |

---

## 4) Struktur Knowledge di Folder Ini

Lihat `README.md` untuk peta baca & aturan maintenance. Ringkas:

| File | Sifat |
|---|---|
| `README.md` | Entry point global — mulai dari sini |
| `core_principles.md` | Prinsip universal, lintas project, stabil |
| `examples/core_index.md` | **Profile aplikasi StockFlow (file ini)** — per-project |
| `core_status.md` | Status ✅/⬜, tanggal verifikasi, changelog — volatil |
| `core_auth_api.md` | Kontrak detail core 1 (StockFlow) |
| `core_role_permission_api.md` | Kontrak detail core 2 (StockFlow) |
| `core_crud_api.md` | Kontrak detail core 3+4 (StockFlow) |
| `core_stepup_security_api.md` | Kontrak detail core 5 (StockFlow) |
