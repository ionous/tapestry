package eph

import (
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

// returns true if newly added
func (d *Domain) AddOpposite(oneWord, otherWord string) (err error) {
	return d.opposites.AddPair(oneWord, otherWord)
}

// use the domain rules ( and hierarchy ) to find the opposite of the passed word
// the way opposites are defined, there can be more than one opposite word for a given singular word.
// in that case, attempts to pick one.
func (d *Domain) FindOpposite(word string) (ret string, err error) {
	if e := VisitTree(d, func(dep Dependency) (err error) {
		scope := dep.(*Domain)
		if other, ok := scope.opposites.FindOpposite(word); ok {
			ret, err = other, Visited
		}
		return
	}); e != nil && e != Visited {
		err = e
	}
	return
}

// while it'd probably be faster to do this while we assemble,
// keep this assembly separate from the writing produces nicer code and tests.
func (c *Catalog) WriteOpposites(w Writer) (err error) {
	if deps, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			d := dep.Leaf().(*Domain)
			for _, p := range d.opposites {
				def := d.GetDefinition(MakeKey("opposite", p.one))
				if e := w.Write(mdl.Opposite, d.name, p.one, p.other, def.at); e != nil {
					err = errutil.Append(err, domainError{d.name, e})
				}
			}
		}
	}
	return
}

func (op *EphOpposites) Phase() assert.Phase { return assert.PluralPhase }

func (op *EphOpposites) Weave(k assert.Assertions) (err error) {
	return k.AssertOpposite(op.Opposite, op.Word)
}

func (ctx *Context) AssertOpposite(opOpposite, opWord string) (err error) {
	d, at := ctx.d, ctx.at
	if oneWord, ok := UniformString(opOpposite); !ok {
		err = InvalidString(opOpposite)
	} else if otherWord, ok := UniformString(opWord); !ok {
		err = InvalidString(opWord)
	} else if ok, e := d.RefineDefinition(MakeKey("opposite", oneWord), at, otherWord); e != nil {
		err = e
	} else if ok {
		err = d.AddOpposite(oneWord, otherWord)
	}
	return
}
