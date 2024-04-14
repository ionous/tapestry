package flex

import (
	"io"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/support/files"
)

// a error, a mapping, or sequence
type subDoc any

// read a tell document from a channel.
// this posts the final document, or any error, to out;
// returns a channel to put document runes into.
func newAsyncDoc(out chan<- subDoc, includeComments bool) chan<- rune {
	in := make(ReverseReader)
	go func() {
		if content, e := files.ReadTellRunes(in, includeComments); e != nil {
			out <- e
		} else {
			out <- content
		}
	}()
	return in
}

// a rune reader that pulls from a channel.
// uses a -1 rune to indicate eof.
type ReverseReader chan rune

func (c ReverseReader) ReadRune() (ret rune, size int, err error) {
	if q, ok := <-c; q == -1 || !ok {
		err = io.EOF
	} else {
		ret, size = q, utf8.RuneLen(q)
	}
	return
}
