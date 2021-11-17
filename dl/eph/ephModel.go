package eph

import (
	"sort"
	"strings"

	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// database/sql like interface
type Writer interface {
	Write(q string, args ...interface{}) error
}

// domain name and materialized parents separated by commas
var mdl_domain = tables.Insert("mdl_domain", "domain", "path")

// a plural word ("many") can have at most one singular definition per domain
// ie. "people" and "persons" are valid plurals of "person",
// but "people" as a singular can only be defined as "person" ( not "cat" )
var mdl_plural = tables.Insert("mdl_plural", "domain", "many", "one")

// domains should be in "most" core to least order
// each line should have all the dependencies it needs
func writeDomains(out Writer, ds AllDomains) (err error) {
	// we *try* as much as possible to keep the order stableish
	sorted := make([]string, 0, len(ds))
	for _, d := range ds {
		if len(d.at) == 0 {
			err = errutil.Append(err, errutil.New("domain never declared", d.name))
		} else {
			sorted = append(sorted, d.name)
			d.status = Unresolved // for the resolution callbacks to trigger each time we writeDomains
		}
	}
	if err == nil {
		sort.Strings(sorted)
		for _, n := range sorted {
			d := ds[n]
			if e := d.resolveCb(func(d *Domain) (err error) {
				deps := d.Resolved()
				sort.Strings(deps) // sort for some amount of consistency
				ls := strings.Join(deps, ",")
				if e := out.Write(mdl_domain, d.name, ls); e != nil {
					err = errutil.New("domain", d.name, "couldn't output", e)
				}
				return
			}); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}
