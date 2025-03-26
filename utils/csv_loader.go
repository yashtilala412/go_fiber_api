package utils

import (
	"encoding/csv"
	"os"
	"sync"

	"go.uber.org/zap"
)

// Global variable to store cached data
var (
	cachedData [][]string
	cacheOnce  sync.Once
	cacheLock  sync.RWMutex
)

// LoadCSV reads CSV data into memory once (caching) and skips faulty rows
func LoadCSV(logger *zap.Logger, filePath string) ([][]string, error) {
	logger.Info("Attempting to open CSV file", zap.String("filePath", filePath)) // Debug log

	cacheOnce.Do(func() {
		if filePath == "" {
			logger.Error("CSV file path is empty!")
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			logger.Error("Failed to open CSV file", zap.String("filePath", filePath), zap.Error(err))
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
			zap.String("filePath", filePath),
			zap.Int("totalRows", len(data)),
			zap.Int("skippedRows", skippedRows))
	})

	// Return cached data
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	return cachedData, nil
}
