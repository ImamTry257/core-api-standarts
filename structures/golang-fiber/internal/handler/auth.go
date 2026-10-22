package handler

import (
	"core-api-playbook-starter/internal/dto"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler — Core 1 Auth (STUB penunjuk, bukan implementasi).
// Sembilan endpoint wajib mengikuti core_auth_api.md §2.1–§2.9 — request, pesan,
// dan kode status persis dari kontrak; stub ini hanya menandai "isi di sini".
type AuthHandler struct {
	// authUC usecase.AuthUseCase — aktifkan saat usecase diisi (pola learning-auth-services).
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func notImplemented(section string) dto.APIResponse {
	return dto.Fail("belum diimplementasi — lihat core_auth_api.md " + section)
}

// Login — core_auth_api.md §2.1.
// Request: {email, password (sha256 hex 64 char)}; gagal kredensial → pesan generik
// "email or password is wrong" (jangan bocorkan akun ada/tidak). Sukses → sesi PARTIAL_AUTH (§1.5).
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(notImplemented("§2.1"))
}

// VerifyOtp — core_auth_api.md §2.2: kode OTP → sesi FULL_AUTH (§1.5) + token final.
func (h *AuthHandler) VerifyOtp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(notImplemented("§2.2"))
}

// ResendOtp — core_auth_api.md §2.3: kirim ulang OTP login (ikut batas di kontrak).
func (h *AuthHandler) ResendOtp(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(notImplemented("§2.3"))
}

// CheckActivation — core_auth_api.md §2.4: cek status aktivasi akun.
func (h *AuthHandler) CheckActivation(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(notImplemented("§2.4"))
}

// ConfirmActivation — core_auth_api.md §2.5: konfirmasi kode aktivasi.
func (h *AuthHandler) ConfirmActivation(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(notImplemented("§2.5"))
}

// ResendActivation — core_auth_api.md §2.6: endpoint RIIL-nya di group users:
// POST /api/v1/dashboard/users/{id}/resend-activation (FULL_AUTH + permission
// users,resend_activation) — BUKAN /auth/activation/resend (tidak ada di backend).
// Didaftarkan di router.go bersama rute users, bukan di group /auth.
func (h *AuthHandler) ResendActivation(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(notImplemented("§2.6"))
}

// ForgotPassword — core_auth_api.md §2.7: minta link/kode reset via email.
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(notImplemented("§2.7"))
}

// ResetPassword — core_auth_api.md §2.8: tukar token reset dengan password baru.
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(notImplemented("§2.8"))
}

// Refresh — core_auth_api.md §2.9: POST /auth/refresh-token — refresh token
// sekali pakai (dirotasi) → access token baru; hanya sesi FULL_AUTH (§1.5).
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(notImplemented("§2.9"))
}
