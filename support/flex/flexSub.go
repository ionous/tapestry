package flex

import (
	"github.com/ionous/tell/charm"
	"github.com/ionous/tell/runes"
)

// hands the rune that ended the document, plus the content
type AfterDocument func(q rune, content any) charm.State

func DecodeSubDoc(notify AfterDocument, includeComments bool) (ret charm.State) {
	var indent int
	return charm.Step(
		// determine the indentation of the first line of the tell document
		findIndent(2, &indent),
		charm.Statement("startSubDoc", func(q rune) (ret charm.State) {
			// async routine receives runes via runes
			// and posts final results to out
			out := make(chan subDoc)
			runes := newAsyncDoc(out, includeComments)
			dec := subDocDecoder{
				out:    out,
				runes:  runes,
				indent: indent,
				notify: notify,
			}
			return dec.NewRune(q)
		}))
}

type subDocDecoder struct {
	out    chan subDoc
	runes  chan<- rune
	indent int
	notify AfterDocument
}

// send the pending document and unhandled rune to the after document handler
func (n *subDocDecoder) finishDoc(q rune) (ret charm.State) {
	close(n.runes)
	return n.notify(q, <-n.out)
}

// send runes to the document
func (n *subDocDecoder) NewRune(q rune) (ret charm.State) {
	switch q {
	case runes.Eof:
		ret = n.finishDoc(q)

	default:
		select {
		// check if the *last* rune ended the document
		// ( ex. with an error )
		case content := <-n.out:
			ret = n.notify(q, content)

		// or, send the new rune into the reader
		case n.runes <- q:
			// if it was a newline, on the next line,
			// we want to eat whitespace until we match the original indent
			if q != runes.Newline {
				ret = n
			} else {
				ret = n.matchIndent()
			}
		}
	}
	return
}

// assume we're just past a newline
// waits until we've reached an equal indent then passes control to after;
// treats everything other than a space (or newline) as unhandled.
// assumes indent is at least 1.
func (n *subDocDecoder) matchIndent() charm.State {
	var spaces int
	return charm.Self("matchIndent", func(matchIndent charm.State, q rune) (ret charm.State) {
		switch q {
		default: // unhandled, use whatever doc data we have
			ret = n.finishDoc(q)
		case runes.Newline:
			spaces = 0
			ret = matchIndent
		case runes.Space:
			if spaces = spaces + 1; spaces < n.indent {
				ret = matchIndent // keep reading spaces
			} else {
				ret = n // matchedIndent, decode more of the doc
			}
		}
		return
	})
}
