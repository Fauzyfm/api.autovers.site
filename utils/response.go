package utils


import "github.com/gofiber/fiber/v2"

// JSONError → gunakan untuk return error
func JSONError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"message": message})
}

// JSONSuccess → gunakan untuk return success response
func JSONSuccess(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(data)
}