
# URL Shortener (Golang)

A robust, modular URL shortener service written in Go, featuring API key authentication, Redis-backed storage with per-request TTL control, and an extensible architecture.

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
  - [POST /shorten](#post-shorten)
  - [GET /{code}](#get-code)
- [TTL Configuration](#ttl-configuration)
- [Testing](#testing)
- [Extending](#extending)
- [Redis CLI Reference](#redis-cli-reference)

---

## Features

- Shorten URLs to unique 6-character alphanumeric codes
- Redirect short codes to original URLs (HTTP 302)
- RESTful API:
  - `POST /shorten` — protected by API key
  - `GET /{code}` — public redirect
- Redis-backed storage with configurable TTL
  - Global default TTL set at startup (default: 24 hours)
  - Per-request TTL override via `ttl_seconds` JSON field
- API key authentication middleware (`X-API-Key` header)
- Modular, extensible design — swap storage backends via the `Storage` interface
- Comprehensive unit tests (mock storage, no Redis required for handler tests)

---

## Project Structure

```
url-shortner/
├── main.go                # Entry point, HTTP server setup
├── handlers/              # HTTP request handlers
│   └── api.go
├── middleware/            # Authentication middleware
│   └── auth.go
├── shortener/             # Code generation logic
│   └── shortener.go
├── storage/               # Storage interface & Redis implementation
│   ├── storage.go
│   └── redis.go
├── tests/                 # Unit & integration tests
│   ├── handlers_test.go
│   └── shortener_test.go
├── go.mod
└── README.md
```

---

## Prerequisites

- Go 1.25.7+
- Redis server

---

## Redis Setup

```bash
# Install
sudo apt update && sudo apt install redis-server

# Start
sudo systemctl start redis-server

# Verify
sudo systemctl status redis-server
```

The application connects to Redis at `localhost:6379` by default.

---

## Build

```bash
go build -o url-shortner main.go
```

---

## Run

```bash
export URL_SHORTENER_API_KEY="your-secret-key"

# Using go run
go run main.go

# Or the compiled binary
./url-shortner
```

The server listens on `:8080`.

---

## Authentication

All requests to `POST /shorten` must include a valid API key:

| Header | Value |
|--------|-------|
| `X-API-Key` | value of `URL_SHORTENER_API_KEY` env var |

Missing or invalid keys receive `HTTP 401 Unauthorized`.

---

## API Usage

### POST /shorten

Shorten a URL. Requires the `X-API-Key` header.

**Request body**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `url` | string | ✅ | The URL to shorten (must start with `http`) |
| `ttl_seconds` | integer | ❌ | Custom TTL for this entry in seconds. `0` or omitted uses the server default (24 h). Must be non-negative. |

**Example — default TTL**
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-key" \
  -d '{"url": "https://example.com"}'
# Response: {"code":"abc123"}
```

**Example — custom TTL (1 hour)**
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-key" \
  -d '{"url": "https://example.com", "ttl_seconds": 3600}'
# Response: {"code":"xyz789"}
```

**Responses**

| Status | Meaning |
|--------|---------|
| `200 OK` | `{"code": "<short-code>"}` |
| `400 Bad Request` | Invalid URL or negative `ttl_seconds` |
| `401 Unauthorized` | Missing or invalid API key |
| `405 Method Not Allowed` | Non-POST request |

---

### GET /{code}

Redirect to the original URL.

```bash
curl -v http://localhost:8080/abc123
# HTTP/1.1 302 Found
# Location: https://example.com
```

| Status | Meaning |
|--------|---------|
| `302 Found` | Redirect to original URL |
| `404 Not Found` | Code does not exist or has expired |

---

## TTL Configuration

URLs stored in Redis expire automatically. TTL can be controlled at two levels:

| Level | How to set | Default |
|-------|-----------|---------|
| **Global default** | Pass `ttlSeconds` to `NewRedisStorage(addr, ttlSeconds)` in `main.go` | 86400 s (24 h) |
| **Per request** | Include `"ttl_seconds": <n>` in the `POST /shorten` body | Falls back to global default when `0` or omitted |

---

## Testing

Handler and storage unit tests use a mock storage and do not require a running Redis instance:

```bash
go test ./tests/... -v
```

Integration tests (`TestShortenHandler`, `TestRedirectHandler`, `TestRedisStorage`) require Redis at `localhost:6379`.

---

## Extending

- **New storage backend** — implement the `Storage` interface in `storage/storage.go`:
  ```go
  type Storage interface {
      Shorten(url string, ttlSeconds int64) string
      Resolve(code string) (string, bool)
  }
  ```
- **Custom global TTL** — pass a second argument to `NewRedisStorage`:
  ```go
  storage.NewRedisStorage("localhost:6379", 7200) // 2-hour default
  ```
- **Additional middleware** — add rate limiting, JWT auth, or request logging in the `middleware` package.

---

## Redis CLI Reference

Inspect and manage stored URLs directly:

```bash
redis-cli KEYS '*'          # List all short codes
redis-cli GET <code>        # Get original URL for a code
redis-cli SET <code> <url>  # Manually add/update a mapping
redis-cli DEL <code>        # Delete a mapping
redis-cli TTL <code>        # Check remaining TTL (seconds)
```

---

## License

MIT
