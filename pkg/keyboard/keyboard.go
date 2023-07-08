package keyboard

import (
	"strings"
	"sync"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/types"
	"golang.org/x/exp/slices"
)

const (
	defaultEventChanSize = 256
	defaultListenerSize  = 256
	defaultPressedSize   = 256
)

var (
	// eventChan is a channel to receive keyboard events.
	eventChan chan types.KeyboardEvent
	// listeners is a map of hotkey listeners.
	listeners map[string]*hotkeyListener
	// pressed is a map of pressed keys.
	pressed [defaultPressedSize]bool
)

// init initializes variables.
func init() {
	reset()
}

// reset resets all variables.
func reset() {
	eventChan = make(chan types.KeyboardEvent, defaultEventChanSize)
	listeners = make(map[string]*hotkeyListener, defaultListenerSize)
	pressed = [defaultPressedSize]bool{}
}

// Hotkey represents a hotkey.
type Hotkey []types.VKCode

// hotkeyListener listens the hotkey.
type hotkeyListener struct {
	// mu is a mutex to protect the hotkey.
	mu sync.Mutex
	// justPressed is true if the hotkey is just pressed.
	justPressed bool
	// hotkey is a hotkey to listen.
	hotkey Hotkey
	// callback is a callback function which is called when the hotkey is pressed.
	callback func()
}

// NewHotkey creates a new hotkey.
func NewHotkey(keys ...types.VKCode) Hotkey {
	slices.Sort(keys)
	return Hotkey(keys)
}

// String returns a string representation of the hotkey.
func (h Hotkey) String() string {
	b := strings.Builder{}
	for i, key := range h {
		if i > 0 {
			b.WriteString("+")
		}
		b.WriteString(key.String())
	}
	return b.String()
}

// Pressed returns true if the hotkey is pressed.
func (h Hotkey) Pressed() bool {
	for _, key := range h {
		if !pressed[key] {
			return false
		}
	}
	return true
}

// newHotkeyListener creates a new hotkey listener.
func newHotkeyListener(hotkey Hotkey, callback func()) *hotkeyListener {
	if callback == nil {
		callback = func() {}
	}
	return &hotkeyListener{
		hotkey:   hotkey,
		callback: callback,
	}
}

// Notify calls the callback function if the hotkey is pressed.
func (l *hotkeyListener) Notify() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.hotkey.Pressed() {
		l.justPressed = false
		return
	}
	if l.justPressed {
		return
	}
	l.justPressed = true
	l.callback()
}

// RegisterHotkey registers the hotkey.
func RegisterHotkey(hotkey Hotkey, callback func()) {
	listeners[hotkey.String()] = newHotkeyListener(hotkey, callback)
}

// UnregisterHotkey unregisters the hotkey.
func UnregisterHotkey(hotkey Hotkey) {
	delete(listeners, hotkey.String())
}

// handleEvent handles the keyboard event.
func handleEvent(event types.KeyboardEvent) {
	switch event.Message {
	case types.WM_KEYDOWN, types.WM_SYSKEYDOWN:
		pressed[event.VKCode] = true
		for _, listener := range listeners {
			listener.Notify()
		}
	case types.WM_KEYUP, types.WM_SYSKEYUP:
		pressed[event.VKCode] = false
	}
}

// Start starts the keyboard hook.
func Start() error {
	go func() {
		for event := range eventChan {
			handleEvent(event)
		}
	}()
	return keyboard.Install(nil, eventChan)
}

// Stop stops the keyboard hook.
func Stop() error {
	err := keyboard.Uninstall()
	if err != nil {
		return err
	}
	close(eventChan)
	reset()
	return nil
}
