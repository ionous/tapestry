package eph

import (
	"strings"

	"github.com/ionous/errutil"
)

type ResolvedDomains []*Domain

// for each domain in the passed list, output its full ancestry tree ( or just its parents )
func (ds *ResolvedDomains) WriteDomains(w Writer, fullTree bool) (err error) {
	for _, d := range *ds {
		if deps, e := d.GetDependencies(); e != nil {
			err = errutil.Append(err, e)
		} else if e := w.Write(mdl_domain, d.name, strings.Join(deps.Ancestors(fullTree), ",")); e != nil {
			err = errutil.Append(err, errutil.New("domain", d.name, "couldn't write", e))
		}
	}
	return
}

// for each domain, determine the kinds that it defined
func (ds *ResolvedDomains) ResolveKinds(cat *Catalog) (ret ResolvedKinds, err error) {
	var out ResolvedKinds
	for _, d := range *ds {
		if e := d.kinds.ResolveKinds(&out); e != nil {
			err = e
			break
		}
	}
	return
}

// walk the domains and run the commands remaining in their queues
// this is a little squirly because it needs cat for rivals and ephemera processing
// but the resolved domains itself comes from the catalog..
func (ds *ResolvedDomains) ProcessDomains(cat *Catalog, phaseActions PhaseActions) (err error) {
	for _, d := range *ds {
		if e := cat.checkRivals(d); e != nil {
			err = e
			break
		} else {
			for phase, ephlist := range d.phases {
				for _, el := range ephlist {
					if e := el.Eph.Catalog(cat, d, el.At); e != nil {
						err = errutil.Append(err, e)
					}
				}
				if err != nil {
					break
				} else if act, ok := phaseActions[Phase(phase)]; ok {
					if e := act(cat, d); e != nil {
						err = e
						break

					}
				}
			}
		}
	}
	return
}
