package eph

import (
	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	domains         map[string]*Domain
	processing      DomainStack
	plurals         PluralTable
	phase           Phase
	resolvedDomains DependencyTable
}

// use the domain rules ( and hierarchy ) to turn the passed plural into its singular form
func (c *Catalog) Singularize(domain, plural string) (ret string, err error) {
	if explict, e := c.plurals.FindSingular((*catDependencyFinder)(c), domain, plural); e != nil {
		err = e
	} else if len(explict) > 0 {
		ret = explict
	} else {
		ret = inflect.Singularize(plural)
	}
	return
}

// used by ephemera during assembly to record some piece of information
// that would cause problems it were specified differently elsewhere.
// ex. some in game password specified as the word "secret" in one place, but "mongoose" somewhere else.
func (c *Catalog) AddDefinition(domain, cat, at, key, value string) (err error) {
	if d, e := c.GetDependentDomains(domain); e != nil {
		err = e
	} else {
		key := CategoryKey{cat, key}
		err = CheckConflicts(d.AllAncestors(), (*catArtifactFinder)(c), key, at, value)
	}
	return
}

// setting fields is a bit .... painful.
// traits want to set the same value to multiple keys so we place the key list last..
func (c *Catalog) AddFields(currDomain *Domain, namedKind string, field FieldDefinition) (err error) {
	var tgtKind *Kind // add the field to this
	var ks []string   // ancestors of the kind
	for {
		if dk, e := FindScopedKind(c, currDomain.name, namedKind); e != nil {
			err = e
			break
		} else {
			currDomain = dk.Domain
			currKind := dk.Kind
			if e := field.CheckConflict(currKind); e != nil {
				err = KindError{currKind.name, currDomain.name, e}
				break
			} else {
				// handle if the current kind was the first kind.
				if tgtKind == nil {
					tgtKind = currKind
					// we will need the kind hierarchy to search parents
					if pks, e := currKind.reqs.GetDependencies(); e != nil {
						err = e
						break // our ancestor kinds
					} else {
						ks = pks.Ancestors()
					}
				}
			}
			// move up to the next kind.
			if pop := len(ks) - 1; pop < 0 {
				break
			} else {
				namedKind, ks = ks[pop], ks[:pop]
				// loop!
			}
		}
	}
	// if everything worked out store definition in our first kind.
	if err == nil && tgtKind != nil {
		field.AddToKind(tgtKind)
	}

	return
}

func (c *Catalog) GetDependentDomains(n string) (ret Dependents, err error) {
	if d, ok := c.domains[n]; !ok {
		err = errutil.New("unknown domains have no dependencies")
	} else {
		ret, err = d.reqs.Resolve(d.name, (*catDependencyFinder)(c))
	}
	return
}

// return the uniformly named domain ( if it exists )
func (c *Catalog) GetDomain(n string) (*Domain, bool) {
	d, ok := c.domains[n]
	return d, ok
}

// return the uniformly named domain ( creating it if necessary )
func (c *Catalog) EnsureDomain(n string) (ret *Domain) {
	if d, ok := c.domains[n]; ok {
		ret = d
	} else {
		d = &Domain{name: n /*, finder: c*/}
		if c.domains == nil {
			c.domains = map[string]*Domain{n: d}
		} else {
			c.domains[n] = d
		}
		ret = d
	}
	return
}

// creates domains, suspends all other ephemera until the domains are resolved.
func (c *Catalog) AddEphemera(ephAt EphAt) (err error) {
	if d, ok := c.processing.Top(); !ok {
		err = errutil.New("no domain")
	} else if currPhase, phase := c.phase, ephAt.Eph.Phase(); currPhase > phase {
		err = errutil.New("unexpected phase")
	} else if phase == DomainPhase {
		err = ephAt.Eph.Assemble(c, d, ephAt.At)
	} else {
		d.phases[phase] = append(d.phases[phase], ephAt)
	}
	return
}

// work out the hierarchy of all the domains, and return them in a list.
// the list has the "shallowest" domains first, and the most derived ( "deepest" ) domains last.
func (c *Catalog) ResolveDomains() (ret DependencyTable, err error) {
	if len(c.resolvedDomains) != 0 {
		ret = c.resolvedDomains
	} else {
		out := make([]Dependents, 0, len(c.domains))
		// walk all domains in the map
		for n, d := range c.domains {
			if len(d.at) == 0 {
				err = errutil.Append(err, errutil.New("domain never declared", d.name))
			} else if dep, e := c.GetDependentDomains(n); e != nil {
				err = errutil.Append(err, e)
			} else {
				out = append(out, dep)
			}
		}
		if err == nil {
			c.resolvedDomains = out
			ret = out
			ret.SortTable()
		}
	}
	return
}

// for given the passed domain hierarchy, determine the kinds that it defined
// func (c *Catalog) ResolveKinds(ds DependencyTable) (ret ResolvedKinds, err error) {
// 	var out ResolvedKinds
// 	for _, n := range ds {
// 		if d, ok := c.GetDomain(n.Name()); !ok {
// 			err = errutil.Fmt("unknown domain %q", n.Name())
// 			break
// 		} else {
// 			for _, k := range d.kinds {
// 				kf := catKindFinder{c, d}
// 				if res, e := kind.reqs.Resolve(kind, &kf); e != nil {
// 					err = errutil.Append(err, e)
// 				} else if ps := res.Parents(); len(ps) > 1 {
// 					e := errutil.Fmt("kind %q should have at most one parent (has: %v)", kn, ps)
// 					err = errutil.Append(err, e)
// 				} else {
// 					kinds := res.Ancestors()
// 					*out = append(*out, ResolvedKind{kn, kinds})
// 				}
// 			}
// 		}
// 	}
// 	return
// }

// type catKindFinder struct {
// 	c      *Catalog
// 	domain *Domain
// }

// func (kf *catKindFinder) GetRequirements(name string) (ret *Requires, err error) {
// 	// the interface for this assumes its one flat pool
// 	// when really, every kind we find should move up its domain path upwards
// 	// perhaps this does mean we should be returning "Node(s)" so that the nodes can implement search themselves
// 	if res, e := FindScopedKind(c, kf.domain.name, name); e != nil {
// 		err = e
// 	} else {
// 		ret = &res.Kind.reqs
// 	}
// 	return
// }

// walk the domains and run the commands remaining in their queues
func (c *Catalog) ProcessDomains(phaseActions PhaseActions) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, deps := range ds {
			if e := c.AssembleDomain(deps, phaseActions); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func (c *Catalog) AssembleDomain(deps Dependents, phaseActions PhaseActions) (err error) {
	if d, ok := c.GetDomain(deps.Name()); !ok {
		err = errutil.New("unknown domain", deps.Name())
	} else {
		c.processing = DomainStack{d} // so ephemera can add other ephemera
		if e := c.checkRivals(deps); e != nil {
			err = e
		} else {
			for w, ephlist := range d.phases {
				for _, el := range ephlist {
					if e := el.Eph.Assemble(c, d, el.At); e != nil {
						err = errutil.Append(err, e)
					}
				}
				if err != nil {
					break
				} else if act, ok := phaseActions[Phase(w)]; ok {
					if e := act(c, d); e != nil {
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
func (c *Catalog) checkRivals(res Dependents) (err error) {
	if parents := res.Parents(); len(parents) > 1 {
		def := make(Artifacts) // start with nothing and merge in to check for artifacts
		for _, p := range parents {
			if d, ok := c.domains[p]; ok {
				if e := def.Merge(d.defs); e != nil {
					err = DomainError{p, e}
					break
				}
			}
		}
	}
	return
}

// private helper to make the catalog compatible with the ArtifactFinder ( for domains )
type catArtifactFinder Catalog

func (c *catArtifactFinder) GetArtifacts(name string) (ret *Artifacts, okay bool) {
	if d, ok := c.domains[name]; ok {
		ret, okay = &d.defs, true
	}
	return
}

// private helper to make the catalog compatible with the DependencyFinder ( for domains )
type catDependencyFinder Catalog

func (c *catDependencyFinder) GetRequirements(name string) (ret *Requires, okay bool) {
	if d, ok := c.domains[name]; ok {
		ret, okay = &d.reqs, true
	}
	return
}
