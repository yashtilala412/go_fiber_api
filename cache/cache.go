package cache

import (
	"sync"
	"time"
)

var cacheMap = make(map[string]interface{})
var mutex = sync.RWMutex{}

// Set stores a value in the cache
func Set(key string, value interface{}) {
	mutex.Lock()
	cacheMap[key] = value
	mutex.Unlock()
}

// Get retrieves a value from the cache
func Get(key string) (interface{}, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	value, found := cacheMap[key]
	return value, found
}

// ClearCache clears all cache data
func ClearCache() {
	mutex.Lock()
	defer mutex.Unlock()
	cacheMap = make(map[string]interface{})
}

// AutoClearCache clears cache every X seconds (for development only)
func AutoClearCache(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			ClearCache()
		}
	}()
}
