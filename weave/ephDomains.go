package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

type DomainFinder interface {
	GetDomain(n string) (*Domain, bool)
}

type Domain struct {
	name          string
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

// have all parent domains been processed?
func (d *Domain) isReadyForProcessing() bool {
	return nil == d.visit(func(uses *Domain) (err error) {
		if d != uses && uses.currPhase < TempSplit {
			err = errutil.New("break")
		}
		return
	})
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
func (d *Domain) Resolve() (ret []string, err error) {
	c := d.catalog // we shouldnt have to worry about dupes, because in theory we didnt add them.
	if rows, e := c.db.Query(`select uses from domain_tree 
		where base = ?1 order by dist desc`, d.name); e != nil {
		err = e
	} else {
		ret, err = tables.ScanStrings(rows)
	}
	return
}

func (d *Domain) GetDefinition(key keyType) Definition {
	return d.defs[key.hash]
}

// temp... hopefullly.
func (d *Domain) visit(visit func(d *Domain) error) (err error) {
	cat := d.catalog
	if tree, e := d.Resolve(); e != nil {
		err = e
	} else {
		for _, el := range tree {
			if p, ok := cat.GetDomain(el); !ok {
				err = errutil.Fmt("unexpected domain %q", el)
				break
			} else if e := visit(p); e != nil {
				err = e
				break
			}
		}
	}
	return
}

// used by ephemera during assembly to record some piece of information
// that would cause problems it were specified differently elsewhere.
// ex. an in-game password specified as the word "secret" in one place, but "mongoose" somewhere else.
func (d *Domain) AddDefinition(key keyType, at, value string) (err error) {
	if e := d.visit(func(scope *Domain) (err error) {
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

func (d *Domain) runPhase(ctx *Weaver) (err error) {
	w := ctx.phase
	d.currPhase = w // hrmm
	// note: even if there are no ephemera... in a given phase..
	// there can still be rivals and other results to process
	allowDupes := w != assert.AncestryPhase
	if e := d.checkRivals(allowDupes); e != nil {
		err = e
	} else {
		// don't "range" over the phase data since the contents can change during traversal.
		// tbd: have "Schedule" immediately execute the statement if in the correct phase?
		for len(d.scheduling[w]) > 0 {
			els := d.scheduling[w]
			next := els[0]
			d.scheduling[w] = els[1:]

			if e := next.call(ctx); e != nil {
				if !errors.Is(e, mdl.Duplicate) {
					err = errutil.Append(err, e)
				} else if d.catalog.warn != nil {
					d.catalog.warn(e)
				}
			}
		}
		d.currPhase++
	}
	return
}

// used by assembler to check that domains with multiple parents don't contain conflicting information.
// ex. "plane: a flying vehicle" and "plane: a woodworking tool" both included by some child domain.
//
// NOTE: we're not copying single parents, or even remembering this merged set:
// instead we're crawling up the conflicts in AddDefinition...
func (d *Domain) checkRivals(allowDupes bool) (err error) {
	// start with nothing and merge in to check for artifacts
	a := make(Artifacts)
	return d.visit(func(scope *Domain) (err error) {
		if scope != d {
			if defs := scope.defs; len(defs) > 0 {
				if e := a.Merge(defs, allowDupes, func(warn error) {
					LogWarning(domainError{scope.name, warn})
				}); e != nil {
					err = domainError{scope.name, e}
				}
			}
		}
		return
	})
}
