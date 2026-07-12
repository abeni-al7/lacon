import { FileArchive } from 'lucide-react'

export default function Header() {
  return (
    <header className="text-center py-12 px-4 animate-fade-in">
      <div className="flex items-center justify-center gap-3 mb-4">
        <div className="w-12 h-12 rounded-xl bg-[var(--color-primary)] flex items-center justify-center">
          <FileArchive size={28} className="text-white" />
        </div>
        <h1 className="text-4xl sm:text-5xl font-bold tracking-tight text-[var(--color-foreground)]">
          Lacon
        </h1>
      </div>
      <p className="text-lg sm:text-xl text-[var(--color-foreground)] opacity-70 max-w-md mx-auto leading-relaxed">
        Compress and decompress files using <span className="font-semibold opacity-100">Huffman coding</span> — fast, lossless, and entirely client-to-server.
      </p>
    </header>
  )
}