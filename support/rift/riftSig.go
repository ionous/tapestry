package rift

import (
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// parses a dictionary key of ascii words separated by, and terminating with, a colon.
// the words must start with a letter, but can contain spaces and underscores.
// ex. `a:`, `a:b:`, `and:more complex:keys_like_this:`
func NewSignature(hist *History, indent int, writeBack func(string) error) charm.State {
	var sig Signature
	return hist.PushIndent(indent, &sig, func() (err error) {
		if str, e := sig.getSignature(); e != nil {
			err = e
		} else {
			err = writeBack(str)
		}
		return
	})
}

type Signature struct {
	out     strings.Builder
	pending []rune
}

// first character of the signature must be a letter
// subsequent characters of words can be letters, numbers, spaces, or "connectors" (underscore)
// colons separate word parts
func (sig *Signature) NewRune(r rune) (ret charm.State) {
	switch {
	case (r == Space || r == Newline) && !sig.isKeyPending():
		break // done

	case unicode.IsLetter(r):
		sig.append(r)
		ret = sig

	case r == SignatureSeparator:
		if sig.isKeyPending() {
			sig.append(r) // the signature includes the separator
			sig.flushWord()
			ret = sig
		}

	case r == Space || r == SignatureConnector || unicode.IsDigit(r):
		sig.append(r)
		ret = sig
	}
	return
}

func (sig *Signature) getSignature() (ret string, err error) {
	if len(sig.pending) > 0 {
		err = errutil.New("Signature must end with a colon")
	} else {
		ret = sig.out.String()
	}
	return
}

func (sig *Signature) isKeyPending() bool {
	return len(sig.pending) > 0
}

func (sig *Signature) flushWord() {
	sig.out.WriteString(string(sig.pending))
	sig.pending = sig.pending[0:0]
}

func (sig *Signature) append(r rune) {
	sig.pending = append(sig.pending, r)
}
