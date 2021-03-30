package print

import (
	"bytes"
	"strings"

	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt/writer"
)

// Parens buffers writer.Output, grouping a stream of writes.
// Close adds the closing paren.
func Parens() *BracketSpanner {
	return &BracketSpanner{open: '(', close: ')'}
}

type BracketSpanner struct {
	Spanner     // inside the brackets: write with spaces
	open, close rune
}

func (p *BracketSpanner) ChunkOutput() writer.ChunkOutput {
	return p.WriteChunk
}

func (p *BracketSpanner) WriteChunk(c writer.Chunk) (ret int, err error) {
	if c.IsClosed() {
		if p.Len() > 0 {
			p.buf.WriteRune(p.close)
		}
	} else {
		if p.buf.Len() > 0 {
			ret, err = p.Spanner.WriteChunk(c)
		} else {
			var buf bytes.Buffer
			ret, err = c.WriteTo(&buf)
			// wrote something locally? prepend it with the open.
			if buf.Len() > 0 {
				p.buf.WriteRune(p.open)
				buf.WriteTo(&p.buf)
			}
		}
	}
	return
}

// Carriage filters writer.Output, creating a new line after every write
func Carriage(out writer.Output) writer.Output {
	f := &Filter{
		Rest: func(c writer.Chunk) (ret int, err error) {
			ret, err = c.WriteTo(out)
			out.WriteRune('\n')
			return
		},
	}
	return writer.ChunkOutput(f.WriteChunk)
}

// Capitalize filters writer.Output, capitalizing the first string.
func Capitalize(out writer.Output) writer.Output {
	f := &Filter{
		First: func(c writer.Chunk) (int, error) {
			cap := lang.Capitalize(c.String())
			return out.WriteString(cap)
		},
		Rest: func(c writer.Chunk) (int, error) {
			return c.WriteTo(out)
		},
	}
	return writer.ChunkOutput(f.WriteChunk)
}

// TitleCase filters writer.Output, capitalizing every write.
func TitleCase(out writer.Output) writer.Output {
	f := &Filter{
		Rest: func(c writer.Chunk) (int, error) {
			cap := lang.Capitalize(c.String())
			return out.WriteString(cap)
		},
	}
	return writer.ChunkOutput(f.WriteChunk)
}

// Lowercase filters writer.Output, lowering every string.
func Lowercase(out writer.Output) writer.Output {
	f := &Filter{
		Rest: func(c writer.Chunk) (int, error) {
			cap := strings.ToLower(c.String())
			return out.WriteString(cap)
		},
	}
	return writer.ChunkOutput(f.WriteChunk)
}

// Slash filters writer.Output, separating writes with a slash.
func Slash(out writer.Output) writer.Output {
	f := &Filter{
		First: func(c writer.Chunk) (int, error) {
			return c.WriteTo(out)
		},
		Rest: func(c writer.Chunk) (int, error) {
			out.WriteString(" /")
			return c.WriteTo(out)
		},
	}
	return writer.ChunkOutput(f.WriteChunk)
}

// Tag - surrounds a block of text with an html-like element.
// It establishes a new "context" for writing --
// any existing "span", etc. does not apply to the contents of the tag.
func Tag(out writer.Output, tag string) writer.OutputCloser {
	var buf bytes.Buffer
	writeTag(&buf, tag, true)
	f := &Filter{
		Rest: func(c writer.Chunk) (int, error) {
			return c.WriteTo(&buf)
		},
		Last: func(int) (err error) {
			writeTag(&buf, tag, false)
			_, err = buf.WriteTo(out)
			return
		},
	}
	return writer.ChunkOutput(f.WriteChunk)
}

func writeTag(w writer.Output, tag string, open bool) {
	w.WriteRune('<')
	if !open {
		w.WriteRune('/')
	}
	w.WriteString(tag)
	w.WriteRune('>')
}
