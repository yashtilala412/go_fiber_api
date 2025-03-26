// LoadAppData loads App & Review CSV data once and caches it
package services

import (
	"sync"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/config"
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
func LoadAppData(logger *zap.Logger, cfg *config.Config) ([][]string, [][]string, error) {
	once.Do(func() {
		logger.Info("Loading CSV data...")

		// Load main app data CSV
		data, err := utils.LoadCSV(logger, cfg.CSVFilePath)
		if err != nil {
			logger.Error("Failed to load app CSV", zap.Error(err))
			loadErr = err
			return
		}

		// Load review CSV data
		review, err := utils.LoadCSV(logger, cfg.ReviewCSVPath)
		if err != nil {
			logger.Error("Failed to load review CSV", zap.Error(err))
			loadErr = err
			return
		}

		appData = data
		reviewData = review

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
func GetReview() [][]string {
	return reviewData
}
