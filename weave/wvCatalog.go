package weave

import (
	"database/sql"
	"log"
	"strconv"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
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

	processing SceneStack

	run rt.Runtime    // custom runtime for running macros
	db  *tables.Cache // for domain processing, rival testing; tbd: move to mdl entirely?
}

type StepFunction func(weaver.Phase) (bool, error)

func NewCatalog(db *sql.DB) *Catalog {
	return NewCatalogWithWarnings(db, nil, nil)
}

func NewCatalogWithWarnings(db *sql.DB, run rt.Runtime, warn func(error)) *Catalog {
	if run == nil {
		// fix: this only happens for tests; unwind the calls so that's explicit
		dec := qdb.DecodeNone("unsupported decoder")
		if q, e := qdb.NewQueries(db, dec); e != nil {
			panic(e)
		} else {
			run = qna.NewRuntime(q)
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

func (cat *Catalog) GetRuntime() rt.Runtime {
	return cat.run
}

// return the uniformly named domain ( if it exists )
func (cat *Catalog) GetDomain(n string) (*Domain, bool) {
	d, ok := cat.domains[n]
	return d, ok
}

// return the uniformly named domain ( if it exists )
func (cat *Catalog) CurrentScene() (ret string) {
	if d, ok := cat.processing.Top(); ok {
		ret = d.Name()
	}
	return
}

// walk the domains and run the commands remaining in their queues
func (cat *Catalog) AssembleCatalog() (err error) {
	for {
		if len(cat.processing) > 0 {
			err = errutil.New("mismatched begin/end scene")
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

// used for jess:
// runs the passed function until it returns true or errors
// if currently processing, the first step will execute next phase.
func (cat *Catalog) Step(cb StepFunction) (err error) {
	if d, ok := cat.processing.Top(); !ok {
		err = errutil.New("step has no scene")
	} else if z := d.currPhase; z < 0 {
		err = errutil.Fmt("step scene %q already woven", d.name)
	} else {
		d.steps = append(d.steps, cb)
	}
	return
}

// run the passed function now or in the future.
func (cat *Catalog) Schedule(when weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) (err error) {
	return cat.SchedulePos(nil, when, cb)
}

// run the passed function now or in the future.
// pass the line number of the operation.
// tbd: maybe it would be better if the ops carried their full source pos
// currently they just have line number.
func (cat *Catalog) SchedulePos(t typeinfo.Markup, when weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) (err error) {
	if d, ok := cat.processing.Top(); !ok {
		err = errutil.New("schedule has no active scene")
	} else {
		pos := mdl.MakeSource(t)
		err = d.schedule(pos, when, cb)
	}
	return
}

func (cat *Catalog) EnsureScene(name string) (ret *Domain) {
	// find or create the domain
	n := inflect.Normalize(name)
	if d, ok := cat.domains[n]; ok {
		ret = d
	} else {
		d = &Domain{name: n, cat: cat}
		cat.pendingDomains = append(cat.pendingDomains, d)
		cat.domains[n] = d
		ret = d
	}
	return
}

func (cat *Catalog) PushScene(d *Domain, at mdl.Source) (ret *mdl.Pen, err error) {
	if d.currPhase != 0 {
		err = errutil.New("trying to start an in-progress scene", d.name)
	} else {
		pen := cat.PinPos(d.name, at)
		cd := cat.CurrentScene()
		if len(cd) > 0 {
			// we're implicitly dependent on the currently running domain
			// noting that tests dont always have a current domain
			// but really it shouldn't be a requirement anyway.
			err = pen.AddDependency(cd)
		}
		if err == nil {
			cat.processing.Push(d)
			ret = pen
		}
	}
	return
}

func (cat *Catalog) PopScene() {
	if _, ok := cat.processing.Top(); !ok {
		panic("mismatched PopScene")
	} else {
		cat.processing.Pop()
	}
}

func (cat *Catalog) findRivals() (err error) {
	var rivals error
	if e := findRivals(cat.db, func(group, domain, key, value, at string) (_ error) {
		rivals = errutil.Append(rivals, errutil.Fmt("%w in domain %q at %q for %s %q",
			mdl.ErrConflict, domain, at, group, value))
		return
	}); e != nil {
		err = e
	} else {
		err = rivals
	}
	return
}
