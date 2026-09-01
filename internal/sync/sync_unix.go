//go:build !windows

package sync

import "os/exec"

func hideWindow(cmd *exec.Cmd) {
}
