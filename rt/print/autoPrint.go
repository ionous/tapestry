package print

import (
	"unicode"

	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/web/text"
)

// NewAutoWriter accepts incoming text chunks and writes them to target writing newlines at the end of sentences.
func NewAutoWriter(w writer.Output) writer.ChunkOutput {
	out := writer.ChunkWriter{text.Html2Text(w)}
	return func(c writer.Chunk) (int, error) {
		n, e := c.WriteTo(out)
		if last, _ := c.DecodeLastRune(); unicode.Is(unicode.Sentence_Terminal, last) {
			out.WriteRune('\n')
		}
		return n, e
	}
}
