import type { SyncResult } from '../hooks/use-sync'

type Props = {
	vaultPath: string
	syncing: boolean
	result: SyncResult | null
	onSync: () => void
	onPickFolder: () => void
}

export function SyncPanel({ vaultPath, syncing, result, onSync, onPickFolder }: Props) {
	const isError = result?.status === 'error'
	const isSkipped = result?.status === 'skipped'

	const background = isError
		? 'rgba(248,113,113,0.06)'
		: isSkipped
			? 'rgba(250,204,21,0.06)'
			: 'rgba(74,222,128,0.06)'

	const border = isError
		? 'rgba(248,113,113,0.15)'
		: isSkipped
			? 'rgba(250,204,21,0.15)'
			: 'rgba(74,222,128,0.15)'

	const color = isError ? '#f87171' : isSkipped ? '#facc15' : '#4ade80'

	return (
		<div
			style={{
				padding: '24px 20px',
				display: 'flex',
				flexDirection: 'column',
				gap: '24px',
			}}
		>
			<div
				style={{
					display: 'flex',
					flexDirection: 'column',
					gap: '8px',
				}}
			>
				<span
					style={{
						fontSize: '0.6rem',
						letterSpacing: '0.1em',
						color: 'rgba(233,230,221,0.25)',
					}}
				>
					VAULT PATH
				</span>

				<button
					onClick={onPickFolder}
					style={{
						background: 'rgba(233,230,221,0.04)',
						border: '1px solid rgba(233,230,221,0.08)',
						borderRadius: '8px',
						padding: '10px 12px',
						textAlign: 'left',
						cursor: 'pointer',
						fontFamily: 'inherit',
						transition: 'border-color 120ms',
						display: 'flex',
						alignItems: 'center',
						justifyContent: 'space-between',
						gap: '8px',
					}}
					onMouseEnter={event => {
						event.currentTarget.style.borderColor = 'rgba(233,230,221,0.18)'
					}}
					onMouseLeave={event => {
						event.currentTarget.style.borderColor = 'rgba(233,230,221,0.08)'
					}}
				>
					<span
						style={{
							fontSize: '0.75rem',
							color: vaultPath ? '#e9e6dd' : 'rgba(233,230,221,0.25)',
							overflow: 'hidden',
							textOverflow: 'ellipsis',
							whiteSpace: 'nowrap',
						}}
					>
						{vaultPath || 'Select vault folder…'}
					</span>

					<span
						style={{
							fontSize: '0.65rem',
							color: 'rgba(233,230,221,0.30)',
							flexShrink: 0,
						}}
					>
						↗
					</span>
				</button>
			</div>

			<button
				onClick={onSync}
				disabled={syncing || !vaultPath}
				style={{
					background: syncing ? 'rgba(233,230,221,0.04)' : 'rgba(233,230,221,0.08)',
					border: '1px solid rgba(233,230,221,0.10)',
					borderRadius: '8px',
					padding: '12px',
					cursor: syncing || !vaultPath ? 'not-allowed' : 'pointer',
					fontFamily: 'inherit',
					fontSize: '0.8rem',
					color: syncing || !vaultPath ? 'rgba(233,230,221,0.25)' : '#e9e6dd',
					transition: 'background 120ms, color 120ms',
					letterSpacing: '0.04em',
				}}
				onMouseEnter={event => {
					if (!syncing && vaultPath) {
						event.currentTarget.style.background = 'rgba(233,230,221,0.12)'
					}
				}}
				onMouseLeave={event => {
					event.currentTarget.style.background = syncing
						? 'rgba(233,230,221,0.04)'
						: 'rgba(233,230,221,0.08)'
				}}
			>
				{syncing ? 'Syncing…' : 'Sync now'}
			</button>

			{result && (
				<div
					style={{
						padding: '10px 12px',
						borderRadius: '8px',
						background,
						border: `1px solid ${border}`,
						display: 'flex',
						flexDirection: 'column',
						gap: '4px',
					}}
				>
					<span
						style={{
							fontSize: '0.75rem',
							color,
						}}
					>
						{result.message}
					</span>

					{result.timestamp && (
						<span
							style={{
								fontSize: '0.6rem',
								color: 'rgba(233,230,221,0.25)',
							}}
						>
							{result.timestamp}
						</span>
					)}
				</div>
			)}
		</div>
	)
}
