package controllers

import (
	"net/http"
	"strconv"
	"time"

	"cineverse/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Admin: Add Movie
func AdminAddMovie(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var m models.Movie
		if err := c.ShouldBind(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&m).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"movie": m})
	}
}

// Admin: Create show
func AdminAddShow(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			MovieID uint    `json:"movie_id" form:"movie_id"`
			Hall    string  `json:"hall" form:"hall"`
			Start   string  `json:"start_time" form:"start_time"` // RFC3339 or custom parse
			Seats   int     `json:"seats_total" form:"seats_total"`
			Price   float64 `json:"price" form:"price"`
		}
		if err := c.ShouldBind(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// parse time
		t, err := time.Parse(time.RFC3339, payload.Start)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time, use RFC3339 format", "example": "2025-01-02T15:04:05Z"})
			return
		}
		show := models.Show{
			MovieID:     payload.MovieID,
			Hall:        payload.Hall,
			StartTime:   t,
			SeatsTotal:  payload.Seats,
			SeatsBooked: 0,
			Price:       payload.Price,
		}
		if err := db.Create(&show).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"show": show})
	}
}

// Admin: List all bookings (with optional status filter)
func AdminListBookings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bookings []models.Booking
		status := c.Query("status")
		query := db.Preload("User").Preload("Show").Preload("Show.Movie")
		if status != "" {
			query = query.Where("status = ?", status)
		}
		if err := query.Find(&bookings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"bookings": bookings})
	}
}

// Admin: Update booking status
func AdminUpdateBookingStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)
		var payload struct {
			Status string `json:"status"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var booking models.Booking
		if err := db.First(&booking, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
			return
		}
		booking.Status = payload.Status
		if err := db.Save(&booking).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"booking": booking})
	}
}
