package router

import (
	"go-generate-sk/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/generate-sk", handler.GenerateSK)
}
