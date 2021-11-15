package eph

import (
	"strconv"
	"testing"

	"github.com/kr/pretty"
)

type Blob string

func (b Blob) Catalog(c *Catalog, d *Domain, at string) {}

func TestDomains_Flat(t *testing.T) {
	flat := []EphBeginDomain{{
		Name:     "a",
		Requires: []string{"b", "d"},
	}, {
		Name:     "b",
		Requires: []string{"c", "d"},
	}, {
		Name:     "c",
		Requires: []string{"d", "e"},
	}, {
		Name:     "e",
		Requires: []string{"d"},
	}}
	if ds, e := flatten(flat); e != nil {
		t.Fatal(e)
	} else if out, e := ds.Resolve(); e != nil {
		t.Fatal(e)
	} else {
		var names []string
		for _, d := range out {
			names = append(names, d.name)
		}
		if diff := pretty.Diff(names, []string{"g", "d", "e", "c", "b", "a"}); len(diff) > 0 {
			t.Fatal(names, diff)
		}
	}
}

func flatten(flat []EphBeginDomain) (_ Domains, err error) {
	var cat Catalog // the catalog processing requires a global (root) domain.
	cat.domains.processing.Push(cat.domains.GetDomain("g"))
	for i, el := range flat {
		if e := cat.AddEphemera(EphAt{At: strconv.Itoa(i), Eph: &el}); e != nil {
			err = e
			break
		} else if e := cat.AddEphemera(EphAt{At: strconv.Itoa(i), Eph: &EphEndDomain{Name: el.Name}}); e != nil {
			err = e
			break
		}
	}
	return cat.domains, err
}

func TestDomains_Circ(t *testing.T) {
	flat := []EphBeginDomain{{
		Name:     "a",
		Requires: []string{"b", "d"},
	}, {
		Name:     "b",
		Requires: []string{"c", "d"},
	}, {
		Name:     "c",
		Requires: []string{"d", "e"},
	}, {
		Name:     "d",
		Requires: []string{"a"},
	}}

	if ds, e := flatten(flat); e != nil {
		t.Fatal(e)
	} else if _, e := ds.Resolve(); e == nil {
		t.Fatal("expected circular error detection")
	} else {
		t.Log("ok", e)
	}
}
