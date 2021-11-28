package eph

import "github.com/ionous/errutil"

type DomainKind struct {
	*Domain
	*Kind
}

func FindScopedKind(df DomainFinder, domain, kind string) (ret DomainKind, err error) {
	var ds []string
	for {
		if d, ok := df.GetDomain(domain); !ok {
			err = errutil.New("unknown domain", domain)
			break
		} else if k, ok := d.kinds[kind]; ok {
			ret = DomainKind{d, k}
			break
		} else {
			// first time through, get the list of ancestor domains
			if ds == nil {
				if deps, e := d.GetDependencies(); e != nil {
					err = e
				} else {
					ds = deps.Ancestors()
				}
			}
			// pop off the next domain name from the list of ancestors
			if pop := len(ds) - 1; pop < 0 {
				err = errutil.New("unknown kind", kind)
				break // no more domains
			} else {
				domain, ds = ds[pop], ds[:pop]
				// loop!
			}
		}
	}
	return
}
