package rift

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
)

type Mapping struct {
	doc    *Document
	sig    Signature
	values MapValues
}

// maybe doc is a factory even?
func NewMapping(doc *Document, indent int, writeBack func(vs MapValues) error) charm.State {
	n := &Mapping{doc: doc}
	return doc.PushIndent(indent, n, func() (err error) {
		// see if there was a value-less key in the pipeline
		// ex. "signature:<eof>"
		if sig, e := n.sig.getSignature(); e != nil {
			err = e
		} else {
			if len(sig) > 0 {
				n.values.Add(sig, nil)
			}
			err = writeBack(n.values)
		}
		return
	})
}

func (n *Mapping) NewRune(first rune) charm.State {
	return charm.RunStep(first, &n.sig, charm.Statement("after sig", func(space rune) (ret charm.State) {
		if sig, e := n.sig.getSignature(); e != nil {
			ret = charm.Error(e)
		} else {
			// add a nil placeholder value
			// alt: could trigger the pop() write every time
			// rather than just the last time ( ie. pop before returning the indented state )
			n.values.Add(sig, nil)
			// note: the end of a signature is indicated by a colon and a space;
			// we have to pass that to parseCollection because it requires at least one space.
			ret = charm.RunState(space, parseCollection(n.doc, func(val any) (_ error) {
				n.values[len(n.values)-1].Value = val // = n.values.Append(sig, val)
				return
			}))
		}
		return
	}))
}
