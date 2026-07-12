import { Analytics } from "@vercel/analytics/react"
import { useState } from 'react'
import Header from './components/Header'
import ModeSelector from './components/ModeSelector'
import FileUploader from './components/FileUploader'
import CompressionInfo from './components/CompressionInfo'
import Footer from './components/Footer'
import { useCompression } from './hooks/useCompression'
import { ArrowRight, RotateCcw } from 'lucide-react'

function App() {
  const [mode, setMode] = useState<'encode' | 'decode'>('encode')
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const { isProcessing, error, result, process, reset } = useCompression()

  const handleModeChange = (newMode: 'encode' | 'decode') => {
    if (isProcessing) return
    setMode(newMode)
    setSelectedFile(null)
    reset()
  }

  const handleSubmit = () => {
    if (selectedFile && !isProcessing) {
      process(selectedFile, mode)
    }
  }

  const actionLabel = mode === 'encode' ? 'Compress' : 'Decompress'

  return (
    <>
      <Header />

      <main className="flex-1 w-full max-w-lg mx-auto px-4 pb-16 space-y-5">
        <ModeSelector mode={mode} onChange={handleModeChange} disabled={isProcessing} />

        <FileUploader
          mode={mode}
          onFileSelect={setSelectedFile}
          isProcessing={isProcessing}
          selectedFile={selectedFile}
        />

        {/* Submit CTA */}
        {selectedFile && !isProcessing && !result && (
          <div className="flex justify-center animate-slide-up">
            <button
              type="button"
              onClick={handleSubmit}
              className="
                group relative inline-flex items-center gap-2.5 px-8 py-3.5 rounded-xl font-semibold
                text-base text-white transition-all duration-200 overflow-hidden
                active:scale-[0.97]
              "
            >
              {/* Gradient background */}
              <div className={`
                absolute inset-0 rounded-xl transition-opacity duration-200
                ${mode === 'encode'
                  ? 'bg-gradient-to-r from-[var(--color-primary)] to-[var(--color-primary-light)]'
                  : 'bg-gradient-to-r from-[var(--color-accent)] to-[var(--color-accent-light)]'
                }
              `} />
              {/* Hover overlay */}
              <div className="absolute inset-0 rounded-xl bg-white/0 group-hover:bg-white/10 transition-colors duration-200" />
              {/* Shadow */}
              <div className={`
                absolute inset-0 rounded-xl blur-xl opacity-30 transition-opacity duration-200
                ${mode === 'encode'
                  ? 'bg-[var(--color-primary)]'
                  : 'bg-[var(--color-accent)]'
                }
              `} />
              {/* Content */}
              <span className="relative z-10 flex items-center gap-2.5">
                {actionLabel}
                <ArrowRight size={18} className="group-hover:translate-x-0.5 transition-transform duration-200" />
              </span>
            </button>
          </div>
        )}

        {/* Processing state */}
        {isProcessing && (
          <div className="flex flex-col items-center gap-4 py-10 animate-fade-in">
            <div className="relative">
              <div className="w-12 h-12 rounded-full border-3 border-[var(--color-border)] border-t-[var(--color-primary)] animate-spin" />
              <div className="absolute inset-0 rounded-full animate-pulse-ring border border-[var(--color-primary)]/20" />
            </div>
            <div className="text-center space-y-1">
              <p className="text-sm font-medium text-[var(--color-foreground)]">
                {mode === 'encode' ? 'Compressing file...' : 'Decompressing file...'}
              </p>
              <p className="text-xs text-[var(--color-foreground)] opacity-40">
                This should only take a moment
              </p>
            </div>
          </div>
        )}

        {/* Error state */}
        {error && (
          <div className="animate-scale-in">
            <div className="rounded-2xl border border-[var(--color-destructive)]/20 bg-[var(--color-destructive)]/5 p-5">
              <div className="flex items-start gap-3">
                <div className="w-7 h-7 rounded-lg bg-[var(--color-destructive)]/10 flex items-center justify-center shrink-0">
                  <span className="text-[var(--color-destructive)] font-bold text-sm">!</span>
                </div>
                <div className="flex-1 min-w-0">
                  <h4 className="font-semibold text-sm text-[var(--color-destructive)]">Processing failed</h4>
                  <p className="text-sm text-[var(--color-foreground)] opacity-70 mt-1 leading-relaxed">{error}</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Result */}
        {result && (
          <CompressionInfo
            originalSize={result.originalSize}
            resultSize={result.resultSize}
            filename={result.filename}
            mode={mode}
          />
        )}

        {/* Reset button */}
        {result && (
          <div className="flex justify-center animate-slide-up">
            <button
              type="button"
              onClick={() => {
                setSelectedFile(null)
                reset()
              }}
              className="
                inline-flex items-center gap-2 px-6 py-3 rounded-xl text-sm font-medium
                text-[var(--color-foreground)] opacity-50 hover:opacity-90
                bg-[var(--color-card)] border border-[var(--color-border)]
                hover:border-[var(--color-primary)]/30 hover:shadow-sm
                transition-all duration-200 active:scale-[0.98]
              "
            >
              <RotateCcw size={15} />
              Process another file
            </button>
          </div>
        )}
      </main>

      <Footer />
      <Analytics />
    </>
  )
}

export default App