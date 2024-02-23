package mdl

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"github.com/ionous/errutil"
)

// used by fields to map field name to referenced class ( if any )
// future: wrap the fieldset with a "promise" to avoid looking up the same info repeatedly
type fieldCache map[string]kindInfo

type fieldHandler func(kid, cls kindInfo, field string, aff affine.Affinity) error

func (c *fieldCache) writeFields(pen *Pen, target kindInfo, fields []FieldInfo, call fieldHandler) (err error) {
	for _, f := range fields {
		if e := f.validate(); e != nil {
			err = e
		} else if cls, e := c.getClass(pen, f); e != nil {
			err = e
		} else {
			e := call(target, cls, f.Name, f.Affinity)
			if e := eatDuplicates(pen.warn, e); e != nil {
				err = e
			} else if f.Init != nil {
				e := pen.addDefaultValue(target, f.Name, f.Init)
				if e := eatDuplicates(pen.warn, e); e != nil {
					err = e
				}
			}
			if err != nil {
				err = errutil.Fmt("%w trying to write field %q", err, f.Name)
				break
			}
		}
	}
	return
}

func (c *fieldCache) store(name string, cls kindInfo) {
	if *c == nil {
		*c = make(fieldCache)
	}
	(*c)[name] = cls
}

// for patterns, waits to create the pattern after all fields are known
// which ensures that "extend pattern" (to add locals) happens after define pattern (for parameters and locals)
func (c *fieldCache) getClass(pen *Pen, field FieldInfo) (ret kindInfo, err error) {
	if clsName := field.getClass(); len(clsName) > 0 {
		if a, ok := (*c)[clsName]; ok {
			ret = a
		} else if cls, e := pen.findOptionalKind(clsName); e != nil {
			err = errutil.Fmt("%w trying to find field %q", e, field.Name)
		} else {
			c.store(clsName, cls)
			ret = cls
		}
	}
	return
}
