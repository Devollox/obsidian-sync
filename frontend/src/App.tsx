import { useState } from 'react'
import { SettingsPanel } from './components/SettingsPanel'
import { SyncPanel } from './components/SyncPanel'
import { TitleBar } from './components/TitleBar'
import { useSettings, useSync } from './hooks/use-sync'
import './style.css'

type Tab = 'sync' | 'settings'

export default function App() {
	const [tab, setTab] = useState<Tab>('sync')
	const { result, syncing, sync } = useSync()
	const { settings, saving, save } = useSettings()

	const pickFolder = async () => {
		const path = await window.go.main.App.PickFolder()
		if (path && settings) save({ ...settings, vault_path: path })
	}

	return (
		<div
			style={{
				width: '100%',
				height: '100vh',
				background: '#0a0a0c',
				display: 'flex',
				flexDirection: 'column',
				border: '1px solid rgba(233,230,221,0.08)',
				borderRadius: '12px',
				overflow: 'hidden',
				position: 'relative',
			}}
		>
			<TitleBar
				onMinimize={() => window.go.main.App.MinimizeWindow()}
				onHide={() => window.go.main.App.HideWindow()}
				onQuit={() => window.go.main.App.QuitApp()}
			/>

			<div style={{ height: '1px', background: 'rgba(233,230,221,0.06)', flexShrink: 0 }} />

			<div style={{ display: 'flex', padding: '0 20px', gap: '0', flexShrink: 0 }}>
				{(['sync', 'settings'] as Tab[]).map(t => (
					<button
						key={t}
						onClick={() => setTab(t)}
						style={{
							background: 'none',
							border: 'none',
							borderBottom: `1px solid ${tab === t ? 'rgba(233,230,221,0.50)' : 'transparent'}`,
							padding: '10px 0',
							marginRight: '20px',
							cursor: 'pointer',
							fontFamily: 'inherit',
							fontSize: '0.65rem',
							letterSpacing: '0.08em',
							color: tab === t ? 'rgba(233,230,221,0.80)' : 'rgba(233,230,221,0.25)',
							transition: 'color 150ms, border-color 150ms',
						}}
					>
						{t.toUpperCase()}
					</button>
				))}
			</div>

			<div style={{ height: '1px', background: 'rgba(233,230,221,0.06)', flexShrink: 0 }} />

			<div className='scroll-y' style={{ flex: 1 }}>
				{tab === 'sync' && settings && (
					<SyncPanel
						vaultPath={settings.vault_path}
						syncing={syncing}
						result={result}
						onSync={sync}
						onPickFolder={pickFolder}
					/>
				)}
				{tab === 'settings' && settings && (
					<SettingsPanel
						settings={settings}
						saving={saving}
						onChange={save}
						onPickFolder={pickFolder}
					/>
				)}
				{!settings && (
					<div
						style={{
							flex: 1,
							display: 'flex',
							alignItems: 'center',
							justifyContent: 'center',
							fontSize: '0.7rem',
							color: 'rgba(233,230,221,0.20)',
							padding: '40px',
						}}
					>
						Loading…
					</div>
				)}
			</div>
		</div>
	)
}
