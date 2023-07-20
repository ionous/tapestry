package weave

import (
	"database/sql"
	"errors"
	"log"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grokdb"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
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

	*mdl.Modeler

	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	macros      macroReg
	autoCounter Counters
	Env         Environ

	domainNouns map[domainNoun]*ScopedNoun

	gdb  grokdb.Source
	warn func(error)
}

type domainNoun struct{ domain, noun string }

func NewCatalog(db *sql.DB) *Catalog {
	return NewCatalogWithWarnings(db, nil, nil)
}

func (cat *Catalog) eatDuplicates(e error) (err error) {
	if e == nil || !errors.Is(e, mdl.Duplicate) {
		err = e
	} else if cat.warn != nil {
		cat.warn(e)
	}
	return
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
	var logerr mdl.Log
	if warn != nil {
		logerr = func(fmt string, parts ...any) {
			warn(errutil.Fmt(fmt, parts...))
		}
	}
	m, e := mdl.NewModelerWithWarnings(db, logerr)
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
		warn:        warn, // compat for scoped noun
	}
}

func (cat *Catalog) SetSource(x string) {
	cat.cursor = x
}

func (cat *Catalog) Runtime() rt.Runtime {
	return cat.run
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
		err = errutil.Fmt("circular or unknown domain %q", first.name)
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
			//
			for p := Phase(0); p <= RequireAll; p++ {
				ctx := Weaver{Catalog: cat, Domain: d, Phase: p, Runtime: cat.run}
				if e := d.runPhase(&ctx); e != nil {
					err = e
					break
				} else if e := d.flush(true); e != nil {
					err = e
					break
				}
			}
			// play out suspended actions:
			if err == nil && len(d.suspended) > 0 {
				err = d.flush(false)
			}
			//
			cat.processing.Pop()
			if err == nil {
				d.currPhase = -1 // all done.
				ret = d
			}
		}
	}
	return
}

// calls to schedule() between begin/end domain write to this newly declared domain.
func (cat *Catalog) DomainStart(name string, requires []string) (err error) {
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

func (cat *Catalog) DomainEnd() (err error) {
	if _, ok := cat.processing.Top(); !ok {
		err = errutil.New("unexpected domain ending when there's no domain")
	} else {
		cat.processing.Pop()
	}
	return
}

// note: values are written per *noun* not per domain....
// todo: this should be removed once the values can be written directly into the db
// see writeValues()
func (cat *Catalog) AddNounValue(opNoun, opField string, opPath []string, value literal.LiteralValue) error {
	return cat.Schedule(RequireNames, func(ctx *Weaver) (err error) {
		d, at := ctx.Domain, ctx.At
		if noun, ok := UniformString(opNoun); !ok {
			err = InvalidString(opNoun)
		} else if field, ok := UniformString(opField); !ok {
			err = InvalidString(opField)
		} else if path, e := UniformStrings(opPath); e != nil {
			err = e
		} else if n, e := d.GetClosestNoun(noun); e != nil {
			err = e
		} else {
			err = n.WriteValue(at, field, path, value)
		}
		return
	})
}

// fix: at some point it'd be nice to write values as they are generated
// the basic idea i think would be to write each field AND sub-record path individually
// and, on write, do a test to ensure the path is meaningful,
// and that no "directory value" value exists for any sub path
// ex. "a.b.c" is okay, so long as there's no record stored at "a.b" directly.
// the runtime would change the way it reconstitutes values using multiple rows
func (cat *Catalog) writeValues() (err error) {
Loop:
	for _, n := range cat.domainNouns {
		if rv := n.localRecord; rv.isValid() {
			for _, fv := range rv.rec.Fields {
				// pin each time because the values might have multiple sources
				if e := cat.Modeler.Pin(n.domain.name, rv.at).AddValue(n.name, fv.Field, fv.Value); e != nil {
					err = e
					break Loop
				}
			}
		}
	}
	return
}

// NewCounter generates a unique string, and uses local markup to try to create a stable one.
// instead consider "PreImport" could be used to write a key into the markup if one doesnt already exist.
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

func (cat *Catalog) Schedule(when Phase, what func(*Weaver) error) (err error) {
	if d, ok := cat.processing.Top(); !ok {
		err = errutil.New("unknown top level domain")
	} else {
		err = d.schedule(cat.cursor, when, what)
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

	if d.currPhase < 0 || d.currPhase >= RequireDependencies {
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
			err = cat.Modeler.Pin(d.name, at).AddDependency("")
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
							if e := cat.Modeler.Pin(n, at).AddDependency(dep); e != nil {
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
