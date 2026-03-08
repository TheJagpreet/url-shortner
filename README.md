

# URL Shortener (Golang)

A robust, modular URL shortener service written in Go, featuring API key authentication, Redis-backed storage, and extensible architecture.

---

## Table of Contents
- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Redis Setup](#redis-setup)
- [Build](#build)
- [Run](#run)
- [Authentication](#authentication)
- [API Usage](#api-usage)
- [Testing](#testing)
- [Extending](#extending)
- [Storage Details](#storage-details)

---

## Features

- Shorten URLs to unique codes (random, 6-character alphanumeric)
- Redirect short codes to original URLs (HTTP 302)
- RESTful API:  
     - `POST /shorten` (protected by API key)
     - `GET /{code}` (redirect)
- Redis-backed storage (default, with TTL support)
- API Key authentication middleware
- Unit and integration tests
- Modular, extensible design (easy to add new storage or features)
- In-memory storage option (via shortener package)
- Configurable Redis TTL (default 24h)
- Example Redis CLI commands for manual inspection

---

## Project Structure

```
url-shortner/
├── main.go                # Entry point, HTTP server setup
├── handlers/              # API endpoint handlers
│   └── api.go
├── middleware/            # Authentication middleware
│   └── auth.go
├── shortener/             # URL shortening logic (in-memory)
│   └── shortener.go
├── storage/               # Storage interface & Redis implementation
│   ├── storage.go
│   └── redis.go
├── tests/                 # Unit & integration tests
│   ├── handlers_test.go
│   └── shortener_test.go
├── go.mod                 # Go module definition
└── README.md              # Project documentation
```

---

## Prerequisites

- Go 1.25.7+
- Redis server (for storage backend)

---

## Redis Setup

Install Redis:
```bash
sudo apt update
sudo apt install redis-server
```
Start Redis server:
```bash
sudo systemctl start redis-server
```
Check Redis status:
```bash
sudo systemctl status redis-server
```
By default, the app connects to Redis at `localhost:6379`. Ensure Redis is running before starting the app or running tests.

---

## Build

Build the application:
```bash
go build -o url-shortner main.go
```

---

## Run

Set the API key as an environment variable before running:
```bash
export URL_SHORTENER_API_KEY="your-secret-key"
```

Run the server:
```bash
go run main.go
```
Or use the built binary:
```bash
./url-shortner
```

---

## Authentication

- The `/shorten` endpoint requires an API key via the `X-API-Key` header.
- The API key is set via the `URL_SHORTENER_API_KEY` environment variable.
- If the key is missing or incorrect, requests are rejected with HTTP 401.

---

## API Usage

### POST /shorten

Shorten a URL (requires API key):
```bash
curl -X POST http://localhost:8080/shorten \
           -H "Content-Type: application/json" \
           -H "X-API-Key: your-secret-key" \
           -d '{"url": "https://example.com"}'
# Response: {"code": "abc123"}
```

### GET /{code}

Redirect to original URL:
```bash
curl -v http://localhost:8080/abc123
# Response: HTTP 302 redirect to https://example.com
```

---

## Testing

Run all tests (Redis must be running):
```bash
go test ./tests -v
```
- Tests cover API key authentication, URL shortening, redirection, and Redis storage.

---

## Extending

- Add analytics, rate limiting, or advanced authentication (JWT, OAuth) as needed.
- Swap storage backend by implementing the `Storage` interface.
- Adjust Redis TTL via `NewRedisStorage(addr, ttlSeconds)`.

---

## Storage Details

- URLs are stored in Redis with a default TTL of 24 hours.
- Storage interface allows for easy extension to other backends.
- Example Redis CLI commands:
     - List all keys: `redis-cli KEYS '*'`
     - Get value: `redis-cli GET <code>`
     - Set value: `redis-cli SET <code> <url>`
     - Delete key: `redis-cli DEL <code>`
     - Check TTL: `redis-cli TTL <code>`

---

## License
MIT
