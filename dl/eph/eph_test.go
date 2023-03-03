package eph

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

type testOut []string

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
	out       []Ephemera
	noShuffle bool
}

func dd(names ...string) []string {
	return names
}

func (dt *domainTest) makeDomain(names []string, add ...Ephemera) {
	n, req := names[0], names[1:]
	if !dt.noShuffle {
		// shuffle the incoming ephemera
		rand.Shuffle(len(add), func(i, j int) { add[i], add[j] = add[j], add[i] })
		// shuffle the order of domain dependencies
		rand.Shuffle(len(req), func(i, j int) { req[i], req[j] = req[j], req[i] })
	}
	dt.out = append(dt.out, &EphBeginDomain{
		Name:     n,
		Requires: req,
	})
	dt.out = append(dt.out, add...)
	dt.out = append(dt.out, &EphEndDomain{
		Name: n,
	})
}

func (dt *domainTest) addToCat(cat *Catalog) (err error) {
	for _, el := range dt.out {
		if e := cat.AddEphemera("x", el); e != nil {
			err = e
			break
		}
	}
	return
}

// relation, kind, cardinality, otherKinds
func newRelation(r, k, c, o string) *EphRelations {
	var card EphCardinality
	switch c {
	case tables.ONE_TO_ONE:
		card = EphCardinality{EphCardinality_OneOne_Opt, &OneOne{Kind: k, OtherKind: o}}
	case tables.ONE_TO_MANY:
		card = EphCardinality{EphCardinality_OneMany_Opt, &OneMany{Kind: k, OtherKinds: o}}
	case tables.MANY_TO_ONE:
		card = EphCardinality{EphCardinality_ManyOne_Opt, &ManyOne{Kinds: k, OtherKind: o}}
	case tables.MANY_TO_MANY:
		card = EphCardinality{EphCardinality_ManyMany_Opt, &ManyMany{Kinds: k, OtherKinds: o}}
	default:
		panic("unknown cardinality")
	}
	return &EphRelations{
		Rel:         r,
		Cardinality: card,
	}
}

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Value: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Value: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Value: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Value: s} }

func buildAncestors(dt domainTest) (ret *Catalog, err error) {
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		err = e
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryActions,
		FieldPhase:    FieldActions,
	}); e != nil {
		err = e
	} else {
		ret = &cat
	}
	return
}

func buildNouns(dt domainTest) (ret *Catalog, err error) {
	var cat Catalog
	if e := dt.addToCat(&cat); e != nil {
		err = e
	} else if e := cat.AssembleCatalog(PhaseActions{
		AncestryPhase: AncestryActions,
		FieldPhase:    FieldActions,
		NounPhase:     NounActions,
	}); e != nil {
		err = e
	} else {
		ret = &cat
	}
	return
}
