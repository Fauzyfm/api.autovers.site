package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte(getJwtSecret())

func getJwtSecret() string {
	secret := os.Getenv("SECRET_JWT_AUTOVERS")
	if secret == "" {
		secret = "SECRET_JWT_AUTOVERS" // fallback
	}
	return secret
}

type JwtClaims struct {
	Email    string `json:"email"`
	UserName string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type VerificationClaims struct {
	Email   string `json:"email"`
	Purpose string `json:"purpose"`
	jwt.RegisteredClaims
}

// Generate JWT token
func GenerateToken(email, username, role string) (string, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback ke UTC jika timezone tidak tersedia
		loc = time.UTC
	}
	
	claims := &JwtClaims{
		Email:    email,
		UserName: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().In(loc).Add(time.Hour * 1)), // 1 jam
			IssuedAt:  jwt.NewNumericDate(time.Now().In(loc)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

// Parse & verify token
func ParseToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func GenerateVerificationToken(email string) (string, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback ke UTC jika timezone tidak tersedia
		loc = time.UTC
	}

	claims := &VerificationClaims{
		Email:   email,
		Purpose: "email_verification",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().In(loc).Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().In(loc)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

func ParseVerificationToken(tokenString string) (*VerificationClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&VerificationClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*VerificationClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if claims.Purpose != "email_verification" {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// Generate Reset Password Token
func GenerateResetPasswordToken(email string) (string, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback ke UTC jika timezone tidak tersedia
		loc = time.UTC
	}

	claims := &VerificationClaims{
		Email:   email,
		Purpose: "password_reset", // Purpose untuk membedakan dari email_verification
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().In(loc).Add(30 * time.Minute)), // 30 menit
			IssuedAt:  jwt.NewNumericDate(time.Now().In(loc)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}

// Parse Reset Password Token
func ParseResetPasswordToken(tokenString string) (*VerificationClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&VerificationClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*VerificationClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if claims.Purpose != "password_reset" {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
