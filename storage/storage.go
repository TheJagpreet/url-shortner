package storage

// Storage defines the interface for URL storage backends
// Extend this for Redis, in-memory, or other implementations

type Storage interface {
	Shorten(url string) string
	Resolve(code string) (string, bool)
}
