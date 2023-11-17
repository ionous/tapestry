package cin

import (
	"strings"
	"unicode"

	"github.com/ionous/errutil"
)

type Signature struct {
	Name   string
	Params []Parameter
}

// helper for debug printing
func (s *Signature) DebugString() string {
	var b strings.Builder
	b.WriteString(s.Name)
	for i, p := range s.Params {
		if i == 0 {
			b.WriteRune(':')
		} else {
			b.WriteRune(',')
		}
		b.WriteString(p.DebugString())
	}
	b.WriteRune(':')
	return b.String()
}

// given a key such as "Command noun:trait:change choice:"
// separate out the command name and parameter labels
func ReadSignature(key string) (ret Signature, err error) {
	var sig sigReader
	if e := sig.readSig(key); e != nil {
		err = e
	} else {
		ret = Signature{sig.cmd, sig.params}
	}
	return
}

type sigReader struct {
	cmd       string
	params    []Parameter // argument names
	currLabel string
	buf       runeBuffer
}

type Parameter struct {
	Label  string
	Choice string // optional
}

func (p *Parameter) DebugString() string {
	var b strings.Builder
	if l := p.Label; len(l) == 0 {
		b.WriteRune('_')
	} else {
		b.WriteString(l)
	}
	if len(p.Choice) > 0 {
		b.WriteRune(' ')
		b.WriteString(p.Choice)
	}
	return b.String()
}

func (p *Parameter) String() string {
	out := p.Label
	if len(p.Choice) > 0 {
		out = out + " " + p.Choice
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
		s.params = append(s.params, Parameter{}) // blank, unlabeled
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
			err = errutil.New("arguments should always end with a separator", s.cmd, s.currLabel)
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
func (s *sigReader) flush() Parameter {
	var out Parameter
	// nothing accumulated? then our parameter is anonymous
	// ( that's totally fine for our first param )
	if str := s.buf.unbuffer(); len(str) > 0 {
		// if we dont have a label yet, str is the label
		if len(s.currLabel) == 0 {
			out.Label = str
		} else {
			// otherwise, str is the choice
			out.Label, out.Choice = s.currLabel, str
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
