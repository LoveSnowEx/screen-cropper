package handler

import (
	"io"
	"net/http"

	"github.com/LoveSnowEx/screen-cropper/pkg/screen"
)

// ScreenHandler is a http.Handler that serves the screen image.
type ScreenHandler struct {
	screen screen.Screen
}

// NewScreenHandler creates a new ScreenHandler.
func NewScreenHandler(screen screen.Screen) http.Handler {
	return &ScreenHandler{screen: screen}
}

// ServeHTTP serves the screen image.
func (s *ScreenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/screen.png" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	if _, err := io.Copy(w, s.screen.Reader()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
