import { useState, useCallback } from 'react'

interface CompressionResult {
  originalSize: number
  resultSize: number
  filename: string
}

interface UseCompressionReturn {
  isProcessing: boolean
  error: string | null
  result: CompressionResult | null
  process: (file: File, mode: 'encode' | 'decode') => Promise<void>
  reset: () => void
}

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export function useCompression(): UseCompressionReturn {
  const [isProcessing, setIsProcessing] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [result, setResult] = useState<CompressionResult | null>(null)

  const process = useCallback(async (file: File, mode: 'encode' | 'decode') => {
    setIsProcessing(true)
    setError(null)
    setResult(null)

    const originalSize = file.size

    try {
      const formData = new FormData()
      formData.append('file', file)

      const endpoint = mode === 'encode' ? '/encode' : '/decode'
      const response = await fetch(`${API_URL}${endpoint}`, {
        method: 'POST',
        body: formData,
      })

      if (!response.ok) {
        const contentType = response.headers.get('content-type')
        if (contentType?.includes('application/json')) {
          const errData = await response.json()
          throw new Error(errData.error || `Request failed with status ${response.status}`)
        }
        throw new Error(`Request failed with status ${response.status}`)
      }

      const blob = await response.blob()
      const contentDisposition = response.headers.get('content-disposition')
      let filename = mode === 'encode' ? `${file.name}.lacon` : `${file.name}.decoded`

      if (contentDisposition) {
        const match = contentDisposition.match(/filename="?(.+?)"?$/)
        if (match) {
          filename = match[1]
        }
      }

      // Trigger download
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = filename
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)

      setResult({
        originalSize,
        resultSize: blob.size,
        filename,
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An unexpected error occurred')
    } finally {
      setIsProcessing(false)
    }
  }, [])

  const reset = useCallback(() => {
    setError(null)
    setResult(null)
  }, [])

  return { isProcessing, error, result, process, reset }
}