import { ArrowDown, FileText, Archive, Percent } from 'lucide-react'

interface CompressionInfoProps {
  originalSize: number
  resultSize: number
  filename: string
  mode: 'encode' | 'decode'
}

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

export default function CompressionInfo({ originalSize, resultSize, filename, mode }: CompressionInfoProps) {
  const ratio = originalSize > 0 ? ((1 - resultSize / originalSize) * 100).toFixed(1) : '0.0'
  const isSmaller = resultSize < originalSize

  return (
    <div className="animate-fade-in">
      <div className="rounded-2xl border border-[var(--color-border)] bg-[var(--color-muted)] p-6">
        <div className="flex items-center gap-2 mb-5">
          <div className="w-8 h-8 rounded-lg bg-[var(--color-primary)]/10 flex items-center justify-center">
            <ArrowDown size={18} className="text-[var(--color-primary)]" />
          </div>
          <h3 className="font-semibold text-[var(--color-foreground)]">
            {mode === 'encode' ? 'Compression Complete' : 'Decompression Complete'}
          </h3>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div className="rounded-xl bg-[var(--color-background)] p-4 border border-[var(--color-border)]">
            <div className="flex items-center gap-2 mb-2">
              <FileText size={16} className="text-[var(--color-foreground)] opacity-50" />
              <span className="text-xs font-medium uppercase tracking-wider text-[var(--color-foreground)] opacity-50">
                Original
              </span>
            </div>
            <p className="text-xl font-bold text-[var(--color-foreground)]">
              {formatSize(originalSize)}
            </p>
          </div>

          <div className="rounded-xl bg-[var(--color-background)] p-4 border border-[var(--color-border)]">
            <div className="flex items-center gap-2 mb-2">
              <Archive size={16} className="text-[var(--color-foreground)] opacity-50" />
              <span className="text-xs font-medium uppercase tracking-wider text-[var(--color-foreground)] opacity-50">
                {mode === 'encode' ? 'Compressed' : 'Decompressed'}
              </span>
            </div>
            <p className="text-xl font-bold text-[var(--color-foreground)]">
              {formatSize(resultSize)}
            </p>
          </div>

          <div className="rounded-xl bg-[var(--color-background)] p-4 border border-[var(--color-border)]">
            <div className="flex items-center gap-2 mb-2">
              <Percent size={16} className="text-[var(--color-foreground)] opacity-50" />
              <span className="text-xs font-medium uppercase tracking-wider text-[var(--color-foreground)] opacity-50">
                Ratio
              </span>
            </div>
            <p className={`text-xl font-bold ${isSmaller ? 'text-green-600 dark:text-green-400' : 'text-[var(--color-foreground)]'}`}>
              {isSmaller ? `${ratio}% smaller` : ratio === '0.0' ? '—' : `${Math.abs(parseFloat(ratio))}% larger`}
            </p>
          </div>
        </div>

        <p className="mt-4 text-sm text-[var(--color-foreground)] opacity-60 text-center">
          File <span className="font-medium opacity-80">{filename}</span> has been downloaded automatically.
        </p>
      </div>
    </div>
  )
}