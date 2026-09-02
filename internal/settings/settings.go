package settings

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Settings struct {
	VaultPath     string `json:"vault_path"`
	Interval      int    `json:"interval"`
	AutoSync      bool   `json:"auto_sync"`
	DailySync     bool   `json:"daily_sync"`
	SyncOnStartup bool   `json:"sync_on_startup"`
	Autostart     bool   `json:"autostart"`
	StartHidden   bool   `json:"start_hidden"`
}

func DefaultSettings() *Settings {
	return &Settings{
		VaultPath:     "",
		Interval:      30,
		AutoSync:      false,
		DailySync:     false,
		SyncOnStartup: true,
		Autostart:     false,
		StartHidden:   false,
	}
}

func configDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "Obsync"), nil
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "settings.json"), nil
}

func Load() (*Settings, error) {
	path, err := configPath()
	if err != nil {
		return DefaultSettings(), nil
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return DefaultSettings(), nil
	}

	if err != nil {
		return nil, err
	}

	s := DefaultSettings()

	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}

	if s.Interval <= 0 {
		s.Interval = 30
	}

	return s, nil
}

func Save(s *Settings) error {
	if s == nil {
		return errors.New("settings are nil")
	}

	dir, err := configDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	if s.Interval <= 0 {
		s.Interval = 30
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(dir, "settings.json")

	return os.WriteFile(path, data, 0o600)
}
