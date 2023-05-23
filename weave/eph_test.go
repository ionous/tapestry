package weave

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/tables"
	"git.sr.ht/~ionous/tapestry/test/testdb"
	"git.sr.ht/~ionous/tapestry/weave/eph"
	"github.com/ionous/errutil"
)

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Value: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Value: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Value: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Value: s} }

type Warnings []error

// override the global eph warning function
// returns a defer-able function which:
// 1. restores the warning function; and,
// 2. raises a Fatal error if there are any unhandled warnings.
func (w *Warnings) catch(t *testing.T) func() {
	was := LogWarning
	LogWarning = func(e any) {
		(*w) = append((*w), e.(error))
	}
	return func() {
		if len(*w) > 0 {
			t.Fatal("unhandled warnings", *w)
		}
		LogWarning = was
	}
}

// return the warnings as a raw list, clear all stored errors.
func (w *Warnings) all() (ret []error) {
	ret, (*w) = (*w), nil
	return ret
}

// remove and return the first warning, or error if there are none left.
func (w *Warnings) shift() (err error) {
	if cnt := len(*w); cnt == 0 {
		err = errutil.New("out of warnings")
	} else {
		err, (*w) = (*w)[0], (*w)[1:]
	}
	return
}

type domainTest struct {
	name string
	// queues the commands so we can makeDomain without worrying about error handling
	// also allows shuffling the declarations (within a single domain)
	queue     []eph.Ephemera
	noShuffle bool
	db        *sql.DB
	cat       *Catalog
}

func (dt *domainTest) Open(name string) *sql.DB {
	if dt.db == nil {
		dt.db = testdb.Open(name, testdb.Memory, "")
	}
	return dt.db
}

func (dt *domainTest) Close() {
	if dt.db != nil {
		dt.db.Close()
		dt.db = nil
	}
}

func dd(names ...string) []string {
	return names
}

func (dt *domainTest) makeDomain(names []string, add ...eph.Ephemera) {
	n, req := names[0], names[1:]
	if !dt.noShuffle {
		// shuffle the incoming ephemera
		rand.Shuffle(len(add), func(i, j int) { add[i], add[j] = add[j], add[i] })
		// shuffle the order of domain dependencies
		rand.Shuffle(len(req), func(i, j int) { req[i], req[j] = req[j], req[i] })
	}
	dt.queue = append(dt.queue, &eph.BeginDomain{
		Name:     n,
		Requires: req,
	})
	dt.queue = append(dt.queue, add...)
	dt.queue = append(dt.queue, &eph.EndDomain{
		Name: n,
	})
}

func (dt *domainTest) Assemble() (ret *Catalog, err error) {
	for _, el := range dt.queue {
		if e := el.Assert(dt.cat); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		err = dt.cat.AssembleCatalog()
	}
	if errs := dt.cat.Errors; len(errs) > 0 {
		err = errutil.New(err, errs)
	}
	return dt.cat, err
}

// relation, kind, cardinality, otherKinds
func newRelation(r, k, c, o string) *eph.Relations {
	var card eph.Cardinality
	switch c {
	case tables.ONE_TO_ONE:
		card = &eph.OneOne{Kind: k, OtherKind: o}
	case tables.ONE_TO_MANY:
		card = &eph.OneMany{Kind: k, OtherKinds: o}
	case tables.MANY_TO_ONE:
		card = &eph.ManyOne{Kinds: k, OtherKind: o}
	case tables.MANY_TO_MANY:
		card = &eph.ManyMany{Kinds: k, OtherKinds: o}
	default:
		panic("unknown cardinality")
	}
	return &eph.Relations{
		Rel:         r,
		Cardinality: card,
	}
}

func (dt *domainTest) readAssignments() ([]string, error) {
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
}

// an array of dependencies
func (dt *domainTest) readDomain(n string) (ret []string, err error) {
	if ds, e := tables.QueryStrings(dt.db,
		`select uses from domain_tree 
		where base = ?1 
		order by dist desc`,
		n); e != nil {
		err = e
	} else {
		// the domain tree includes the domain itself
		// the tests dont expect that.
		ret = ds[:len(ds)-1]
	}
	return
}

func (dt *domainTest) readDomains() ([]string, error) {
	return tables.QueryStrings(dt.db,
		`select domain from domain_order`)
}

// domain, kind, field name, affinity, subtype
func (dt *domainTest) readFields() ([]string, error) {
	return tables.QueryStrings(dt.db, `
  select mk.domain ||':'|| mk.kind ||':'|| mf.field  ||':'|| 
         mf.affinity  ||':'|| coalesce(mt.kind, '') 
	from mdl_field mf 
	join mdl_kind mk 
		on(mf.kind = mk.rowid)
	left join mdl_kind mt 
		on(mf.type = mt.rowid)`)
}

// domain, input, serialized program
func (dt *domainTest) readGrammar() ([]string, error) {
	return tables.QueryStrings(dt.db, `
  select domain ||':'|| name ||':'|| prog
	from mdl_grammar`)
}

// domain, kind, materialized path
func (dt *domainTest) readKinds() (ret []string, err error) {
	type kind struct {
		id           int
		domain, kind string
	}
	var kinds []kind
	var k kind
	if e := tables.QueryAll(dt.db,
		`select rowid, domain, kind
		from mdl_kind
		order by rowid
		`, func() (_ error) {
			kinds = append(kinds, k)
			return
		}, &k.id, &k.domain, &k.kind); e != nil {
		err = e
	} else {
		for _, k := range kinds {
			// just to be confusing, this is the opposite order of KindOfAncestors
			// root is on the right here.
			if path, e := tables.QueryStrings(dt.db,
				`select mk.kind
				from mdl_kind ks
				join mdl_kind mk
					-- is Y (is their name) a part of X (our path)
					on instr(',' || ks.path,
									 ',' || mk.rowid || ',' )
				where ks.rowid = ?1
				order by mk.rowid desc`, k.id); e != nil {
				err = e
				break
			} else {
				row := fmt.Sprintf("%s:%s:%s", k.domain, k.kind, strings.Join(path, ","))
				ret = append(ret, row)
			}
		}
	}
	return
}

// domain, kind, name, serialized initialization
func (dt *domainTest) readLocals() ([]string, error) {
	// read kind assignments for those that have initializers
	// the original order was per domain, and then in dependency order withing that domain
	//  wonder if itd be enough to simply go by kind id since the order of writing should follow
	return tables.QueryStrings(dt.db, `
	select mk.domain ||':'|| mk.kind ||':'|| mf.field ||':'|| ma.value
	from mdl_assign ma
	join mdl_field mf 
		on(ma.field = mf.rowid)
	join mdl_kind mk 
		on(mf.kind = mk.rowid)
	order by mk.rowid, mf.rowid, ma.rowid`)
}

// domain, noun, name, rank
func (dt *domainTest) readNames() ([]string, error) {
	return tables.QueryStrings(dt.db, `
	select my.domain ||':'|| mn.noun ||':'|| my.name ||':'|| my.rank
	from mdl_name my
	join mdl_noun mn
		on(my.noun = mn.rowid)
	order by mn.rowid, my.rowid`)
}

// domain, noun, kind
func (dt *domainTest) readNouns() ([]string, error) {
	return tables.QueryStrings(dt.db, `
  select mn.domain ||':'|| mn.noun ||':'|| mk.kind
	from mdl_noun mn
	join mdl_kind mk
		on(mn.kind = mk.rowid)
	order by mk.rowid, mn.rowid`)
}

// domain, oneWord, otherWord
func (dt *domainTest) readOpposites() ([]string, error) {
	return tables.QueryStrings(dt.db, `
  select domain ||':'|| oneWord ||':'|| otherWord
	from mdl_rev mp 
	order by domain, oneWord`)
}

// domain, relation, noun, other noun
// original order was domain, and alpha relation name
func (dt *domainTest) readPairs() ([]string, error) {
	return tables.QueryStrings(dt.db, `
	with domain_rank(domain, rank) as (
		select *, row_number() over (order by 1) 
		from domain_order
	)
	select mp.domain ||':'|| mk.kind ||':'|| one.noun ||':'|| other.noun
	from domain_rank dr
	join mdl_pair mp 
	on(mp.domain = dr.domain)
	join mdl_kind mk
		on(mp.relKind = mk.rowid)
	join mdl_noun one
		on(oneNoun = one.rowid)
	join mdl_noun other
		on(otherNoun = other.rowid)
	order by dr.rank, mk.kind`)
}

// domain, pattern, labels, result field
func (dt *domainTest) readPatterns() ([]string, error) {
	return tables.QueryStrings(dt.db, `
  select mk.domain ||':'|| mk.kind ||':'|| labels ||':'|| result
	from mdl_pat mp
	join mdl_kind mk 
		on(mp.kind = mk.rowid)`)
}

// domain, many, one
func (dt *domainTest) readPlurals() ([]string, error) {
	return tables.QueryStrings(dt.db, `
  select domain ||':'|| many ||':'|| one
	from mdl_plural`)
}

// domain, relation, one kind, other kind, cardinality
func (dt *domainTest) readRelations() ([]string, error) {
	return tables.QueryStrings(dt.db, `
  select mk.domain ||':'|| mk.kind ||':'|| one.kind ||':'|| other.kind ||':'|| mr.cardinality
  from mdl_rel mr 
  join mdl_kind mk
  	on(relKind = mk.rowid)
	join mdl_kind one
		on(oneKind = one.rowid)
	join mdl_kind other
		on(otherKind = other.rowid)`)
}

// domain, pattern, target, phase, filter, prog
func (dt *domainTest) readRules() ([]string, error) {
	return tables.QueryStrings(dt.db, `
  select mr.domain ||':'|| mk.kind ||':'|| coalesce(mt.kind, '') ||':'|| phase ||':'|| filter ||':'|| prog
  from mdl_rule mr 
  join mdl_kind mk
  	on(mr.kind = mk.rowid)
  left join mdl_kind mt 
  	on(mr.target = mt.rowid)`)
}

// domain, noun, field, value
func (dt *domainTest) readValues() ([]string, error) {
	return tables.QueryStrings(dt.db, `
	select mn.domain ||':'|| mn.noun ||':'|| mf.field ||':'|| mv.value
	from mdl_value mv 
	join mdl_noun mn 
		on(mv.noun = mn.rowid)
	join mdl_field mf	
		on(mv.field = mf.rowid)`)
}
