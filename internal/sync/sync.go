package sync

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Status string

const (
	StatusIdle    Status = "idle"
	StatusSyncing Status = "syncing"
	StatusDone    Status = "done"
	StatusSkipped Status = "skipped"
	StatusError   Status = "error"
)

type SyncResult struct {
	Status    Status `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type Syncer struct {
	ctx context.Context
}

func NewSyncer() *Syncer {
	return &Syncer{}
}

func (s *Syncer) Startup(ctx context.Context) {
	s.ctx = ctx
}

func git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	hideWindow(cmd)

	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func alreadySyncedToday(vaultPath string) bool {
	out, err := git(vaultPath, "log", "-1", "--format=%s")
	if err != nil {
		return false
	}

	today := time.Now().Format("2006-01-02")
	return strings.HasPrefix(out, "sync: "+today)
}

func (s *Syncer) Sync(vaultPath string, dailyMode bool) (*SyncResult, error) {
	ts := time.Now().Format("15:04:05")

	if dailyMode && alreadySyncedToday(vaultPath) {
		return &SyncResult{
			Status:    StatusSkipped,
			Message:   "Already synced today",
			Timestamp: ts,
		}, nil
	}

	return s.sync(vaultPath, ts, false)
}

func (s *Syncer) SyncDaily(vaultPath string) (*SyncResult, error) {
	ts := time.Now().Format("15:04:05")

	if alreadySyncedToday(vaultPath) {
		return &SyncResult{
			Status:    StatusSkipped,
			Message:   "Already synced today",
			Timestamp: ts,
		}, nil
	}

	return s.sync(vaultPath, ts, false)
}

func (s *Syncer) SyncOnStartup(vaultPath string, dailyMode bool) (*SyncResult, error) {
	ts := time.Now().Format("15:04:05")

	if dailyMode {
		if alreadySyncedToday(vaultPath) {
			return &SyncResult{
				Status:    StatusSkipped,
				Message:   "Already synced today",
				Timestamp: ts,
			}, nil
		}

		return s.sync(vaultPath, ts, false)
	}

	return s.sync(vaultPath, ts, true)
}

func (s *Syncer) sync(vaultPath, ts string, pullOnly bool) (*SyncResult, error) {
	stashed := false

	out, err := git(vaultPath, "status", "--porcelain")
	if err != nil {
		return &SyncResult{
			Status:    StatusError,
			Message:   "status: " + out,
			Timestamp: ts,
		}, nil
	}

	if out != "" {
		out, err = git(vaultPath, "stash", "push", "-u", "-m", "obsync-autostash")
		if err != nil {
			return &SyncResult{
				Status:    StatusError,
				Message:   "stash: " + out,
				Timestamp: ts,
			}, nil
		}

		stashed = true
	}

	out, err = git(vaultPath, "pull", "--rebase")
	if err != nil {
		if stashed {
			popOut, popErr := git(vaultPath, "stash", "pop")
			if popErr != nil {
				return &SyncResult{
					Status:    StatusError,
					Message:   "pull: " + out + "; stash pop: " + popOut,
					Timestamp: ts,
				}, nil
			}
		}

		return &SyncResult{
			Status:    StatusError,
			Message:   "pull: " + out,
			Timestamp: ts,
		}, nil
	}

	if stashed {
		out, err = git(vaultPath, "stash", "pop")
		if err != nil {
			return &SyncResult{
				Status:    StatusError,
				Message:   "stash pop: " + out,
				Timestamp: ts,
			}, nil
		}
	}

	if pullOnly {
		return &SyncResult{
			Status:    StatusDone,
			Message:   "Pulled latest changes",
			Timestamp: ts,
		}, nil
	}

	out, err = git(vaultPath, "add", ".")
	if err != nil {
		return &SyncResult{
			Status:    StatusError,
			Message:   "add: " + out,
			Timestamp: ts,
		}, nil
	}

	commitMsg := fmt.Sprintf(
		"sync: %s",
		time.Now().Format("2006-01-02 15:04:05"),
	)

	out, err = git(vaultPath, "commit", "-m", commitMsg)
	if err != nil {
		if strings.Contains(out, "nothing to commit") {
			return &SyncResult{
				Status:    StatusDone,
				Message:   "Nothing to commit",
				Timestamp: ts,
			}, nil
		}

		return &SyncResult{
			Status:    StatusError,
			Message:   "commit: " + out,
			Timestamp: ts,
		}, nil
	}

	out, err = git(vaultPath, "push")
	if err != nil {
		return &SyncResult{
			Status:    StatusError,
			Message:   "push: " + out,
			Timestamp: ts,
		}, nil
	}

	return &SyncResult{
		Status:    StatusDone,
		Message:   "Vault synced successfully",
		Timestamp: ts,
	}, nil
}
