package settings

import (
	"encoding/json"
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

func configPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exe), "settings.json"), nil
}

func Load() (*Settings, error) {
	path, err := configPath()
	if err != nil {
		return DefaultSettings(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultSettings(), nil
	}

	s := DefaultSettings()
	if err := json.Unmarshal(data, s); err != nil {
		return DefaultSettings(), nil
	}

	return s, nil
}

func Save(s *Settings) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
