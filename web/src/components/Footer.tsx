import { GitFork, FileCode, Server } from 'lucide-react'

export default function Footer() {
  return (
    <footer className="mt-auto text-center py-10 px-4">
      <div className="w-full max-w-lg mx-auto mb-6">
        <div className="h-px bg-gradient-to-r from-transparent via-[var(--color-border)] to-transparent" />
      </div>
      <div className="flex flex-wrap items-center justify-center gap-x-5 gap-y-2 text-sm text-[var(--color-foreground)] opacity-45">
        <span className="inline-flex items-center gap-1.5">
          <FileCode size={14} />
          Built with Go
        </span>
        <span className="hidden sm:inline w-1 h-1 rounded-full bg-current" />
        <span className="inline-flex items-center gap-1.5">
          <Server size={14} />
          Huffman Coding
        </span>
        <span className="hidden sm:inline w-1 h-1 rounded-full bg-current" />
        <a
          href="https://github.com/abeni-al7/lacon"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-1.5 hover:opacity-70 transition-opacity"
        >
          <GitFork size={14} />
          GitHub
        </a>
        <span className="hidden sm:inline w-1 h-1 rounded-full bg-current" />
        <span>MIT License</span>
      </div>
    </footer>
  )
}