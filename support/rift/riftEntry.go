package rift

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// represents the "contents" of an entry
type riftEntry struct {
	doc                  *Document
	pendingValue         pendingValue
	addsValue            func(val any, comment string) error
	header, buffer, tail CommentBuffer
	bufferedLines        int
	depth                int
}

// pop parser states up to the current indentation level
func (ent *riftEntry) popToIndent() charm.State {
	return ent.doc.popToIndent()
}

func (ent *riftEntry) finalizeEntry() (err error) {
	if val, e := ent.pendingValue.FinalizeValue(); e != nil {
		err = e
	} else {
		var w strings.Builder
		w.WriteString(ent.buffer.String())
		head := ent.header.String()
		tail := ent.tail.String()
		if len(head) > 0 || len(tail) > 0 {
			w.WriteRune(HTab)
			w.WriteString(head)
			w.WriteRune(HTab)
			w.WriteString(tail)
		}
		err = ent.addsValue(val, w.String())
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
		ent.header.buf.Reset()
	}
	return
}

// parses contents and loops (by popping) after its done
func ContentsLoop(ent *riftEntry) charm.State {
	return charm.Step(Contents(ent),
		charm.Self("after entry", func(afterEntry charm.State, r rune) (ret charm.State) {
			switch r {
			case Newline:
				ret = NextIndent(ent.popToIndent)
			}
			return
		}))
}

// contents appear after a collection marker
// [-] <inline spaces> ( <inline comment> | <value> )
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
			if ent.doc.Col >= ent.depth {
				ret = ReadComment(&ent.buffer, contents)
			}

		case Newline:
			ret = NextIndent(func() (ret charm.State) {
				if at := ent.doc.Col; at >= ent.depth {
					ret = BufferRegion(ent, at)
				} else {
					ret = ent.popToIndent()
				}
				return
			})

		default:
			if ent.doc.Col >= ent.depth {
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
					switch at := ent.doc.Col; {
					case at == depth:
						// at the same indent, stick where we're at.
						ret = buffering
					case ent.bufferedLines == 1 && at > depth:
						// the ideal multiline comment a single comment followed by indented lines
						ret = IndentedComment(ent, at)
					default:
						ret = ent.popToIndent()
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
				if at := ent.doc.Col; at == depth {
					ret = indentedComment
				} else {
					ret = HeaderRegion(ent, at)
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
				if ent.doc.Col == depth {
					ret = headerComment
				} else {
					ret = ent.popToIndent()
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
				if at := ent.doc.Col; at >= ent.depth {
					ret = TrailingComment(ent, at)
				} else {
					ret = ent.popToIndent()
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
			ent.tail.WriteRune(Newline)
			ret = ReadComment(&ent.tail, NextIndent(func() (ret charm.State) {
				if at := ent.doc.Col; at >= ent.depth {
					ret = NestedComment(ent, at)
				} else {
					ret = ent.popToIndent()
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
			ent.tail.WriteRune(Newline)
			ret = ReadComment(&ent.tail, NextIndent(func() (ret charm.State) {
				if ent.doc.Col == ent.depth {
					ret = nested
				} else {
					ret = ent.popToIndent()
				}
				return
			}))
			return
		}
		return
	})
}
