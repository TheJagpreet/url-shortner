package storage

// Storage defines the interface for URL storage backends
// Extend this for Redis, in-memory, or other implementations

type Storage interface {
	// Shorten stores the given URL and returns a short code.
	// ttlSeconds sets the expiry for this entry; 0 means use the backend default.
	Shorten(url string, ttlSeconds int64) string
	Resolve(code string) (string, bool)
}
