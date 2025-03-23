package service

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"strings"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/cache"
)

// App struct represents an app entry
type App struct {
	Name  string
	Price string
}

// FetchAllApps reads from the cache or CSV file and returns filtered apps
func FetchAllApps(limit, page int, price string) ([]string, error) {
	// Check if data is already cached
	cachedData, found := cache.Get("apps")
	if found {
		log.Println("Serving from cache")
		return applyPagination(cachedData.([]string), limit, page), nil
	}

	// Open CSV file
	file, err := os.Open("storage/googleplaystore.csv")
	if err != nil {
		log.Println("Error opening CSV file:", err)
		return nil, errors.New("failed to open CSV file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var apps []string
	lineNumber := 0

	for {
		record, err := reader.Read()
		if err != nil {
			// Stop reading on EOF
			if err.Error() == "EOF" {
				break
			}
			log.Printf("Skipping invalid row at line %d: %v\n", lineNumber, err)
			continue
		}

		lineNumber++
		if lineNumber == 1 {
			continue // Skip the header row
		}

		// Ensure the record has enough columns before accessing index 7
		if len(record) < 8 {
			log.Printf("Skipping row %d: expected at least 8 columns, got %d\n", lineNumber, len(record))
			continue
		}

		// Filter apps by price if specified
		if price == "" || strings.TrimSpace(record[7]) == price {
			apps = append(apps, record[0]) // App name is in the first column
		}
	}

	// Store in cache
	cache.Set("apps", apps)

	// Return paginated data
	return applyPagination(apps, limit, page), nil
}

// applyPagination applies pagination logic
func applyPagination(data []string, limit, page int) []string {
	start := (page - 1) * limit
	if start >= len(data) {
		return []string{}
	}
	end := start + limit
	if end > len(data) {
		end = len(data)
	}
	return data[start:end]
}
