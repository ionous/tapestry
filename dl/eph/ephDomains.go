package eph

import (
	"git.sr.ht/~ionous/iffy/tables/mdl"
	"github.com/ionous/errutil"
)

type DomainFinder interface {
	GetDomain(n string) (*Domain, bool)
}

type Domain struct {
	Requires      // other domains this needs ( can have multiple direct parents )
	catalog       *Catalog
	currPhase     Phase // lift into some "ProcessingDomain" structure?
	phases        [NumPhases]PhaseData
	kinds         ScopedKinds
	nouns         ScopedNouns
	checks        asmChecks
	resolvedKinds cachedTable
	resolvedNouns cachedTable
	pairs         PluralPairs
	rules         map[string]Rulesets  // pattern name to rules for that pattern
	relatives     map[string]Relatives // relation name to pairs of nouns
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

func (d *Domain) GetDefinition(phase Phase, key string) (ret string) {
	defs := d.phases[phase].defs
	return defs[key].value
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
func (d *Domain) Assemble(w Phase, flags PhaseFlags) (err error) {
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
		err = ds.WriteTable(w, mdl.Domain, true)
	}
	return
}

// EphBeginDomain
func (op *EphBeginDomain) Phase() Phase { return DomainPhase }

//
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
func (op *EphEndDomain) Phase() Phase { return DomainPhase }

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
