package shortener

import (
	"math/rand"
	"sync"
	"time"
)

const CodeLength = 6

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type Shortener struct {
	mu   sync.RWMutex
	urls map[string]string // code -> original URL
}

func NewShortener() *Shortener {
	return &Shortener{
		urls: make(map[string]string),
	}
}

func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, CodeLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *Shortener) GenerateCode() string {
	return GenerateCode()
}

func (s *Shortener) Shorten(url string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	code := s.GenerateCode()
	for {
		if _, exists := s.urls[code]; !exists {
			break
		}
		code = s.GenerateCode()
	}
	s.urls[code] = url
	return code
}

func (s *Shortener) Resolve(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.urls[code]
	return url, ok
}
