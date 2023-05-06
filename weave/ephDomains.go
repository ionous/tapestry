package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/weave/assert"

	"github.com/ionous/errutil"
)

type DomainFinder interface {
	GetDomain(n string) (*Domain, bool)
}

type Domain struct {
	Requires      // other domains this needs ( can have multiple direct parents )
	catalog       *Catalog
	defs          Artifacts
	kinds         ScopedKinds // the kinds declared in this domain
	nouns         ScopedNouns
	checks        asmChecks
	resolvedKinds cachedTable
	resolvedNouns cachedTable
	rules         map[string]Rulesets  // pattern name to rules for that pattern
	relatives     map[string]Relatives // relation name to pairs of nouns

	// a domain that's fully processed will be in some final state
	currPhase assert.Phase

	// scheduled assertions
	scheduling [assert.NumPhases][]memento
}

type memento struct {
	cb func(*Weaver) error
	at string
}

func (op *memento) call(ctx *Weaver) error {
	ctx.at = op.at
	return op.cb(ctx)
}

func (cat *Catalog) Schedule(when assert.Phase, what func(*Weaver) error) (err error) {
	if d, ok := cat.processing.Top(); !ok {
		err = errutil.New("unknown top level domain")
	} else {
		err = d.schedule(cat.cursor, when, what)
	}
	return
}

func (d *Domain) schedule(at string, when assert.Phase, what func(*Weaver) error) (err error) {
	if d.currPhase > when {
		err = errutil.Fmt("unexpected phase(%s) for %q", when, d.name)
	} else /*if when == d.currPhase {
		err = what(&ctx)
	} else */{
		d.scheduling[when] = append(d.scheduling[when], memento{what, at})
	}
	return
}

// return the domain hierarchy: the ancestors ending just before the domain itself.
// direct parents may not be contiguous ( depending on whether their ancestors overlap. )
func (d *Domain) Resolve() (ret Dependencies, err error) {
	if len(d.at) == 0 {
		err = domainError{d.name, errutil.New("never defined", d.name)}
	} else if ds, e := d.resolve(d, (*catDependencyFinder)(d.catalog)); e != nil {
		err = domainError{d.name, e}
	} else {
		ret = ds
	}
	return
}

func (d *Domain) GetDefinition(key keyType) Definition {
	return d.defs[key.hash]
}

// used by ephemera during assembly to record some piece of information
// that would cause problems it were specified differently elsewhere.
// ex. an in-game password specified as the word "secret" in one place, but "mongoose" somewhere else.
func (d *Domain) AddDefinition(key keyType, at, value string) (err error) {
	if e := VisitTree(d, func(dep Dependency) (err error) {
		scope := dep.(*Domain)
		if e := scope.defs.CheckConflict(key, value); e != nil {
			err = domainError{scope.name, e}
		}
		return
	}); e != nil {
		err = e
	} else {
		if d.defs == nil {
			d.defs = make(Artifacts)
		}
		d.defs[key.hash] = Definition{key: key, at: at, value: value}
	}
	return
}

// add a definition that can be overridden in subsequent domains.
// returns "okay" if the refinement was added ( ex. not duplicated )
func (d *Domain) RefineDefinition(key keyType, at, value string) (okay bool, err error) {
	var de domainError
	var conflict *Conflict
	if e := d.AddDefinition(key, at, value); e == nil {
		okay = true
	} else if !errors.As(e, &de) || !errors.As(de.Err, &conflict) {
		err = e // some unknown error?
	} else {
		switch conflict.Reason {
		case Redefined:
			// redefined definitions are only a problem in the same domain.
			// ( ie. we allow subdomains to reset / override the plurals )
			if d.name == de.Domain {
				err = e
			} else {
				okay = true
				// FIX! see Domain.AddDefinition
				// the earlier "AddDefinition" doesnt actually add it because this is a redefinition
				// *but* we actually do want that information....
				if d.defs == nil {
					d.defs = make(Artifacts)
				}
				d.defs[key.hash] = Definition{key: key, at: at, value: value}
				LogWarning(e) // even though its okay, let the user know.
			}
		case Duplicated:
			// duplicated definitions are all okay;
			// but if its in a derived domain: let the user know.
			if de.Domain != d.name {
				LogWarning(e)
			}
		default:
			err = e // some unknown conflict?
		}
	}
	return
}

// the domain should have been resolved already.
func (d *Domain) runPhase(ctx *Weaver) (err error) {
	if ds, e := d.Resolve(); e != nil {
		err = e
	} else {
		w := ctx.phase
		d.currPhase = w // hrmm
		// note: even if there are no ephemera... in a given phase..
		// there can still be rivals and other results to process
		allowDupes := w != assert.AncestryPhase
		if e := d.checkRivals(ds, allowDupes); e != nil {
			err = e
		} else {
			// don't "range" over the phase data since the contents can change during traversal.
			// tbd: have "Schedule" immediately execute the statement if in the correct phase?
			for len(d.scheduling[w]) > 0 {
				els := d.scheduling[w]
				next := els[0]
				d.scheduling[w] = els[1:]
				if e := next.call(ctx); e != nil {
					err = errutil.Append(err, e)
				}
			}
			d.currPhase++
		}
	}
	return
}

// used by assembler to check that domains with multiple parents don't contain conflicting information.
// ex. "plane: a flying vehicle" and "plane: a woodworking tool" both included by some child domain.
//
// NOTE: we're not copying single parents, or even remembering this merged set:
// instead we're crawling up the conflicts in AddDefinition...
func (d *Domain) checkRivals(ds Dependencies, allowDupes bool) (err error) {
	// start with nothing and merge in to check for artifacts
	if deps := ds.Parents(); len(deps) > 1 {
		a := make(Artifacts)
		for _, dep := range deps {
			scope := dep.(*Domain) // allow this to panic
			if defs := scope.defs; len(defs) > 0 {
				if e := a.Merge(defs, allowDupes); e != nil {
					err = domainError{scope.name, e}
					break
				}
			}
		}
	}
	return
}

func (d *Domain) addDependencies(reqs []string) (err error) {
	if d.status == xResolved {
		err = errutil.Fmt("domain %q already resolved, can't add new dependencies", d.name)
	} else {
		for _, req := range reqs {
			if d.name == req {
				err = errutil.Fmt("circular reference: %s can't depend on itself", req)
				break
			} else {
				d.Requires.AddRequirement(req)
			}
		}
	}
	return
}
