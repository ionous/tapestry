package rift

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// represents the "contents" of an entry
type riftEntry struct {
	Collection
	pendingValue         pendingValue
	buffer, header, tail strings.Builder
	bufferedLines        int
	depth                int
}

func (ent *riftEntry) finalizeEntry() (err error) {
	c := ent.Collection
	if val, e := ent.pendingValue.FinalizeValue(); e != nil {
		err = e
	} else {
		w := c.CommentWriter()
		w.WriteString(ent.buffer.String())
		head := ent.header.String()
		tail := ent.tail.String()
		if len(head) > 0 || len(tail) > 0 {
			w.WriteRune(HTab)
			w.WriteString(head)
			w.WriteRune(HTab)
			w.WriteString(tail)
		}
		// fix: modify the signature to write the comment at the same time?
		err = c.writeValue(val)
	}
	return
}

// fix: can this be made internal
// even evaluating the pendingValue might be better
func (ent *riftEntry) writeHeader() (ret string, err error) {
	if ent.bufferedLines > 1 {
		err = errutil.New("ambiguous multiline comment.")
	} else {
		ret = ent.header.String()
		ent.header.Reset()
	}
	return
}

// parses contents and loops (by popping) after its done
func ContentsLoop(ent *riftEntry) charm.State {
	return charm.Step(Contents(ent),
		charm.Self("after entry", func(afterEntry charm.State, r rune) (ret charm.State) {
			switch r {
			case Newline:
				doc := ent.Document()
				ret = NextIndent(doc.popToIndent)
			}
			return
		}))
}

// contents appear after a collection marker:
//
// [-] <inline spaces> ( <inline comment> | <value> )
//
//	<buffered comment>
//	   <indented additional lines>
//	<header comment>
//	<value>
func Contents(ent *riftEntry) charm.State {
	return charm.Self("contents", func(contents charm.State, r rune) (ret charm.State) {
		switch r {
		case Space:
			ret = contents

		case Hash:
			// these use >= so that content can appear at column zero in documents
			if doc := ent.Document(); doc.Col >= ent.depth {
				ret = ReadComment(ent.CommentWriter(), contents)
			}

		case Newline:
			ret = NextIndent(func() (ret charm.State) {
				if doc := ent.Document(); doc.Col >= ent.depth {
					ret = BufferRegion(ent, doc.Col)
				} else {
					ret = doc.popToIndent()
				}
				return
			})

		default:
			if doc := ent.Document(); doc.Col >= ent.depth {
				ret = ValueOfEntry(ent, r)
			}
		}
		return
	})
}

// we are at the start of a line where comment buffering might occur.
func BufferRegion(ent *riftEntry, depth int) charm.State {
	return charm.Self("buffering", func(buffering charm.State, r rune) (ret charm.State) {
		switch r {

		// possibly a value at the same depth as the buffering section
		default:
			ret = ValueOfEntry(ent, r)

		// after a completely empty line: move to the header region.
		// FIX -- should be checking for nesting
		// newline should just loop
		case Newline:
			ret = HeaderRegion(ent, depth)

			// read comment, and search for next indent
		case Hash:
			ret = ReadComment(&ent.buffer,
				NextIndent(func() (ret charm.State) {
					ent.bufferedLines++
					switch doc := ent.Document(); {
					case doc.Col == depth:
						// at the same indent, stick where we're at.
						ret = buffering
					case ent.bufferedLines == 1 && doc.Col > depth:
						// the ideal multiline comment a single comment followed by indented lines
						ret = IndentedComment(ent, doc.Col)
					default:
						ret = doc.popToIndent()
					}
					return
				}))
			return
		}
		return
	})
}

// anything at the same indent can be a continuing comment
// anything at a different indent can be the header or value.
func IndentedComment(ent *riftEntry, depth int) (ret charm.State) {
	out := &ent.buffer
	return charm.Self("indented comment", func(indentedComment charm.State, r rune) (ret charm.State) {
		switch r {
		case Hash:
			ret = ReadComment(out, NextIndent(func() (ret charm.State) {
				// loop if we're at our indent, or start treating future comments as value header comments.
				if doc := ent.Document(); doc.Col == depth {
					ret = indentedComment
				} else {
					ret = HeaderRegion(ent, doc.Col)
				}
				return
			}))

		}
		return
	})
}

// after the buffering section
// collect the header and get ready to pass it to the value.
func HeaderRegion(ent *riftEntry, depth int) (ret charm.State) {
	out := &ent.header
	return charm.Self("header comment", func(headerComment charm.State, r rune) (ret charm.State) {
		switch r {
		case Hash:
			ret = ReadComment(out, NextIndent(func() (ret charm.State) {
				if doc := ent.Document(); doc.Col == depth {
					ret = headerComment
				} else {
					ret = doc.popToIndent()
				}
				return
			}))

		default:
			ret = ValueOfEntry(ent, r)
		}
		return
	})
}

// at the start of a rune which might be a value:
// reads that value and any trailing comment describing it.
func ValueOfEntry(ent *riftEntry, r rune) (ret charm.State) {
	// dont bother trying to read a value if it wasn't meant to be.
	if r != Newline && r != Space {
		ret = charm.RunState(r, charm.Step(NewValue(ent),
			InlineComment(ent)))
	}
	return
}

// fix/share: almost exactly the same as the padding contents...
func InlineComment(ent *riftEntry) (ret charm.State) {
	return charm.Self("trailing comments", func(inline charm.State, r rune) (ret charm.State) {
		switch r {
		case Space: // eat spaces on the line after the value
			ret = inline
		case Hash: // an inline comment? read it; loop to us to handle the newline.
			ret = ReadComment(&ent.tail, inline)
		case Newline: // a newline ( regardless of whether there was a comment )
			ret = NextIndent(func() (ret charm.State) {
				// on the following line, trailing comments can appear at or deeper than the entry.
				if doc := ent.Document(); doc.Col >= ent.depth {
					ret = TrailingComment(ent, doc.Col)
				} else {
					ret = doc.popToIndent()
				}
				return
			})
		}
		return
	})
}

// an optional comment can appear on the first line after a value
// starts on something other than whitespace
func TrailingComment(ent *riftEntry, depth int) charm.State {
	return charm.Self("trailing", func(trailing charm.State, r rune) (ret charm.State) {
		switch r {
		case Hash: // nested comments can appear at or deeper than the entry
			ret = ReadComment(&ent.tail, NextIndent(func() (ret charm.State) {
				if doc := ent.Document(); doc.Col >= ent.depth {
					ret = NestedComment(ent, doc.Col)
				} else {
					ret = doc.popToIndent()
				}
				return
			}))
			return
		}
		return
	})
}

// nested comments are fixed at the passed depth
// starts on something other than whitespace
func NestedComment(ent *riftEntry, depth int) charm.State {
	return charm.Self("nested", func(nested charm.State, r rune) (ret charm.State) {
		switch r {
		case Hash: // loop at the same depth
			ret = ReadComment(&ent.tail, NextIndent(func() (ret charm.State) {
				if doc := ent.Document(); doc.Col == ent.depth {
					ret = nested
				} else {
					ret = doc.popToIndent()
				}
				return
			}))
			return
		}
		return
	})
}
