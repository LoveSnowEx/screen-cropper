package main

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/LoveSnowEx/screen-cropper/pkg/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed assets/icon.ico
var icon []byte

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.bindSystray()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
