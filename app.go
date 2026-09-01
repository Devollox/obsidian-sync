package main

import (
	"context"
	"time"

	gosync "local/obsync/internal/sync"
	"local/obsync/internal/settings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx    context.Context
	syncer *gosync.Syncer
	ticker *time.Ticker
	stop   chan struct{}
}

func NewApp() *App {
	return &App{
		syncer: gosync.NewSyncer(),
		stop:   make(chan struct{}),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.syncer.Startup(ctx)

	s, _ := settings.Load()
	if s.StartHidden {
		runtime.Hide(ctx)
	}
	if s.SyncOnStartup && s.VaultPath != "" {
		go a.syncer.SyncOnStartup(s.VaultPath)
	}
	if s.AutoSync && s.VaultPath != "" {
		a.startAutoSync(s)
	}
}

func (a *App) startAutoSync(s *settings.Settings) {
	if a.ticker != nil {
		a.ticker.Stop()
		select {
		case a.stop <- struct{}{}:
		default:
		}
	}

	a.ticker = time.NewTicker(time.Duration(s.Interval) * time.Minute)
	daily := s.DailySync

	go func() {
		for {
			select {
			case <-a.ticker.C:
				cfg, _ := settings.Load()
				if !cfg.AutoSync || cfg.VaultPath == "" {
					continue
				}
				if daily {
					a.syncer.SyncDaily(cfg.VaultPath)
				} else {
					a.syncer.Sync(cfg.VaultPath)
				}
			case <-a.stop:
				return
			}
		}
	}()
}

func (a *App) Sync() (*gosync.SyncResult, error) {
	s, _ := settings.Load()
	if s.VaultPath == "" {
		return &gosync.SyncResult{
			Status:  gosync.StatusError,
			Message: "Vault path is not set",
		}, nil
	}
	return a.syncer.Sync(s.VaultPath)
}

func (a *App) GetSettings() (*settings.Settings, error) {
	s, err := settings.Load()
	if err != nil {
		return nil, err
	}
	s.Autostart = settings.GetAutostart()
	return s, nil
}

func (a *App) SaveSettings(s *settings.Settings) error {
	if err := settings.SetAutostart(s.Autostart); err != nil {
		s.Autostart = settings.GetAutostart()
	}
	if err := settings.Save(s); err != nil {
		return err
	}
	if s.AutoSync && s.VaultPath != "" {
		a.startAutoSync(s)
	} else if a.ticker != nil {
		a.ticker.Stop()
	}
	return nil
}

func (a *App) PickFolder() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Obsidian Vault",
	})
}

func (a *App) HideWindow() {
	runtime.Hide(a.ctx)
}

func (a *App) MinimizeWindow() {
	runtime.WindowMinimise(a.ctx)
}

func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}
