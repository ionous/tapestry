package eph

import "strconv"

type testOut map[string][]outEl

type outEl []interface{}

func (x *testOut) Write(q string, args ...interface{}) (err error) {
	if *x == nil {
		*x = make(map[string][]outEl)
	}
	els := (*x)[q]
	els = append(els, args)
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
	dt.out = append(dt.out, &EphBeginDomain{
		Name:     names[0],
		Requires: names[1:],
	})
	dt.out = append(dt.out, add...)
	dt.out = append(dt.out, &EphEndDomain{
		Name: names[0],
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
