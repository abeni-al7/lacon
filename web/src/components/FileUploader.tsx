import { useRef, useState, useCallback } from 'react'
import { Upload, Lock, Unlock, X, File, CheckCircle2 } from 'lucide-react'

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
          transition-all duration-300 overflow-hidden
          ${isDragging
            ? 'border-[var(--color-primary)] bg-[var(--color-primary)]/5 scale-[1.02] shadow-[var(--shadow-glow)]'
            : selectedFile
              ? 'border-[var(--color-primary)]/40 bg-[var(--color-card)] hover:border-[var(--color-primary)]/60'
              : 'border-[var(--color-border)] bg-[var(--color-card)] hover:border-[var(--color-primary)]/40 hover:bg-[var(--color-muted)] hover:shadow-md'
          }
          ${isProcessing ? 'pointer-events-none opacity-60' : ''}
        `}
      >
        {/* Drag overlay glow */}
        {isDragging && (
          <div className="absolute inset-0 rounded-2xl bg-gradient-to-br from-[var(--color-primary)]/5 to-transparent pointer-events-none" />
        )}

        <div className="relative p-8 sm:p-10 text-center">
          {selectedFile ? (
            <div className="flex flex-col items-center gap-4">
              <div className="relative">
                <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-[var(--color-primary)]/10 to-[var(--color-primary)]/5 flex items-center justify-center border border-[var(--color-primary)]/10">
                  <File size={30} className="text-[var(--color-primary)]" />
                </div>
                <div className="absolute -top-1 -right-1 w-5 h-5 rounded-full bg-[var(--color-success)] flex items-center justify-center shadow-sm">
                  <CheckCircle2 size={12} className="text-white" />
                </div>
              </div>
              <div className="space-y-1">
                <p className="font-semibold text-[var(--color-foreground)] text-base">{selectedFile.name}</p>
                <p className="text-sm text-[var(--color-foreground)] opacity-50">{formatSize(selectedFile.size)}</p>
              </div>
              {/* File size visual bar */}
              <div className="w-full max-w-[200px] h-1.5 rounded-full bg-[var(--color-border)] overflow-hidden">
                <div
                  className="h-full rounded-full bg-gradient-to-r from-[var(--color-primary)] to-[var(--color-primary-light)] transition-all duration-500"
                  style={{ width: `${Math.min(100, (selectedFile.size / 10485760) * 100)}%` }}
                />
              </div>
              {!isProcessing && (
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation()
                    onFileSelect(null as unknown as File)
                  }}
                  className="inline-flex items-center gap-1.5 text-sm text-[var(--color-foreground)] opacity-40 hover:opacity-70 transition-opacity"
                >
                  <X size={14} />
                  Remove file
                </button>
              )}
            </div>
          ) : (
            <div className="flex flex-col items-center gap-4">
              <div className={`
                w-16 h-16 rounded-2xl flex items-center justify-center transition-all duration-300
                ${isDragging
                  ? 'bg-[var(--color-primary)]/10 scale-110'
                  : 'bg-[var(--color-muted)]'
                }
              `}>
                {mode === 'encode'
                  ? <Lock size={28} className={`transition-colors duration-300 ${isDragging ? 'text-[var(--color-primary)]' : 'text-[var(--color-primary)] opacity-60'}`} />
                  : <Unlock size={28} className={`transition-colors duration-300 ${isDragging ? 'text-[var(--color-accent)]' : 'text-[var(--color-accent)] opacity-60'}`} />
                }
              </div>
              <div className="space-y-1">
                <p className="font-semibold text-[var(--color-foreground)]">
                  {isDragging ? 'Release to drop your file' : 'Drop a file here or click to browse'}
                </p>
                <p className="text-sm text-[var(--color-foreground)] opacity-50">
                  {mode === 'encode' ? 'Any plain text file supported' : 'Upload a .lacon compressed file'}
                </p>
              </div>
              <div className={`
                inline-flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium transition-all duration-200
                ${isDragging
                  ? 'bg-[var(--color-primary)] text-white shadow-md shadow-[var(--color-primary)]/20'
                  : 'bg-[var(--color-muted)] text-[var(--color-foreground)] opacity-60 border border-[var(--color-border)]'
                }
              `}>
                <Upload size={14} />
                <span>Choose file to {actionLabel.toLowerCase()}</span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}