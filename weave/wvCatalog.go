package weave

import (
	"database/sql"
	"log"

	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grokdb"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	*mdl.Modeler     // db output
	Env          Env // generic storage for command processing

	SuspendSchedule int // hack until promises are available

	// cleanup? seems redundant to have three views of the same domain
	domains        map[string]*Domain
	processing     DomainStack
	pendingDomains []*Domain

	cursor string        // current source position
	run    rt.Runtime    // custom runtime for running macros
	gdb    grokdb.Source // english parser
	db     *tables.Cache // for domain processing, rival testing; tbd: move to mdl entirely?
}

type ScheduledCallback func(*Weaver) error

func NewCatalog(db *sql.DB) *Catalog {
	return NewCatalogWithWarnings(db, nil, nil)
}

func NewCatalogWithWarnings(db *sql.DB, run rt.Runtime, warn func(error)) *Catalog {
	if run == nil {
		dec := decoder.DecodeNone("unsupported decoder")
		if e := tables.CreateAll(db); e != nil {
			panic(e)
		} else if qx, e := qdb.NewQueries(db, false); e != nil {
			panic(e)
		} else {
			run = qna.NewRuntime(
				log.Writer(),
				qx,
				dec,
			)
		}
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
		Env:     make(Env),
		cursor:  "x", // fix
		db:      cache,
		domains: make(map[string]*Domain),
		Modeler: m,
		run:     run,
		gdb:     gdb,
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

		if e := cat.run.ActivateDomain(d.name); e != nil {
			err = e
		} else if e := cat.findRivals(); e != nil {
			err = e
		} else {
			cat.processing.Push(d)
			//
			for p := Phase(0); p <= RequireAll; p++ {
				w := Weaver{Catalog: cat, Domain: d.name, Phase: p, Runtime: cat.run}
				if e := d.runPhase(&w); e != nil {
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
	if d, e := cat.addDomain(name, cat.cursor, requires...); e != nil {
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

func (cat *Catalog) Schedule(when Phase, what ScheduledCallback) (err error) {
	if d, ok := cat.processing.Top(); !ok {
		err = errutil.New("unknown top level domain")
	} else {
		err = d.schedule(cat.cursor, when, what)
	}
	return
}

// return the uniformly named domain ( creating it if necessary )
func (cat *Catalog) addDomain(name, at string, reqs ...string) (ret *Domain, err error) {
	// find or create the domain
	n := inflect.Normalize(name)
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
				// check for circular references:
				if req := inflect.Normalize(req); n == req {
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
						if e := cat.Modeler.Pin(n, at).AddDependency(req); e != nil {
							err = e
							break
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
