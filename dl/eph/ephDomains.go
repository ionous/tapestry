package eph

import (
	"errors"
	"git.sr.ht/~ionous/tapestry/imp/assert"

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
	plural        PluralPairs
	opposites     OppositePairs
	rules         map[string]Rulesets  // pattern name to rules for that pattern
	relatives     map[string]Relatives // relation name to pairs of nouns
}

type PhaseData struct {
	eph []EphAt
}

func (d *Domain) Resolve() (ret Dependencies, err error) {
	if len(d.at) == 0 {
		err = DomainError{d.name, errutil.New("never defined")}
	} else if ds, e := d.resolve(d, (*catDependencyFinder)(d.catalog)); e != nil {
		err = DomainError{d.name, e}
	} else {
		ret = ds
	}
	return
}

func (d *Domain) AddEphemera(at string, ep Ephemera) (err error) {
	if currPhase, phase := d.currPhase, ep.Phase(); currPhase > phase {
		err = errutil.New("unexpected phase")
	} else {
		// fix? consider all ephemera in a flat slice ( again ) scanning by phase instead of partitioning.
		// that way we dont need all the separate lists and we can append....
		els := d.phases[phase]
		els.eph = append(els.eph, EphAt{
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
			err = DomainError{scope.name, e}
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
	var de DomainError
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
func (d *Domain) Assemble(w assert.Phase, flags PhaseFlags) (err error) {
	if ds, e := d.GetDependencies(); e != nil {
		err = e
	} else {
		// note: even if there are no ephemera... in a given phase..
		// there can still be rivals and other results to process
		d.currPhase = w // hrmmm...
		if e := d.checkRivals(w, ds, !flags.NoDuplicates); e != nil {
			err = e
		} else {
			// don't "range" over the phase data since the contents can change during traversal.
			// fix: if we were merging in the definitions we wouldnt have to walk upwards...
			// what's best? see note in checkRivals()
			for i := 0; i < len(d.phases[w].eph); i++ {
				op := d.phases[w].eph[i]
				if e := op.Eph.Assemble(d.catalog, d, op.At); e != nil {
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
func (d *Domain) checkRivals(phase assert.Phase, ds Dependencies, allowDupes bool) (err error) {
	// start with nothing and merge in to check for artifacts
	if deps := ds.Parents(); len(deps) > 1 {
		a := make(Artifacts)
		for _, dep := range deps {
			scope := dep.(*Domain) // allow this to panic
			if defs := scope.defs; len(defs) > 0 {
				if e := a.Merge(defs, allowDupes); e != nil {
					err = DomainError{scope.name, e}
					break
				}
			}
		}
	}
	return
}

func (c *Catalog) WriteDomains(w Writer) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		err = ds.WriteTable(w, mdl.Domain, true)
	}
	return
}

// EphBeginDomain
func (op *EphBeginDomain) Phase() assert.Phase { return assert.DomainPhase }

func (op *EphBeginDomain) Assemble(c *Catalog, _nil *Domain, at string) (err error) {
	if n, ok := UniformString(op.Name); !ok {
		err = InvalidString(op.Name)
	} else if reqs, e := UniformStrings(op.Requires); e != nil {
		err = e // transform all the names first to determine any errors
	} else if d, e := c.EnsureDomain(n, at, reqs...); e != nil {
		err = e
	} else {
		c.processing.Push(d)
	}
	return
}

// EphEndDomain
func (op *EphEndDomain) Phase() assert.Phase { return assert.DomainPhase }

// pop the most recent domain
func (op *EphEndDomain) Assemble(c *Catalog, _nil *Domain, at string) (err error) {
	// we expect it's the current domain, the parent of this command, that's the one ending
	if n, ok := UniformString(op.Name); !ok && len(op.Name) > 0 {
		err = InvalidString(op.Name)
	} else if d, ok := c.processing.Top(); !ok {
		err = errutil.New("unexpected domain ending when there's no domain")
	} else if n != d.name && len(op.Name) > 0 {
		err = errutil.New("unexpected domain ending, requested", op.Name, "have", d.name)
	} else {
		c.processing.Pop()
	}
	return
}
