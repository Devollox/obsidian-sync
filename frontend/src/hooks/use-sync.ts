import { useEffect, useRef, useState } from 'react'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export type SyncResult = {
	status: 'idle' | 'syncing' | 'done' | 'skipped' | 'error'
	message: string
	timestamp: string
}

export type AppSettings = {
	vault_path: string
	interval: number
	auto_sync: boolean
	daily_sync: boolean
	sync_on_startup: boolean
	autostart: boolean
	start_hidden: boolean
}

declare global {
	interface Window {
		go: {
			main: {
				App: {
					Sync: () => Promise<SyncResult>
					GetLastSyncResult: () => Promise<SyncResult | null>
					GetSettings: () => Promise<AppSettings>
					SaveSettings: (s: AppSettings) => Promise<void>
					PickFolder: () => Promise<string>
					HideWindow: () => void
					MinimizeWindow: () => void
					QuitApp: () => void
				}
			}
		}
	}
}

export function useSync() {
	const [result, setResult] = useState<SyncResult | null>(null)
	const [syncing, setSyncing] = useState(false)

	useEffect(() => {
		const off = EventsOn('sync:result', (data: SyncResult) => {
			setSyncing(false)
			setResult(data)
		})

		window.go.main.App.GetLastSyncResult()
			.then(lastResult => {
				if (lastResult) {
					setResult(lastResult)
				}
			})
			.catch(() => {})

		return () => off()
	}, [])

	const sync = async () => {
		setSyncing(true)

		try {
			const syncResult = await window.go.main.App.Sync()
			setResult(syncResult)
		} catch (error) {
			setResult({
				status: 'error',
				message: String(error),
				timestamp: '',
			})
		} finally {
			setSyncing(false)
		}
	}

	return {
		result,
		syncing,
		sync,
	}
}

export function useSettings() {
	const [settings, setSettings] = useState<AppSettings | null>(null)
	const [saving, setSaving] = useState(false)
	const saveTimer = useRef<ReturnType<typeof setTimeout> | null>(null)

	useEffect(() => {
		window.go.main.App.GetSettings().then(setSettings)
	}, [])

	useEffect(() => {
		return () => {
			if (saveTimer.current) {
				clearTimeout(saveTimer.current)
			}
		}
	}, [])

	const save = (next: AppSettings) => {
		setSettings(next)

		if (saveTimer.current) {
			clearTimeout(saveTimer.current)
		}

		saveTimer.current = setTimeout(async () => {
			setSaving(true)

			try {
				await window.go.main.App.SaveSettings(next)
			} finally {
				setSaving(false)
			}
		}, 400)
	}

	return {
		settings,
		saving,
		save,
	}
}
