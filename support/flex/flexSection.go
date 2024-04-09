package flex

import (
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

func MakeSection(r Unreader) Section {
	return Section{runes: r, lastError: io.EOF}
}

// even on error returns true;
// callers can see internal errors via ReadRune.
func (k *Section) NextSection() bool {
	var done bool
	if k.lastError == io.EOF {
		// waiting for start?
		if k.line == 0 {
			k.lastError = nil
		} else {
			// try to read to see if it was a fake EOF
			if _, _, e := k.runes.ReadRune(); e != nil {
				k.lastError = e
				done = e == io.EOF
			} else if e := k.runes.UnreadRune(); e != nil {
				k.lastError = e
			} else {
				k.lastError = nil
				k.StartingLine = k.line
			}
		}
	}
	return !done
}

func (k *Section) ReadRune() (r rune, n int, err error) {
	if k.lastError != nil {
		err, k.lastError = k.lastError, nil
	} else if k.dashing > 0 {
		k.dashing--
		r, n = dash, 1
	} else {
		// not dashing, read:
		if r, n, err = k.runes.ReadRune(); err == nil {
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
						k.lastError = io.EOF
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
