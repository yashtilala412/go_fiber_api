package cache

import (
	"encoding/csv"
	"os"
	"sync"
	"time"

	"git.pride.improwised.dev/Onboarding-2025/Yash-Tilala/config"
)

type CacheItem struct {
	Data       [][]string //csv content store in memory
	Expiration int64      //cache expiration
}

type CSVCache struct {
	data map[string]CacheItem
	mu   sync.RWMutex
}

func NewCsvCache() *CSVCache {
	return &CSVCache{
		data: make(map[string]CacheItem),
	}
}

var GlobalCSVCache = NewCsvCache()

// LoadCSVToCache reads a CSV file and stores it in memory
func (c *CSVCache) LoadCSVToCache(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	//  Use expiration time from config
	expiration := time.Now().Add(config.CacheExpiration).Unix()

	// Store in cache
	c.data[filename] = CacheItem{Data: records, Expiration: expiration}

	return nil
}

// GetCSVData retrieves data from cache
func (c *CSVCache) GetCSVData(filename string) ([][]string, bool) {
	c.mu.RLock()
	item, exists := c.data[filename]
	c.mu.RUnlock()

	if !exists || time.Now().Unix() > item.Expiration {
		return nil, false // Expired or not found
	}
	return item.Data, true
}
