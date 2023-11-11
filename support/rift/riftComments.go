package rift

import (
	"strings"
)

// the bits needed from go's missing string builder / generic writer interface
type CommentWriter interface {
	WriteRune(r rune) (int, error)
}

func KeepCommentWriter() CommentBlock {
	return CommentBlock{keepComments: true}
}

func DiscardCommentWriter() CommentBlock {
	return CommentBlock{keepComments: false}
}

// signature for functions which create comment blocks
type CommentFactory func() CommentBlock

// holds comments for a collection.
// fyi: don't copy a comments block with content
// ( re: strings.Builder zero value )
type CommentBlock struct {
	keepComments bool
	comments     strings.Builder
}

// implements Collection for aggregation
func (b *CommentBlock) CommentWriter() (ret CommentWriter) {
	if b.keepComments {
		ret = &b.comments
	} else {
		ret = nullCommentWriter
	}
	return
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

func (w *CommentBuffer) WriteRune(r rune) (_ int, _ error) {
	if r == Space {
		w.spaces++
	} else {
		w.spaces = 0
	}
	return w.buf.WriteRune(r)
}

func (w *CommentBuffer) String() string {
	str := w.buf.String()
	return str[:len(str)-w.spaces]
}
