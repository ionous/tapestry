package weave

import (
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

func (c *Catalog) WriteNames(m mdl.Modeler) error {
	return forEachNoun(c, func(n *ScopedNoun) (err error) {
		{
			const ofs = -1 // aliases are forced first, in order of declaration.
			for i, a := range n.aliases {
				at := n.aliasat[i]
				if e := m.Name(n.domain.name, n.name, a, ofs, at); e != nil {
					err = e
					break
				}
			}
		}
		if err == nil {
			for ofs, name := range n.Names() {
				if e := m.Name(n.domain.name, n.name, name, ofs, n.at); e != nil {
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

func (cat *Catalog) AssertNounKind(opNoun, opKind string) error {
	return cat.Schedule(assert.NounPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		_, noun := d.StripDeterminer(opNoun)
		_, kind := d.StripDeterminer(opKind)

		if name, ok := UniformString(noun); !ok {
			err = InvalidString(opNoun)
		} else if kn, ok := UniformString(kind); !ok {
			err = InvalidString(opKind)
		} else if k, ok := d.findPluralKind(kn); !ok {
			return errutil.Fmt("unknown kind %q for noun %q at %q", opKind, opNoun, at)
		} else {
			if e := cat.writer.Noun(d.name, name, k, at); e != nil {
				err = e
			} else {
				noun := d.EnsureNoun(name, at)
				noun.AddRequirement(k)
				noun.UpdateFriendlyName(opNoun)
			}
		}
		return
	})
}
