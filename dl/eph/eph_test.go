package eph

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

func catchWarnings(out *[]error) func() {
	was := LogWarning
	LogWarning = func(e error) {
		*out = append(*out, e)
	}
	return func() {
		LogWarning = was
	}
}

type domainTest struct {
	out []Ephemera
}

func dd(names ...string) []string {
	return names
}

func (dt *domainTest) makeDomain(names []string, add ...Ephemera) {
	n, req := names[0], names[1:]
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
	cat.processing.Push(cat.EnsureDomain("g", "global"))
	for i, el := range dt.out {
		if e := cat.AddEphemera(EphAt{At: strconv.Itoa(i), Eph: el}); e != nil {
			err = e
			break
		}
	}
	return
}
