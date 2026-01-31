package main

import (
	"belajar-go-fiber/config"
	"belajar-go-fiber/handlers"
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

	// ⭐ HEALTH CHECK / ROOT ENDPOINT
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "API running",
			"version": "1.0.0",
		})
	})

	// ⭐ INITIALIZE DATABASE
	config.InitDatabase()

	// ⭐ PUBLIC ROUTES (tidak memerlukan authentication)
	routes.AuthRoutes(app)

	// ⭐ PROTECTED ROUTES (hanya /auth/me dan /auth/logout yang perlu authentication)
	app.Get("/auth/me", middlewares.ProtectRoute(), middlewares.RequireRole("admin", "user"), handlers.MeHandler)
	app.Post("/auth/logout", middlewares.ProtectRoute(), middlewares.RequireRole("admin", "user"), handlers.LogoutHandler)

	app.Listen(":8080")
}
