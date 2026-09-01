//go:build !darwin

package tray

import (
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed icon.ico
var icon []byte

type Handlers struct {
	OnShow func()
	OnQuit func()
}

func Run(h Handlers) {
	systray.Run(func() {
		systray.SetIcon(icon)
		systray.SetTitle("Obsync")
		systray.SetTooltip("Obsync")

		mShow := systray.AddMenuItem("Show", "Show the window")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the application")

		go func() {
			for {
				select {
				case <-mShow.ClickedCh:
					if h.OnShow != nil {
						h.OnShow()
					}
				case <-mQuit.ClickedCh:
					systray.Quit()
					if h.OnQuit != nil {
						h.OnQuit()
					}
				}
			}
		}()
	}, func() {})
}

func Stop() {
	systray.Quit()
}
