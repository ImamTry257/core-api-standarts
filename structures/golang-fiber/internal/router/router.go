package router

import (
	"core-api-playbook-starter/internal/handler"
	"core-api-playbook-starter/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Register — pola route group mengikuti learning-api-gateway/internal/router/auth_router.go:
// prefix /api/v1/dashboard, group auth terbuka, group protected dibungkus JWT.
func Register(app *fiber.App) {
	// ===== Core 1: AUTH ===== (core_auth_api.md §2) — semua rute terdaftar,
	// handler-nya masih stub 501 sampai diisi.
	auth := app.Group("/api/v1/dashboard/auth")
	h := handler.NewAuthHandler()

	// Terbuka (belum butuh token):
	auth.Post("/login", h.Login)                         // §2.1
	auth.Post("/verify-code", h.VerifyOtp)               // §2.2
	auth.Post("/resend-code", h.ResendOtp)               // §2.3
	auth.Post("/activation/check", h.CheckActivation)    // §2.4
	auth.Post("/activation/confirm", h.ConfirmActivation) // §2.5
	auth.Post("/forgot/password", h.ForgotPassword)      // §2.7
	auth.Post("/reset/password", h.ResetPassword)        // §2.8

	// Hanya sesi FULL_AUTH — middleware dipasang PER ROUTE (bukan group.Use):
	// di Fiber v2, middleware grup terdaftar setelah rute login sehingga melekat ke
	// seluruh prefix (path tak terdaftar ikut 401, padahal kontrak §3 → 404).
	auth.Post("/refresh-token", middleware.JWTAuth(true), h.Refresh) // §2.9 — path persis kontrak (bukan /refresh)

	// §2.6 — endpoint RIIL ada di group users (bukan /auth/activation/resend —
	// path blueprint itu tidak ada di backend):
	app.Post("/api/v1/dashboard/users/:id/resend-activation",
		middleware.JWTAuth(true),
		middleware.Permission("users", "resend_activation"),
		h.ResendActivation)

	// ===== Core 2: RBAC ===== (core_role_permission_api.md §5)
	// Contoh: app.Get("/api/v1/dashboard/users",
	//     middleware.JWTAuth(true), middleware.Permission("users", "read"), handler.ListUsers)
	// Route Role-Permission & challenge memakai group terproteksi role.strict:SUPERADMIN → core_role_permission_api.md §2–§4.

	// ===== Core 3+4: CRUD + PAGINATION ===== (core_crud_api.md §1.2 & §4)
	// Contoh per domain:
	//   app.Post("/api/v1/dashboard/example",     JWTAuth, Permission("example","create"), handler.Store)   → 201
	//   app.Get("/api/v1/dashboard/example",      JWTAuth, Permission("example","read"),   handler.Index)    → envelope {contents,page,size,totalData}
	//   app.Get("/api/v1/dashboard/example/{id}", JWTAuth, Permission("example","read"),   handler.Show)     → 200
	//   app.Put("/api/v1/dashboard/example/{id}", JWTAuth, Permission("example","update"), handler.Update)   → 200
	//   app.Delete("/api/v1/dashboard/example/{id}", JWTAuth, Permission("example","delete"), handler.Destroy)→ 200 data:null
	//   app.Get("/api/v1/dashboard/example/names", JWTAuth, Permission("example","read"), handler.Names)     → lookup {id,name} terpaginasi

	// ===== Core 5: STEP-UP SECURITY ===== (core_stepup_security_api.md)
	// Route: /role-permissions/security-policies, /role-permissions/request-actions (+/consume),
	//        /auth/role/seal, /auth/role/verify — proteksi JWTAuth + permission/role.strict
	//        (lihat §1.1 kontrak core 5 untuk proteksi masing-masing).
}
