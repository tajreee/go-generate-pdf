package main

import (
	"go-generate-sk/config"
	"go-generate-sk/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.GetConfig()

	// config.ConnectDB()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		// AllowOrigins:     "*",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
		// AllowCredentials: true,
	}))

	router.SetupRoutes(app)
	app.Listen(":3000")

}
