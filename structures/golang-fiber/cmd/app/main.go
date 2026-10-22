package main

import (
	"core-api-playbook-starter/internal/config"
	"core-api-playbook-starter/internal/handler"
	"core-api-playbook-starter/internal/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		// Bentuk error middleware: {"error", "code"} — core_auth_api.md §3
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error(), "code": "INTERNAL_ERROR"})
		},
	})

	app.Get("/health", handler.Health)

	// Route Core 1 (auth) + catatan route Core 2/3/4/5 — internal/router/router.go
	router.Register(app)

	if err := app.Listen(":" + cfg.Port); err != nil {
		panic(err)
	}
}
