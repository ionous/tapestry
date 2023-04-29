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

type testOut []string

func (log *testOut) Results() []string {
	return (*log)[1:]
}

// ignores the category for simpler testing
func (log *testOut) Write(cat string, args ...interface{}) (err error) {
	if len(*log) == 0 {
		err = errutil.New("testOut not initialized")
	} else if (*log)[0] == cat {
		if argcnt, catcnt := len(args), strings.Count(cat, "?"); cat != "" && catcnt != argcnt {
			err = errutil.Fmt("not enough parameters. want %d have %d", catcnt, argcnt)
		} else {
			var b strings.Builder
			for i, arg := range args {
				if i > 0 {
					b.WriteRune(':')
				}
				b.WriteString(fmt.Sprint(arg))
			}
			(*log) = append((*log), b.String())
		}
	}
	return
}

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

func (dt *domainTest) Dequeue() (err error) {
	for _, el := range dt.queue {
		if e := el.Assert(dt.cat); e != nil {
			err = e
			break
		}
	}
	return
}

func (dt *domainTest) Assemble() (ret *Catalog, err error) {
	err = dt.Dequeue()
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

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Value: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Value: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Value: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Value: s} }
