package print

import (
	"io"
	"unicode"

	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/web/markup"
)

// NewAutoWriter takes incoming text chunks and writes them to w adding newlines at the end of sentences.
func NewAutoWriter(w io.Writer) writer.ChunkOutput {
	out := writer.ChunkWriter{markup.ToText(w)}
	return func(c writer.Chunk) (int, error) {
		n, e := c.WriteTo(out)
		if last, _ := c.DecodeLastRune(); unicode.Is(unicode.Sentence_Terminal, last) {
			out.WriteRune('\n')
		}
		return n, e
	}
}

func NewLineWriter(w io.Writer) writer.ChunkOutput {
	out := writer.ChunkWriter{w}
	return func(c writer.Chunk) (int, error) {
		n, e := c.WriteTo(out)
		if last, _ := c.DecodeLastRune(); unicode.Is(unicode.Sentence_Terminal, last) {
			out.WriteRune('\n')
		}
		return n, e
	}
}
