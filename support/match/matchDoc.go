package match

import (
	"io"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/tell/charm"
	"github.com/ionous/tell/runes"
)

// handle the parsed document.
// the document data also includes the unprocessed content which ended the document.
// ( ex. deindentation )
type AfterDocument func(AsyncDoc) charm.State

// reads a document via channels
// ( which allows reading a (sub) document to become a state in a larger document )
type AsyncDoc struct {
	// the final document ( or error if file.ReadTellRunes failed )
	Content any
	indent  mismatchedIndent
}

// Sub-documents are defined by their indentation level.
// And, on each new line they have to collect enough whitespace
// to determine whether the line is part of their content.
// If the line has a lesser indent, the doc ends,
// but it still has the whitespace it collected (which the parent doc needs.)
// ParseUnhandledContent() sends that whitespace to the passed state.
func (doc AsyncDoc) ParseUnhandledContent(n charm.State) charm.State {
	for range doc.indent.spaces {
		n = n.NewRune(runes.Space)
	}
	return n.NewRune(doc.indent.unhandled)
}

// ofs is line offset for error reporting
func decodeDoc(includeComments bool, lineOfs int, notify AfterDocument) charm.State {
	var indent int
	return charm.Step(
		// determine the indentation of the first line of the tell document
		findIndent(2, &indent),
		charm.Statement("startSubDoc", func(q rune) (ret charm.State) {
			// async routine pulls runes from the 'runes' channel.
			out := make(chan AsyncDoc)
			async := tellDocDecoder{
				out:    out,
				runes:  newAsyncDoc(out, lineOfs, indent, includeComments),
				notify: notify,
			}
			return async.NewRune(q)
		}))
}

type tellDocDecoder struct {
	out    chan AsyncDoc
	runes  chan rune
	notify AfterDocument
}

// send runes to the document
func (n *tellDocDecoder) NewRune(q rune) (ret charm.State) {
	select {
	// check if the previous rune ended the document
	// (ex. on error, or on a return to an earlier indent level )
	case res := <-n.out:
		// the result has the *previous* rune, but not the current rune.
		// after handling all the collected data, send the new state the new rune.
		if next := n.notify(res); next != nil {
			ret = next.NewRune(q)
		}

	// send the rune to the background tell reader.
	case n.runes <- q:
		if q != charm.Eof {
			// keep looping after every normal rune.
			ret = n
		} else {
			// if we've sent the eof: wait until the async document has finished.
			// ( it will have the eof we just sent in its content,
			// - so we don't need to manually send it to the next state. )
			res := <-n.out
			ret = n.notify(res)
		}
	}
	return
}

// the existing tell document reader expects a "RuneReader"
// it wants to pull values at its own pace.
// however, the charm states only get runes one by one.
// so, this creates a channel that all of the runes can post to.
// assumes we start already indented to 'indent'
func newAsyncDoc(out chan<- AsyncDoc, ofs, indent int, includeComments bool) chan rune {
	in := channelReader{
		indent: indent,
		runes:  make(chan rune),
		// spaces is zero, because we start at the right indentation
	}
	go func() {
		if content, e := files.ReadTellRunes(&in, files.Ofs{Line: ofs}, includeComments); e != nil {
			out <- AsyncDoc{Content: e}
		} else {
			out <- AsyncDoc{Content: content, indent: in.ending}
		}
	}()
	return in.runes
}

// a rune reader that pulls from a channel.
// uses a -1 rune to indicate eof.
type channelReader struct {
	indent int
	runes  chan rune
	ending mismatchedIndent
	spaces int // ranges from indent down to 0
}

func (n *channelReader) ReadRune() (ret rune, size int, err error) {
	for {
		// read from the channel;
		// close happens only if ReadRune() returns an error ( ex. EOF )
		// ( which causes files.ReadTellRunes to fail, which causes newAsyncDoc to return )
		q := <-n.runes
		if n.ending.unhandled != 0 {
			// this should never happen, and is unrecoverable.
			// after we set n.ending below, files.ReadTellRunes and newAsyncDoc should exit.
			panic("unexpectedly reading after failure")
		} else {
			// start matching indent after newlines by continuing to read from the channel
			// ( even in we haven't reached the proper indent. )
			if q == runes.Newline {
				n.spaces = n.indent
				ret, size = q, utf8.RuneLen(q)
				break
			} else if n.spaces == 0 {
				// after we have eaten enough spaces
				// return the rune as is.
				ret, size = q, utf8.RuneLen(q)
				break
			} else if q != runes.Space {
				// not enough spaces is an error... or the end of the document.
				// we can't know which.
				n.ending = mismatchedIndent{n.indent - n.spaces, q}
				err = io.EOF
				break
			} else {
				n.spaces--
				continue
			}
		}
	}
	return
}

// a series of left aligned spaces followed by some unexpected character
// ( the character was expected to be another space, but wasn't )
type mismatchedIndent struct {
	spaces    int
	unhandled rune
}

func (m mismatchedIndent) Error() string {
	return "mismatched indent"
}
