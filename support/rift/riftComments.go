package rift

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
)

type CommentWriter interface {
	WriteRune(r rune) (int, error)
	WriteString(s string) (int, error)
}
type NullComments struct{}

var nullComments NullComments

func (NullComments) WriteRune(r rune) (_ int, _ error)     { return }
func (NullComments) WriteString(s string) (_ int, _ error) { return }

// read everything until the end of the line as a comment.
func ReadComment(out CommentWriter, eol charm.State) charm.State {
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

type CommentBlock struct {
	keepComments bool
	comments     strings.Builder
}

func (b *CommentBlock) Comments() (ret CommentWriter) {
	if b.keepComments {
		ret = &b.comments
	} else {
		ret = nullComments
	}
	return
}
