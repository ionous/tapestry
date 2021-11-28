package eph

import (
	"errors"
	"strings"

	"github.com/ionous/errutil"
)

const (
	AspectKinds = "aspect"
	RecordKinds = "record"
)

// name of a kind to assembly info
// ready after phase Ancestry
type Kinds map[string]*Kind

type Kind struct {
	name   string
	reqs   Requires // references to ancestors ( at most it can have one direct parent )
	traits []traitDef
	fields []fieldDef
}

type KindError struct {
	Kind, Domain string
	Err          error
}

func (n KindError) Error() string {
	return errutil.Sprint(n.Err, n.Kind, "in", n.Domain)
}
func (n KindError) Unwrap() error {
	return n.Err
}

// we only allow each kind to be given a single parent ( via resolve )
// but we still have to determine what the hierarchy is ( so we reuse the same bits as domain )
func (ks *Kinds) AddKind(k, p string) {
	if *ks == nil {
		*ks = make(Kinds)
	}
	kind, ok := (*ks)[k]
	if !ok {
		kind = &Kind{name: k}
		(*ks)[k] = kind
	}
	if len(p) > 0 {
		kind.reqs.AddRequirement(p)
	}
}

// distill a tree of kinds into a set of names and their hierarchy
func (ks *Kinds) ResolveKinds(out *ResolvedKinds) (err error) {
	for kind, info := range *ks {
		if res, e := info.reqs.Resolve(kind, (*kindOfDependencies)(ks)); e != nil {
			err = errutil.Append(err, e)
		} else if ps := res.Parents(); len(ps) > 1 {
			e := errutil.Fmt("kind %q should have at most one parent (has: %v)", kind, ps)
			err = errutil.Append(err, e)
		} else {
			kinds := res.Ancestors()
			*out = append(*out, ResolvedKind{kind, kinds})
		}
	}
	return
}

// private helper to make the kinds compatible with the DependencyFinder  ( for resolve )
type kindOfDependencies Kinds

func (c kindOfDependencies) GetRequirements(name string) (ret *Requires, okay bool) {
	if k, ok := c[name]; ok {
		ret, okay = &k.reqs, true
	}
	return
}

func (el *EphKinds) Phase() Phase { return AncestryPhase }

func (el *EphKinds) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if singleKind, e := c.Singularize(d.name, strings.TrimSpace(el.Kinds)); e != nil {
		err = e
	} else if newKind, ok := UniformString(singleKind); !ok {
		err = InvalidString(el.Kinds)
	} else if parentKind, ok := UniformString(el.From); !ok {
		err = InvalidString(el.From)
	} else {
		err = addKind(c, d, at, newKind, parentKind)
	}
	return
}

func addKind(c *Catalog, d *Domain, at string, newKind, parentKind string) (err error) {
	if e := c.AddDefinition(d.name, mdl_kind, at, newKind, parentKind); e != nil {
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
