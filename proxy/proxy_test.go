package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyHandler_CacheMissThenHit(t *testing.T) {
	// Simula el servidor origen
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"hello"}`))
	}))
	defer origin.Close()

	ClearCache()

	handler := ProxyHandler(origin.URL)

	// 1st request (expected MISS)
	req1 := httptest.NewRequest("GET", "/greeting", nil)
	rec1 := httptest.NewRecorder()
	handler.ServeHTTP(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rec1.Code)
	}
	if rec1.Header().Get("X-Cache") != "MISS" {
		t.Errorf("Expected X-Cache: MISS, got %s", rec1.Header().Get("X-Cache"))
	}
	body1, _ := io.ReadAll(rec1.Body)
	if string(body1) != `{"message":"hello"}` {
		t.Errorf("Unexpected body: %s", body1)
	}

	// Second request (expected HIT)
	req2 := httptest.NewRequest("GET", "/greeting", nil)
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req2)

	if rec2.Header().Get("X-Cache") != "HIT" {
		t.Errorf("Expected X-Cache: HIT, got %s", rec2.Header().Get("X-Cache"))
	}
	body2, _ := io.ReadAll(rec2.Body)
	if string(body2) != `{"message":"hello"}` {
		t.Errorf("Unexpected body: %s", body2)
	}
}
