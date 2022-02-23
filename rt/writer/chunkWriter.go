package writer

import (
	"io"
	"unicode/utf8"
)

// ChunkWriter - adapter a regular io.Writer into a chunk writer
type ChunkWriter struct{ io.Writer }

func (c ChunkWriter) WriteByte(b byte) error {
	_, e := c.Write([]byte{b})
	return e
}

func (c ChunkWriter) WriteRune(q rune) (int, error) {
	return WriteRune(c.Writer, q)
}

func (c ChunkWriter) WriteString(s string) (ret int, err error) {
	return WriteString(c.Writer, s)
}

// WriteString helper mimics io.WriteString, invoking w.WriteString if it implements it
// otherwise falling back to WriteRune, if it implements it or Write otherwise.
// using io.WriteString is generally enough because the formatting layers use ChunkWriter which use this.
func WriteString(w io.Writer, s string) (ret int, err error) {
	if ws, ok := w.(io.StringWriter); ok {
		ret, err = ws.WriteString(s)
	} else if rw, ok := w.(RuneWriter); ok {
		for _, q := range s {
			if i, e := rw.WriteRune(q); e != nil {
				err = e
				break
			} else {
				ret += i
			}
		}
	} else {
		ret, err = w.Write([]byte(s))
	}
	return
}

// WriteRune helper mimics io.WriteString, this invokes RuneWriter.WriteRune if w implements it.
// Otherwise w.Write gets called with the bytes necessary to encode the specified rune.
// Returns the number of bytes written.
func WriteRune(w io.Writer, q rune) (ret int, err error) {
	if rw, ok := w.(RuneWriter); ok {
		ret, err = rw.WriteRune(q)
	} else {
		var p [utf8.UTFMax]byte
		n := utf8.EncodeRune(p[:], q)
		ret, err = w.Write(p[:n])
	}
	return
}

type RuneWriter interface {
	WriteRune(q rune) (int, error)
}
