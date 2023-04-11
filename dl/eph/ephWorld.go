package eph

import "git.sr.ht/~ionous/tapestry/lang"

// PATCH:
// this implements "Model"
// but really it should be "Runtime"
// and ideally it would be built one domain at a time.
// the Catalog ( etc ) would write into it and appropriate things would read from it.
// currently, the catalog has the data in memory -- so its implemented here instead.

func (c *Catalog) PluralOf(word string) (ret string) {
	if d, ok := c.processing.Top(); ok {
		if a, e := d.Pluralize(word); e == nil {
			ret = a
		}
	}
	if len(ret) == 0 {
		ret = lang.Pluralize(word)
	}
	return
}

func (c *Catalog) SingularOf(word string) (ret string) {
	if d, ok := c.processing.Top(); ok {
		if a, e := d.Singularize(word); e == nil {
			ret = a
		}
	}
	if len(ret) == 0 {
		ret = lang.Singularize(word)
	}
	return
}

func (c *Catalog) OppositeOf(word string) (ret string) {
	if d, ok := c.processing.Top(); ok {
		if a, e := d.FindOpposite(word); e == nil {
			ret = a
		}
	}
	if len(ret) == 0 {
		ret = word
	}
	return
}
