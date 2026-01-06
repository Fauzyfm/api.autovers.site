package middlewares

import (
	"belajar-go-fiber/utils"

	"github.com/gofiber/fiber/v2"
)

// RequireRole - Middleware untuk check apakah user role sesuai dengan allowed roles
// Cara pakai:
//
//	app.Get("/endpoint", RequireRole("admin", "user"), handler)
//	app.Get("/admin-only", RequireRole("admin"), handler)
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil role dari context (sudah di-set oleh ProtectRoute middleware)
		userRole, ok := c.Locals("role").(string)
		if !ok || userRole == "" {
			return utils.JSONError(c, fiber.StatusUnauthorized, "Unauthorized - role not found")
		}

		// 2. Check apakah user role ada di allowed roles
		isRoleAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isRoleAllowed = true
				break
			}
		}

		// 3. Jika role tidak diizinkan, return forbidden
		if !isRoleAllowed {
			return utils.JSONError(c, fiber.StatusForbidden, "Forbidden - insufficient permissions")
		}

		// 4. Lanjut ke handler
		return c.Next()
	}
}
