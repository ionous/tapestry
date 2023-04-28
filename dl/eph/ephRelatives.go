package eph

import (
	"sort"

	"git.sr.ht/~ionous/tapestry/imp/assert"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

func (c *Catalog) WritePairs(w Writer) (err error) {
	if _, e := c.ResolveNouns(); e != nil {
		err = e
	} else if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, deps := range ds {
			d := deps.Leaf().(*Domain)
			// pairs are stored by relation name
			// we sort to give some consistency, the order shouldnt really matter.
			names := make([]string, 0, len(d.relatives))
			for k := range d.relatives {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, relName := range names {
				if e := writePairs(w, d, relName, d.relatives[relName].pairs); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

func writePairs(w Writer, d *Domain, relName string, rs []Relative) (err error) {
	// note: we dont have to test the existence of the kinds and nouns, assembly has already done that
	// sometimes, though... its helpful for testing.
	/*if rel, ok := d.GetKind(relName); !ok {
		err = errutil.New("couldnt find relation", relName)
	} else */sort.Slice(rs, func(i, j int) (less bool) {
		a, b := rs[i], rs[j]
		switch {
		case a.firstNoun < b.firstNoun:
			less = true
		case a.firstNoun == b.firstNoun:
			less = a.secondNoun < b.secondNoun
		}
		return
	})
	for _, p := range rs {
		/*if n1, ok := d.GetNoun(p.firstNoun); !ok {
			err = errutil.New("couldnt find first noun", p.firstNoun)
			break
		} else if n1, ok := d.GetNoun(p.secondNoun); !ok {
			err = errutil.New("couldnt find second noun", p.secondNoun)
			break
		} else*/if e := w.Write(mdl.Pair, d.name, relName, p.firstNoun, p.secondNoun, p.at); e != nil {
			err = e
			break
		}
	}
	return
}

type Relatives struct {
	pairs []Relative
}

type Relative struct {
	firstNoun, secondNoun, at string
}

func (rs *Relatives) AddPair(a, b, at string) {
	rs.pairs = append(rs.pairs, Relative{a, b, at})
}

func (op *EphRelatives) Phase() assert.Phase { return assert.RelativePhase }

func (op *EphRelatives) Weave(k assert.Assertions) (err error) {
	return k.AssertRelative(op.Rel, op.Noun, op.OtherNoun)
}

// validate that the pattern for the rule exists then add the rule to the *current* domain
// ( rules are de/activated based on domain, they can be part some child of the domain where the pattern was defined. )
func (ctx *Context) AssertRelative(opRel, opNoun, opOtherNoun string) (err error) {
	d, at := ctx.d, ctx.at
	if name, ok := UniformString(opRel); !ok {
		err = InvalidString(opRel)
	} else if rel, ok := d.GetPluralKind(name); !ok || !rel.HasAncestor(kindsOf.Relation) {
		err = errutil.Fmt("unknown or invalid relation %q", opRel)
	} else if card := rel.domain.GetDefinition(MakeKey("rel", rel.name, "card")); len(card.value) == 0 {
		err = errutil.Fmt("unknown or invalid cardinality for %q", opRel)
	} else if first, e := getClosestNoun(d, opNoun); e != nil {
		err = e
	} else if second, e := getClosestNoun(d, opOtherNoun); e != nil {
		err = e
	} else {
		var addPair bool
		switch card.value {
		case tables.ONE_TO_ONE:
			// when one-to-one, the meaning of the two columns is the same
			// and sorting the names so that first is less than second simplifies testing for uniqueness
			if first.name > second.name {
				first, second = second, first
			}
			addPair, err = relate(d, rel, first.name, at, second.name)

		case tables.ONE_TO_MANY:
			// one parent to many children; so given second noun ( a child ) there is only one valid first noun ( a parent )
			addPair, err = relate(d, rel, second.name, at, first.name)

		case tables.MANY_TO_ONE:
			// many children to one parent; so given first noun ( a child ) there is only one valid second noun( a parent )
			addPair, err = relate(d, rel, first.name, at, second.name)

		case tables.MANY_TO_MANY:
			uniquePair := first.name + second.name
			addPair, err = relate(d, rel, uniquePair, at, uniquePair)
		default:
			err = errutil.Fmt("unknown or invalid cardinality %q for %q", card.value, opRel)
		}
		//
		if err == nil && addPair {
			if d.relatives == nil {
				d.relatives = make(map[string]Relatives)
			}
			pairs := d.relatives[rel.name]
			pairs.AddPair(first.name, second.name, at)
			d.relatives[rel.name] = pairs
		}
	}
	return
}

func relate(d *Domain, rel *ScopedKind, key, at, other string) (okay bool, err error) {
	return d.RefineDefinition(MakeKey("rel", rel.name, key), at, other)
}
