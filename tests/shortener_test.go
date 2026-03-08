package urlshortner_test

import (
	"testing"
	"url-shortner/storage"
)

func TestRedisStorage(t *testing.T) {
	store := storage.NewRedisStorage("localhost:6379")
	url := "https://example.org"
	code := store.Shorten(url, 0) // 0 uses the default TTL
	if len(code) == 0 {
		t.Errorf("Expected non-empty code")
	}
	resolved, ok := store.Resolve(code)
	if !ok || resolved != url {
		t.Errorf("RedisStorage failed to resolve code")
	}
}
