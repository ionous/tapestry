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
func (dt *domainTest) readDomain(n string) ([]string, error) {
	return tables.QueryStrings(dt.db,
		`select uses from domain_tree where base = ?1 order by dist desc`,
		n)
}

func (dt *domainTest) readDomains() ([]string, error) {
	return tables.QueryStrings(dt.db,
		`select domain from domain_order`)
}

// domain, kind, field name, affinity, subtype
func (dt *domainTest) readFields() ([]string, error) {
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
}

// domain, input, serialized program
func (dt *domainTest) readGrammar() ([]string, error) {
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
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
		from mdl_kind mk`, func() (_ error) {
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
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
}

// domain, noun, name, rank
func (dt *domainTest) readNames() ([]string, error) {
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
}

// domain, noun, kind
func (dt *domainTest) readNouns() ([]string, error) {
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
}

// domain, oneWord, otherWord
func (dt *domainTest) readOpposites() ([]string, error) {
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| oneWord ||':'|| otherWord
from mdl_rev mp 
order by domain, oneWord`)
}

// domain, relation, noun, other noun
func (dt *domainTest) readPairs() ([]string, error) {
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
}

// domain, pattern, labels
func (dt *domainTest) readPatterns() ([]string, error) {
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
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
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
}

// domain, pattern, target, phase, filter, prog
func (dt *domainTest) readRules() ([]string, error) {
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
}

// domain, noun, field, value
func (dt *domainTest) readValues() ([]string, error) {
	return tables.QueryStrings(dt.db, `
    select domain ||':'|| many ||':'|| one
from mdl_plural`)
}
