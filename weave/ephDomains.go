package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/weave/assert"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

type DomainFinder interface {
	GetDomain(n string) (*Domain, bool)
}

type Domain struct {
	Requires      // other domains this needs ( can have multiple direct parents )
	catalog       *Catalog
	defs          Artifacts
	kinds         ScopedKinds
	nouns         ScopedNouns
	checks        asmChecks
	resolvedKinds cachedTable
	resolvedNouns cachedTable
	opposites     OppositePairs
	rules         map[string]Rulesets  // pattern name to rules for that pattern
	relatives     map[string]Relatives // relation name to pairs of nouns

	//
	currPhase assert.Phase

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
		err = errutil.New("no top domain")
	} else {
		err = d.Schedule(cat.cursor, when, what)
	}
	return
}

func (d *Domain) Schedule(at string, when assert.Phase, what func(*Weaver) error) (err error) {
	if currPhase, phase := d.currPhase, when; currPhase > phase {
		err = errutil.New("unexpected phase")
	} else {
		d.scheduling[when] = append(d.scheduling[when], memento{
			what, at,
		})
	}
	return
}

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
func (d *Domain) AssembleDomain(w assert.Phase, flags PhaseFlags) (err error) {
	if ds, e := d.GetDependencies(); e != nil {
		err = e
	} else {
		d.currPhase = w // hrmm
		// note: even if there are no ephemera... in a given phase..
		// there can still be rivals and other results to process
		if e := d.checkRivals(ds, !flags.NoDuplicates); e != nil {
			err = e
		} else {
			ctx := Weaver{d: d, phase: w, Runtime: d.catalog.run}
			// don't "range" over the phase data since the contents can change during traversal.
			// fix: if we were merging in the definitions we wouldnt have to walk upwards...
			// what's best? see note in checkRivals()
			for i := 0; i < len(d.scheduling[w]); i++ {
				if e := d.scheduling[w][i].call(&ctx); e != nil {
					err = errutil.Append(err, e)
				}
			}
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

// for each domain in the passed list, output its full ancestry tree ( or just its parents )
func (d *Domain) WriteDomain(w Writer) (err error) {
	if dep, e := d.Resolve(); e != nil {
		err = e
	} else {
		name, row, at := d.Name(), dep.Strings(true), d.OriginAt()
		err = w.Write(mdl.Domain, name, row, at)
	}
	return
}
