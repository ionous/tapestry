package weave

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	*mdl.Modeler     // db output
	Env          Env // generic storage for command processing

	// these three all refer to domain... cleanup?
	domains        map[string]*Domain // all domains
	pendingDomains []*Domain          // unprocessed domains
	sceneStack     SceneStack         // tracks scene begins

	run rt.Runtime    // custom runtime for running macros
	db  *tables.Cache // for domain processing, rival testing; tbd: move to mdl entirely?
}

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
		logerr = func(str string, parts ...any) {
			warn(fmt.Errorf(str, parts...))
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

// return name of the scene currently being woven.
// fix: i dont think this should be exposed
// its only for jess.
func (cat *Catalog) CurrentScene() (ret string) {
	if d, ok := cat.sceneStack.Top(); ok {
		ret = d.Name()
	}
	return
}

// walk the domains and run the commands remaining in their queues
func (cat *Catalog) AssembleCatalog() (err error) {
	for {
		if len(cat.sceneStack) > 0 {
			err = errors.New("mismatched begin/end scene")
			break
		} else if len(cat.pendingDomains) == 0 {
			break
		} else if e := cat.assembleNext(); e != nil {
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

func (cat *Catalog) assembleNext() (err error) {
	found := -1 // tbd: a better way?
	for i := 0; i < len(cat.pendingDomains); i++ {
		next := cat.pendingDomains[i]
		if e := cat.isDomainReady(next.name); e == nil {
			found = i
			break
		} else if e != errNotReady {
			err = e
			break
		}
	}
	if found < 0 && err == nil {
		first := cat.pendingDomains[0]
		err = fmt.Errorf("circular or unknown domain %q", first.name)
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
			// sets up the domain stack in case new scheduling comes in.
			log.Printf("weaving %s...\n", d.name)
			cat.sceneStack.Push(d)
			err = d.runAll()
			cat.sceneStack.Pop()
			if err == nil {
				pos := compact.Source{File: "initialization", Line: -1}
				err = d.finalizeDomain(cat.PinPos(d.name, pos))
			}
		}
	}
	return
}

func (cat *Catalog) Schedule(when weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) error {
	pos := compact.Source{}
	return cat.SchedulePos(pos, when, cb)
}

// run the passed function now or in the future.
func (cat *Catalog) ScheduleCmd(key typeinfo.Markup, when weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) error {
	pos := compact.MakeSource(key.GetMarkup(false))
	return cat.SchedulePos(pos, when, cb)
}

// run the passed function now or in the future.
func (cat *Catalog) SchedulePos(pos compact.Source, when weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) error {
	if when == 0 {
		panic("xxx")
	}
	d, _ := cat.sceneStack.Top()
	return d.proc.schedule(when, func(now weaver.Phase) error {
		pen := cat.Modeler.PinPos(d.name, pos)
		w := localWeaver{d, pen}
		return cb(w, cat.run)
	})
}

// used for jess:
// runs the passed function until it returns true or errors
// if currently processing, the first step will execute next phase.
func (cat *Catalog) Step(cb func(weaver.Phase) (bool, error)) error {
	d, _ := cat.sceneStack.Top()
	return d.proc.schedule(0, func(now weaver.Phase) (err error) {
		if ok, e := cb(now); e != nil {
			err = e
		} else if !ok {
			err = fmt.Errorf("%w step", weaver.ErrMissing)
		}
		return
	})
}

// find or create the named domain
func (cat *Catalog) EnsureScene(name string) (ret *Domain) {
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

func (cat *Catalog) SceneBegin(name string, at compact.Source, exe []rt.Execute) (ret *mdl.Pen, err error) {
	if d := cat.EnsureScene(name); d.finished {
		err = fmt.Errorf("trying to begin a scene %s that already finished weaving", d.name)
	} else {
		// fix: maybe there should be a specific non Pen method for create?
		pen := cat.PinPos(d.name, at)
		if e := pen.AddScene(d.name, at.Comment); e != nil {
			err = e
		} else if cd, ok := cat.sceneStack.Top(); ok {
			// we're implicitly dependent on the currently running domain
			// noting that tests dont always have a current domain
			// but really it shouldn't be a requirement anyway.
			err = pen.AddDependency(cd.name)
		}
		if err == nil {
			d.pos = at
			d.startup = exe
			cat.sceneStack.Push(d)
			ret = pen
		}
	}
	return
}

func (cat *Catalog) SceneEnd() {
	if _, ok := cat.sceneStack.Top(); !ok {
		panic("mismatched SceneEnd")
	} else {
		cat.sceneStack.Pop()
	}
}

func (cat *Catalog) findRivals() (err error) {
	if res, e := findRivals(cat.db); e != nil {
		err = e
	} else if len(res) > 0 {
		err = fmt.Errorf("%w in %w", mdl.ErrConflict, rivalErrorList(res))
	}
	return
}

var errNotReady = errors.New("not ready")

// determine if all parent domains have been processed.
// fix: it should be possible to ask for a list of all such domains
// rather then ask for domains one at a time and filter.
func (cat *Catalog) isDomainReady(scene string) (err error) {
	if rows, e := cat.db.Query(`
		select uses from domain_tree 
		where base = ?1 
		and uses != ?1
		order by dist desc`, scene); e != nil {
		err = e
	} else {
		var uses string
		err = tables.ScanAll(rows, func() (err error) {
			if d, ok := cat.domains[uses]; !ok {
				err = errNotReady
			} else if !d.finished {
				err = errNotReady
			}
			return
		}, &uses)
	}
	return
}
