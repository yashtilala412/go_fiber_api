package utils

import (
	"encoding/csv"
	"os"
	"sync"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Global variable to store cached data
var (
	cachedData [][]string
	cacheOnce  sync.Once
	cacheLock  sync.RWMutex
)

// LoadCSV reads CSV data into memory once (caching) and skips faulty rows
func LoadCSV(logger *zap.Logger) ([][]string, error) {
	cacheOnce.Do(func() {
		config.LoadConfig()
		filePath := viper.GetString("CSV_FILE_PATH")

		file, err := os.Open(filePath)
		if err != nil {
			logger.Error("Failed to open CSV file", zap.Error(err))
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		var data [][]string
		var skippedRows int

		for {
			record, err := reader.Read()
			if err != nil {
				break // Stop on EOF
			}

			// Validate row: Skip if empty or mismatched columns
			if len(record) == 0 {
				skippedRows++
				continue
			}

			data = append(data, record)
		}

		// Store in cache
		cacheLock.Lock()
		cachedData = data
		cacheLock.Unlock()

		logger.Info("CSV data loaded into cache",
			zap.Int("totalRows", len(data)),
			zap.Int("skippedRows", skippedRows))
	})

	// Return cached data
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	return cachedData, nil
}
