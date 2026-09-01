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
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func alreadySyncedToday(vaultPath string) bool {
	out, err := git(vaultPath, "log", "-1", "--format=%s")
	if err != nil {
		return false
	}
	today := time.Now().Format("2006-01-02")
	return strings.Contains(out, "sync: "+today)
}

func (s *Syncer) Sync(vaultPath string) (*SyncResult, error) {
	ts := time.Now().Format("15:04:05")
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

func (s *Syncer) SyncOnStartup(vaultPath string) (*SyncResult, error) {
	ts := time.Now().Format("15:04:05")
	return s.sync(vaultPath, ts, true)
}

func (s *Syncer) sync(vaultPath, ts string, pullOnly bool) (*SyncResult, error) {
	stashed := false

	out, _ := git(vaultPath, "status", "--porcelain")
	hasChanges := out != ""

	if hasChanges {
		if _, err := git(vaultPath, "stash"); err == nil {
			stashed = true
		}
	}

	if _, err := git(vaultPath, "pull", "--rebase"); err != nil {
		if stashed {
			git(vaultPath, "stash", "pop")
		}
		out2, _ := git(vaultPath, "pull", "--rebase")
		if !strings.Contains(out2, "Already up to date") {
			return &SyncResult{
				Status:    StatusError,
				Message:   "pull: " + out2,
				Timestamp: ts,
			}, nil
		}
	}

	if stashed {
		git(vaultPath, "stash", "pop")
	}

	if pullOnly {
		return &SyncResult{
			Status:    StatusDone,
			Message:   "Pulled latest changes",
			Timestamp: ts,
		}, nil
	}

	if _, err := git(vaultPath, "add", "."); err != nil {
		return &SyncResult{Status: StatusError, Message: "add failed", Timestamp: ts}, nil
	}

	commitMsg := fmt.Sprintf("sync: %s", time.Now().Format("2006-01-02 15:04:05"))
	out3, err := git(vaultPath, "commit", "-m", commitMsg)
	if err != nil {
		if strings.Contains(out3, "nothing to commit") {
			return &SyncResult{
				Status:    StatusDone,
				Message:   "Nothing to commit",
				Timestamp: ts,
			}, nil
		}
		return &SyncResult{Status: StatusError, Message: "commit: " + out3, Timestamp: ts}, nil
	}

	if out4, err := git(vaultPath, "push"); err != nil {
		return &SyncResult{Status: StatusError, Message: "push: " + out4, Timestamp: ts}, nil
	}

	return &SyncResult{
		Status:    StatusDone,
		Message:   "Vault synced successfully",
		Timestamp: ts,
	}, nil
}
