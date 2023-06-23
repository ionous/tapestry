package weave

import (
	"database/sql"
	"errors"
	"log"

	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
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
	writer         *mdl.Modeler
	run            rt.Runtime
	db             *tables.Cache
	warn           func(error)

	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	macros      macroReg
	autoCounter Counters
	env         Environ

	domainNouns map[domainNoun]*ScopedNoun
}

// request future processing from the catalog's importer.
type Schedule interface {
	Schedule(*Catalog) error
}

type domainNoun struct{ domain, noun string }

func NewCatalog(db *sql.DB) *Catalog {
	return NewCatalogWithWarnings(db, nil)
}

func NewCatalogWithWarnings(db *sql.DB, warn func(error)) *Catalog {
	qx, e := qdb.NewQueryx(db)
	if e != nil {
		panic(e)
	}
	m, e := mdl.NewModeler(db)
	if e != nil {
		panic(e)
	}
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
		db:          tables.NewCache(db),
		domains:     make(map[string]*Domain),
		writer:      m,
		run: qna.NewRuntimeOptions(
			log.Writer(),
			qx,
			qna.DecodeNone("unsupported decoder"),
			qna.NewOptions()),
	}
}

func (k *Catalog) SetSource(x string) {
	k.cursor = x
}

// return the uniformly named domain ( if it exists )
func (c *Catalog) GetDomain(n string) (*Domain, bool) {
	d, ok := c.domains[n]
	return d, ok
}

func (c *Catalog) EndDomain() (err error) {
	if _, ok := c.processing.Top(); !ok {
		err = errutil.New("unexpected domain ending when there's no domain")
	} else {
		c.processing.Pop()
	}
	return
}

// calls to schedule() between begin/end domain write to this newly declared domain.
func (c *Catalog) BeginDomain(name string, requires []string) (err error) {
	if n, ok := UniformString(name); !ok {
		err = InvalidString(name)
	} else if reqs, e := UniformStrings(requires); e != nil {
		err = e // transform all the names first to determine any errors
	} else if d, e := c.addDomain(n, c.cursor, reqs...); e != nil {
		err = e
	} else {
		c.processing.Push(d)
	}
	return
}

func (c *Catalog) ensureDomain(n string) (ret *Domain) {
	// find or create the domain
	if d, ok := c.domains[n]; ok {
		ret = d
	} else {
		d = &Domain{name: n, catalog: c}
		c.pendingDomains = append(c.pendingDomains, d)
		c.domains[n] = d
		ret = d
	}
	return
}

// log if the error is a duplicate;
// only return non-duplicate, non-nil errors
func (c *Catalog) eatDuplicates(e error) (err error) {
	if !errors.Is(e, mdl.Duplicate) {
		err = e
	} else if c.warn != nil {
		c.warn(e)
	}
	return
}

// return the uniformly named domain ( creating it if necessary )
func (c *Catalog) addDomain(n, at string, reqs ...string) (ret *Domain, err error) {
	d := c.ensureDomain(n)
	if d.currPhase >= assert.RequireDependencies {
		err = errutil.New("can't add new dependencies to parent domains", d.name)
	} else {
		// domains are implicitly dependent on their parent domain
		if p, ok := c.processing.Top(); ok {
			reqs = append(reqs, p.name)
		}
		// probably asking for  trouble:
		// the tests have no top level domain (tapestry) the way weave does
		// we still need them to wind up in the table eventually...
		if len(reqs) == 0 {
			err = c.writer.Domain(d.name, "", at)
		} else {
			for _, req := range reqs {
				if dep, ok := UniformString(req); !ok {
					err = errutil.New("invalid name", req)
					break
				} else if e := c.findDomainCycles(d.name, dep); e != nil {
					err = e
					break
				} else {
					e := c.writer.Domain(d.name, dep, at)
					if e := c.eatDuplicates(e); e != nil {
						err = e
						break
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

// check before inserting a reference to avoid circularity;
// errors if one was detected.
func (c *Catalog) findDomainCycles(n, req string) (err error) {
	if n == req {
		err = errutil.Fmt("circular reference: %q can't depend on itself", n)
	} else {
		var exists bool
		if e := c.db.QueryRow(
			`select 1 
		from domain_tree 
		where base = ?1
		and uses = ?2
		and base != uses`, req, n).Scan(&exists); e != nil && e != sql.ErrNoRows {
			err = e
		} else if exists {
			err = errutil.Fmt("circular reference: %q requires %q", req, n)
		}
	}
	return
}

func (c *Catalog) findRivals() (err error) {
	var rivals error
	if e := findRivals(c.db, func(group, domain, key, value, at string) (_ error) {
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

// exposed for testing:
// tbd: maybe this could pull in the newly relevant domains;
// ( ex. use domain = active count )
// currently it happens after the domains have been activated
// and therefore compares everything to everything each time.
// note: fields don't have rivals because they all exist in the same domain as their owner kind.
func findRivals(db tables.Querier, onConflict func(group, domain, key, value, at string) error) (err error) {
	if rows, e := db.Query(`
	with active_grammar as (
		select mg.*
		from mdl_grammar mg 
		join active_domains
		using (domain)
	),

	active_facts as (
		select mx.*
		from mdl_fact mx
		join active_domains
		using (domain)
	)
	
	select 'fact', a.domain, a.at, a.fact, a.value
		from active_facts as a 
		join active_facts as b 
			using(fact)
		where a.domain != b.domain 
		and a.value != b.value
	union all

	select 'kind', a.domain, a.at, a.kind, ''
		from active_kinds as a 
		join active_kinds as b 
			using(name)
		where a.domain != b.domain
	union all

	select 'grammar', a.domain, a.at, a.name, ''
		from active_grammar as a 
		join active_grammar as b 
			using(name)
		where a.domain != b.domain 
		and a.prog != b.prog
	union all

	select 'opposite', a.domain, a.at, a.oneWord, a.otherWord 
		from active_rev as a 
		join active_rev as b 
			using(oneWord)
		where a.domain != b.domain 
		and a.otherWord != b.otherWord
	union all

	select 'plural', a.domain, a.at, a.many, a.one
		from active_plurals as a 
		join active_plurals as b 
			using(many)
		where a.domain != b.domain
		and a.one != b.one
	`); e != nil {
		err = errutil.New("database error", e)
	} else {
		var group, domain, key, value string
		var at sql.NullString
		if e := tables.ScanAll(rows, func() error {
			return onConflict(group, domain, key, value, at.String)
		}, &group, &domain, &at, &key, &value); e != nil && e != sql.ErrNoRows {
			err = e
		}
	}
	return
}
