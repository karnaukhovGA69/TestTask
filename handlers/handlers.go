package handlers

import (
	"encoding/json"
	"main/internal"
	"net/http"
	"strings"
)

type Handler struct {
	service *internal.Service
}

func NewHandler(service *internal.Service) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) GetHandler(rw http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimSpace(r.PathValue("shortURL"))
	if shortURL == "" {
		http.Error(rw, "пустой URL", http.StatusBadRequest)
		return
	}

	originalURL, err := h.service.GetLongURL(shortURL)
	if err != nil {
		http.Error(rw, "Не найден", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{"URL": originalURL})
}

type URL struct {
	Url string `json:"url"`
}

func (h *Handler) PostHandler(rw http.ResponseWriter, r *http.Request) {
	var url URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(rw, "не получилось распарсить URL", http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.CreateShortURL(url.Url)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	fullShortURL := "http://" + r.Host + "/" + shortURL
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]string{"shortURL": fullShortURL})

}
