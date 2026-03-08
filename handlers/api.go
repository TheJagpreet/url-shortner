package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"url-shortner/storage"
)

type API struct {
	Storage storage.Storage
}

func NewAPI(storage storage.Storage) *API {
	return &API{Storage: storage}
}

func (a *API) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL        string `json:"url"`
		TTLSeconds int64  `json:"ttl_seconds"`
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !strings.HasPrefix(req.URL, "http") {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.TTLSeconds < 0 {
		http.Error(w, "ttl_seconds must be non-negative", http.StatusBadRequest)
		return
	}
	code := a.Storage.Shorten(req.URL, req.TTLSeconds)
	resp := map[string]string{"code": code}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (a *API) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	url, ok := a.Storage.Resolve(code)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
