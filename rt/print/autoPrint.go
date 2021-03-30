package print

import (
	"unicode"

	"git.sr.ht/~ionous/iffy/rt/writer"
	"git.sr.ht/~ionous/iffy/web/text"
)

// NewAutoWriter accepts incoming text chunks and writes them to target writing newlines at the end of sentences.
func NewAutoWriter(w writer.Output) writer.ChunkOutput {
	out := writer.ChunkWriter{text.Html2Text(w)}
	return func(c writer.Chunk) (int, error) {
		n, e := c.WriteTo(out)
		if last, _ := c.DecodeLastRune(); unicode.Is(unicode.Terminal_Punctuation, last) {
			out.WriteRune('\n')
		}
		return n, e
	}
}
