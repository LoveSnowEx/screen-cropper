package systray

import (
	"context"

	"fyne.io/systray"
)

var (
	done = make(chan struct{}, 1)
)

// Config is the configuration for the systray
type Config struct {
	Title     string
	Tooltip   string
	IconBytes []byte
}

// Run the systray and block until systray is closed
func Run(ctx context.Context, cfg Config) {
	onReady := func() {
		systray.SetIcon(cfg.IconBytes)
		systray.SetTitle(cfg.Title)
		systray.SetTooltip(cfg.Tooltip)

		mQuit := systray.AddMenuItem("Quit", "Quit the app")
		go func() {
			<-mQuit.ClickedCh
			systray.Quit()
		}()
	}

	onExit := func() {
		close(done)
	}

	go func() {
		<-ctx.Done()
		systray.Quit()
	}()

	systray.Run(onReady, onExit)
}

// Done returns a channel that closed when the systray is closed
func Done() <-chan struct{} {
	return done
}
