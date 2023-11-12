package rift

import (
	"strings"
)

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
