import { Lock, Unlock } from 'lucide-react'

interface ModeSelectorProps {
  mode: 'encode' | 'decode'
  onChange: (mode: 'encode' | 'decode') => void
  disabled?: boolean
}

export default function ModeSelector({ mode, onChange, disabled }: ModeSelectorProps) {
  return (
    <div className="flex items-center justify-center gap-2 p-1.5 rounded-xl bg-[var(--color-muted)] border border-[var(--color-border)] w-fit mx-auto animate-fade-in">
      <button
        type="button"
        disabled={disabled}
        onClick={() => onChange('encode')}
        className={`
          flex items-center gap-2 px-5 py-2.5 rounded-lg text-sm font-medium transition-all duration-200
          ${mode === 'encode'
            ? 'bg-[var(--color-primary)] text-white shadow-sm'
            : 'text-[var(--color-foreground)] opacity-60 hover:opacity-90'
          }
          disabled:opacity-50 disabled:cursor-not-allowed
        `}
      >
        <Lock size={16} />
        Encode
      </button>
      <button
        type="button"
        disabled={disabled}
        onClick={() => onChange('decode')}
        className={`
          flex items-center gap-2 px-5 py-2.5 rounded-lg text-sm font-medium transition-all duration-200
          ${mode === 'decode'
            ? 'bg-[var(--color-accent)] text-white shadow-sm'
            : 'text-[var(--color-foreground)] opacity-60 hover:opacity-90'
          }
          disabled:opacity-50 disabled:cursor-not-allowed
        `}
      >
        <Unlock size={16} />
        Decode
      </button>
    </div>
  )
}