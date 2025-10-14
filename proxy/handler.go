package proxy

import (
	"encoding/json"
	"io"
	"net/http"
)

type CacheStats struct {
	Count int      `json:"count"`
	Keys  []string `json:"keys"`
}

type ClearResponse struct {
	Message string `json:"message"`
}

func ProxyHandler(origin string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cacheKey := r.URL.Path + "?" + r.URL.RawQuery

		if cached, ok := GetFromCache(cacheKey); ok {
			w.Header().Set("X-Cache", "HIT")
			w.Write(cached)
			return
		}

		resp, err := http.Get(origin + r.URL.Path + "?" + r.URL.RawQuery)
		if err != nil {
			http.Error(w, "Error contacting origin", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading origin response", http.StatusInternalServerError)
			return
		}

		SetCache(cacheKey, body)

		for k, v := range resp.Header {
			for _, val := range v {
				w.Header().Add(k, val)
			}
		}
		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}
}
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	cacheLock.RLock()
	defer cacheLock.RUnlock()

	keys := make([]string, 0, len(cache))
	for k := range cache {
		keys = append(keys, k)
	}

	stats := CacheStats{
		Count: len(cache),
		Keys:  keys,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func ClearCacheHandler(w http.ResponseWriter, r *http.Request) {
	ClearCache()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ClearResponse{
		Message: "Cache cleared successfully",
	})
}
