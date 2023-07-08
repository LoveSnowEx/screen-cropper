package main

import (
	"context"
	"fmt"
	"log"

	_ "embed"

	"github.com/LoveSnowEx/screen-cropper/pkg/keyboard"
	"github.com/LoveSnowEx/screen-cropper/pkg/screen"
	"github.com/LoveSnowEx/screen-cropper/pkg/systray"
	"github.com/moutend/go-hook/pkg/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed assets/icon.ico
var icon []byte

// App struct
type App struct {
	ctx    context.Context
	screen screen.Screen
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		screen: screen.NewScreen(screen.MaxBound()),
	}
}

// bindSystray binds the systray to the app
func (a *App) bindSystray() {
	// Run the systray
	go systray.Run(a.ctx, systray.Config{
		Title:     "Screen Cropper",
		Tooltip:   "Screen Cropper",
		IconBytes: icon,
	})

	// Quit the app when the systray is closed
	go func() {
		<-systray.Done()
		runtime.Quit(a.ctx)
	}()
}

// bindKeyboard binds the keyboard to the app
func (a *App) bindKeyboard() {
	// Register the hotkey
	keyboard.RegisterHotkey(keyboard.NewHotkey(types.VK_LCONTROL, types.VK_LMENU, types.VK_A), func() {
		runtime.EventsEmit(a.ctx, "capture")
	})
	// Start the keyboard
	if err := keyboard.Start(); err != nil {
		log.Fatal(err)
	}
	// Stop the keyboard when the app closes
	go func() {
		<-a.ctx.Done()
		if err := keyboard.Stop(); err != nil {
			log.Fatal(err)
		}
	}()
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.bindSystray()
	a.bindKeyboard()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	go a.screen.Capture()
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
