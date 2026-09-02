import type { AppSettings } from '../hooks/use-sync'

type Props = {
	settings: AppSettings
	saving: boolean
	onChange: (s: AppSettings) => void
	onPickFolder: () => void
}

export function SettingsPanel({ settings, saving, onChange }: Props) {
	const set = <K extends keyof AppSettings>(key: K, value: AppSettings[K]) =>
		onChange({ ...settings, [key]: value })

	return (
		<div style={{ padding: '20px', display: 'flex', flexDirection: 'column', gap: '4px' }}>
			<div
				style={{
					fontSize: '0.6rem',
					letterSpacing: '0.1em',
					color: 'rgba(233,230,221,0.25)',
					marginBottom: '12px',
				}}
			>
				SETTINGS {saving && <span style={{ color: 'rgba(233,230,221,0.20)' }}>· saving…</span>}
			</div>

			<ToggleRow
				label='Auto sync'
				desc='Automatically sync on interval'
				value={settings.auto_sync}
				onChange={v => set('auto_sync', v)}
			/>
			<ToggleRow
				label='Sync on startup'
				desc='Pull changes when the app launches'
				value={settings.sync_on_startup}
				onChange={v => set('sync_on_startup', v)}
			/>
			{settings.auto_sync && (
				<>
					<ToggleRow
						label='Daily mode'
						desc='Only one sync per day — skips if already synced today'
						value={settings.daily_sync}
						onChange={v => set('daily_sync', v)}
					/>

					{!settings.daily_sync && (
						<div style={{ padding: '8px 0 4px' }}>
							<div
								style={{
									fontSize: '0.6rem',
									letterSpacing: '0.08em',
									color: 'rgba(233,230,221,0.25)',
									marginBottom: '8px',
								}}
							>
								INTERVAL
							</div>
							<div style={{ display: 'flex', gap: '6px', flexWrap: 'wrap' }}>
								{[5, 15, 30, 60].map(v => (
									<button
										key={v}
										onClick={() => set('interval', v)}
										style={{
											padding: '5px 12px',
											borderRadius: '6px',
											border: `1px solid ${settings.interval === v ? 'rgba(233,230,221,0.30)' : 'rgba(233,230,221,0.08)'}`,
											background:
												settings.interval === v ? 'rgba(233,230,221,0.10)' : 'transparent',
											color: settings.interval === v ? '#e9e6dd' : 'rgba(233,230,221,0.35)',
											fontSize: '0.72rem',
											cursor: 'pointer',
											fontFamily: 'inherit',
											transition: 'all 120ms',
										}}
									>
										{v}m
									</button>
								))}
							</div>
						</div>
					)}
				</>
			)}

			<ToggleRow
				label='Autostart'
				desc='Launch on system startup'
				value={settings.autostart}
				onChange={v => set('autostart', v)}
			/>

			<ToggleRow
				label='Start hidden'
				desc='Hide to tray on launch'
				value={settings.start_hidden}
				onChange={v => set('start_hidden', v)}
			/>
		</div>
	)
}

function ToggleRow({
	label,
	desc,
	value,
	onChange,
}: {
	label: string
	desc: string
	value: boolean
	onChange: (v: boolean) => void
}) {
	return (
		<div
			style={{
				display: 'flex',
				alignItems: 'center',
				justifyContent: 'space-between',
				padding: '10px 0',
				borderBottom: '1px solid rgba(233,230,221,0.07)',
			}}
		>
			<div style={{ display: 'flex', flexDirection: 'column', gap: '2px' }}>
				<span style={{ fontSize: '0.8rem', color: '#e9e6dd' }}>{label}</span>
				<span style={{ fontSize: '0.62rem', color: 'rgba(233,230,221,0.30)' }}>{desc}</span>
			</div>
			<button
				onClick={() => onChange(!value)}
				style={{
					width: '36px',
					height: '20px',
					borderRadius: '10px',
					border: 'none',
					background: value ? 'rgba(74,222,128,0.70)' : 'rgba(233,230,221,0.12)',
					cursor: 'pointer',
					position: 'relative',
					transition: 'background 200ms',
					flexShrink: 0,
				}}
			>
				<span
					style={{
						position: 'absolute',
						top: '2px',
						left: value ? '18px' : '2px',
						width: '16px',
						height: '16px',
						borderRadius: '50%',
						background: value ? '#fff' : 'rgba(233,230,221,0.50)',
						transition: 'left 200ms',
					}}
				/>
			</button>
		</div>
	)
}
