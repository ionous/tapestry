package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// eats ascii whitespace, tracking indent
// in yaml, spaces are indents and after indentation tabs are allowed.
// that seems mildly interesting for end of line alignment of comments
// but sticking to no tabs at all seems even better.
type Whitespace struct {
	Indent   int
	Lines    int
	optional bool // by default whitespace is required.
}

func OptionalWhitespace() *Whitespace {
	return &Whitespace{optional: true}
}

func (p *Whitespace) StateName() string {
	return "whitespace"
}

func (p *Whitespace) IsEmpty() bool {
	return p.Indent+p.Lines == 0
}

func (p *Whitespace) GetSpacing() (retDepth, retLines int) {
	retDepth, retLines = p.Indent, p.Lines
	return
}

// first character of the signature must be a letter
func (p *Whitespace) NewRune(r rune) (ret charm.State) {
	if r == charm.Eof {
		p.Lines++ // treat it as an new
	} else if r == Newline {
		p.Lines++
		p.Indent = 0
		ret = p
	} else if r == Space {
		p.Indent++
		ret = p
	} else if !p.optional && r != charm.Eof && p.IsEmpty() {
		e := errutil.New("expected whitespace")
		ret = charm.Error(e)
	}
	return
}
