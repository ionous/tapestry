package rift

import (
	"strings"
)

// the bits needed from go's missing string builder / generic writer interface
type CommentWriter interface {
	WriteRune(r rune) (int, error)
}

// a "null device" for CommentWriter. it eats all writes.
type NullCommentWriter struct{}

var nullCommentWriter NullCommentWriter

func (NullCommentWriter) WriteRune(r rune) (_ int, _ error) { return }

// a strings builder which trims write spaces
type CommentBuffer struct {
	buf    strings.Builder
	spaces int
}

func (w *CommentBuffer) Len() int {
	return w.buf.Len() - w.spaces
}

func (w *CommentBuffer) WriteRune(r rune) (_ int, _ error) {
	if r == Space {
		w.spaces++
	} else {
		if r < Space && w.spaces > 0 {
			s := w.String() // trim
			w.buf.Reset()
			w.buf.WriteString(s)
		}
		w.spaces = 0
	}
	return w.buf.WriteRune(r)
}

func (w *CommentBuffer) String() string {
	str := w.buf.String()
	return str[:len(str)-w.spaces]
}
