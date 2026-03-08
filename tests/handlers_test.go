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

// mockStorage is an in-memory Storage implementation used for unit tests.
type mockStorage struct {
	urls         map[string]string
	lastTTL      int64
	codeSequence []string
	callCount    int
}

func newMockStorage(codes ...string) *mockStorage {
	return &mockStorage{
		urls:         make(map[string]string),
		codeSequence: codes,
	}
}

func (m *mockStorage) Shorten(url string, ttlSeconds int64) string {
	code := "mockcode"
	if m.callCount < len(m.codeSequence) {
		code = m.codeSequence[m.callCount]
	}
	m.callCount++
	m.lastTTL = ttlSeconds
	m.urls[code] = url
	return code
}

func (m *mockStorage) Resolve(code string) (string, bool) {
	url, ok := m.urls[code]
	return url, ok
}

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

func TestShortenHandlerTTL(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
		wantTTL    int64
	}{
		{
			name:       "custom TTL is forwarded to storage",
			body:       map[string]interface{}{"url": "https://example.com", "ttl_seconds": 3600},
			wantStatus: http.StatusOK,
			wantTTL:    3600,
		},
		{
			name:       "omitted ttl_seconds defaults to 0 (uses backend default)",
			body:       map[string]interface{}{"url": "https://example.com"},
			wantStatus: http.StatusOK,
			wantTTL:    0,
		},
		{
			name:       "negative ttl_seconds is rejected",
			body:       map[string]interface{}{"url": "https://example.com", "ttl_seconds": -1},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := newMockStorage("abc123")
			api := handlers.NewAPI(mock)
			b, _ := json.Marshal(tc.body)
			req := httptest.NewRequest("POST", "/shorten", bytes.NewReader(b))
			w := httptest.NewRecorder()
			api.ShortenHandler(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("Expected status %d, got %d", tc.wantStatus, w.Code)
			}
			if tc.wantStatus == http.StatusOK && mock.lastTTL != tc.wantTTL {
				t.Errorf("Expected TTL %d passed to storage, got %d", tc.wantTTL, mock.lastTTL)
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	store := storage.NewRedisStorage("localhost:6379")
	api := handlers.NewAPI(store)
	url := "https://test.com"
	code := store.Shorten(url, 0) // 0 uses the default TTL
	req := httptest.NewRequest("GET", "/"+code, nil)
	w := httptest.NewRecorder()
	api.RedirectHandler(w, req)
	if w.Code != http.StatusFound {
		t.Errorf("Expected 302 Found, got %d", w.Code)
	}
}

