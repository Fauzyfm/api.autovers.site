package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// ConfigureCORS - Configure CORS middleware untuk allow frontend
// Bisa di-set via environment variable ALLOWED_ORIGINS
func ConfigureCORS() fiber.Handler {
	// Get allowed origins dari environment variable atau default
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")

	// Default allowed origins (development)
	defaultOrigins := []string{
		"http://localhost:3000",
	}

	var allowedOrigins string
	if allowedOriginsEnv != "" {
		allowedOrigins = allowedOriginsEnv
	} else {
		allowedOrigins = strings.Join(defaultOrigins, ",")
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins, // Frontend domains yang boleh akses
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		ExposeHeaders:    "Content-Length,X-JSON-Response",
		AllowCredentials: true, // ⭐ PENTING untuk cookies/auth
		MaxAge:           300,  // Pre-flight cache 5 menit
	})
}
