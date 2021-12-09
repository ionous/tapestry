package eph

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
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
	LogWarning = func(e error) {
		(*w) = append((*w), e)
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
	// fix: it would be nice to not need an explicit top domain.
	cat.processing.Push(cat.EnsureDomain("g", "x"))
	for _, el := range dt.out {
		if d, ok := cat.processing.Top(); !ok {
			err = errutil.New("no top domain")
			break
		} else if e := d.AddEphemera(EphAt{At: "x", Eph: el}); e != nil {
			err = e
			break
		}
	}
	return
}

// kind, parent
func addKinds(out *[]Ephemera, kps ...string) {
	for i, cnt := 0, len(kps); i < cnt; i += 2 {
		k, p := kps[i], kps[i+1]
		*out = append(*out, &EphKinds{Kinds: k, From: p})
	}
}

// kind, name, affinity(key), class
func addFields(out *[]Ephemera, knacs ...string) {
	for i, cnt := 0, len(knacs); i < cnt; i += 4 {
		k, n, a, c := knacs[i], knacs[i+1], knacs[i+2], knacs[i+3]
		*out = append(*out, &EphFields{
			Kinds: k, Name: n, Affinity: Affinity{a}, Class: c,
		})
	}
}

// noun, kind
func addNouns(out *[]Ephemera, nks ...string) {
	for i, cnt := 0, len(nks); i < cnt; i += 2 {
		n, k := nks[i], nks[i+1]
		*out = append(*out, &EphNouns{Noun: n, Kind: k})
	}
}

// noun(string), field(string), value(literal)
func addValues(out *[]Ephemera, nfvs ...interface{}) {
	for i, cnt := 0, len(nfvs); i < cnt; i += 3 {
		n, f, v := nfvs[i].(string), nfvs[i+1].(string), nfvs[i+2].(literal.LiteralValue)
		*out = append(*out, &EphValues{Noun: n, Field: f, Value: v})
	}
}

// relation, kind, cardinality, otherKinds
func addRelations(out *[]Ephemera, rkcos ...string) {
	for i, cnt := 0, len(rkcos); i < cnt; i += 4 {
		r, k, c, o := rkcos[i], rkcos[i+1], rkcos[i+2], rkcos[i+3]
		var card EphCardinality
		switch c {
		case tables.ONE_TO_ONE:
			card = EphCardinality{EphCardinality_OneOne_Opt, &OneOne{k, o}}
		case tables.ONE_TO_MANY:
			card = EphCardinality{EphCardinality_OneMany_Opt, &OneMany{k, o}}
		case tables.MANY_TO_ONE:
			card = EphCardinality{EphCardinality_ManyOne_Opt, &ManyOne{k, o}}
		case tables.MANY_TO_MANY:
			card = EphCardinality{EphCardinality_ManyMany_Opt, &ManyMany{k, o}}
		default:
			panic("unknown cardinality")
		}
		*out = append(*out, &EphRelations{
			Rel:         r,
			Cardinality: card,
		})
	}
}

// add noun, stem/rel, otherNoun ephemera
func addRelatives(out *[]Ephemera, nros ...string) {
	for i, cnt := 0, len(nros); i < cnt; i += 3 {
		n, r, o := nros[i], nros[i+1], nros[i+2]
		*out = append(*out, &EphRelatives{
			Noun:      n,
			Rel:       r,
			OtherNoun: o,
		})
	}
}

func B(b bool) rt.BoolEval          { return &literal.BoolValue{b} }
func I(n int) rt.NumberEval         { return &literal.NumValue{float64(n)} }
func F(n float64) rt.NumberEval     { return &literal.NumValue{n} }
func T(s string) *literal.TextValue { return &literal.TextValue{s} }
