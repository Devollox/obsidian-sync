# Obsync

<div align="center">
  <img width="416" height="355" alt="image" src="https://github.com/user-attachments/assets/dfb56efc-c7d3-4a16-ab83-2fc3e9e74359" />
</div>

##

Lightweight desktop app for automated Obsidian vault synchronization with GitHub, built with [Wails](https://wails.io/) + React + TypeScript.

## Features

- Manual sync with one click
- Auto sync on configurable interval (5 / 15 / 30 / 60 min)
- Daily mode — one sync per day, skips if already synced today
- Sync on startup — pulls latest changes on every launch
- Vault folder picker
- Autostart on system startup
- Start hidden — hide to tray on launch
- System tray with Show / Quit

## Setup

### 1. Install Git

Make sure Git is installed and available in PATH.

```
git --version
```

### 2. Create a GitHub repository

Go to [github.com/new](https://github.com/new) and create a new **private** repository. Do not initialize it with any files.

### 3. Initialize your Obsidian vault as a Git repository

Open a terminal in your vault folder and run:

```bash
git init
git remote add origin https://github.com/YOUR_USERNAME/YOUR_REPO.git
git branch -M main
git add .
git commit -m "init: vault"
git push -u origin main
```

### 4. Set up authentication

For HTTPS (recommended):

```bash
git config --global credential.helper manager
```

This stores your credentials via Git Credential Manager. On first push, a browser window will open to authenticate with GitHub.

For SSH — add your key to GitHub and use `git@github.com:YOUR_USERNAME/YOUR_REPO.git` as the remote.

### 5. Select vault in the app

Open obsync, go to **Sync** tab, click the vault path field and select your vault folder.

Hit **Sync now** to verify everything works.

### 6. Configure auto sync (optional)

Go to **Settings** tab:

- Enable **Sync on startup** to pull changes every time you open the app
- Enable **Auto sync** and pick an interval for background sync
- Enable **Daily mode** to limit syncing to once per day

## How sync works

Each sync run does the following:

1. `git stash` — stashes any unstaged changes
2. `git pull --rebase` — pulls remote changes
3. `git stash pop` — restores stashed changes
4. `git add .` — stages everything
5. `git commit -m "sync: YYYY-MM-DD HH:MM:SS"` — commits with timestamp
6. `git push` — pushes to remote

Daily mode checks the last commit message. If it contains today's date (`sync: 2026-09-01`), the sync is skipped.

## Build

```powershell
wails build
```

## Stack

- [Wails v2](https://wails.io/) — Go + WebView desktop framework
- React 18 + TypeScript
