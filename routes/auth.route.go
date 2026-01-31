package routes

import (
	"belajar-go-fiber/handlers"

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

// ProtectedRoutes - Function ini tidak digunakan lagi karena protected routes
// sekarang didaftarkan langsung di main.go untuk menghindari konflik middleware
func ProtectedRoutes(app fiber.Router) {
	// DEPRECATED: Protected routes sekarang di main.go
	// Ini dibiarkan untuk compatibility jika ada referensi lain
}
