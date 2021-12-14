package eph

import (
	"github.com/ionous/errutil"
)

var NounPhaseActions = PhaseAction{
	Do: func(d *Domain) error {
		_, e := d.ResolveNouns()
		return e
	},
}

func (c *Catalog) WriteNouns(w Writer) error {
	return forEachNoun(c, func(d *Domain, k *ScopedKind, n *ScopedNoun) (err error) {
		return w.Write(mdl_noun, d.name, n.name, k.name, n.at)
	})
}

func (c *Catalog) WriteNames(w Writer) error {
	return forEachNoun(c, func(d *Domain, k *ScopedKind, n *ScopedNoun) (err error) {
		{
			const ofs = -1 // aliases are forced first, in order of declaration.
			for i, a := range n.aliases {
				at := n.aliasat[i]
				if e := w.Write(mdl_name, d.name, n.name, a, ofs, at); e != nil {
					err = e
					break
				}
			}
		}
		if err == nil {
			for ofs, name := range n.Names() {
				if e := w.Write(mdl_name, d.name, n.name, name, ofs, n.at); e != nil {
					err = e
					break
				}
			}
		}
		return
	})
}

func forEachNoun(c *Catalog, it func(*Domain, *ScopedKind, *ScopedNoun) error) (err error) {
	if ns, e := c.ResolveNouns(); e != nil {
		err = e
	} else {
		for _, ndep := range ns {
			n := ndep.Leaf().(*ScopedNoun)
			if k, e := n.Kind(); e != nil {
				err = errutil.Append(err, e)
			} else if e := it(n.domain, k, n); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

type NounError struct {
	Noun string
	Err  error
}

func (n NounError) Error() string {
	return errutil.Sprintf("%v for noun %q", n.Err, n.Noun)
}
func (n NounError) Unwrap() error {
	return n.Err
}

// name of a noun to assembly info
func (op *EphNouns) Phase() Phase { return NounPhase }

// noun, kind
func (op *EphNouns) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if name, ok := UniformString(op.Noun); !ok {
		err = InvalidString(op.Noun)
	} else if kn, ok := UniformString(op.Kind); !ok {
		err = InvalidString(op.Kind)
	} else if k, ok := d.GetPluralKind(kn); !ok {
		err = errutil.New("unknown kind", k)
	} else {
		noun := d.EnsureNoun(name, at)
		// we can only add requirements to the noun in the same domain that it was declared
		// if in a different domain: the nouns have to match up
		if noun.domain == d {
			noun.AddRequirement(k.name)
		} else if !noun.HasAncestor(k.name) {
			err = NounError{name, errutil.Fmt("can't redefine parent as %q", op.Kind)}
		} else {
			e := errutil.New("duplicate noun definition at", at)
			LogWarning(NounError{name, e})
		}
	}
	return
}
