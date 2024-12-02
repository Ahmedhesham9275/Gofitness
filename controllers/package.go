package controllers

import (
	"fitnesshub/database"
	"fitnesshub/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePackage(c *gin.Context) {
	var input models.Package
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	input.UserID = userID.(uint)

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Package"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func GetPackages(c *gin.Context) {
	var Packages []models.Package
	if err := database.DB.Find(&Packages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Packages"})
		return
	}

	c.JSON(http.StatusOK, Packages)
}

func GetPackage(c *gin.Context) {
	var Package models.Package
	if err := database.DB.First(&Package, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}

	c.JSON(http.StatusOK, Package)
}

func UpdatePackage(c *gin.Context) {
	var Package models.Package
	if err := database.DB.First(&Package, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	if Package.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this Package"})
		return
	}

	var input models.Package
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Package.Title = input.Title
	Package.Description = input.Description
	if err := database.DB.Save(&Package).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Package"})
		return
	}

	c.JSON(http.StatusOK, Package)
}

func DeletePackage(c *gin.Context) {
	var Package models.Package
	if err := database.DB.First(&Package, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}

	// Retrieve the user ID from the context
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	// Check if the current user is the owner of the Package
	if Package.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this Package"})
		return
	}

	if err := database.DB.Delete(&Package).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Package"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
