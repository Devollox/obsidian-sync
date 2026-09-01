type Props = {
	onMinimize: () => void
	onHide: () => void
	onQuit: () => void
}

export function TitleBar({ onMinimize, onHide, onQuit }: Props) {
	return (
		<div
			className='drag'
			style={{
				height: '40px',
				display: 'flex',
				alignItems: 'center',
				justifyContent: 'space-between',
				padding: '0 10px 0 16px',
				flexShrink: 0,
			}}
		>
			<span
				style={{ fontSize: '0.65rem', letterSpacing: '0.08em', color: 'rgba(233,230,221,0.30)' }}
			>
				Obsync
			</span>

			<div className='no-drag' style={{ display: 'flex', gap: '4px', alignItems: 'center' }}>
				<WinBtn onClick={onMinimize} title='Minimize'>
					–
				</WinBtn>
				<WinBtn onClick={onHide} title='Hide to tray'>
					□
				</WinBtn>
				<WinBtn onClick={onQuit} title='Quit' danger>
					×
				</WinBtn>
			</div>
		</div>
	)
}

function WinBtn({
	onClick,
	title,
	children,
	danger,
}: {
	onClick: () => void
	title: string
	children: React.ReactNode
	danger?: boolean
}) {
	return (
		<button
			onClick={onClick}
			title={title}
			style={{
				width: '22px',
				height: '22px',
				borderRadius: '50%',
				border: 'none',
				background: danger ? 'rgba(248,113,113,0.12)' : 'rgba(233,230,221,0.07)',
				color: danger ? 'rgba(248,113,113,0.65)' : 'rgba(233,230,221,0.35)',
				fontSize: '0.9rem',
				cursor: 'pointer',
				display: 'flex',
				alignItems: 'center',
				justifyContent: 'center',
				lineHeight: 1,
				transition: 'background 120ms, color 120ms',
				fontFamily: 'inherit',
			}}
			onMouseEnter={e => {
				e.currentTarget.style.background = danger
					? 'rgba(248,113,113,0.28)'
					: 'rgba(233,230,221,0.14)'
				e.currentTarget.style.color = danger ? '#f87171' : 'rgba(233,230,221,0.80)'
			}}
			onMouseLeave={e => {
				e.currentTarget.style.background = danger
					? 'rgba(248,113,113,0.12)'
					: 'rgba(233,230,221,0.07)'
				e.currentTarget.style.color = danger ? 'rgba(248,113,113,0.65)' : 'rgba(233,230,221,0.35)'
			}}
		>
			{children}
		</button>
	)
}
