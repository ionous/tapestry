package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

type Mapping struct {
	h          *History
	currentSig Signature
	values     MapValues
}

// maybe h is a factory even?
func NewMapping(h *History, indent int, writeBack func(vs MapValues) error) charm.State {
	n := &Mapping{h: h}
	return h.PushIndent(indent, n, func() (err error) {
		// see if there was a value-less key in the pipeline
		// ex. "signature:<eof>"
		if sig, e := n.currentSig.getSignature(); e != nil {
			err = e
		} else {
			vs := n.values
			if len(sig) > 0 {
				vs = vs.Append(sig, nil)
			}
			writeBack(vs)
		}
		return
	})
}

func (n *Mapping) NewRune(first rune) charm.State {
	return charm.RunStep(first, &n.currentSig, charm.Statement("after sig", func(r rune) charm.State {
		startingIndent := n.h.CurrentIndent()
		// padding is the space between the dash and any value
		return charm.RunState(r, RequireSpaces("padding", startingIndent, func(padding int) (ret charm.State) {
			switch {
			// if the indent is less or equal,
			// than we're a new line and the value was null
			case padding <= startingIndent:
				ret = n.h.PopIndent(padding)
			default:
				// an increased indentation means a value:
				ret = charm.Step(NewValue(n.h, padding, func(val any) (err error) {
					if sig, e := n.currentSig.getSignature(); e != nil {
						err = e
					} else if len(sig) == 0 {
						err = errutil.New("missing signature") // this shouldnt be possible
					} else {
						n.values = n.values.Append(sig, val)
						n.currentSig = Signature{} // reset
					}
					return
				}), charm.Statement("after value", func(r rune) charm.State {
					// after value we require a newline:
					return charm.RunState(r, RequireLines("tail", padding, func(tail int) (ret charm.State) {
						switch {
						case tail <= startingIndent:
							ret = n.h.PopIndent(tail)
						default:
							// the only valid increases are collections,
							// and after a collection, we should be at a less or equal indent.
							e := errutil.New("invalid indentation")
							ret = charm.Error(e)
						}
						return
					}))
				}))
			}
			return
		}))
	}))
}
