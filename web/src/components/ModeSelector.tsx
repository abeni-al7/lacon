import { Lock, Unlock } from 'lucide-react'

interface ModeSelectorProps {
  mode: 'encode' | 'decode'
  onChange: (mode: 'encode' | 'decode') => void
  disabled?: boolean
}

export default function ModeSelector({ mode, onChange, disabled }: ModeSelectorProps) {
  return (
    <div className="flex items-center justify-center animate-fade-in">
      <div className="relative inline-flex items-center p-1 rounded-2xl bg-[var(--color-muted)] border border-[var(--color-border)] shadow-sm">
        {/* Active indicator */}
        <div
          className={`
            absolute top-1 bottom-1 w-[calc(50%-4px)] rounded-xl transition-all duration-300 ease-[cubic-bezier(0.34,1.56,0.64,1)]
            ${mode === 'encode'
              ? 'left-1 bg-gradient-to-r from-[var(--color-primary)] to-[var(--color-primary-light)] shadow-md shadow-[var(--color-primary)]/20'
              : 'left-[calc(50%+2px)] bg-gradient-to-r from-[var(--color-accent)] to-[var(--color-accent-light)] shadow-md shadow-[var(--color-accent)]/20'
            }
          `}
        />

        <button
          type="button"
          disabled={disabled}
          onClick={() => onChange('encode')}
          className={`
            relative z-10 flex items-center gap-2.5 px-6 py-3 rounded-xl text-sm font-semibold transition-all duration-200
            ${mode === 'encode'
              ? 'text-white'
              : 'text-[var(--color-foreground)] opacity-60 hover:opacity-90'
            }
            disabled:opacity-50 disabled:cursor-not-allowed
          `}
        >
          <Lock size={16} className={mode === 'encode' ? 'text-white/90' : ''} />
          Encode
        </button>

        <button
          type="button"
          disabled={disabled}
          onClick={() => onChange('decode')}
          className={`
            relative z-10 flex items-center gap-2.5 px-6 py-3 rounded-xl text-sm font-semibold transition-all duration-200
            ${mode === 'decode'
              ? 'text-white'
              : 'text-[var(--color-foreground)] opacity-60 hover:opacity-90'
            }
            disabled:opacity-50 disabled:cursor-not-allowed
          `}
        >
          <Unlock size={16} className={mode === 'decode' ? 'text-white/90' : ''} />
          Decode
        </button>
      </div>
    </div>
  )
}