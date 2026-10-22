# Structure — Laravel

Peta struktur project Laravel yang menerapkan 5 core. **Sumber:** `fe-be-stockflow` (struktur nyata; `fe-be-flow` memakai pola sama). Kontraknya di 4 file `core_*_api.md`.

---

## 1) Tree

```
<app-laravel>/
├── routes/
│   ├── api.php                  ← seluruh route API + middleware (auth.jwt, permission:<r>,<a>, role.strict)
│   ├── web.php
│   └── console.php
├── app/
│   ├── Http/
│   │   ├── Controllers/
│   │   │   ├── Controller.php              ← base class (abstract, polos) — dipakai `extends`
│   │   │   └── Api/V1/
│   │   │       ├── Dashboard/   ← domain dashboard: AuthController, UserController, RolePermissionController,
│   │   │       │                   SecurityPolicyController, RequestActionController, CategoryController, dll.
│   │   │       └── Ops/         ← varian ops (OpsAuthController, dst.) — varian Core 1
│   │   └── Middleware/
│   │       ├── AuthJwtMiddleware.php        ← Core 1 (Bearer → FULL_AUTH)
│   │       ├── PermissionMiddleware.php     ← Core 2 (permission <resource>,<action>)
│   │       ├── RoleMiddleware.php           ← Core 2
│   │       ├── RoleStrictMiddleware.php     ← Core 2 (role persis, mis. SUPERADMIN)
│   │       └── OpsAuthJwtMiddleware.php     ← varian Core 1 (ops)
│   ├── Services/
│   │   ├── AuthService.php                  ← Core 1 (state sesi, OTP, activation, refresh)
│   │   └── RolePermissionSecurityService.php← Core 2+5 (challenge, payloadDigest)
│   ├── Support/
│   │   ├── ApiResponse.php                  ← Core 3/4 (3 bentuk response + envelope)
│   │   ├── AppError.php                     ← error bisnis ber-code
│   │   ├── JwtService.php                   ← Core 1 (terbit/parse token)
│   │   ├── OtpService.php                   ← Core 1 (generate OTP)
│   │   └── OpsJwtService.php                ← varian Core 1
│   ├── Models/                               ← User, Role, Permission, RolePermission, SecurityPolicy(+Target),
│   │                                           RequestAction, Category, Supplier, Product, dst.
│   ├── Jobs/                                 ← kirim email (activation/OTP/reset)
│   └── Providers/
├── database/migrations/
├── config/ · composer.json · artisan
```

---

## 2) Pemetaan Folder → Core

| Core | File di struktur ini | Rujukan kontrak |
|---|---|---|
| **1 Auth** | `routes/api.php` (prefix `/auth`), `Controllers/Api/V1/Dashboard/AuthController`, `Middleware/AuthJwtMiddleware`, `Services/AuthService`, `Support/{JwtService,OtpService}`, `Jobs` (email OTP/activation) | `core_auth_api.md` §1–§3 |
| **2 RBAC** | `Middleware/{Permission,Role,RoleStrict}Middleware`, `Controllers/.../RolePermissionController`, `RolePermissionSecurityController`, `Services/RolePermissionSecurityService`, `Models/{Role,Permission,RolePermission}` | `core_role_permission_api.md` §2–§5 |
| **3 CRUD** | `Controllers/.../<Domain>Controller` (6 operasi), `Support/ApiResponse` (pesan baku & kode status) | `core_crud_api.md` §1–§3, §6 |
| **4 Pagination** | `Support/ApiResponse::success` (envelope `{contents,page,size,totalData}`) + query `search,page,size` di tiap `index()` | `core_crud_api.md` §4 |
| **5 Step-up** | `Controllers/.../SecurityPolicyController`, `RequestActionController`, `RolePermissionSecurityController` (challenge), `AuthController::sealRole/verifyRoleSeal`, `Models/{SecurityPolicy(+Target),RequestAction,RolePermissionSecurityChallenge}` | `core_stepup_security_api.md` §1–§5 |

**Catatan:**
- Route `permission:<resource>,<action>` didefinisikan **per route** di `routes/api.php` (lihat pola `Route::middleware([...])->prefix(...)`).
- Pola response SEMUA lewat `Support/ApiResponse` — 3 bentuk (`success` / `error` / `middlewareError`) sesuai `core_auth_api.md` §1.1.
- Varian `Ops/` = Core 1 untuk role gudang (login NIP) — terpisah dari dashboard.

---

## 3) Stub

Ada di `stubs/` (penunjuk berkomentar, **bukan kode aktif**):

| Stub | Untuk core | Yang ditandai |
|---|---|---|
| `Controller.php` | — | base class `App\Http\Controllers\Controller` untuk `extends` (tanpa logika) |
| `AuthController.php` | 1 | 9 endpoint auth (login → OTP → activation → forgot/reset → refresh) dengan rujukan §2.1–§2.9 + pemetaan group route |
| `AuthJwtMiddleware.php` | 1 | alur Bearer → FULL_AUTH + bentuk error `{error,code}` |
| `PermissionMiddleware.php` | 2 | `<resource>,<action>`, SUPERADMIN bypass, default deny 403 |
| `ApiResponse.php` | 3/4 | 3 bentuk response + envelope pagination |
| `ExampleCrudController.php` | 3/4 | skeleton 6 operasi dengan pesan & kode status baku |
| `ExampleRequestActionController.php` | 5 | lifecycle policy → intent → consume |
