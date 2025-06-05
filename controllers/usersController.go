package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/cheeszy/go-crud/initializers"
	"github.com/cheeszy/go-crud/models"
	"github.com/cheeszy/go-crud/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	token, err := services.GenerateToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	user.VerificationToken = token
	user.VerificationExpiresAt = time.Now().Add(15 * time.Minute) // Token valid for 24 hours

	initializers.DB.Save(&user)

	go services.SendVerificationEmail(user.Email, token)

	c.JSON(http.StatusOK, gin.H{"message": "User registered. Please check your email to verify your account."})
}

func Users(c *gin.Context) {
	var users []models.User
	initializers.DB.Find(&users)

	//Respond with them
	c.JSON(200, gin.H{
		"Users": users,
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := initializers.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// password checking (misal bcrypt)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// IsVerified checking
	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Please verify your email.",
		})
		return
	}

	// create JWT token using user ID as sub
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,                               // user ID as subject
		"exp": time.Now().Add(time.Hour * 24).Unix(), // token expired 24 hours later
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// return token to client
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{
		"text": "Bye",
	})
}

func VerifyEmail(c *gin.Context) {

	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
	}

	var user models.User
	if err := initializers.DB.Where("verification_token = ?", token).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or token is invalid"})
		return
	}

	if time.Now().After(user.VerificationExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token has expired"})
		return
	}

	// Update the user's IsVerified status
	user.IsVerified = true
	user.VerificationExpiresAt = time.Time{}
	user.VerificationToken = ""
	user.ResendCount = 0

	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func ResendVerificationEmail(c *gin.Context) {

	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := initializers.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already verified"})
		return
	}

	const maxResendCount = 3
	const cooldownPeriod = 24 * time.Hour

	if time.Since(user.VerificationExpiresAt) > cooldownPeriod {
		user.ResendCount = 0
	}

	if user.ResendCount >= maxResendCount {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Maximum resend limit reached. Try again 24 hours later."})
		return
	}

	token, err := services.GenerateToken(32)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	expiry := time.Now().Add(15 * time.Minute)
	user.VerificationToken = token

	user.VerificationExpiresAt = expiry
	user.ResendCount++
	user.LastVerificationSentAt = time.Now()

	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	go services.SendVerificationEmail(user.Email, token)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Verification email resent successfully",
		"expires_at":      user.VerificationExpiresAt,
		"resend_count":    user.ResendCount,
		"resend_limit":    maxResendCount,
		"remaining_quota": maxResendCount - user.ResendCount,
	})

}
