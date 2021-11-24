package eph

import (
	"errors"

	"github.com/ionous/errutil"
)

// a mapping of a kind to its ancestors that can be resolved into a flat list of kinds
type Kinds map[string]*Requires

// we only allow each kind to be given a single parent ( via resolve )
// but we still have to determine what the hierarchy is ( so we reuse the same bits as domain )
func (ks *Kinds) AddKind(k, p string) {
	if *ks == nil {
		*ks = make(Kinds)
	}
	kind, ok := (*ks)[k]
	if !ok {
		kind = new(Requires)
		(*ks)[k] = kind
	}
	if len(p) > 0 {
		kind.AddRequirement(p)
	}
}

// distill a tree of kinds into a set of names and their hierarchy
func (ks *Kinds) ResolveKinds(out *ResolvedKinds) (err error) {
	for kind, deps := range *ks {
		if res, e := deps.Resolve(kind, (*kindFinder)(ks)); e != nil {
			err = errutil.Append(err, e)
		} else if res.NumParents() > 1 {
			e := errutil.Fmt("kind %q should have at most one parent (has: %v)", kind, res.Ancestors(false))
			err = errutil.Append(err, e)
		} else {
			kinds := res.Ancestors(true)
			*out = append(*out, ResolvedKind{kind, kinds})
		}
	}
	return
}

// private helper to make the kinds compatible with the DependencyFinder ( for resolve )
type kindFinder Kinds

func (c kindFinder) GetRequirements(name string) (ret *Requires, okay bool) {
	ret, okay = c[name]
	return
}

func (el *EphKinds) Phase() Phase { return AncestryPhase }

func (el *EphKinds) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if kinds, ok := UniformString(el.Kinds); !ok {
		err = InvalidString(el.Kinds)
	} else if parentKind, ok := UniformString(el.Kind); !ok {
		err = InvalidString(el.Kind)
	} else if newKind, e := c.Singularize(d.name, kinds); e != nil {
		err = e
	} else if e := c.CheckConflict(d.name, mdl_kind, at, newKind, parentKind); e != nil {
		var de DomainError
		var conflict *Conflict
		if !errors.As(e, &de) || !errors.As(de.Err, &conflict) || conflict.Reason != Duplicated {
			err = e
		} else if de.Domain != d.name {
			LogWarning(e) // log duplicated definitions in subsequent domains
		}
	} else {
		d.kinds.AddKind(newKind, parentKind)
	}
	return
}
