package urlshortner_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"url-shortner/handlers"
	"url-shortner/middleware"
	"url-shortner/storage"
)

func TestShortenHandler(t *testing.T) {
	store := storage.NewRedisStorage("localhost:6379")
	api := handlers.NewAPI(store)
	reqBody := map[string]string{"url": "https://test.com"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(b))
	w := httptest.NewRecorder()
	// Set valid API key header
	os.Setenv("URL_SHORTENER_API_KEY", "test-key")
	req.Header.Set("X-API-Key", "test-key")
	mux := http.NewServeMux()
	mux.Handle("/shorten", middleware.APIKeyAuth(http.HandlerFunc(api.ShortenHandler)))
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}

	// Test missing/invalid API key
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/shorten", bytes.NewReader(b))
	mux.ServeHTTP(w2, req2)
	if w2.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", w2.Code)
	}
}

func TestRedirectHandler(t *testing.T) {
	store := storage.NewRedisStorage("localhost:6379")
	api := handlers.NewAPI(store)
	url := "https://test.com"
	code := store.Shorten(url)
	req := httptest.NewRequest("GET", "/"+code, nil)
	w := httptest.NewRecorder()
	api.RedirectHandler(w, req)
	if w.Code != http.StatusFound {
		t.Errorf("Expected 302 Found, got %d", w.Code)
	}
}
