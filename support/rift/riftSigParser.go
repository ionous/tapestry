package rift

import (
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// parses a dictionary key of ascii words separated by, and terminating with, a colon.
// the words must start with a letter, but can contain spaces and underscores.
// ex. `a:`, `a:b:`, `and:more complex:keys_like_this:`
type SigParser struct {
	// tbd: break into parts? maybe slices of runes?
	runes     charm.Runes
	separator bool
}

// func (p *SigParser) Reset() string {
// 	r := p.runes.String()
// 	p.runes = charm.Runes{}
// 	return r
// }

func (p *SigParser) Signature() (ret string, err error) {
	if !p.separator {
		err = errutil.New("Signatures must end with a colon")
	} else {
		ret = p.runes.String()
	}
	return
}

// first character of the signature must be a letter
func (p *SigParser) NewRune(r rune) (ret charm.State) {
	if unicode.IsLetter(r) {
		ret = p.runes.Accept(r, charm.Statement("sig head", p.body)) // loop...
	}
	return
}

// subsequent characters can be letters or colons
func (p *SigParser) body(r rune) (ret charm.State) {
	okay := false
	if r == SignatureSeparator && !p.separator {
		okay = true
		p.separator = true
	} else if isSigWord(r) {
		okay = true
		p.separator = false
	}
	if okay {
		ret = p.runes.Accept(r, charm.Statement("sig body", p.body)) // loop...
	}
	return
}
