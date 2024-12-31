package handlers

import (
	// "log"
	"errors"
	"net/http"
	// "strconv"
	"time"

	"gorm.io/gorm"

	"github.com/alexbeattie/golangone/models"
	"github.com/alexbeattie/golangone/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
	db      *gorm.DB
}

func NewHandler(service *services.Service, db *gorm.DB) *Handler {
	return &Handler{
		service: service,
		db:      db,
	}
}

func (h *Handler) GetUserPreferences(c *gin.Context) {
    deviceID := c.Param("userId") // Using deviceID instead of userID
    
    var prefs models.UserPreferences
    result := h.db.Where("user_id = ?", deviceID).First(&prefs)
    
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            // Return default preferences if none found
            prefs = models.UserPreferences{
                UserID:          deviceID,
                ShowAddress:     true,
                ShowEngineHours: true,
                ShowOdometer:    true,
                ShowVin:         true,
                ShowSpeed:       true,
                ShowHeading:     true,
                ShowBattery:     true,
                ShowSatellites:  true,
                ShowLastUpdate:  true,
            }
            c.JSON(http.StatusOK, prefs)
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch preferences"})
        return
    }
    
    c.JSON(http.StatusOK, prefs)
}

func (h *Handler) UpdateUserPreferences(c *gin.Context) {
    userId := c.Param("userId")
    var preferences models.UserPreferences
    
    if err := c.ShouldBindJSON(&preferences); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    preferences.UserID = userId

    // Start transaction
    tx := h.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Try to find existing preferences
    var existingPrefs models.UserPreferences
    result := tx.Where("user_id = ?", userId).First(&existingPrefs)
    
    if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    if result.Error == nil {
        // Update existing
        preferences.ID = existingPrefs.ID
        if err := tx.Save(&preferences).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
            return
        }
    } else {
        // Create new
        if err := tx.Create(&preferences).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create preferences"})
            return
        }
    }

    // Commit transaction
    tx.Commit()

    c.JSON(http.StatusOK, preferences)
}


func (h *Handler) GetDevices(c *gin.Context) {
	devices, err := h.service.FetchDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"devices": devices})
}
func (h *Handler) GetDeviceInfo(c *gin.Context) {
    // You can add query params handling if needed
    deviceInfo, err := h.service.FetchDeviceInfo(nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, deviceInfo)
}
func (h *Handler) GetDriveStopRoute(c *gin.Context) {
    deviceID := c.Query("device_id")
    if deviceID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
        return
    }

    fromStr := c.Query("dt_tracker_from")
    toStr := c.Query("dt_tracker_to")
    stopDuration := c.DefaultQuery("stop_duration", "5m0s")

    from, err := time.Parse(time.RFC3339, fromStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format"})
        return
    }

    to, err := time.Parse(time.RFC3339, toStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format"})
        return
    }

    routeData, err := h.service.FetchDriveStopRoute(deviceID, from, to, stopDuration)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, routeData)
}