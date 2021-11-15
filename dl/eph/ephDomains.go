package eph

import (
	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
)

type Domain struct {
	name, at string
	inflect  inflect.Ruleset
	deps     domainList
}

// is it better to store pointers or names?
type domainList []*Domain

type Domains struct {
	pool       domainPool
	processing domainStack
	ephemera   EphList
}

// for now, duplicates are okay.
func (dl *domainList) add(d *Domain) {
	(*dl) = append(*dl, d)
}
func (dl *domainList) remove(d *Domain) {
	end := len(*dl) - 1
	for i := end; i >= 0; i-- {
		if (*dl)[i] == d {
			(*dl)[i] = (*dl)[end]
			(*dl) = (*dl)[:end]
			break
		}
	}
}
func (dl *domainList) contains(d *Domain) (okay bool) {
	for _, el := range *dl {
		if el == d {
			okay = true
			break
		}
	}
	return
}

func (d *Domain) resolve(resolved *domainList, unresolved *domainList) (err error) {
	unresolved.add(d)
	for _, edge := range d.deps {
		if !resolved.contains(edge) {
			if unresolved.contains(edge) {
				err = errutil.New("Circular reference detected:", d.name, "->", edge.name)
				break
			} else if e := edge.resolve(resolved, unresolved); e != nil {
				err = errutil.New(e, "->", edge.name)
				break
			}
		}
	}
	if err == nil {
		resolved.add(d)
		unresolved.remove(d)
	}
	return
}

func (ds *Domains) Resolve() (ret domainList, err error) {
	var unresolved domainList
	for _, d := range ds.pool {
		if !ret.contains(d) {
			if e := d.resolve(&ret, &unresolved); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func (ds *Domains) GetDomain(n string) (ret *Domain) {
	if d, ok := ds.pool[n]; ok {
		ret = d
	} else {
		d = &Domain{name: n}
		if ds.pool == nil {
			ds.pool = domainPool{n: d}
		} else {
			ds.pool[n] = d
		}
		ret = d
	}
	return
}

type domainPool map[string]*Domain

// domainStack - keep track of the current block while marshaling.
// ( so that end block can be called. )
type domainStack []*Domain

func (k *domainStack) Push(b *Domain) {
	(*k) = append(*k, b)
}

func (k *domainStack) Top() (ret *Domain, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, okay = (*k)[end], true
	}
	return
}

// return false if empty
func (k *domainStack) Pop() (ret *Domain, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, (*k) = (*k)[end], (*k)[:end]
		okay = true
	}
	return
}
