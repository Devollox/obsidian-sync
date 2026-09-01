//go:build darwin

package tray

type Handlers struct {
	OnShow func()
	OnQuit func()
}

func Run(h Handlers) {}
func Stop()          {}
