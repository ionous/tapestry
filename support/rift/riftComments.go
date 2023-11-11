package rift

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
)

// the bits needed from go's missing string builder / generic writer interface
type CommentWriter interface {
	WriteRune(r rune) (int, error)
	WriteString(s string) (int, error)
}

func KeepCommentWriter() CommentBlock {
	return CommentBlock{keepComments: true}
}

func DiscardCommentWriter() CommentBlock {
	return CommentBlock{keepComments: false}
}

type CommentFactory func() CommentBlock

// read everything until the end of the line as a comment.
func ReadComment(out CommentWriter, eol charm.State) charm.State {
	out.WriteRune(Hash)
	return charm.Self("read comment", func(self charm.State, r rune) (ret charm.State) {
		if r == Newline {
			ret = eol
		} else {
			out.WriteRune(r)
			ret = self
		}
		return
	})
}

// dont copy a comments block with content
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

type NullCommentWriter struct{}

var nullCommentWriter NullCommentWriter

func (NullCommentWriter) WriteRune(r rune) (_ int, _ error)     { return }
func (NullCommentWriter) WriteString(s string) (_ int, _ error) { return }
