package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/weave/eph"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

type DomainFinder interface {
	GetDomain(n string) (*Domain, bool)
}

type Domain struct {
	Requires      // other domains this needs ( can have multiple direct parents )
	catalog       *Catalog
	currPhase     assert.Phase // lift into some "ProcessingDomain" structure?
	defs          Artifacts
	phases        [assert.NumPhases]PhaseData
	kinds         ScopedKinds
	nouns         ScopedNouns
	checks        asmChecks
	resolvedKinds cachedTable
	resolvedNouns cachedTable
	opposites     OppositePairs
	rules         map[string]Rulesets  // pattern name to rules for that pattern
	relatives     map[string]Relatives // relation name to pairs of nouns
}

type PhaseData struct {
	eph []ephAt
}

// At
type ephAt struct {
	At  string       `if:"label=at,type=text"`
	Eph eph.Ephemera `if:"label=eph"`
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

func (d *Domain) QueueEphemera(at string, ep eph.Ephemera) (err error) {
	if currPhase, phase := d.currPhase, assert.AliasPhase; /* FIIIIIIIIIXS ep.Phase()*/ currPhase > phase {
		err = errutil.New("unexpected phase")
	} else {
		// fix? consider all ephemera in a flat slice ( again ) scanning by phase instead of partitioning.
		// that way we dont need all the separate lists and we can append....
		els := d.phases[phase]
		els.eph = append(els.eph, ephAt{
			At:  at,
			Eph: ep,
		})
		d.phases[phase] = els
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

// the domain is resolved already.
func (d *Domain) AssembleDomain(w assert.Phase, flags PhaseFlags) (err error) {
	if ds, e := d.GetDependencies(); e != nil {
		err = e
	} else {
		// note: even if there are no ephemera... in a given phase..
		// there can still be rivals and other results to process
		d.currPhase = w // hrmmm...
		if e := d.checkRivals(ds, !flags.NoDuplicates); e != nil {
			err = e
		} else {
			ctx := Context{c: d.catalog, d: d, phase: w}
			// don't "range" over the phase data since the contents can change during traversal.
			// fix: if we were merging in the definitions we wouldnt have to walk upwards...
			// what's best? see note in checkRivals()
			for i := 0; i < len(d.phases[w].eph); i++ {
				op := d.phases[w].eph[i]
				ctx.at = op.At
				if e := op.Eph.Weave(&ctx); e != nil {
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
