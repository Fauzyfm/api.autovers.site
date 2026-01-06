package routes

import (
	"belajar-go-fiber/handlers"
	"belajar-go-fiber/middlewares"

	"github.com/gofiber/fiber/v2"
)

// AuthRoutes - Public routes (Register, Login, Verify Email, Forgot Password, Reset Password)
func AuthRoutes(app *fiber.App) {
	app.Post("/auth/register", handlers.RegisterHandler)
	app.Post("/auth/login", handlers.LoginHandler)
	app.Get("/auth/verify", handlers.VerificationEmailHandler)
	app.Post("/auth/forgot-password", handlers.ForgotPasswordHandler)
	app.Post("/auth/reset-password", handlers.ResetPasswordHandler)
}

// ProtectedRoutes - Routes yang memerlukan authentication
func ProtectedRoutes(app fiber.Router) {
	// Hanya admin dan user yang bisa akses
	app.Get("/auth/me", middlewares.RequireRole("admin", "user"), handlers.MeHandler)
	app.Post("/auth/logout", middlewares.RequireRole("admin", "user"), handlers.LogoutHandler)
	// Tambahkan route protected lainnya di sini
	// Contoh: app.Get("/users/profile", middlewares.RequireRole("admin", "user"), handlers.GetProfile)
	// Contoh admin-only: app.Get("/admin/dashboard", middlewares.RequireRole("admin"), handlers.AdminDashboard)
}
