import type { SyncResult } from '../hooks/use-sync'

type Props = {
  vaultPath: string
  syncing: boolean
  result: SyncResult | null
  onSync: () => void
  onPickFolder: () => void
}

export function SyncPanel({ vaultPath, syncing, result, onSync, onPickFolder }: Props) {
  return (
    <div style={{ padding: '24px 20px', display: 'flex', flexDirection: 'column', gap: '24px' }}>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
        <span style={{ fontSize: '0.6rem', letterSpacing: '0.1em', color: 'rgba(233,230,221,0.25)' }}>
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
          onMouseEnter={e => (e.currentTarget.style.borderColor = 'rgba(233,230,221,0.18)')}
          onMouseLeave={e => (e.currentTarget.style.borderColor = 'rgba(233,230,221,0.08)')}
        >
          <span style={{
            fontSize: '0.75rem',
            color: vaultPath ? '#e9e6dd' : 'rgba(233,230,221,0.25)',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            whiteSpace: 'nowrap',
          }}>
            {vaultPath || 'Select vault folder…'}
          </span>
          <span style={{ fontSize: '0.65rem', color: 'rgba(233,230,221,0.30)', flexShrink: 0 }}>↗</span>
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
        onMouseEnter={e => {
          if (!syncing && vaultPath) e.currentTarget.style.background = 'rgba(233,230,221,0.12)'
        }}
        onMouseLeave={e => {
          e.currentTarget.style.background = syncing ? 'rgba(233,230,221,0.04)' : 'rgba(233,230,221,0.08)'
        }}
      >
        {syncing ? 'Syncing…' : 'Sync now'}
      </button>

      {result && (
        <div style={{
          padding: '10px 12px',
          borderRadius: '8px',
          background: result.status === 'error'
            ? 'rgba(248,113,113,0.06)'
            : 'rgba(74,222,128,0.06)',
          border: `1px solid ${result.status === 'error' ? 'rgba(248,113,113,0.15)' : 'rgba(74,222,128,0.15)'}`,
          display: 'flex',
          flexDirection: 'column',
          gap: '4px',
        }}>
          <span style={{
            fontSize: '0.75rem',
            color: result.status === 'error' ? '#f87171' : '#4ade80',
          }}>
            {result.message}
          </span>
          {result.timestamp && (
            <span style={{ fontSize: '0.6rem', color: 'rgba(233,230,221,0.25)' }}>
              {result.timestamp}
            </span>
          )}
        </div>
      )}
    </div>
  )
}
