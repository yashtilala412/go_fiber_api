package services

import (
	"errors"
	"sync"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/utils"
	"go.uber.org/zap"
)

var (
	appData    [][]string
	reviewData [][]string
	once       sync.Once
	loadErr    error
)

// LoadAppData loads CSV data once and caches it
func LoadAppData(logger *zap.Logger) ([][]string, [][]string, error) {
	once.Do(func() {
		logger.Info("Loading CSV data...")

		data, err := utils.LoadCSV(logger) // Ensure LoadCSV uses the logger
		if err != nil {
			logger.Error("Failed to load CSV data", zap.Error(err))
			loadErr = err
			return
		}

		if len(data) > 0 {
			appData = data[:len(data)/2]    // First half as appData
			reviewData = data[len(data)/2:] // Second half as reviewData
		} else {
			loadErr = errors.New("CSV data is empty")
		}

		logger.Info("CSV data successfully loaded",
			zap.Int("appRows", len(appData)),
			zap.Int("reviewRows", len(reviewData)),
		)
	})

	return appData, reviewData, loadErr
}

// GetApps retrieves cached app data
func GetApps() [][]string {
	return appData
}
