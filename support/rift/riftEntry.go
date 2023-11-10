package rift

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// represents the "contents" of an entry
type riftEntry struct {
	Collection
	pendingValue   pendingValue
	buffer, header strings.Builder
	bufferedLines  int
	depth          int
}

func (ent *riftEntry) finalizeEntry() (err error) {
	c := ent.Collection
	if val, e := ent.pendingValue.FinalizeValue(); e != nil {
		err = e
	} else {
		c.Comments().WriteString(ent.buffer.String())
		c.Comments().WriteString(ent.header.String())
		// fix: modify the signature to write the comment at the same time?
		err = c.WriteValue(val)
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
				ret = NextIndent(doc.Pop)
			}
			return
		}))
}

// the contents exist after a collection marker:
//
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
			if doc := ent.Document(); doc.Col >= ent.depth {
				ret = ReadComment(ent.Comments(), contents)
			}
		case Newline:
			ret = NextIndent(func() (ret charm.State) {
				if doc := ent.Document(); doc.Col >= ent.depth {
					ret = BufferRegion(ent, doc.Col)
				} else {
					ret = doc.Pop()
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
						ret = doc.Pop()
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
					ret = doc.Pop()
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
			charm.Self("trailing comments", func(tail charm.State, r rune) (ret charm.State) {
				// fix ... need to implement this:
				switch r {
				case Space:
					ret = tail
				}
				return
			})))
	}
	return
}
