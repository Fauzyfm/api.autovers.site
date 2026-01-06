package repositories

import (
	"belajar-go-fiber/config"
	"belajar-go-fiber/models"
)

func VerifyUserByEmail(email string) error {
	return config.DB.Model(&models.Users{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"activeUser":        true,
			"verificationToken": "true",
		}).Error
}

// Find user by email
func FindUserByEmail(email string) (*models.Users, error) {
	var user models.Users
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create user in DB
func CreateUser(user *models.Users) error {
	return config.DB.Create(user).Error
}

// Check if email already exists
func IsEmailRegistered(email string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Users{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Find user by email with active status
func FindUserByEmailWithActiveStatus(email string) (*models.Users, error) {
	var user models.Users
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update inactive user (re-register case)
func UpdateInactiveUser(email string, user *models.Users) error {
	return config.DB.Model(&models.Users{}).
		Where("email = ? AND \"activeUser\" = ?", email, false).
		Updates(map[string]interface{}{
			"userName":          user.UserName,
			"noHandphone":       user.NoHandphone,
			"password":          user.Password,
			"verificationToken": user.VerificationToken,
		}).Error
}

// SaveResetPasswordToken - Save reset password token ke field verificationToken
func SaveResetPasswordToken(userID, token string) error {
	return config.DB.Model(&models.Users{}).
		Where("id = ?", userID).
		Update("verificationToken", token).Error
}

// UpdatePasswordAndClearToken - Update password dan clear reset token
func UpdatePasswordAndClearToken(userID, hashedPassword string) error {
	return config.DB.Model(&models.Users{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password":          hashedPassword,
			"verificationToken": "true", // Clear token setelah reset
		}).Error
}
