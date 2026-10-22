package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTAuth — Core 1: proteksi sesi FULL_AUTH (core_auth_api.md §1.3 & §3).
// Error middleware memakai bentuk {"error", "code"} — BERBEDA dari error bisnis {status,message,data}.
//
// STUB: header sudah diperiksa, VERIFIKASI TOKEN MASIH PLACEHOLDER —
// wajib diisi (parse + cek exp + isi Locals) sebelum dipakai produksi.
func JWTAuth(required bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
				"code":  "AUTH_MISSING_HEADER",
			})
		}
		if !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid format header",
				"code":  "AUTH_INVALID_HEADER_FORMAT",
			})
		}

		token := strings.TrimPrefix(header, "Bearer ")

		// TODO: verifikasi signature & expiry (pola SetJWT/VerifyJWT —
		// learning-auth-services/pkg/utils/jwt.go). Bila gagal → 401
		// {"error": "invalid or expired token", "code": "AUTH_INVALID_OR_EXPIRED_TOKEN"}.
		// Setelah lolos: c.Locals("user_id", uid); c.Locals("role", role).
		_ = token
		if required {
			// placeholder: lolos tanpa verifikasi — GANTI dengan cek token asli.
		}

		return c.Next()
	}
}
