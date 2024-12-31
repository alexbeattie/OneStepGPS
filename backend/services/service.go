// services/service.go
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
	"github.com/alexbeattie/golangone/config"
	"github.com/alexbeattie/golangone/models"
)

type Service struct {
	db     *gorm.DB
	config *config.Config
	client *http.Client
}

// NewService returns a new instance of Service.
//
// The DB and Config parameters are required and used to initialize the service.
// The HTTP client is initialized with a timeout of 10 seconds.
func NewService(db *gorm.DB, config *config.Config) *Service {
	return &Service{
		db:     db,
		config: config,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *Service) GetPreferences(userID uint) (*models.UserPreferences, error) {
	var preferences models.UserPreferences
	if err := s.db.First(&preferences, userID).Error; err != nil {
		return nil, err
	}
	return &preferences, nil
}

func (s *Service) UpdatePreferences(preferences *models.UserPreferences) error {
	return s.db.Save(preferences).Error
}

func (s *Service) FetchDevices() ([]models.Device, error) {
	url := fmt.Sprintf("https://track.onestepgps.com/v3/api/public/device?latest_point=true&api-key=%s",
		s.config.OneStepGPSAPIKey)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch devices: %w", err)
	}
	defer resp.Body.Close()

	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.ResultList, nil
}

func (s *Service) FetchDeviceInfo(params map[string]string) (*models.DeviceInfoResponse, error) {
    url := fmt.Sprintf("https://track.onestepgps.com/v3/api/public/device-info?lat_lng=1&api-key=%s",
        s.config.OneStepGPSAPIKey)
    
    resp, err := s.client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch device info: %w", err)
    }
    defer resp.Body.Close()

    var response models.DeviceInfoResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode device info response: %w", err)
    }

    return &response, nil
}
func (s *Service) FetchDriveStopRoute(deviceID string, fromTime, toTime time.Time, stopDuration string) (*models.DriveStopResponse, error) {
    url := fmt.Sprintf("https://track.onestepgps.com/v3/api/public/route/drive-stop?device_id=%s&dt_tracker_from=%s&dt_tracker_to=%s&stop_duration=%s&return_points=true&max_return_points=999",
        deviceID,
        fromTime.Format(time.RFC3339),
        toTime.Format(time.RFC3339),
        stopDuration)
    
    url = fmt.Sprintf("%s&api-key=%s", url, s.config.OneStepGPSAPIKey)
    
    resp, err := s.client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch drive-stop route: %w", err)
    }
    defer resp.Body.Close()

    var response models.DriveStopResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode drive-stop response: %w", err)
    }

    return &response, nil
}