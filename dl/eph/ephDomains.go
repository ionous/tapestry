package eph

import (
	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
)

type Resolution int

const (
	Unresolved Resolution = iota
	Processing
	Resolved
	Errored
)

type Domain struct {
	name, at     string
	originalName string
	inflect      inflect.Ruleset
	deps         DomainList
	eph          EphList
	status       Resolution
	resolved     UniqueNames
	err          error
}

type DomainList []*Domain

// for now, duplicates are okay.
func (dl *DomainList) add(d *Domain) {
	if !dl.contains(d) {
		(*dl) = append(*dl, d)
	}
}

func (dl *DomainList) contains(d *Domain) (okay bool) {
	for _, el := range *dl {
		if el == d {
			okay = true
			break
		}
	}
	return
}

func (d *Domain) Resolved() []string {
	return []string(d.resolved)
}

// Recursively determine the domain's dependency list
func (d *Domain) Resolve(newlyResolved func(*Domain) error) (err error) {
	switch d.status {
	case Resolved:
		// ignore things that are already resolved

	case Processing:
		d.status, d.err = Errored, errutil.New("Circular reference detected:", d.name)
		err = d.err

	case Unresolved:
		d.status = Processing
		var res UniqueNames
		for _, dep := range d.deps {
			if e := dep.Resolve(newlyResolved); e != nil {
				d.status, d.err = Errored, errutil.New(e, "->", d.name)
				err = d.err
				break
			} else {
				res.AddName(dep.name)
				for _, n := range dep.resolved {
					res.AddName(n)
				}
			}
		}
		if err == nil {
			d.status = Resolved
			d.resolved = res
			//
			if newlyResolved != nil {
				if e := newlyResolved(d); e != nil {
					d.status, d.err = Errored, e
					err = d.err
				}
			}
		}
	default:
		if e := d.err; e != nil {
			err = e
		} else {
			err = errutil.New("Unknown error processing", d.name)
		}
	}
	return
}

type AllDomains map[string]*Domain

// DomainStack - keep track of the current block while marshaling.
// ( so that end block can be called. )
type DomainStack []*Domain

func (k *DomainStack) Push(b *Domain) {
	(*k) = append(*k, b)
}

func (k *DomainStack) Top() (ret *Domain, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, okay = (*k)[end], true
	}
	return
}

// return false if empty
func (k *DomainStack) Pop() (ret *Domain, okay bool) {
	if end := len(*k) - 1; end >= 0 {
		ret, (*k) = (*k)[end], (*k)[:end]
		okay = true
	}
	return
}
