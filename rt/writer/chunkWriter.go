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
	return WriteRune(c, q)
}

func (c ChunkWriter) WriteString(s string) (ret int, err error) {
	for _, q := range s {
		if i, e := c.WriteRune(q); e != nil {
			err = e
			break
		} else {
			ret += i
		}
	}
	return
}

func WriteRune(w io.Writer, q rune) (int, error) {
	var p [utf8.UTFMax]byte
	n := utf8.EncodeRune(p[:], q)
	return w.Write(p[:n])
}
