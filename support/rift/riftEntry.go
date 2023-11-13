package rift

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// represents the "contents" of an entry
type riftEntry struct {
	doc             *Document
	pendingValue    pendingValue
	addsValue       func(val any, comment string) error
	pad, head, tail CommentBuffer
	bufferedLines   int
	depth           int
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
		pad := ent.pad.String()
		head := ent.head.String()
		tail := ent.tail.String()
		w.WriteString(pad)
		if len(pad) > 0 && len(head) > 0 {
			w.WriteRune(Newline)
		}
		if len(head) > 0 || len(tail) > 0 {
			w.WriteString(head) // padding
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
		ret = ent.head.String()
		ent.head.Reset()
	}
	return
}

// parses contents and loops (by popping) after its done
func ContentsLoop(ent *riftEntry) charm.State {
	return charm.Step(Contents(ent),
		charm.Self("after entry", func(afterEntry charm.State, r rune) (ret charm.State) {
			switch r {
			case Newline:
				ret = NextIndent(ent.doc, nil)
			}
			return
		}))
}

// Content appears after a collection marker:
// ( noting the collection marker for a document is the start of a file )
//
//	  <spaces> <padding comment>
//		 <header comment>
//		 <value> <spaces> <inline comment>
//	  <trailing comment>
func Contents(ent *riftEntry) charm.State {
	return charm.Self("contents", func(contents charm.State, r rune) (ret charm.State) {
		switch r {
		case Space:
			ret = contents
		case Newline: // a blank line with no contents is the header.
			ret = NextIndent(ent.doc, func(at int) (ret charm.State) {
				if at >= ent.depth {
					ret = HeaderRegion(ent, at, NextValue(ent))
				}
				return
			})
		case Hash: // a hash is the entry comment ( aka padding )
			if at := ent.doc.Col; at >= ent.depth {
				ret = ReadPadding(ent, at)
			}
		default:
			if ent.doc.Col >= ent.depth {
				ret = ReadValue(ent, r)
			}
		}
		return
	})
}

// starts with the comment hash
func ReadPadding(ent *riftEntry, depth int) charm.State {
	doc := ent.doc
	return ReadComment(&ent.pad, NextIndent(doc, func(at int) (ret charm.State) {
		switch {
		case at == depth:
			// the same indent means switch to header
			ret = HeaderRegion(ent, depth, NextValue(ent))
		case at > depth:
			// a deeper indent means nesting ( after nesting, use the header at the original depth )
			doc.Push(depth, HeaderRegion(ent, depth, NextValue(ent)))
			ret = NestedComment(doc, &ent.pad)
		}
		return
	}))
}

// we're at the start of ... something
// could be a value or a comment.
// fix: the reason we cant have the caller use "Step"
// is that pop doesnt really. it doesnt break out of the parent / child
func HeaderRegion(ent *riftEntry, depth int, next charm.State) charm.State {
	return charm.Self("header", func(header charm.State, r rune) (ret charm.State) {
		switch r {
		default:
			ret = next.NewRune(r)
		case Hash:
			ret = ReadComment(&ent.head, header)
		case Newline:
			ret = NextIndent(ent.doc, func(at int) (ret charm.State) {
				switch {
				case at == depth:
					ret = ContinueHeader(ent, depth)
				case at > depth:
					ent.doc.Push(depth, NextValue(ent))
					ret = NestedComment(ent.doc, &ent.head)
				}
				return
			})
		}
		return
	})
}

// subsequent lines of the header are all at the value's indent
// keep reading comments at that indent until there is a value.
func ContinueHeader(ent *riftEntry, depth int) charm.State {
	return charm.Self("second header", func(header charm.State, r rune) (ret charm.State) {
		switch r {
		default:
			ret = ReadValue(ent, r)
		case Hash:
			ent.head.WriteRune(Newline)
			ret = ReadComment(&ent.head, header)
		case Newline:
			ret = MaintainIndent(ent.doc, header, depth)
		}
		return
	})
}

// at the start of a rune which might be a value:
// reads that value and any trailing comment describing it.
func ReadValue(ent *riftEntry, r rune) (ret charm.State) {
	// dont bother trying to read a value if it wasn't meant to be.
	if r != Newline && r != Space {
		ret = charm.RunState(r, NextValue(ent))
	}
	return
}

func NextValue(ent *riftEntry) (ret charm.State) {
	return charm.Step(NewValue(ent), InlineComment(ent))
}

// these are comments to the right of a known value.
func InlineComment(ent *riftEntry) (ret charm.State) {
	inlineIndent := -1
	return charm.Self("inline comment", func(loop charm.State, r rune) (ret charm.State) {
		switch r {
		case Space: // eat spaces on the line after the value
			ret = loop
		case Hash: // an inline comment? read it; loop to us to handle the newline.
			inlineIndent = ent.doc.Col
			ret = ReadComment(&ent.tail, loop)
		case Newline: // a newline ( regardless of whether there was a comment )
			ret = NextIndent(ent.doc, func(at int) (ret charm.State) {
				// the trailing comment indent cant be deeper than its inline comment.
				if (at >= ent.depth) && (inlineIndent < 0 || at <= inlineIndent) {
					// when trailing comments are right aligned with the indent comment
					// use nesting, otherwise use normal newlines.
					var sep rune
					if inlineIndent == at {
						sep = Nestline
					} else {
						sep = Newline
					}
					ret = charm.RunState(r, TrailingComment(ent, at, inlineIndent, sep))
				}
				return
			})
		}
		return
	})
}

// an optional comment can appear on the first line after a value
// starts on something other than whitespace
// at the indent we want to stick with.
func TrailingComment(ent *riftEntry, wantIndent, inlineIndent int, sep rune) charm.State {
	return charm.Self("trailing comments", func(loop charm.State, r rune) (ret charm.State) {
		switch r {
		case Hash:
			ent.tail.WriteRune(sep)
			ret = ReadComment(&ent.tail, loop)
		case Newline:
			ret = MaintainIndent(ent.doc, loop, wantIndent)
		}
		return
	})
}
