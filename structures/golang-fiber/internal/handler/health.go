package handler

import (
	"core-api-playbook-starter/internal/dto"

	"github.com/gofiber/fiber/v2"
)

// Health — endpoint contoh yang BERJALAN penuh: memperlihatkan bentuk envelope sukses
// {status:true, message, data} (core_auth_api.md §1.1).
func Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(dto.Success("OK", fiber.Map{
		"service": "core-api-playbook-starter",
	}))
}
