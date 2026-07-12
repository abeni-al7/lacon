import { FileArchive, Zap } from 'lucide-react'

export default function Header() {
  return (
    <header className="text-center pt-16 pb-8 px-4 animate-fade-in-down">
      <div className="flex items-center justify-center gap-3 mb-5">
        <div className="relative">
          <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-[var(--color-primary)] to-[var(--color-primary-light)] flex items-center justify-center shadow-lg shadow-[var(--color-primary)]/20">
            <FileArchive size={30} className="text-white" />
          </div>
          <div className="absolute -top-1 -right-1 w-5 h-5 rounded-full bg-[var(--color-accent)] flex items-center justify-center shadow-sm">
            <Zap size={10} className="text-white" fill="white" />
          </div>
        </div>
        <div>
          <h1 className="text-4xl sm:text-5xl font-extrabold tracking-tight">
            <span className="bg-gradient-to-r from-[var(--color-primary)] to-[var(--color-accent)] bg-clip-text text-transparent">
              Lacon
            </span>
          </h1>
        </div>
      </div>
      <p className="text-base sm:text-lg text-[var(--color-foreground)] opacity-65 max-w-lg mx-auto leading-relaxed">
        Compress and decompress files using{' '}
        <span className="font-semibold opacity-90">Huffman coding</span> — fast, lossless, and entirely server-side.
      </p>
      <div className="flex items-center justify-center gap-2 mt-5">
        <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium bg-[var(--color-primary)]/10 text-[var(--color-primary)] border border-[var(--color-primary)]/15">
          <span className="w-1.5 h-1.5 rounded-full bg-[var(--color-primary)] animate-pulse" />
          Fully server-side
        </span>
        <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium bg-[var(--color-accent)]/10 text-[var(--color-accent)] border border-[var(--color-accent)]/15">
          <span className="w-1.5 h-1.5 rounded-full bg-[var(--color-accent)]" />
          Lossless
        </span>
      </div>
    </header>
  )
}