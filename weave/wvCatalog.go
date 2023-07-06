package weave

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/grokdb"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	domains        map[string]*Domain
	processing     DomainStack
	pendingDomains []*Domain
	Errors         []error
	cursor         string
	run            rt.Runtime
	db             *tables.Cache
	warn           func(error)

	*mdl.Modeler

	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	macros      macroReg
	autoCounter Counters
	Env         Environ

	domainNouns map[domainNoun]*ScopedNoun

	gdb grokdb.Source
}

type domainNoun struct{ domain, noun string }

func NewCatalog(db *sql.DB) *Catalog {
	return NewCatalogWithWarnings(db, nil, nil)
}

func NewCatalogWithWarnings(db *sql.DB, run rt.Runtime, warn func(error)) *Catalog {
	if run == nil {
		qx, e := qdb.NewQueryx(db)
		if e != nil {
			panic(e)
		}
		run = qna.NewRuntimeOptions(
			log.Writer(),
			qx,
			qna.DecodeNone("unsupported decoder"),
			qna.NewOptions(),
		)
	}
	m, e := mdl.NewModeler(db)
	if e != nil {
		panic(e)
	}
	cache := tables.NewCache(db)
	gdb := grokdb.NewSource(cache)
	// fix: this needs cleanup
	// write should be called modeler
	// initial cursor should be set externally or passed in
	// what should be public for Catalog?
	// no panics on creation... etc.
	return &Catalog{
		warn:        warn,
		macros:      make(macroReg),
		oneTime:     make(map[string]bool),
		domainNouns: make(map[domainNoun]*ScopedNoun),
		autoCounter: make(Counters),
		cursor:      "x", // fix
		db:          cache,
		domains:     make(map[string]*Domain),
		Modeler:     m,
		run:         run,
		gdb:         gdb,
	}
}

func (cat *Catalog) SetSource(x string) {
	cat.cursor = x
}

func (cat *Catalog) Runtime() rt.Runtime {
	return cat.run
}

func (cat *Catalog) Warn(e error) {
	if cat.warn != nil {
		cat.warn(e)
	}
}

// return the uniformly named domain ( if it exists )
func (cat *Catalog) GetDomain(n string) (*Domain, bool) {
	d, ok := cat.domains[n]
	return d, ok
}

// walk the domains and run the commands remaining in their queues
func (cat *Catalog) AssembleCatalog() (err error) {
	var ds []*Domain
	for {
		if len(cat.processing) > 0 {
			err = errutil.New("mismatched begin/end domain")
			break
		} else if len(cat.pendingDomains) == 0 {
			break
		} else if was, e := cat.assembleNext(); e != nil {
			err = e
			break
		} else {
			ds = append(ds, was)
		}
	}
	if err == nil {
		err = cat.writeValues()
	}
	return
}

func (cat *Catalog) assembleNext() (ret *Domain, err error) {
	found := -1 // tbd: a better way?
	for i := 0; i < len(cat.pendingDomains); i++ {
		next := cat.pendingDomains[i]
		if ok, e := next.isReadyForProcessing(); e != nil {
			err = e
			break
		} else if ok {
			found = i
			break
		}
	}
	if found < 0 && err == nil {
		first := cat.pendingDomains[0]
		err = errutil.New("circular or unknown domain %q", first.name)
	}
	if err == nil {
		// chop this one out, then process
		d := cat.pendingDomains[found]
		cat.pendingDomains = append(cat.pendingDomains[:found], cat.pendingDomains[found+1:]...)

		if _, e := cat.run.ActivateDomain(d.name); e != nil {
			err = e
		} else if e := cat.findRivals(); e != nil {
			err = e
		} else {
			cat.processing.Push(d)
			for p := assert.Phase(0); p <= assert.RequireAll; p++ {
				ctx := Weaver{Catalog: cat, Domain: d, Phase: p, Runtime: cat.run}
				if e := d.runPhase(&ctx); e != nil {
					err = e
					break
				}
			}
			cat.processing.Pop()
			if err == nil {
				ret = d
			}
		}
	}
	return
}

func (cat *Catalog) AssertAlias(opShortName string, opAliases ...string) error {
	return cat.Schedule(assert.RequireNouns, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		if shortName, ok := UniformString(opShortName); !ok {
			err = errutil.New("invalid name", opShortName)
		} else if n, e := d.getClosestNoun(shortName); e != nil {
			err = e
		} else {
			for _, a := range opAliases {
				if a, ok := UniformString(a); !ok {
					err = errutil.Append(err, InvalidString(a))
				} else {
					err = cat.AddName(d.name, n.name, a, -1, at)
				}
			}
		}
		return
	})
}

// if a parent kind was specified, make the kid dependent on it.
// note: a singular to plural (if needed ) gets handled by the dependency resolver's kindFinder and GetPluralKind()
func (cat *Catalog) AssertAncestor(opKind, opAncestor string) error {
	return cat.Schedule(assert.RequirePlurals, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		// tbd: are the determiners of kinds useful for anything?
		if _, kind := d.UniformDeterminer(opKind); len(kind) == 0 {
			err = InvalidString(kind)
		} else if _, ancestor := d.UniformDeterminer(opAncestor); len(ancestor) == 0 && len(opAncestor) > 0 {
			err = InvalidString(opAncestor)
		} else {
			err = cat.Schedule(assert.RequireDeterminers, func(ctx *Weaver) error {
				e := cat.AddKind(d.name, kind, ancestor, at)
				return cat.eatDuplicates(e)
			})
		}
		return
	})
}

// generates traits and adds them to a custom aspect kind.
func (cat *Catalog) AssertAspectTraits(opAspects string, opTraits []string) error {
	// uses the ancestry phase because it generates kinds ( one per aspect. )
	return cat.Schedule(assert.RequireDeterminers, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		if aspect, ok := UniformString(opAspects); !ok {
			err = InvalidString(opAspects)
		} else if traits, e := UniformStrings(opTraits); e != nil {
			err = e
		} else {
			e := cat.AddKind(d.name, aspect, kindsOf.Aspect.String(), at)
			if e := cat.eatDuplicates(e); e != nil {
				err = e
			} else if len(traits) > 0 {
				err = d.schedule(at, assert.RequireResults, func(ctx *Weaver) error {
					return cat.AddAspect(d.name, aspect, at, traits)
				})
			}
		}
		return
	})
}

func (cat *Catalog) AssertCheck(opName string, prog []rt.Execute, expect literal.LiteralValue) error {
	// uses domain phase, because it needs to ensure a domain exists
	return cat.Schedule(assert.RequireAll, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		if name, ok := UniformString(opName); !ok {
			err = InvalidString(opName)
		} else {
			err = cat.AddCheck(d.name, name, expect, prog, at)
		}
		return
	})
}

func (cat *Catalog) AssertDefinition(path ...string) error {
	return cat.Schedule(assert.RequireAll, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		if end := len(path) - 1; end <= 0 {
			err = errutil.New("path too short", path)
		} else {
			key, value := strings.Join(path[:end], "/"), path[end]
			err = cat.AddFact(d.name, key, value, at)
		}
		return
	})
}

// calls to schedule() between begin/end domain write to this newly declared domain.
func (cat *Catalog) AssertDomainStart(name string, requires []string) (err error) {
	if n, ok := UniformString(name); !ok {
		err = InvalidString(name)
	} else if reqs, e := UniformStrings(requires); e != nil {
		err = e // transform all the names first to determine any errors
	} else if d, e := cat.addDomain(n, cat.cursor, reqs...); e != nil {
		err = e
	} else {
		cat.processing.Push(d)
	}
	return
}

func (cat *Catalog) AssertDomainEnd() (err error) {
	if _, ok := cat.processing.Top(); !ok {
		err = errutil.New("unexpected domain ending when there's no domain")
	} else {
		cat.processing.Pop()
	}
	return
}

func (cat *Catalog) AssertField(kind, field, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.RequireResults, func(ctx *Weaver) error {
		d, at := ctx.Domain, ctx.At
		return addField(ctx, kind, field, class, func(kind, field, class string) (err error) {
			// shortcut: if we specify a field name for a record and no class, we'll expect the class to be the name.
			if len(class) == 0 && isRecordAffinity(aff) {
				class = field
			}
			e := cat.AddMember(d.name, kind, field, aff, class, at)
			if e := cat.eatDuplicates(e); e != nil {
				err = e
			} else if init != nil {
				err = cat.AddDefault(d.name, kind, field, init)
			}
			return
		})
	})
}

// jump/skip/hop	{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}
func (cat *Catalog) AssertGrammar(opName string, prog *grammar.Directive) error {
	return cat.Schedule(assert.RequireRules /*GrammarPhase*/, func(ctx *Weaver) error {
		d, at := ctx.Domain, ctx.At
		return cat.AddGrammar(d.name, opName, prog, at)
	})
}

func (cat *Catalog) AssertNounKind(opNoun, opKind string) error {
	return cat.Schedule(assert.RequireDefaults, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		_, name := d.StripDeterminer(opNoun)
		if noun, ok := UniformString(name); !ok {
			err = InvalidString(opNoun)
		} else if _, kind := d.UniformDeterminer(opKind); len(kind) == 0 {
			err = InvalidString(opKind)
		} else if e := cat.AddNoun(d.name, noun, kind, at); e != nil {
			err = e
		} else {
			cat.domainNouns[domainNoun{d.name, noun}] = &ScopedNoun{domain: d, name: noun}
			err = d.makeNames(noun, name, at)
		}
		return
	})
}

// note: values are written per *noun* not per domain....
func (cat *Catalog) AssertNounValue(opNoun, opField string, opPath []string, opValue literal.LiteralValue) error {
	return cat.Schedule(assert.RequireNames, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		if noun, ok := UniformString(opNoun); !ok {
			err = InvalidString(opNoun)
		} else if field, ok := UniformString(opField); !ok {
			err = InvalidString(opField)
		} else if path, e := UniformStrings(opPath); e != nil {
			err = e
		} else if noun, e := d.getClosestNoun(noun); e != nil {
			err = e
		} else if n, ok := cat.domainNouns[domainNoun{noun.domain, noun.name}]; !ok {
			err = errutil.Fmt("unexpected noun %q in domain %q", noun.name, noun.domain)
		} else if rv, e := n.recordValues(at); e != nil {
			err = e
		} else if value := opValue; value == nil {
			err = errutil.New("null value", opNoun, opField)
		} else {
			return rv.writeValue(noun.name, at, field, path, value)
		}
		return
	})
}

// fix: at some point it'd be nice to write values as they are generated
// the basic idea i think would be to write each field AND sub-record path individually
// and, on write, do a test to ensure the path is meaningful,
// and that no "directory value" value exists for any sub path
// ex. "a.b.c" is okay, so long as there's no record stored at "a.b" directly.
// the runtime would change the way it reconstitutes values to handle all that.
func (cat *Catalog) writeValues() (err error) {
Loop:
	for _, n := range cat.domainNouns {
		if rv := n.localRecord; rv.isValid() {
			for _, fv := range rv.rec.Fields {
				if e := cat.AddValue(n.domain.name, n.name, fv.Field, fv.Value, rv.at); e != nil {
					err = e
					break Loop
				}
			}
		}
	}
	return
}

func (cat *Catalog) AssertOpposite(opOpposite, opWord string) error {
	return cat.Schedule(assert.RequireDependencies, func(ctx *Weaver) (err error) {
		if a, ok := UniformString(opOpposite); !ok {
			err = InvalidString(opOpposite)
		} else if b, ok := UniformString(opWord); !ok {
			err = InvalidString(opWord)
		} else {
			err = cat.AddOpposite(ctx.Domain.name, a, b, cat.cursor)
		}
		return
	})
}

// writes a definition of kindName?args=arg1,arg2,arg3
func (cat *Catalog) AssertParam(kind, field, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.RequireAncestry, func(ctx *Weaver) error {
		d, at := ctx.Domain, ctx.At
		return addField(ctx, kind, field, class, func(kind, field, class string) (err error) {
			if init != nil {
				err = errutil.New("parameters don't currently support initial values")
			} else {
				e := cat.AddKind(d.name, kind, kindsOf.Pattern.String(), at)
				if e := cat.eatDuplicates(e); e != nil {
					err = e
				} else {
					err = cat.AddParameter(d.name, kind, field, aff, class, at)
				}
			}
			return
		})
	})
}

// add to the plurals to the database and ( maybe ) remember the plural for the current domain's set of rules
// not more than one singular per plural ( but the other way around is fine. )
func (cat *Catalog) AssertPlural(opSingular, opPlural string) error {
	return cat.Schedule(assert.RequireDependencies, func(ctx *Weaver) (err error) {
		if plural, ok := UniformString(opPlural); !ok {
			err = InvalidString(opPlural)
		} else if singular, ok := UniformString(opSingular); !ok {
			err = InvalidString(opSingular)
		} else {
			err = cat.AddPlural(ctx.Domain.name, plural, singular, cat.cursor)
		}
		return
	})
}

func (cat *Catalog) AssertResult(kind, field, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.RequireParameters, func(ctx *Weaver) error {
		d, at := ctx.Domain, ctx.At
		return addField(ctx, kind, field, class, func(kind, field, class string) (err error) {
			if init != nil {
				err = errutil.New("return values don't currently support initial values")
			} else {
				e := cat.AddKind(d.name, kind, kindsOf.Pattern.String(), at)
				if e := cat.eatDuplicates(e); e != nil {
					err = e
				} else {
					err = cat.AddResult(d.name, kind, field, aff, class, at)
				}
			}
			return
		})
	})
}

func (cat *Catalog) AssertRelation(opRel, a, b string, amany, bmany bool) error {
	// uses ancestry because it defines kinds for each relation
	return cat.Schedule(assert.RequireDeterminers, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		// like aspects, we dont try to singularize these.
		if rel, ok := UniformString(opRel); !ok {
			err = InvalidString(opRel)
		} else if acls, ok := UniformString(a); !ok {
			err = InvalidString(a)
		} else if bcls, ok := UniformString(b); !ok {
			err = InvalidString(b)
		} else if card := makeCard(amany, bmany); len(card) == 0 {
			err = errutil.New("unknown cardinality")
		} else {
			e := cat.AddKind(d.name, rel, kindsOf.Relation.String(), at)
			if e := cat.eatDuplicates(e); e != nil {
				err = e
			} else {
				err = cat.Schedule(assert.RequireResults, func(ctx *Weaver) (err error) {
					return cat.AddRel(d.name, rel, acls, bcls, card, at)
				})
			}
		}
		return
	})
}

// validate that the pattern for the rule exists then add the rule to the *current* domain
// ( rules are de/activated based on domain, they can be part some child of the domain where the pattern was defined. )
func (cat *Catalog) AssertRelative(opRel, opNoun, opOtherNoun string) error {
	return cat.Schedule(assert.RequireNames, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		if noun, ok := UniformString(opNoun); !ok {
			err = InvalidString(opNoun)
		} else if otherNoun, ok := UniformString(opOtherNoun); !ok {
			err = InvalidString(opOtherNoun)
		} else if rel, ok := UniformString(opRel); !ok {
			err = InvalidString(opRel)
		} else if first, e := d.getClosestNoun(noun); e != nil {
			err = e
		} else if second, e := d.getClosestNoun(otherNoun); e != nil {
			err = e
		} else {
			err = cat.AddPair(d.name, rel, first.name, second.name, at)
		}
		return
	})
}

// validate that the pattern for the rule exists then add the rule to the *current* domain
// ( rules are de/activated based on domain, they can be part some child of the domain where the pattern was defined. )
func (cat *Catalog) AssertRule(pattern string, target string, filter rt.BoolEval, flags assert.EventTiming, prog []rt.Execute) error {
	return cat.Schedule(assert.RequireRelatives, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		if name, ok := UniformString(pattern); !ok {
			err = InvalidString(pattern)
		} else if tgt, ok := UniformString(target); len(target) > 0 && !ok {
			err = errutil.Fmt("unknown or invalid target %q for pattern %q", target, pattern)
		} else {
			flags := fromTiming(flags)
			err = cat.AddRule(d.name, name, tgt, flags, filter, prog, at)
		}
		return
	})
}

// NewCounter generates a unique string, and uses local markup to try to create a stable one.
// instead consider  "PreImport" could be used to write a key into the markup if one doesnt already exist.
// and a free function could also extract what it needs from any op's markup.
// ( then Schedule wouldn't need Catalog for counters )
func (cat *Catalog) NewCounter(name string, markup map[string]any) (ret string) {
	// fix: use a special "id" marker instead?
	if at, ok := markup["comment"].(string); ok && len(at) > 0 {
		ret = at
	} else {
		ret = cat.autoCounter.Next(name)
	}
	return
}

func (cat *Catalog) Schedule(when assert.Phase, what func(*Weaver) error) (err error) {
	if d, ok := cat.processing.Top(); !ok {
		err = errutil.New("unknown top level domain")
	} else {
		err = d.schedule(cat.cursor, when, what)
	}
	return
}

// log if the error is a duplicate;
// only return non-duplicate, non-nil errors
func (cat *Catalog) eatDuplicates(e error) (err error) {
	if !errors.Is(e, mdl.Duplicate) {
		err = e
	} else if cat.warn != nil {
		cat.warn(e)
	}
	return
}

// return the uniformly named domain ( creating it if necessary )
func (cat *Catalog) addDomain(n, at string, reqs ...string) (ret *Domain, err error) {
	// find or create the domain
	d, ok := cat.domains[n]
	if !ok {
		d = &Domain{name: n, cat: cat}
		cat.pendingDomains = append(cat.pendingDomains, d)
		cat.domains[n] = d
	}

	if d.currPhase >= assert.RequireDependencies {
		err = errutil.New("can't add new dependencies to parent domains", d.name)
	} else {
		// domains are implicitly dependent on their parent domain
		if p, ok := cat.processing.Top(); ok {
			reqs = append(reqs, p.name)
		}
		// probably asking for  trouble:
		// the tests have no top level domain (tapestry) the way weave does
		// we still need them to wind up in the table eventually...
		if len(reqs) == 0 {
			err = cat.AddDomain(d.name, "", at)
		} else {
			for _, req := range reqs {
				if dep, ok := UniformString(req); !ok {
					err = errutil.New("invalid name", req)
					break
				} else {
					// check for circular references:
					if n == req {
						err = errutil.Fmt("circular reference: %q can't depend on itself", n)
					} else {
						var exists bool
						if e := cat.db.QueryRow(
							`select 1 
						from domain_tree 
						where base = ?1
						and uses = ?2
						and base != uses`, req, n).Scan(&exists); e != nil && e != sql.ErrNoRows {
							err = e
							break
						} else if exists {
							err = errutil.Fmt("circular reference: %q requires %q", req, n)
							break
						} else {
							e := cat.AddDomain(n, dep, at)
							if e := cat.eatDuplicates(e); e != nil {
								err = e
								break
							}
						}
					}
				}
			}
		}
	}
	if err == nil {
		ret = d
	}
	return
}

func (cat *Catalog) findRivals() (err error) {
	var rivals error
	if e := findRivals(cat.db, func(group, domain, key, value, at string) (_ error) {
		rivals = errutil.Append(rivals, errutil.Fmt("%w in domain %q at %q for %s %q",
			mdl.Conflict, domain, at, group, value))
		return
	}); e != nil {
		err = e
	} else {
		err = rivals
	}
	return
}
