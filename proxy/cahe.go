package proxy

import (
	"sync"
)

var (
	cache     = make(map[string][]byte)
	cacheLock sync.RWMutex
)

func GetFromCache(key string) ([]byte, bool) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	val, ok := cache[key]
	return val, ok
}

func SetCache(key string, data []byte) {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cache[key] = data
}

func ClearCache() {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	cache = make(map[string][]byte)
}
