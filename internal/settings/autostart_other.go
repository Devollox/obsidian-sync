//go:build !windows

package settings

func SetAutostart(enabled bool) error { return nil }
func GetAutostart() bool              { return false }
