package eph

import "github.com/ionous/errutil"

// receive ephemera from the importer
type Catalog struct {
	domains    AllDomains
	processing DomainStack
}

func (c *Catalog) GetDomain(n string) (ret *Domain) {
	if d, ok := c.domains[n]; ok {
		ret = d
	} else {
		d = &Domain{name: n}
		if c.domains == nil {
			c.domains = AllDomains{n: d}
		} else {
			c.domains[n] = d
		}
		ret = d
	}
	return
}

// creates domains, suspends all other ephemera until the domains are resolved.
func (c *Catalog) AddEphemera(ephAt EphAt) (err error) {
	if d, ok := c.processing.Top(); !ok {
		err = errutil.New("no domain")
	} else {
		switch el := ephAt.Eph.(type) {
		case *EphEndDomain, *EphBeginDomain:
			err = el.Catalog(c, d, ephAt.At)
		default:
			d.eph.All = append(d.eph.All, ephAt)
		}
	}
	return
}

func (c *Catalog) WriteDomains(out Writer) error {
	return writeDomains(out, c.domains)
}
