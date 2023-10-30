package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
)

// this is the right hand side after a value
type TailParser struct {
	indent    int
	lineCount int
}

func (p *TailParser) StateName() string {
	return "tail"
}

func (p *TailParser) GetTail() (retDepth, retLines int) {
	retDepth, retLines = p.indent, p.lineCount
	return
}

// first character of the signature must be a letter
func (p *TailParser) NewRune(r rune) (ret charm.State) {
	// fix: support cr/lf with a pending newline accumulator
	if r == '\n' {
		p.lineCount++
		p.indent = 0
		ret = p
	} else if r == ' ' {
		p.indent++
		ret = p
	}
	return
}
