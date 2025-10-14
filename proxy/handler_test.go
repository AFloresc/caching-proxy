package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/cache/stats", StatsHandler).Methods("GET")
	r.HandleFunc("/cache/clear", ClearCacheHandler).Methods("POST")
	return r
}

func TestStatsHandler_EmptyCache(t *testing.T) {
	ClearCache() // Ensure clean state

	req := httptest.NewRequest("GET", "/cache/stats", nil)
	rec := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rec.Code)
	}

	var stats CacheStats
	if err := json.NewDecoder(rec.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if stats.Count != 0 {
		t.Errorf("Expected cache count 0, got %d", stats.Count)
	}
}

func TestStatsHandler_WithCache(t *testing.T) {
	ClearCache()
	SetCache("/test", []byte("data"))

	req := httptest.NewRequest("GET", "/cache/stats", nil)
	rec := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(rec, req)

	var stats CacheStats
	if err := json.NewDecoder(rec.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if stats.Count != 1 || stats.Keys[0] != "/test" {
		t.Errorf("Unexpected cache stats: %+v", stats)
	}
}

func TestClearCacheHandler(t *testing.T) {
	SetCache("/clear-me", []byte("bye"))

	req := httptest.NewRequest("POST", "/cache/clear", strings.NewReader(""))
	rec := httptest.NewRecorder()

	router := setupRouter()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rec.Code)
	}

	var resp ClearResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.Message != "Cache cleared successfully" {
		t.Errorf("Unexpected message: %s", resp.Message)
	}

	// Confirm cache is empty
	if _, ok := GetFromCache("/clear-me"); ok {
		t.Errorf("Expected cache to be cleared")
	}
}
