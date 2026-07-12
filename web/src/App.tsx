import { Analytics } from "@vercel/analytics/react"
import { useState } from 'react'
import Header from './components/Header'
import ModeSelector from './components/ModeSelector'
import FileUploader from './components/FileUploader'
import CompressionInfo from './components/CompressionInfo'
import Footer from './components/Footer'
import { useCompression } from './hooks/useCompression'
import { ArrowRight } from 'lucide-react'

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
  const actionColor = mode === 'encode'
    ? 'bg-[var(--color-primary)] hover:bg-[var(--color-primary)]/90'
    : 'bg-[var(--color-accent)] hover:bg-[var(--color-accent)]/90'

  return (
    <>
      <Header />

      <main className="flex-1 w-full max-w-lg mx-auto px-4 pb-12 space-y-6">
        <ModeSelector mode={mode} onChange={handleModeChange} disabled={isProcessing} />

        <FileUploader
          mode={mode}
          onFileSelect={setSelectedFile}
          isProcessing={isProcessing}
          selectedFile={selectedFile}
        />

        {selectedFile && !isProcessing && !result && (
          <div className="flex justify-center animate-fade-in">
            <button
              type="button"
              onClick={handleSubmit}
              className={`
                flex items-center gap-2 px-8 py-3 rounded-xl text-white font-medium
                text-base transition-all duration-200 shadow-sm
                ${actionColor}
                active:scale-[0.98]
              `}
            >
              {actionLabel}
              <ArrowRight size={18} />
            </button>
          </div>
        )}

        {isProcessing && (
          <div className="flex flex-col items-center gap-3 py-8 animate-fade-in">
            <div className="w-10 h-10 rounded-full border-3 border-[var(--color-border)] border-t-[var(--color-primary)] animate-spin" />
            <p className="text-sm text-[var(--color-foreground)] opacity-60">
              {mode === 'encode' ? 'Compressing file...' : 'Decompressing file...'}
            </p>
          </div>
        )}

        {error && (
          <div className="animate-fade-in">
            <div className="rounded-2xl border border-[var(--color-destructive)]/30 bg-[var(--color-destructive)]/5 p-5">
              <div className="flex items-start gap-3">
                <span className="text-[var(--color-destructive)] font-bold text-sm mt-0.5">!</span>
                <div>
                  <h4 className="font-medium text-sm text-[var(--color-destructive)]">Error</h4>
                  <p className="text-sm text-[var(--color-foreground)] opacity-70 mt-1">{error}</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {result && (
          <CompressionInfo
            originalSize={result.originalSize}
            resultSize={result.resultSize}
            filename={result.filename}
            mode={mode}
          />
        )}

        {result && (
          <div className="flex justify-center animate-fade-in">
            <button
              type="button"
              onClick={() => {
                setSelectedFile(null)
                reset()
              }}
              className="px-6 py-2.5 rounded-xl text-sm font-medium text-[var(--color-foreground)] opacity-60 hover:opacity-100 bg-[var(--color-muted)] border border-[var(--color-border)] transition-all duration-200"
            >
              Compress another file
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