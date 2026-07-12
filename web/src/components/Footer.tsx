import { GitFork } from 'lucide-react'

export default function Footer() {
  return (
    <footer className="mt-auto text-center py-8 px-4 border-t border-[var(--color-border)]">
      <div className="flex items-center justify-center gap-5 text-sm text-[var(--color-foreground)] opacity-50">
        <span>Built with Go</span>
        <span className="w-1 h-1 rounded-full bg-current" />
        <a
          href="https://github.com/abeni-al7/lacon"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center gap-1.5 hover:opacity-80 transition-opacity"
        >
          <GitFork size={14} />
          GitHub
        </a>
        <span className="w-1 h-1 rounded-full bg-current" />
        <span>MIT License</span>
      </div>
    </footer>
  )
}