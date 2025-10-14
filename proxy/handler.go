package proxy

import (
	"io"
	"net/http"
)

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
