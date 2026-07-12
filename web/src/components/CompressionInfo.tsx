import { FileText, Archive, Percent, Download, CheckCircle2 } from 'lucide-react'

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
  const isEncode = mode === 'encode'
  const saved = originalSize - resultSize
  const isSignificant = isEncode && isSmaller && parseFloat(ratio) > 5

  return (
    <div className="animate-fade-in-up">
      <div className="rounded-2xl border border-[var(--color-border)] bg-[var(--color-card)] shadow-sm overflow-hidden">
        {/* Header */}
        <div className="p-5 border-b border-[var(--color-border)] bg-gradient-to-r from-[var(--color-muted)] to-transparent">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-[var(--color-success)]/10 to-[var(--color-success)]/5 flex items-center justify-center border border-[var(--color-success)]/10">
              <CheckCircle2 size={20} className="text-[var(--color-success)]" />
            </div>
            <div>
              <h3 className="font-semibold text-[var(--color-foreground)]">
                {isEncode ? 'Compression Complete' : 'Decompression Complete'}
              </h3>
              <p className="text-xs text-[var(--color-foreground)] opacity-50 mt-0.5">
                File has been downloaded automatically
              </p>
            </div>
          </div>
        </div>

        {/* Stats grid */}
        <div className="p-5">
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div className="rounded-xl bg-[var(--color-background)] p-4 border border-[var(--color-border)]">
              <div className="flex items-center gap-2 mb-2">
                <div className="w-7 h-7 rounded-lg bg-[var(--color-primary)]/10 flex items-center justify-center">
                  <FileText size={14} className="text-[var(--color-primary)]" />
                </div>
                <span className="text-xs font-semibold uppercase tracking-wider text-[var(--color-foreground)] opacity-50">
                  Original
                </span>
              </div>
              <p className="text-xl font-bold text-[var(--color-foreground)]">
                {formatSize(originalSize)}
              </p>
            </div>

            <div className="rounded-xl bg-[var(--color-background)] p-4 border border-[var(--color-border)]">
              <div className="flex items-center gap-2 mb-2">
                <div className="w-7 h-7 rounded-lg bg-[var(--color-accent)]/10 flex items-center justify-center">
                  <Archive size={14} className="text-[var(--color-accent)]" />
                </div>
                <span className="text-xs font-semibold uppercase tracking-wider text-[var(--color-foreground)] opacity-50">
                  {isEncode ? 'Compressed' : 'Decompressed'}
                </span>
              </div>
              <p className="text-xl font-bold text-[var(--color-foreground)]">
                {formatSize(resultSize)}
              </p>
            </div>

            <div className="rounded-xl bg-[var(--color-background)] p-4 border border-[var(--color-border)]">
              <div className="flex items-center gap-2 mb-2">
                <div className="w-7 h-7 rounded-lg bg-[var(--color-success)]/10 flex items-center justify-center">
                  <Percent size={14} className="text-[var(--color-success)]" />
                </div>
                <span className="text-xs font-semibold uppercase tracking-wider text-[var(--color-foreground)] opacity-50">
                  Ratio
                </span>
              </div>
              <p className={`text-xl font-bold ${isEncode && isSmaller ? 'text-[var(--color-success)]' : isEncode ? 'text-[var(--color-foreground)]' : 'text-[var(--color-foreground)]'}`}>
                {isEncode && isSmaller ? `${ratio}% saved` : ratio === '0.0' ? '—' : `${Math.abs(parseFloat(ratio))}%`}
              </p>
            </div>
          </div>

          {/* Visual comparison bar (only for encode) */}
          {isEncode && originalSize > 0 && (
            <div className="mt-4 p-4 rounded-xl bg-[var(--color-background)] border border-[var(--color-border)]">
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-medium text-[var(--color-foreground)] opacity-50">Size comparison</span>
                {isSmaller && isSignificant && (
                  <span className="inline-flex items-center gap-1 text-xs font-semibold text-[var(--color-success)]">
                    <Download size={12} />
                    Saved {formatSize(saved)}
                  </span>
                )}
              </div>
              <div className="relative h-3 rounded-full bg-[var(--color-border)] overflow-hidden">
                <div
                  className="absolute inset-y-0 left-0 rounded-full bg-gradient-to-r from-[var(--color-primary)] to-[var(--color-primary-light)] transition-all duration-700 ease-out"
                  style={{ width: `${Math.min(100, (resultSize / originalSize) * 100)}%` }}
                />
              </div>
              <div className="flex justify-between mt-1.5">
                <span className="text-xs text-[var(--color-foreground)] opacity-50">{formatSize(resultSize)}</span>
                <span className="text-xs text-[var(--color-foreground)] opacity-50">{formatSize(originalSize)}</span>
              </div>
            </div>
          )}

          <p className="mt-4 text-xs text-[var(--color-foreground)] opacity-50 text-center">
            File <span className="font-medium opacity-70">{filename}</span> has been downloaded
          </p>
        </div>
      </div>
    </div>
  )
}