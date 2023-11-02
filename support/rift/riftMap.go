package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

type Mapping struct {
	h      *History
	sig    Signature
	values MapValues
}

// maybe h is a factory even?
func NewMapping(h *History, indent int, writeBack func(vs MapValues) error) charm.State {
	n := &Mapping{h: h}
	return h.PushIndent(indent, n, func() (err error) {
		// see if there was a value-less key in the pipeline
		// ex. "signature:<eof>"
		if sig, e := n.sig.getSignature(); e != nil {
			err = e
		} else {
			if len(sig) > 0 {
				n.values = n.values.Append(sig, nil)
			}
			err = writeBack(n.values)
		}
		return
	})
}

func (n *Mapping) NewRune(first rune) charm.State {
	return charm.RunStep(first, &n.sig, charm.Statement("after sig", func(space rune) charm.State {
		// note: the end of a signature is indicated by a colon and a space;
		// we have to pass that to parseCollection because it requires at least one space.
		return charm.RunState(space, parseCollection(n.h, func(val any) (err error) {
			if sig, e := n.sig.getSignature(); e != nil {
				err = e
			} else if len(sig) == 0 {
				err = errutil.New("missing signature") // this shouldnt be possible
			} else {
				n.values = n.values.Append(sig, val)
			}
			return
		}))
	}))
}
