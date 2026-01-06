package services

import (
	"belajar-go-fiber/models"
	"belajar-go-fiber/repositories"
	"belajar-go-fiber/utils"
	"errors"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ==================== PASSWORD SERVICE ====================

// HashPassword - hash password menggunakan bcrypt
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// CheckPasswordHash - bandingkan password dengan hash
func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ==================== REGISTER SERVICE ====================

// RegisterService - handle logic registrasi user
func RegisterService(req *models.RegisterRequest) (models.AuthResponse, error) {
	// Validasi input
	if err := validateRegisterRequest(req); err != nil {
		return models.AuthResponse{}, err
	}

	// Cek email sudah terdaftar
	exists, err := repositories.IsEmailRegistered(req.Email)
	if err != nil {
		return models.AuthResponse{}, errors.New("database error")
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return models.AuthResponse{}, errors.New("failed to hash password")
	}

	// Generate verification token
	tokenVerificationEmail, err := utils.GenerateVerificationToken(req.Email)
	if err != nil {
		return models.AuthResponse{}, errors.New("failed to generate verification token")
	}

	// User data yang akan disimpan
	user := models.Users{
		UserName:          req.UserName,
		Email:             req.Email,
		NoHandphone:       req.NoHandphone,
		Password:          hashedPassword,
		Role:              "user",
		VerificationToken: tokenVerificationEmail,
	}

	// Jika email sudah terdaftar
	if exists {
		// Cek apakah user sudah aktif
		existingUser, err := repositories.FindUserByEmailWithActiveStatus(req.Email)
		if err != nil {
			return models.AuthResponse{}, errors.New("database error")
		}

		// Jika sudah aktif, tidak bisa register ulang
		if existingUser.ActiveUser {
			return models.AuthResponse{}, errors.New("email already registered and verified")
		}

		// Jika belum aktif, update data user yang sebelumnya (re-register)
		if err := repositories.UpdateInactiveUser(req.Email, &user); err != nil {
			return models.AuthResponse{}, errors.New("failed to update user")
		}
	} else {
		// Jika email belum terdaftar, buat user baru
		if err := repositories.CreateUser(&user); err != nil {
			return models.AuthResponse{}, errors.New("failed to create user")
		}
	}

	// Kirim email verifikasi
	if err := sendVerificationEmail(&user); err != nil {
		return models.AuthResponse{}, errors.New("failed to send verification email")
	}

	return models.AuthResponse{
		UserName: req.UserName,
		Message:  "Registration successful, check your email to verify your account",
	}, nil
}

// validateRegisterRequest - validasi request register
func validateRegisterRequest(req *models.RegisterRequest) error {
	if req.UserName == "" || req.Email == "" || req.Password == "" || req.ConfirmPassword == "" || req.NoHandphone == "" {
		return errors.New("all fields are required")
	}

	if req.Password != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	return nil
}

// sendVerificationEmail - kirim email verifikasi
func sendVerificationEmail(user *models.Users) error {
	verifyLink := "https://autovers.site/auth/verify?token=" + user.VerificationToken

	htmlBytes, err := os.ReadFile("templates/email/verify-email.html")
	if err != nil {
		return err
	}

	emailBody := strings.ReplaceAll(
		string(htmlBytes),
		"{{VERIFY_LINK}}",
		verifyLink,
	)

	return utils.SendMail(
		user.Email,
		"Verify your Autovers Account",
		emailBody,
	)
}

// ==================== LOGIN SERVICE ====================

// LoginService - handle logic login user
func LoginService(req *models.LoginRequest) (models.AuthResponse, *models.Users, error) {
	// Validasi input
	if err := validateLoginRequest(req); err != nil {
		return models.AuthResponse{}, nil, err
	}

	// Ambil user dari database
	user, err := repositories.FindUserByEmail(req.Identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.AuthResponse{}, nil, errors.New("invalid credentials")
		}
		return models.AuthResponse{}, nil, errors.New("database error")
	}

	// Cek password
	if !CheckPasswordHash(user.Password, req.Password) {
		return models.AuthResponse{}, nil, errors.New("invalid credentials")
	}

	// Cek apakah user sudah verifikasi email
	if !user.ActiveUser {
		return models.AuthResponse{}, nil, errors.New("please verify your email first")
	}

	return models.AuthResponse{
		UserName: user.UserName,
		Message:  "Login successful",
	}, user, nil
}

// validateLoginRequest - validasi request login
func validateLoginRequest(req *models.LoginRequest) error {
	if req.Identifier == "" || req.Password == "" {
		return errors.New("all fields are required")
	}
	return nil
}

// ==================== VERIFICATION SERVICE ====================

// VerifyEmailService - handle logic verifikasi email
func VerifyEmailService(token string) error {
	if token == "" {
		return errors.New("verification token is required")
	}

	claims, err := utils.ParseVerificationToken(token)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	// Update user sebagai verified
	err = repositories.VerifyUserByEmail(claims.Email)
	if err != nil {
		return errors.New("failed to verify email")
	}

	return nil
}

// ==================== FORGOT PASSWORD SERVICE ====================

// ForgotPasswordService - handle logic forgot password
func ForgotPasswordService(email string) error {
	// Validasi email
	if email == "" {
		return errors.New("email is required")
	}

	// Cek email ada di database
	user, err := repositories.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Jangan reveal apakah email terdaftar atau tidak (security best practice)
			return nil
		}
		return errors.New("database error")
	}

	// Cek user sudah verified
	if !user.ActiveUser {
		return errors.New("please verify your email first")
	}

	// Generate reset password token dengan Purpose: "password_reset"
	resetToken, err := utils.GenerateResetPasswordToken(email)
	if err != nil {
		return errors.New("failed to generate reset token")
	}

	// Save token ke database (gunakan verificationToken field yang ada)
	if err := repositories.SaveResetPasswordToken(user.ID, resetToken); err != nil {
		return errors.New("failed to save reset token")
	}

	// Send reset password email
	if err := sendResetPasswordEmail(user, resetToken); err != nil {
		return errors.New("failed to send reset password email")
	}

	return nil
}

// ResetPasswordService - handle logic reset password
func ResetPasswordService(token, newPassword, confirmPassword string) error {	
	// Validasi token
	if token == "" {
		return errors.New("reset token is required")
	}

	// Validasi password
	if newPassword == "" || confirmPassword == "" {
		return errors.New("all fields are required")
	}

	if newPassword != confirmPassword {
		return errors.New("passwords do not match")
	}

	// Validasi dan parse token dengan purpose "password_reset"
	claims, err := utils.ParseResetPasswordToken(token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	// Cari user dari email di token
	user, err := repositories.FindUserByEmail(claims.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Cek token matches di database
	if user.VerificationToken != token {
		return errors.New("invalid reset token")
	}

	// Hash password baru
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update password dan clear token
	if err := repositories.UpdatePasswordAndClearToken(user.ID, hashedPassword); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

// sendResetPasswordEmail - kirim email reset password
func sendResetPasswordEmail(user *models.Users, resetToken string) error {
	resetLink := "https://autovers.site/auth/reset-password?token=" + resetToken

	htmlBytes, err := os.ReadFile("templates/email/reset-password.html")
	if err != nil {
		return err
	}

	emailBody := strings.ReplaceAll(
		string(htmlBytes),
		"{{RESET_LINK}}",
		resetLink,
	)

	return utils.SendMail(
		user.Email,
		"Reset Your Autovers Password",
		emailBody,
	)
}
