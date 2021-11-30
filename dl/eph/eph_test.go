package eph

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/ionous/errutil"
)

type testOut []string

// ignores the category for simpler testing
func (x *testOut) Write(_cat string, args ...interface{}) (err error) {
	var b strings.Builder
	for i, arg := range args {
		if i > 0 {
			b.WriteRune(':')
		}
		b.WriteString(fmt.Sprint(arg))
	}
	(*x) = append((*x), b.String())
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
	out []Ephemera
}

func dd(names ...string) []string {
	return names
}

func (dt *domainTest) makeDomain(names []string, add ...Ephemera) {
	n, req := names[0], names[1:]
	// shuffle the order of domain dependencies
	rand.Shuffle(len(req), func(i, j int) { req[i], req[j] = req[j], req[i] })
	dt.out = append(dt.out, &EphBeginDomain{
		Name:     n,
		Requires: req,
	})
	dt.out = append(dt.out, add...)
	dt.out = append(dt.out, &EphEndDomain{
		Name: n,
	})
	return
}

func (dt *domainTest) addToCat(cat *Catalog) (err error) {
	// fix: it would be nice to not need an explicit top domain.
	cat.processing.Push(cat.EnsureDomain("g", "global"))
	for i, el := range dt.out {
		if e := cat.AddEphemera(EphAt{At: strconv.Itoa(i), Eph: el}); e != nil {
			err = e
			break
		}
	}
	return
}
