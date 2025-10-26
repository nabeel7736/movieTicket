package controllers

import (
	"cineverse/models"
	"cineverse/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type registerReq struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't hash password"})
			return
		}
		user := models.User{
			FullName: req.FullName,
			Email:    req.Email,
			Password: hashed,
			IsAdmin:  false,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already used or invalid data", "details": err.Error()})
			return
		}
		// Return minimal info
		c.JSON(http.StatusCreated, gin.H{"message": "user registered", "user_id": user.ID})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if !utils.CheckPasswordHash(user.Password, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		token, err := utils.CreateToken(user.ID, user.IsAdmin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token":      token,
			"expires_at": time.Now().Add(72 * time.Hour),
			"user":       gin.H{"id": user.ID, "full_name": user.FullName, "email": user.Email, "is_admin": user.IsAdmin},
		})
	}
}
