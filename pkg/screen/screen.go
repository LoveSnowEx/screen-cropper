package screen

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"sync"

	"github.com/kbinani/screenshot"
)

const defaultBufferSize = 1 << 24 // 16MB

type Screen interface {
	Capture() error
	Reader() io.Reader
}

type screen struct {
	mu    sync.Mutex
	bound image.Rectangle
	buf   []byte
}

// NewScreen returns a new screen
func NewScreen(bound image.Rectangle) Screen {
	return &screen{
		bound: bound,
		buf:   make([]byte, 0, defaultBufferSize),
	}
}

// Capture captures the screen
func (s *screen) Capture() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	img, err := screenshot.CaptureRect(s.bound)
	if err != nil {
		return err
	}
	s.buf = s.buf[:0]
	buf := bytes.NewBuffer(s.buf)
	if err := png.Encode(buf, img); err != nil {
		return err
	}
	s.buf = buf.Bytes()
	return nil
}

// Reader returns a reader for the cached image
func (s *screen) Reader() io.Reader {
	s.mu.Lock()
	defer s.mu.Unlock()
	return bytes.NewReader(s.buf)
}

// MaxBound returns the maximum bounds of all displays
func MaxBound() image.Rectangle {
	minX, minY, maxX, maxY := 0, 0, 0, 0
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		if bounds.Min.X < minX {
			minX = bounds.Min.X
		}
		if bounds.Min.Y < minY {
			minY = bounds.Min.Y
		}
		if bounds.Max.X > maxX {
			maxX = bounds.Max.X
		}
		if bounds.Max.Y > maxY {
			maxY = bounds.Max.Y
		}
	}
	return image.Rect(minX, minY, maxX, maxY)
}
