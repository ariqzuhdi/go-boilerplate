package services

import (
	"errors"
	"os"
	"time"

	"github.com/cheeszy/journaling/dto"
	"github.com/cheeszy/journaling/initializers"
	"github.com/cheeszy/journaling/models"
	"github.com/cheeszy/journaling/repositories"
	"github.com/cheeszy/journaling/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(input dto.RegisterRequest) (string, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	recoveryKey, err := utils.GenerateRecoveryKey()
	if err != nil {
		return "", err
	}

	user := models.User{
		Username:    input.Username,
		Email:       input.Email,
		Password:    string(hashedPassword),
		RecoveryKey: recoveryKey,
	}

	token, err := GenerateToken(32)
	if err != nil {
		return "", err
	}

	user.VerificationToken = token
	user.VerificationExpiresAt = time.Now().Add(15 * time.Minute)

	if err := repositories.CreateUser(initializers.DB, &user); err != nil {
		return "", err
	}

	go SendVerificationEmail(user.Email, token, user.RecoveryKey)

	return recoveryKey, nil
}

func LoginUser(input dto.LoginRequest) (string, int64, error) {
	user, err := repositories.FindUserByEmailOrUsername(initializers.DB, input.Identifier)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		return "", 0, errors.New("Invalid email/username or password")
	}

	if !user.IsVerified {
		return "", 0, errors.New("Please verify your email")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", 0, err
	}

	return tokenString, time.Now().Add(time.Hour * 24).Unix(), nil
}

func VerifyUserEmail(token string) error {
	user, err := repositories.FindUserByVerificationToken(initializers.DB, token)
	if err != nil {
		return errors.New("Invalid token")
	}

	if time.Now().After(user.VerificationExpiresAt) {
		return errors.New("Token expired")
	}

	user.IsVerified = true
	user.VerificationToken = ""
	user.VerificationExpiresAt = time.Time{}
	user.ResendCount = 0

	return repositories.UpdateUser(initializers.DB, user)
}

func ResendVerificationEmail(email string) (map[string]interface{}, error) {
	user, err := repositories.FindUserByEmail(initializers.DB, email)
	if err != nil {
		return nil, errors.New("User not found")
	}

	if user.IsVerified {
		return nil, errors.New("User is already verified")
	}

	const maxResend = 3
	const cooldown = 24 * time.Hour

	if time.Since(user.VerificationExpiresAt) > cooldown {
		user.ResendCount = 0
	}

	if user.ResendCount >= maxResend {
		return nil, errors.New("Max resend limit reached")
	}

	token, _ := GenerateToken(32)
	expiry := time.Now().Add(15 * time.Minute)

	user.VerificationToken = token
	user.VerificationExpiresAt = expiry
	user.ResendCount++
	user.LastVerificationSentAt = time.Now()

	if err := repositories.UpdateUser(initializers.DB, user); err != nil {
		return nil, err
	}

	go SendVerificationEmail(user.Email, token, user.RecoveryKey)

	return map[string]interface{}{
		"message":         "Verification email resent successfully",
		"expires_at":      expiry,
		"resend_count":    user.ResendCount,
		"resend_limit":    maxResend,
		"remaining_quota": maxResend - user.ResendCount,
	}, nil
}
