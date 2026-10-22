# Structure — Go (Fiber)

Peta struktur project Go dengan **Fiber** yang menerapkan 5 core.

**Sumber yang diadopsi:**
- **Layout & pola Fiber** → `learning-api-gateway` (Fiber **v2**): `cmd/app` + `internal/{router,handler,middleware,usecase,dto}`, route group `/api/v1/dashboard/auth`, middleware `JWTAuth(true)` / `RBAC(...)`.
- **Lapisan auth** → `learning-auth-services`: `repository` (interface + impl, gorm) → `usecase` (interface + impl) → `models/{entity,dto}` + `utils/jwt` (`SetJWT`/`VerifyJWT`).

**Catatan jujur (deviasi sengaja):**
1. **Fiber v2** mengikuti rujukan (`learning-api-gateway`); project `upload-service` memakai v3 — bila upgrade, sesuaikan import `github.com/gofiber/fiber/v3`.
2. Envelope list gateway (`{page,totalData,contents}`) **diadaptasi** ke standar core (`{contents,page,size,totalData}` — `core_crud_api.md` §4.1). **Skeleton mengikuti standar, bukan salin gateway.**
3. Gateway membawa `grpcclient` + gRPC auth — skeleton di sini **self-contained** (tanpa gRPC) supaya bisa jalan mandiri; pola layer tetap.
4. Middleware RBAC gateway berbasis **role** → skeleton memakai **permission** `<resource>,<action>` (standar `core_role_permission_api.md` §5), dengan `AUTHZ_PERMISSION_DENIED`.

---

## 1) Tree (skeleton di folder ini)

```
golang-fiber/
├── go.mod                        module core-api-playbook-starter · fiber/v2
├── cmd/app/main.go               bootstrap fiber, error handler {error,code}, daftar route, :PORT
└── internal/
    ├── config/config.go          PORT dari env (default 8080)
    ├── dto/response.go           APIResponse{status,message,data} + APIListResponse{contents,page,size,totalData}
    ├── router/router.go          AuthRoutes (group /api/v1/dashboard/auth) + catatan route Core 2/3/4/5
    ├── handler/
    │   ├── health.go             GET /health → envelope sukses (contoh berjalan)
    │   └── auth.go               AuthHandler stub — 9 endpoint 501, core_auth_api.md §2.1–§2.9
    ├── middleware/
    │   ├── jwt.go                JWTAuth: Bearer → error AUTH_* {error,code} — core_auth_api.md §3
    │   └── rbac.go               Permission("<r>","<a>") default deny AUTHZ_* — core_role_permission_api.md §5
    └── usecase/auth_usecase.go   interface + impl stub 9 metode (pola learning-auth-services)
```

**Jalankan:** `go run ./cmd/app` lalu `GET http://localhost:8080/health`.

---

## 2) Pemetaan Folder → Core

| Core | File di struktur ini | Rujukan kontrak |
|---|---|---|
| **1 Auth** | `middleware/jwt.go`, `handler/auth.go`, `usecase/auth_usecase.go`, `router/router.go` (group `/auth`) — util token mengikuti pola `SetJWT/VerifyJWT` (`learning-auth-services/pkg/utils/jwt.go`) | `core_auth_api.md` §1–§3 |
| **2 RBAC** | `middleware/rbac.go` (`Permission`) — route mutasi wajib diapit `JWTAuth(true)` + `Permission(...)` | `core_role_permission_api.md` §5 |
| **3 CRUD** | `handler/*` per domain + `usecase/*` (interface+impl) — pola 6 operasi | `core_crud_api.md` §1–§3 |
| **4 Pagination** | `dto/response.go` (`APIListResponse`) + parsing `search,page,size` di handler list | `core_crud_api.md` §4 |
| **5 Step-up** | handler/`usecase` security-policy, request-action, role-seal (belum ada stub — rancang mengikuti `core_stepup_security_api.md` §1.4) | `core_stepup_security_api.md` |

---

## 3) Konvensi Pindah dari Laravel ke Fiber

| Laravel | Go (Fiber) |
|---|---|
| `routes/api.php` | `internal/router/*.go` (group + middleware per route) |
| `app/Http/Controllers/...` | `internal/handler/*.go` (1 handler struct per domain) |
| `app/Services/...` | `internal/usecase/*_usecase.go` (interface + impl) |
| `app/Http/Middleware/...` | `internal/middleware/*.go` (`fiber.Handler`) |
| `app/Support/ApiResponse` | `internal/dto/response.go` |
| `app/Models` (Eloquent) | model entity + repository gorm (`pkg/auth/repository` pola learning-auth-services) |

---

## 4) Stub

File skeleton di folder ini **setengah jalan** (bisa `go build` tapi endpoint inti masih `TODO`):
- `GET /health` → **berjalan penuh** (contoh envelope).
- **8 rute group `/auth`** (§2.1–2.5, §2.7–2.9) + `POST /users/:id/resend-activation` (§2.6, +`JWTAuth`+`Permission`) → semua handler stub **501** dengan komentar rujukan kontrak.
- Path refresh persis kontrak: **`/refresh-token`** (§2.9 — bukan `/refresh`); §2.6 sengaja di luar group `/auth` (lihat kontrak: path blueprint `auth/activation/resend` tidak ada di backend).
- `JWTAuth` → sudah menangani bentuk error `AUTH_*`; verifikasi token masih placeholder (**wajib diisi sebelum produksi**).
- `Permission` → default deny 403 `AUTHZ_PERMISSION_DENIED` (siap dipakai di route mutasi).
