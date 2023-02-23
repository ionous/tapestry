package assign

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// the location of a value in the runtime.
type RefValue struct {
	Object, Field string
	Path          DottedPath
}

// a value and its location pulled from an object in the runtime.
// path indicates a further as yet examined value that can be extracted from the root value.
type RootValue struct {
	RootValue g.Value
	RefValue
}

// a string representation of the reference ( mainly for debugging )
func (ref *RefValue) String() string {
	var b strings.Builder
	b.WriteString(ref.Object)
	b.WriteRune('.')
	b.WriteString(ref.Field)
	b.WriteString(ref.Path.String())
	return b.String()
}

// pull the object.field value from the runtime (
func (ref *RefValue) GetRootValue(run rt.Runtime) (ret RootValue, err error) {
	if rootValue, e := run.GetField(ref.Object, ref.Field); e != nil {
		err = e
	} else {
		ret = RootValue{RefValue: *ref, RootValue: rootValue}
	}
	return
}

// let an object know if some of its inner contents has changed
// ( not needed if SetValue is called directly )
func (ref *RefValue) SetDirty(run rt.Runtime) {
	if ref.Object != meta.Variables {
		run.SetField(meta.ObjectDirty, ref.Object, g.True)
	}
}

func (ref *RefValue) SetValue(run rt.Runtime, newValue g.Value) (err error) {
	if len(ref.Path) == 0 {
		// write straight to the named field
		err = run.SetField(ref.Object, ref.Field, newValue)
	} else {
		// write a value ( somewhere ) inside the named field
		if tgt, e := run.GetField(ref.Object, ref.Field); e != nil {
			err = e
		} else {
			// dive into the value, up to the very last field
			last := len(ref.Path) - 1
			for _, at := range ref.Path[:last] {
				if next, e := at.Peek(run, tgt); e != nil {
					err = e
					break
				} else {
					tgt = next
				}
			}
			// then write into that very last field
			if err == nil {
				if e := ref.Path[last].Poke(run, tgt, newValue); e != nil {
					err = e
				} else {
					// manually notify the object, since it can't see the change deep inside some other field.
					ref.SetDirty(run)
				}
			}
		}
	}
	return
}

// unpack a value
func (src *RootValue) GetValue(run rt.Runtime) (ret g.Value, err error) {
	val := src.RootValue
	for i, dot := range src.Path {
		if next, e := dot.Peek(run, val); e != nil {
			err = errutil.New(e, "peeking at part", i)
			break
		} else {
			val = next
		}
	}
	if err == nil {
		ret = val
	}
	return
}

// FIX: convert and warn instead of error on field affinity checks
func (src *RootValue) GetCheckedValue(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := GetValue(run, src); e != nil {
		err = e
	} else if e := safe.Check(v, aff); e != nil {
		err = errutil.New(e, "at", src.RefValue.String())
	} else {
		ret = v
	}
	return
}

func (op *RootValue) GetList(run rt.Runtime) (ret g.Value, err error) {
	if els, e := op.GetValue(run); e != nil {
		err = e
	} else if aff := els.Affinity(); !affine.IsList(aff) {
		err = errutil.New("expected %s was a list, but its a %v", op.RefValue, aff)
	} else {
		ret = els
	}
	return
}
