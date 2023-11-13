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
	b      strings.Builder
	spaces int
}

func (w *CommentBuffer) Reset() {
	w.b.Reset()
	w.spaces = 0
}

func (w *CommentBuffer) Len() int {
	return w.b.Len() - w.spaces
}

func (w *CommentBuffer) WriteRune(r rune) (_ int, _ error) {
	if r == Space {
		w.spaces++
	} else {
		if r < Space && w.spaces > 0 {
			s := w.String() // trim
			w.b.Reset()
			w.b.WriteString(s)
		}
		w.spaces = 0
	}
	return w.b.WriteRune(r)
}

func (w *CommentBuffer) String() string {
	str := w.b.String()
	return str[:len(str)-w.spaces]
}
