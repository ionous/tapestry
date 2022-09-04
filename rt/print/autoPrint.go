package print

import (
	"io"
	"unicode"

	"git.sr.ht/~ionous/tapestry/rt/writer"
)

func NewLineSentences(w io.Writer) writer.ChunkOutput {
	out := writer.ChunkWriter{w}
	return func(c writer.Chunk) (int, error) {
		n, e := c.WriteTo(out)
		if last, _ := c.DecodeLastRune(); unicode.Is(unicode.Sentence_Terminal, last) {
			out.WriteRune('\n')
		}
		return n, e
	}
}
