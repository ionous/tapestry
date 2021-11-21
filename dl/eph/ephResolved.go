package eph

import "github.com/ionous/errutil"

type ResolvedDomains []*Domain

// walk the domains and run the commands remaining in their queues
func (ds ResolvedDomains) ProcessDomains(cat *Catalog) (err error) {
	for _, d := range ds {
		if e := cat.checkRivals(d); e != nil {
			err = e
			break
		} else {
			for _, phase := range d.phases {
				for _, el := range phase {
					if e := el.Eph.Catalog(cat, d, el.At); e != nil {
						err = errutil.Append(err, e)
					}
				}
				if err != nil {
					break
				}
			}
		}
	}
	return
}
