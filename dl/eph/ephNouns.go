package eph

import (
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

// after assembling all the nouns, make sure they can be resolved.
var NounActions = PhaseAction{
	Do: func(d *Domain) error {
		_, e := d.ResolveNouns()
		return e
	},
}

func (c *Catalog) WriteNouns(w Writer) error {
	return forEachNoun(c, func(n *ScopedNoun) (err error) {
		if k, e := n.Kind(); e != nil {
			err = errutil.Append(err, e)
		} else {
			err = w.Write(mdl.Noun, n.domain.name, n.name, k.name, n.at)
		}
		return
	})
}

func (c *Catalog) WriteNames(w Writer) error {
	return forEachNoun(c, func(n *ScopedNoun) (err error) {
		{
			const ofs = -1 // aliases are forced first, in order of declaration.
			for i, a := range n.aliases {
				at := n.aliasat[i]
				if e := w.Write(mdl.Name, n.domain.name, n.name, a, ofs, at); e != nil {
					err = e
					break
				}
			}
		}
		if err == nil {
			for ofs, name := range n.Names() {
				if e := w.Write(mdl.Name, n.domain.name, n.name, name, ofs, n.at); e != nil {
					err = e
					break
				}
			}
		}
		return
	})
}

type nounResolver interface {
	ResolveNouns() (DependencyTable, error)
}

func forEachNoun(c nounResolver, it func(*ScopedNoun) error) (err error) {
	if ns, e := c.ResolveNouns(); e != nil {
		err = e
	} else {
		for _, ndep := range ns {
			n := ndep.Leaf().(*ScopedNoun)
			if e := it(n); e != nil {
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
func (op *EphNouns) Phase() assert.Phase { return assert.NounPhase }

func (op *EphNouns) Weave(k assert.Assertions) (err error) {
	return k.AssertNounKind(op.Noun, op.Kind)
}

// noun, kind
func (op *EphNouns) Assemble(c *Catalog, d *Domain, at string) (err error) {
	_, noun := d.StripDeterminer(op.Noun)
	_, kind := d.StripDeterminer(op.Kind)

	if name, ok := UniformString(noun); !ok {
		err = InvalidString(op.Noun)
	} else if kn, ok := UniformString(kind); !ok {
		err = InvalidString(op.Kind)
	} else if k, ok := d.GetPluralKind(kn); !ok {
		err = errutil.New("unknown kind", op.Kind)
	} else if noun := d.EnsureNoun(name, at); noun.domain == d {
		// we can only add requirements to the noun in the same domain that it was declared
		// if in a different domain: the nouns have to match up
		noun.UpdateFriendlyName(op.Noun)
		noun.AddRequirement(k.name)
	} else if !noun.HasAncestor(k.name) {
		err = NounError{name, errutil.Fmt("can't redefine parent as %q", op.Kind)}
	} else {
		// is this in anyway useful?
		LogWarning(errutil.Sprintf("duplicate noun %s definition at %v", name, at))
	}
	return
}
