package mdl

import (
	"fmt"
)

// used by fields to map field name to referenced class ( if any )
// future: wrap the fieldset with a "promise" to avoid looking up the same info repeatedly
type fieldCache map[string]kindInfo

func (c *fieldCache) store(name string, cls kindInfo) {
	if *c == nil {
		*c = make(fieldCache)
	}
	(*c)[name] = cls
}

// for patterns, waits to create the pattern after all fields are known
// which ensures that "extend pattern" (to add locals) happens after define pattern (for parameters and locals)
func (c *fieldCache) getClass(pen *Pen, field FieldInfo) (ret kindInfo, err error) {
	if clsName := field.getDefaultClass(); len(clsName) > 0 {
		if a, ok := (*c)[clsName]; ok {
			ret = a
		} else if cls, e := pen.findOptionalKind(clsName); e != nil {
			err = fmt.Errorf("%w trying to find field %q", e, field.Name)
		} else {
			c.store(clsName, cls)
			ret = cls
		}
	}
	return
}

func (c *fieldCache) precache(pen *Pen, fields []FieldInfo) (err error) {
	for _, f := range fields {
		if _, e := c.getClass(pen, f); e != nil {
			err = e
			break
		}
	}
	return
}
