//go:build windows

package settings

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const regKey = `Software\Microsoft\Windows\CurrentVersion\Run`
const appName = "obsync"

func SetAutostart(enabled bool) error {
	k, err := registry.OpenKey(
		registry.CURRENT_USER,
		regKey,
		registry.SET_VALUE,
	)
	if err != nil {
		return err
	}
	defer k.Close()

	if !enabled {
		err := k.DeleteValue(appName)

		if errors.Is(err, registry.ErrNotExist) {
			return nil
		}

		return err
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	exe, err = filepath.Abs(exe)
	if err != nil {
		return err
	}

	return k.SetStringValue(appName, `"`+exe+`"`)
}

func GetAutostart() bool {
	k, err := registry.OpenKey(
		registry.CURRENT_USER,
		regKey,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return false
	}
	defer k.Close()

	_, _, err = k.GetStringValue(appName)
	return err == nil
}
