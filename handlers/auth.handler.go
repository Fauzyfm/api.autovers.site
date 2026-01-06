package handlers

import (
	"belajar-go-fiber/models"
	"belajar-go-fiber/services"
	"belajar-go-fiber/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// @Summary Register user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register request"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/register [post]
// RegisterHandler - HTTP handler untuk registrasi
func RegisterHandler(c *fiber.Ctx) error {
	req := new(models.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return utils.JSONError(c, 400, "Registration failed")
	}

	// Panggil service untuk logic bisnis
	response, err := services.RegisterService(req)
	if err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSuccess(c, 201, response)
}

// @Summary Login user
// @Description Login with email/username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
// LoginHandler - HTTP handler untuk login
func LoginHandler(c *fiber.Ctx) error {
	req := new(models.LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return utils.JSONError(c, 400, "Login failed")
	}

	// Panggil service untuk logic bisnis
	response, user, err := services.LoginService(req)
	if err != nil {
		return utils.JSONError(c, 401, err.Error())
	}

	// Set cookie dengan token JWT
	loc, _ := time.LoadLocation("Asia/Jakarta")
	token, err := utils.GenerateToken(user.Email, user.UserName, user.Role)
	if err != nil {
		return utils.JSONError(c, 500, "Failed to generate token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   false, // ubah ke true di production (https)
		Expires:  time.Now().In(loc).Add(1 * time.Hour),
	})

	return utils.JSONSuccess(c, 200, response)
}

// @Summary Get current user info
// @Description Get authenticated user information
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserInfo
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/me [get]
// MeHandler - HTTP handler untuk ambil info user dari context (middleware)
// ⭐ CATATAN: Route ini sudah di-protect oleh middleware, jadi user sudah authenticated
func MeHandler(c *fiber.Ctx) error {
	// Ambil data dari context yang sudah di-set oleh middleware
	email := c.Locals("email").(string)
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)

	return utils.JSONSuccess(c, 200, models.UserInfo{
		Email:    email,
		Username: username,
		Role:     role,
	})
}

// @Summary Verify email
// @Description Verify user email with token from email link
// @Tags auth
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/verify [get]
// VerificationEmailHandler - HTTP handler untuk verifikasi email
func VerificationEmailHandler(c *fiber.Ctx) error {
	token := c.Query("token")

	// Panggil service untuk logic bisnis
	err := services.VerifyEmailService(token)
	if err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSuccess(c, 200, models.MessageResponse{
		Message: "Email verified successfully",
	})
}

// @Summary Logout user
// @Description Clear authentication cookie
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.MessageResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/logout [post]
// LogoutHandler - HTTP handler untuk logout
// ⭐ CATATAN: Route ini sudah di-protect oleh middleware, hanya user authenticated yang bisa logout
func LogoutHandler(c *fiber.Ctx) error {
	// Clear cookie "auth_token" dengan set MaxAge ke -1
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false, // ubah ke true di production (https)
		MaxAge:   -1,    // MaxAge: -1 untuk delete cookie
	})

	return utils.JSONSuccess(c, 200, models.MessageResponse{
		Message: "Logout successful",
	})
}

// @Summary Request password reset
// @Description Send reset password link to email
// @Tags auth
// @Accept json
// @Produce json
// @Param email body map[string]string true "Email address"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/forgot-password [post]
// ForgotPasswordHandler - HTTP handler untuk request reset password
func ForgotPasswordHandler(c *fiber.Ctx) error {
	type ForgotPasswordRequest struct {
		Email string `json:"email"`
	}

	req := new(ForgotPasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return utils.JSONError(c, 400, "Invalid request")
	}

	// Panggil service untuk generate token reset
	err := services.ForgotPasswordService(req.Email)
	if err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSuccess(c, 200, models.MessageResponse{
		Message: "Reset password link sent to your email",
	})
}

// @Summary Reset password
// @Description Reset user password with reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Reset token and new password"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/reset-password [post]
// ResetPasswordHandler - HTTP handler untuk reset password dengan token
func ResetPasswordHandler(c *fiber.Ctx) error {
	type ResetPasswordRequest struct {
		Token           string `json:"token"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	req := new(ResetPasswordRequest)
	if err := c.BodyParser(req); err != nil {
		return utils.JSONError(c, 400, "Invalid request")
	}

	// Panggil service untuk reset password
	err := services.ResetPasswordService(req.Token, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		return utils.JSONError(c, 400, err.Error())
	}

	return utils.JSONSuccess(c, 200, models.MessageResponse{
		Message: "Password reset successfully",
	})
}
