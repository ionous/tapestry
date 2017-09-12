package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// ShadowClass provides a factory for constructing pod-like types.
// Each field in the target is assigned an eval capable of filling that field.
// ShadowClass then implements GetObject to which evalutes the fields and constructs the object.
type ShadowClass struct {
	rtype  r.Type
	fields []FieldIndex
	slots  map[string]_ShadowSlot
}

func Shadow(rtype r.Type) *ShadowClass {
	var fields []FieldIndex
	flatten(rtype, nil, &fields)
	return &ShadowClass{rtype, fields, make(map[string]_ShadowSlot)}
}

// GetObject for a shadow type generates an object from the slots specified.
// It is a constructor.
func (c *ShadowClass) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	obj := run.Emplace(r.New(c.rtype).Interface())
	// walk all the fields we recorded and pass them to the new object
	for k, slot := range c.slots {
		// Unpack evaluates an interface to get its resulting go value.
		if v, e := slot.unpack(run); e != nil {
			err = errutil.New("shadow class", c.rtype, "couldn't unpack", k, e)
			break
		} else if e := obj.SetValue(k, v); e != nil {
			err = errutil.New("shadow class", c.rtype, "couldn't set value", k, e)
			break
		}
	}
	if err == nil {
		ret = obj
	}
	return
}

// Addr: all command interfaces are (normally) implemented as pointers.
// Spec carries around the target element, and has to take its address to make it into a pointer so that it matches the implementation.
// When using constructors, spec uses *ShadowClass as its target.
// We need to return the Value just of ourself.
func (c *ShadowClass) Addr() r.Value {
	return r.ValueOf(c)
}

// Type returns the type of the class for field walking.
// Compatible with reflect.Value
func (c *ShadowClass) Type() r.Type {
	return c.rtype
}

func (c *ShadowClass) NumField() int {
	return len(c.fields)
}

// Field returns the value of the requested field.
// Compatible with reflect.Value
// The spec will provide some type-safety on assignment to this value.
// FIX? one thing this cant handle is setting a state via an enumerated value.
// ex. TriState ( yes, no, maybe ) cmd.Param("yes").Value("true")
func (c *ShadowClass) Field(n int) (ret r.Value) {
	if n < len(c.fields) {
		ret = c.FieldByIndex(c.fields[n])
	}
	return
}

func (c *ShadowClass) FieldByName(n string) (ret r.Value) {
	k := ident.IdOf(n)
	unique.WalkProperties(c.rtype, func(f *r.StructField, idx []int) (done bool) {
		if k == ident.IdOf(f.Name) {
			ret, done = c.FieldByIndex(idx), true
		}
		return
	})
	return
}

func (c *ShadowClass) FieldByIndex(n []int) (ret r.Value) {
	field := c.rtype.FieldByIndex(n)
	// determine what kind of eval can produce the passed type.
	if rtype := evalFromType(field.Type); rtype != nil {
		// create an empty eval for the user to poke into
		rvalue := r.New(rtype).Elem()
		c.slots[field.Name] = _ShadowSlot{rtype, rvalue}
		ret = rvalue
	}
	return
}
