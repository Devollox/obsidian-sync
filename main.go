package main

import (
	"context"
	"embed"
	"log"

	"local/obsync/internal/settings"
	"local/obsync/internal/tray"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v", r)
		}
	}()

	s, err := settings.Load()
	if err != nil {
		log.Printf("settings load error: %v", err)
	}

	app := NewApp()

	err = wails.Run(&options.App{
		Title:          "Obsync",
		Width:          420,
		Height:         360,
		Frameless:      true,
		DisableResize:  true,
		StartHidden:    s.StartHidden,
		HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{
			R: 10,
			G: 10,
			B: 12,
			A: 1,
		},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)

			go tray.Run(tray.Handlers{
				OnShow: func() {
					runtime.Show(ctx)
				},
				OnQuit: func() {
					runtime.Quit(ctx)
				},
			})
		},
		OnBeforeClose: func(ctx context.Context) bool {
			runtime.Hide(ctx)
			return false
		},
		Bind: []interface{}{
			app,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "e0c6c34f-e546-4e1b-9aa0-7c8e5839ce13",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				runtime.Show(app.ctx)
				runtime.WindowUnminimise(app.ctx)
			},
		},
	})

	if err != nil {
		log.Printf("wails error: %v", err)
	}
}
