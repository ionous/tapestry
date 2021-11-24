package eph

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type testOut map[string][]string

func (x *testOut) Write(q string, args ...interface{}) (err error) {
	if *x == nil {
		*x = make(map[string][]string)
	}
	var b strings.Builder
	for i, arg := range args {
		if i > 0 {
			b.WriteRune(':')
		}
		b.WriteString(fmt.Sprint(arg))
	}
	els := (*x)[q]
	els = append(els, b.String())
	(*x)[q] = els
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
	g := cat.EnsureDomain("g")
	g.at = "global"
	cat.processing.Push(g)
	for i, el := range dt.out {
		if e := cat.AddEphemera(EphAt{At: strconv.Itoa(i), Eph: el}); e != nil {
			err = e
			break
		}
	}
	return
}
