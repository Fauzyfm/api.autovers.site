package main

import (
	"belajar-go-fiber/config"
	"belajar-go-fiber/middlewares"
	"belajar-go-fiber/routes"

	_ "belajar-go-fiber/docs" // ⭐ Auto-generated docs

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	app := fiber.New()
	godotenv.Load()

	// ⭐ CORS MIDDLEWARE (harus di awal, sebelum routes)
	app.Use(middlewares.ConfigureCORS())

	// ⭐ SWAGGER DOCUMENTATION UI
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// ⭐ PUBLIC ROUTES
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API running",
		})
	})

	config.InitDatabase()
	routes.AuthRoutes(app)

	// ⭐ PROTECTED ROUTES (menggunakan ProtectRoute middleware)
	protected := app.Group("")
	protected.Use(middlewares.ProtectRoute())
	routes.ProtectedRoutes(protected)

	app.Listen(":8080")
}
