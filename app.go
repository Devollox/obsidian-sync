package main

import (
	"context"
	"sync"
	"time"

	"local/obsync/internal/settings"
	gosync "local/obsync/internal/sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx            context.Context
	syncer         *gosync.Syncer
	ticker         *time.Ticker
	stop           chan struct{}
	lastSyncResult *gosync.SyncResult
	resultMu       sync.RWMutex
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

	s, err := settings.Load()
	if err != nil {
		return
	}

	if s.SyncOnStartup && s.VaultPath != "" {
		go a.runSync(func() (*gosync.SyncResult, error) {
			return a.syncer.SyncOnStartup(
				s.VaultPath,
				s.DailySync,
			)
		})
	}

	if !s.DailySync && s.AutoSync && s.VaultPath != "" {
		a.startAutoSync(s)
	}
}

func (a *App) setLastSyncResult(result *gosync.SyncResult) {
	if result == nil {
		return
	}

	a.resultMu.Lock()
	a.lastSyncResult = result
	a.resultMu.Unlock()
}

func (a *App) GetLastSyncResult() *gosync.SyncResult {
	a.resultMu.RLock()
	defer a.resultMu.RUnlock()

	if a.lastSyncResult == nil {
		return nil
	}

	result := *a.lastSyncResult
	return &result
}

func (a *App) emitSyncResult(result *gosync.SyncResult) {
	if result == nil {
		return
	}

	a.setLastSyncResult(result)

	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "sync:result", result)
	}
}

func (a *App) runSync(fn func() (*gosync.SyncResult, error)) {
	result, err := fn()

	if err != nil {
		a.emitSyncResult(&gosync.SyncResult{
			Status:    gosync.StatusError,
			Message:   err.Error(),
			Timestamp: time.Now().Format("15:04:05"),
		})
		return
	}

	a.emitSyncResult(result)
}

func (a *App) stopAutoSync() {
	if a.ticker == nil {
		return
	}

	a.ticker.Stop()
	a.ticker = nil

	select {
	case a.stop <- struct{}{}:
	default:
	}
}

func (a *App) startAutoSync(s *settings.Settings) {
	a.stopAutoSync()

	if s.DailySync || !s.AutoSync || s.VaultPath == "" {
		return
	}

	if s.Interval <= 0 {
		return
	}

	a.ticker = time.NewTicker(
		time.Duration(s.Interval) * time.Minute,
	)

	go func() {
		for {
			select {
			case <-a.ticker.C:
				cfg, err := settings.Load()
				if err != nil {
					continue
				}

				if cfg.DailySync || !cfg.AutoSync || cfg.VaultPath == "" {
					continue
				}

				a.runSync(func() (*gosync.SyncResult, error) {
					return a.syncer.Sync(cfg.VaultPath, false)
				})

			case <-a.stop:
				return
			}
		}
	}()
}

func (a *App) Sync() (*gosync.SyncResult, error) {
	s, err := settings.Load()
	if err != nil {
		return nil, err
	}

	if s.VaultPath == "" {
		result := &gosync.SyncResult{
			Status:  gosync.StatusError,
			Message: "Vault path is not set",
		}

		a.emitSyncResult(result)
		return result, nil
	}

	result, err := a.syncer.Sync(s.VaultPath, false)
	if err != nil {
		a.emitSyncResult(&gosync.SyncResult{
			Status:    gosync.StatusError,
			Message:   err.Error(),
			Timestamp: time.Now().Format("15:04:05"),
		})
		return nil, err
	}

	a.emitSyncResult(result)
	return result, nil
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

	a.stopAutoSync()

	if !s.DailySync && s.AutoSync && s.VaultPath != "" {
		a.startAutoSync(s)
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
