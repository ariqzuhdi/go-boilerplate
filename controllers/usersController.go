package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/cheeszy/go-crud/initializers"
	"github.com/cheeszy/go-crud/models"
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

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully."})
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
		"sub": user.ID,                               // user ID sebagai subject
		"exp": time.Now().Add(time.Hour * 24).Unix(), // token expired 24 jam
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// return token ke client
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{
		"text": "Bye",
	})
}
