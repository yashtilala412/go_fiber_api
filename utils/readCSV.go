package utils

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

// ReadCSV reads a CSV file and returns its records
func ReadCSV(path string) ([][]string, error) {
	// Get working directory path
	workingDirPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Construct full file path
	filePath := filepath.Join(workingDirPath, path)

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}
