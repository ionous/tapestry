package cin

import (
	"unicode"

	"github.com/ionous/errutil"
)

type sigReader struct {
	cmd       string
	params    []sigParam // argument names
	currLabel string
	buf       runeBuffer
}

type sigParam struct {
	label  string
	choice string // optional
}

func (p *sigParam) String() string {
	out := p.label
	if len(p.choice) > 0 {
		out = out + " " + p.choice
	}
	return out
}

type runeBuffer []rune

const eof = rune(-1)

func (s *sigReader) readSig(str string) (err error) {
	for _, r := range str {
		if e := s.readRune(r); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		err = s.readRune(eof)
	}
	return
}

func (s *sigReader) readRune(r rune) (err error) {
	if len(s.cmd) == 0 {
		s.readCmd(r)
	} else {
		err = s.readParam(r)
	}
	return
}

func (s *sigReader) readCmd(r rune) {
	switch {
	// commands ending with a colon indicate an initial anonymous argument
	case r == ':':
		s.params = append(s.params, sigParam{}) // blank, unlabeled
		fallthrough
	// a space is used to separate a command from its arguments
	// ( and an immediate end of input means there are no arguments )
	case r == ' ' || r == eof:
		s.cmd = s.buf.unbuffer() // we have a full command name now.
	default:
		s.buf.addRune(r)
	}
}

func (s *sigReader) readParam(r rune) (err error) {
	switch {
	case r == eof:
		if len(s.buf)+len(s.currLabel) != 0 {
			err = errutil.New("arguments should always end with a separator")
		}
	// a colon indicates the end of an argument name
	case r == ':':
		s.params = append(s.params, s.flush())
	// spaces in argument names indicate choices ( for an embedded swap value )
	case r == ' ':
		if len(s.currLabel) > 0 {
			err = errutil.New("spaces in argument names indicate choices, and there should be at most one choice")
		} else {
			// remember the label, and start accumulating the choice.
			s.currLabel = s.buf.unbuffer()
		}
	default:
		s.buf.addRune(r)
	}
	return
}

// return (and reset) the pending argument's accumulated label and choice ( if any )
func (s *sigReader) flush() sigParam {
	var out sigParam
	// nothing accumulated? then our parameter is anonymous
	// ( that's totally fine for our first param )
	if str := s.buf.unbuffer(); len(str) > 0 {
		// if we dont have a label yet, str is the label
		if len(s.currLabel) == 0 {
			out.label = str
		} else {
			// otherwise, str is the choice
			out.label, out.choice = s.currLabel, str
			s.currLabel = ""
		}
	}
	return out
}

// camelCase to break_case helper.
func (b *runeBuffer) addRune(r rune) {
	if l := unicode.ToLower(r); len(*b) > 0 && l != r {
		*b = append(*b, '_', l)
	} else {
		*b = append(*b, l)
	}
}

func (b *runeBuffer) unbuffer() string {
	out := string(*b)
	(*b) = (*b)[:0] // reset w/o realloc
	return out
}
