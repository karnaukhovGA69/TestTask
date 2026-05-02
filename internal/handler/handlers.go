package handler

import (
	"encoding/json"
	"errors"
	"main/internal/apperrors"
	"main/internal/service"
	"net/http"
	"strings"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) GetHandler(rw http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimSpace(r.PathValue("shortURL"))
	if shortURL == "" {
		http.Error(rw, apperrors.ErrEmptyURL.Error(), http.StatusBadRequest)
		return
	}

	originalURL, err := h.service.GetLongURL(shortURL)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(rw, "Не найден", http.StatusNotFound)
			return
		}
		if errors.Is(err, apperrors.ErrEmptyURL) {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]string{"URL": originalURL})
}

type URL struct {
	Url string `json:"url"`
}

func (h *Handler) PostHandler(rw http.ResponseWriter, r *http.Request) {
	var url URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(rw, apperrors.ErrBadJSON.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.CreateShortURL(url.Url)
	if err != nil {
		if errors.Is(err, apperrors.ErrEmptyURL) {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fullShortURL := "http://" + r.Host + "/" + shortURL
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]string{"shortURL": fullShortURL})

}
