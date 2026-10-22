package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Permission — Core 2: cek permission "<resource>,<action>" milik role aktif
// (core_principles.md §3 & core_role_permission_api.md §5).
//
// Deviasi sadar: learning-api-gateway memakai RBAC berbasis ROLE — standar core
// memakai PERMISSION per "<resource>.<action>" dengan default deny.
func Permission(resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: 1) role di Locals (dari JWTAuth) → bila SUPERADMIN → c.Next() (bypass).
		//       2) load permission role → cocokkan "<resource>.<action>" huruf kecil.
		//       3) tidak cocok → 403 {"error":"permission denied","code":"AUTHZ_PERMISSION_DENIED"}.
		_ = resource
		_ = action

		// Default deny (stub: selalu ditolak sampai implementasi diisi).
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "permission denied",
			"code":  "AUTHZ_PERMISSION_DENIED",
		})
	}
}
