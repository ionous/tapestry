package print

import (
	"unicode"

	"git.sr.ht/~ionous/iffy/rt/writer"
)

// NewAutoWriter accepts incoming text chunks and writes them to target writing newlines at the end of sentences.
func NewAutoWriter(w writer.Output) writer.ChunkOutput {
	return func(c writer.Chunk) (int, error) {
		n, e := c.WriteTo(w)
		if last, _ := c.DecodeLastRune(); unicode.Is(unicode.Terminal_Punctuation, last) {
			w.WriteRune('\n')
		}
		return n, e
	}
}
