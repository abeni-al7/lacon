import { useRef, useState, useCallback } from 'react'
import { Upload, FileText, Lock, Unlock } from 'lucide-react'

interface FileUploaderProps {
  mode: 'encode' | 'decode'
  onFileSelect: (file: File) => void
  isProcessing: boolean
  selectedFile: File | null
}

export default function FileUploader({ mode, onFileSelect, isProcessing, selectedFile }: FileUploaderProps) {
  const inputRef = useRef<HTMLInputElement>(null)
  const [isDragging, setIsDragging] = useState(false)

  const handleFile = useCallback((file: File) => {
    if (mode === 'decode' && !file.name.endsWith('.lacon')) {
      // Still allow it, but warn
    }
    onFileSelect(file)
  }, [mode, onFileSelect])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
    const file = e.dataTransfer.files[0]
    if (file) handleFile(file)
  }, [handleFile])

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragging(false)
  }, [])

  const handleClick = () => {
    inputRef.current?.click()
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) handleFile(file)
    // Reset so re-selecting the same file triggers onChange
    e.target.value = ''
  }

  const formatSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }

  const actionLabel = mode === 'encode' ? 'Compress' : 'Decompress'
  const acceptedExt = mode === 'decode' ? '.lacon' : undefined

  return (
    <div className="animate-fade-in">
      <input
        ref={inputRef}
        type="file"
        accept={acceptedExt}
        onChange={handleInputChange}
        className="hidden"
      />

      <div
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onClick={isProcessing ? undefined : handleClick}
        className={`
          relative cursor-pointer rounded-2xl border-2 border-dashed
          transition-all duration-200 p-10 text-center
          ${isDragging
            ? 'border-[var(--color-primary)] bg-[var(--color-primary)]/5 scale-[1.02]'
            : selectedFile
              ? 'border-[var(--color-primary)] bg-[var(--color-muted)]'
              : 'border-[var(--color-border)] hover:border-[var(--color-primary)]/50 hover:bg-[var(--color-muted)]'
          }
          ${isProcessing ? 'pointer-events-none opacity-60' : ''}
        `}
      >
        {selectedFile ? (
          <div className="flex flex-col items-center gap-3">
            <div className="w-14 h-14 rounded-xl bg-[var(--color-primary)]/10 flex items-center justify-center">
              <FileText size={28} className="text-[var(--color-primary)]" />
            </div>
            <div>
              <p className="font-medium text-[var(--color-foreground)]">{selectedFile.name}</p>
              <p className="text-sm text-[var(--color-foreground)] opacity-60 mt-0.5">
                {formatSize(selectedFile.size)}
              </p>
            </div>
            {!isProcessing && (
              <button
                type="button"
                onClick={(e) => {
                  e.stopPropagation()
                  onFileSelect(null as unknown as File)
                }}
                className="text-sm text-[var(--color-foreground)] opacity-50 hover:opacity-100 underline underline-offset-2 transition-opacity"
              >
                Choose a different file
              </button>
            )}
          </div>
        ) : (
          <div className="flex flex-col items-center gap-3">
            <div className="w-14 h-14 rounded-xl bg-[var(--color-muted)] flex items-center justify-center">
              {mode === 'encode'
                ? <Lock size={28} className="text-[var(--color-primary)] opacity-70" />
                : <Unlock size={28} className="text-[var(--color-accent)] opacity-70" />
              }
            </div>
            <div>
              <p className="font-medium text-[var(--color-foreground)]">
                {isDragging ? 'Drop file here' : 'Drop a file here or click to browse'}
              </p>
              <p className="text-sm text-[var(--color-foreground)] opacity-50 mt-1">
                {mode === 'encode' ? 'Any file type supported' : 'Upload a .lacon compressed file'}
              </p>
            </div>
            <div className="flex items-center gap-1.5 text-xs text-[var(--color-foreground)] opacity-40 mt-2">
              <Upload size={12} />
              <span>Select file to {actionLabel.toLowerCase()}</span>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}