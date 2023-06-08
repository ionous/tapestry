package weave

import (
	"database/sql"
	"sort"

	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

func (c *Catalog) WritePairs(m mdl.Modeler) (err error) {
	if _, e := c.ResolveNouns(); e != nil {
		err = e
	} else if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, d := range ds {
			// pairs are stored by relation name
			// we sort to give some consistency, the order shouldnt really matter.
			names := make([]string, 0, len(d.relatives))
			for k := range d.relatives {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, relName := range names {
				if e := writePairs(m, d, relName, d.relatives[relName].pairs); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

func writePairs(m mdl.Modeler, d *Domain, relName string, rs []Relative) (err error) {
	sort.Slice(rs, func(i, j int) (less bool) {
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
		if e := m.Pair(d.name, relName, p.firstNoun, p.secondNoun, p.at); e != nil {
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

func (d *Domain) findCardinality(rel string) (ret string, err error) {
	q := d.catalog.db
	if e := q.QueryRow(`
	select cardinality
	from mdl_rel mr 
	join mdl_kind mk
		on (mr.relKind = mk.rowid)
	join domain_tree
		on (uses = domain)
	where base = ?1
	and kind = ?2
	limit 1
	`, d.name, rel).Scan(&ret); e == sql.ErrNoRows {
		err = errutil.Fmt("unknown or invalid cardinality for %q", rel)
	} else {
		err = e
	}
	return
}

// validate that the pattern for the rule exists then add the rule to the *current* domain
// ( rules are de/activated based on domain, they can be part some child of the domain where the pattern was defined. )
func (cat *Catalog) AssertRelative(opRel, opNoun, opOtherNoun string) error {
	return cat.Schedule(assert.RelativePhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if noun, ok := UniformString(opNoun); !ok {
			err = InvalidString(opNoun)
		} else if otherNoun, ok := UniformString(opOtherNoun); !ok {
			err = InvalidString(opOtherNoun)
		} else if name, ok := UniformString(opRel); !ok {
			err = InvalidString(opRel)
		} else if rel, ok := d.findPluralKind(name); !ok {
			err = errutil.Fmt("unknown or invalid relation %q", opRel)
		} else if card, e := d.findCardinality(rel); e != nil {
			err = e
		} else if first, e := d.GetClosestNoun(noun); e != nil {
			err = e
		} else if second, e := d.GetClosestNoun(otherNoun); e != nil {
			err = e
		} else {
			var addPair bool
			switch card {
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
				err = errutil.Fmt("invalid cardinality %q for %q", card, opRel)
			}
			//
			if err == nil && addPair {
				if d.relatives == nil {
					d.relatives = make(map[string]Relatives)
				}
				pairs := d.relatives[rel]
				pairs.AddPair(first.name, second.name, at)
				d.relatives[rel] = pairs
			}
		}
		return
	})
}

func relate(d *Domain, rel, key, at, other string) (okay bool, err error) {
	return d.RefineDefinition(MakeKey("rel", rel, key), at, other)
}
