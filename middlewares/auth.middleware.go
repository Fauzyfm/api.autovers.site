package middlewares

import (
	"belajar-go-fiber/utils"

	"github.com/gofiber/fiber/v2"
)

// ProtectRoute - Middleware untuk protect route yang perlu authentication
// Ambil token dari Cookie "auth_token" dan validasi
func ProtectRoute() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil token dari cookie "auth_token"
		token := c.Cookies("auth_token")

		// 2. Jika tidak ada token
		if token == "" {
			return utils.JSONError(c, fiber.StatusUnauthorized, "Unauthorized - No token provided")
		}

		// 3. Validasi dan parse token
		claims, err := utils.ParseToken(token)
		if err != nil {
			return utils.JSONError(c, fiber.StatusUnauthorized, "Unauthorized - Invalid token")
		}

		// 4. Simpan user info ke context untuk digunakan di handler
		c.Locals("email", claims.Email)
		c.Locals("username", claims.UserName)
		c.Locals("role", claims.Role)

		// 5. Lanjut ke handler berikutnya
		return c.Next()
	}
}

// OptionalAuth - Middleware untuk optional authentication (user bisa ada atau tidak)
// Jika ada token dan valid, simpan info ke context. Jika tidak ada, lanjutkan saja
func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("auth_token")

		// Jika ada token, coba validasi
		if token != "" {
			claims, err := utils.ParseToken(token)
			if err == nil {
				// Token valid, simpan ke context
				c.Locals("email", claims.Email)
				c.Locals("username", claims.UserName)
				c.Locals("role", claims.Role)
				c.Locals("isAuthenticated", true)
			}
		} else {
			// Tidak ada token, tandai belum authenticated
			c.Locals("isAuthenticated", false)
		}

		return c.Next()
	}
}
