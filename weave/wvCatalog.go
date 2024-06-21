package weave

import (
	"database/sql"
	"log"
	"strconv"

	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
	"github.com/ionous/errutil"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	*mdl.Modeler     // db output
	Env          Env // generic storage for command processing

	// cleanup? seems redundant to have three views of the same domain
	domains        map[string]*Domain
	pendingDomains []*Domain

	cursor     string // current file
	processing DomainStack

	run rt.Runtime    // custom runtime for running macros
	db  *tables.Cache // for domain processing, rival testing; tbd: move to mdl entirely?
}

type StepFunction func(weaver.Phase) (bool, error)

func NewCatalog(db *sql.DB) *Catalog {
	return NewCatalogWithWarnings(db, nil, nil)
}

func NewCatalogWithWarnings(db *sql.DB, run rt.Runtime, warn func(error)) *Catalog {
	if run == nil {
		if q, e := qdb.NewQueries(db); e != nil {
			panic(e)
		} else {
			run = qna.NewRuntime(q, nil)
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
	return &Catalog{
		Env:     make(Env),
		db:      tables.NewCache(db),
		domains: make(map[string]*Domain),
		Modeler: m,
		run:     run,
	}
}

// func (cat *Catalog) SetSource(x string) {
// 	cat.cursor = x
// }

func (cat *Catalog) GetRuntime() rt.Runtime {
	return cat.run
}

// return the uniformly named domain ( if it exists )
func (cat *Catalog) GetDomain(n string) (*Domain, bool) {
	d, ok := cat.domains[n]
	return d, ok
}

// return the uniformly named domain ( if it exists )
func (cat *Catalog) CurrentDomain() (ret string) {
	if d, ok := cat.processing.Top(); ok {
		ret = d.Name()
	}
	return
}

// walk the domains and run the commands remaining in their queues
func (cat *Catalog) AssembleCatalog() (err error) {
	for {
		if len(cat.processing) > 0 {
			err = errutil.New("mismatched begin/end domain")
			break
		} else if len(cat.pendingDomains) == 0 {
			break
		} else if _, e := cat.assembleNext(); e != nil {
			err = e
			break
		}
	}
	return
}

func (cat *Catalog) NewCounter(name string) (ret string) {
	next := cat.Env.Inc(name, 1)
	return name + "-" + strconv.Itoa(next)
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
			log.Printf("weaving %s...\n", d.name)
			for z := weaver.Phase(0); z < weaver.NumPhases; z++ {
				if e := d.runPhase(z); e != nil {
					err = e
					break
				}
			}
			cat.processing.Pop()
			if err == nil {
				if e := d.finalizeDomain(); e != nil {
					err = e
				} else {
					d.currPhase = -1 // all done.
					ret = d
				}
			}
		}
	}
	return
}

func (cat *Catalog) BeginFile(name string) (err error) {
	if cur := cat.cursor; len(cur) > 0 {
		err = errutil.New("file already in progress", cur)
	} else {
		cat.cursor = name
	}
	return
}

func (cat *Catalog) EndFile() {
	if len(cat.cursor) == 0 {
		panic("EndFile called but no BeginFile")
	}
	cat.cursor = ""
}

// calls to schedule() between begin/end domain write to this newly declared domain.
// names dont have to be normalized.
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

// run passed function until it returns true or errors
// if currently processing, the first step will execute next phase.
func (cat *Catalog) Step(cb StepFunction) (err error) {
	if d, ok := cat.processing.Top(); !ok {
		err = errutil.New("unknown top level domain")
	} else {
		err = d.step(cat.cursor, cb)
	}
	return
}

// run the passed function now or in the future.
func (cat *Catalog) Schedule(when weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) (err error) {
	if d, ok := cat.processing.Top(); !ok {
		err = errutil.New("unknown top level domain")
	} else {
		err = d.schedule(cat.cursor, when, cb)
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
	if d.currPhase < 0 || d.currPhase > 0 {
		err = errutil.New("can't add new dependencies to parent domains", d.name)
	} else {
		// domains are implicitly dependent on their parent domain
		if p, ok := cat.processing.Top(); ok {
			reqs = append(reqs, p.name)
		}
		// probably asking for trouble:
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
