package eph

import (
	"log"

	"github.com/ionous/errutil"
)

// Catalog - receives ephemera from the importer.
type Catalog struct {
	Writer     // not a huge fan of this here.... hrm...
	domains    AllDomains
	processing DomainStack
}

func (c *Catalog) Warn(e error) {
	log.Println(e) // for now good enough
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

func (c *Catalog) WriteDomains() error {
	return writeDomains(c.Writer, c.domains)
}
