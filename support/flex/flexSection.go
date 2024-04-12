package flex

import (
	"errors"
	"io"
)

// read a structured block until it hits a structured ending.
type Section struct {
	StartingLine int // newline count at start of section
	line         int // track total newlines

	// >0: a number dashes need to be returned to the caller
	// =0: at the start of a line where dashing might occur
	// <0: in the middle of a line, and runes should be handed back transparently.
	dashing   int
	runes     Unreader
	lastError error
}

type Unreader interface {
	io.RuneReader
	UnreadRune() error
}

// return a reader that ends and restarts
// on every dashed divider in a flex document.
// must call "NextSection" to start reading.
func MakeSection(r Unreader) Section {
	return Section{runes: r, lastError: nextSection}
}

var nextSection = errors.New("a call to NextSection was expected")

// valid at the start of a document or the
// after ReadRune() has returned eof; otherwise, panics.
// returns false at the end of a document.
func (k *Section) NextSection() bool {
	var done bool
	if k.lastError != nextSection {
		panic("unexpected call to NextSection")
	} else {
		// waiting for start?
		if k.line == 0 {
			k.lastError = nil
		} else {
			// see if the last ReadRune() returned a real EOF
			if _, _, e := k.runes.ReadRune(); e != nil {
				k.lastError = e    // the document is probably
				done = e == io.EOF // done.
			} else if e := k.runes.UnreadRune(); e != nil {
				k.lastError = e
			} else {
				// it was a faux eof, allow ReadRune() to proceeed.
				k.lastError = nil
				k.StartingLine = k.line - 1
			}
		}
	}
	return !done
}

// read the next rune unless the next rune is a dash.
// this buffers the dashes, counting them to search for an end of section
func (k *Section) ReadRune() (r rune, n int, err error) {
	if k.lastError != nil {
		err, k.lastError = k.lastError, nil
	} else if k.dashing > 0 {
		k.dashing--
		r, n = dash, 1
	} else {
		// not handing out buffered dashes then read; and check the results.
		r, n, err = k.runes.ReadRune()
		switch err {
		case io.EOF:
			// ugly: if its a real eof; we need to allow calls to NextSection
			k.lastError = nextSection
		case nil:
			switch {
			// always check for newlines
			// resets checks for dashes
			case r == newline:
				k.dashing = 0
				k.line++

			// not at start of line; no special handling.
			case k.dashing < 0:
				break

			// at start of line, but not a dash:
			case r != dash:
				k.dashing = -1

			// a dash and at the start of a line:
			default:
				// count how many more dashes there are
				if w, cnt, e := k.readDashes(); e != nil && e != io.EOF {
					err = e
				} else {
					// have exactly the expected dashes and an eol or eof?
					exact := (cnt == dashCount-1) && (e == io.EOF ||
						(e == nil && w == newline))
					if !exact {
						// not exactly the dashes we expect
						// return the initial dash
						// later, return the buffered ones.
						if e == nil {
							err = k.runes.UnreadRune()
						}
						k.dashing = cnt
					} else {
						// our dash count was exact
						// return eof to indicate end of section
						if w == newline {
							k.line++
						}
						k.lastError = nextSection
						err = io.EOF
					}
				}
			}
		}
	}
	return
}

// read until error, or out of dashes.
func (k *Section) readDashes() (w rune, ret int, err error) {
	for {
		if r, _, e := k.runes.ReadRune(); e != nil {
			err = e
			break
		} else {
			if r == dash {
				ret++
			} else {
				w = r
				break
			}
		}
	}
	return
}

const (
	newline    = '\n'
	dash       = '-'
	dashCount  = 3
	dashString = "---"
	space      = ' '
)
