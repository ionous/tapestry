package print

import (
	"io"
	"unicode"

	"git.sr.ht/~ionous/tapestry/rt/writer"
)

// A special writer that looks for trailing full stops, and other terminals
// and writes a line after them.
// https://www.unicode.org/review/pr-23.html
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
