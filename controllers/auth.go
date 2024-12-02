package controllers

import (
	"fitnesshub/database"
	"fitnesshub/models"
	"fitnesshub/utils"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate username
	if len(input.Username) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username cannot be empty"})
		return
	}

	// Validate password
	if len(input.Password) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 4 characters long"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	input.Password = string(hashedPassword)

	if err := database.DB.Create(&input).Error; err != nil {
		// Check if the error is a unique constraint violation
		if isUniqueViolationError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	response := struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}{
		ID:       input.ID,
		Username: input.Username,
	}

	c.JSON(http.StatusCreated, gin.H{"user": response})
}

// isUniqueViolationError checks if the error is a unique constraint violation
func isUniqueViolationError(err error) bool {
	// This regex matches PostgreSQL's unique violation error code as we just have username as unique
	re := regexp.MustCompile(`unique.*constraint`)
	return re.MatchString(err.Error())
}

func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
