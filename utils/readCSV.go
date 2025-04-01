package utils

import (
	"errors"
	"os"
	"path/filepath"
)

// ReadCSV reads a CSV file and returns its records
func ReadCSV(path string) ([]byte, error) {

	filePath := filepath.Join(path)

	// Open file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("failed to open CSV file")
	}

	return fileData, nil
}
