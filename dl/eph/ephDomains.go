package eph

import (
	"github.com/ionous/errutil"
)

type DomainFinder interface {
	GetDomain(n string) (*Domain, bool)
}

type Domain struct {
	Requires  // other domains this needs ( can have multiple direct parents )
	catalog   *Catalog
	currPhase Phase // lift into some "ProcessingDomain" structure?
	phases    [NumPhases]PhaseData
	kinds     ScopedKinds
	nouns     ScopedNouns
	resolvedKinds,
	resolvedNouns cachedTable
	pairs PluralPairs
}

type PhaseData struct {
	eph  []EphAt
	defs Artifacts
}

func (dp *PhaseData) AddDefinition(k string, v Definition) {
	if dp.defs == nil {
		dp.defs = make(Artifacts)
	}
	dp.defs[k] = v
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

func (d *Domain) AddEphemera(ephAt EphAt) (err error) {
	if currPhase, phase := d.currPhase, ephAt.Eph.Phase(); currPhase > phase {
		err = errutil.New("unexpected phase")
	} else if phase == DomainPhase {
		// fix: queue first, and then run?
		err = ephAt.Eph.Assemble(d.catalog, d, ephAt.At)
	} else {
		// fix? consider all ephemera in a flat slice ( again ) scanning by phase instead of partitioning.
		// that way we dont need all the separate lists and we can append....
		els := d.phases[phase]
		els.eph = append(els.eph, ephAt)
		d.phases[phase] = els
	}
	return
}

// used by ephemera during assembly to record some piece of information
// that would cause problems it were specified differently elsewhere.
// ex. some in game password specified as the word "secret" in one place, but "mongoose" somewhere else.
func (d *Domain) AddDefinition(key, at, value string) (err error) {
	if e := VisitTree(d, func(dep Dependency) (err error) {
		scope := dep.(*Domain)
		if e := scope.phases[d.currPhase].defs.CheckConflict(key, value); e != nil {
			err = DomainError{scope.name, e}
		}
		return
	}); e != nil {
		err = e
	} else {
		defs := d.phases[d.currPhase]
		defs.AddDefinition(key, Definition{at: at, value: value})
		d.phases[d.currPhase] = defs
	}
	return
}

// the domain is resolved already.
func (d *Domain) Assemble(phaseActions PhaseActions) (err error) {
	if ds, e := d.GetDependencies(); e != nil {
		err = e
	} else {
		// don't "range" over the phases since the contents can change during traversal.
		for w := 0; w < int(NumPhases); w++ {
			// note: even if there are no ephemera... in a given phase..
			// there can still be rivals and other results to process
			phaseData := d.phases[w]
			currPhase := Phase(w)
			d.currPhase = currPhase // hrmmm...
			act := phaseActions[currPhase]
			if e := d.checkRivals(currPhase, ds, !act.Flags.NoDuplicates); e != nil {
				err = e
				break
			} else {
				// fix: if we were merging in the definitions we wouldnt have to walk upwards...
				// what's best? see note in checkRivals()
				for _, el := range phaseData.eph {
					if e := el.Eph.Assemble(d.catalog, d, el.At); e != nil {
						err = errutil.Append(err, e)
					}
				}
				if err != nil {
					break
				} else if do := act.Do; do != nil {
					if e := do(d); e != nil {
						err = e
						break
					}
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
func (d *Domain) checkRivals(phase Phase, ds Dependencies, allowDupes bool) (err error) {
	// start with nothing and merge in to check for artifacts
	if deps := ds.Parents(); len(deps) > 1 {
		a := make(Artifacts)
		for _, dep := range deps {
			p := dep.(*Domain) // allow this to panic
			if defs := p.phases[phase].defs; len(defs) > 0 {
				if e := a.Merge(defs, allowDupes); e != nil {
					err = DomainError{p.name, e}
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
		err = ds.WriteTable(w, mdl_domain, true)
	}
	return
}

// EphBeginDomain
func (el *EphBeginDomain) Phase() Phase { return DomainPhase }

//
func (el *EphBeginDomain) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if n, ok := UniformString(el.Name); !ok {
		err = InvalidString(el.Name)
	} else if kid, ok := c.GetDomain(n); ok {
		err = errutil.New("domain", n, "at", kid.at, "redeclared", kid.at)
	} else {
		kid := c.EnsureDomain(n, at)
		// add any explicit dependencies
		for _, req := range el.Requires {
			if sub, ok := UniformString(req); !ok {
				err = errutil.Append(err, InvalidString(req))
			} else {
				kid.AddRequirement(sub)
			}
		}
		if err == nil {
			// we are dependent on the parent domain too
			// ( adding it last keeps it closer to the right side of the parent list )
			kid.AddRequirement(d.name)
			c.processing.Push(kid)
		}
	}
	return
}

// EphEndDomain
func (el *EphEndDomain) Phase() Phase { return DomainPhase }

// pop the most recent domain
func (el *EphEndDomain) Assemble(c *Catalog, d *Domain, at string) (err error) {
	// we expect it's the current domain, the parent of this command, that's the one ending
	if n, ok := UniformString(el.Name); !ok {
		err = InvalidString(el.Name)
	} else if n != d.name {
		err = errutil.New("unexpected domain ending, requested", el.Name, "have", d.name)
	} else {
		c.processing.Pop()
	}
	return
}
