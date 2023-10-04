package assign

import (
	"strconv"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// Dot provides access to a value inside another value.
// ex. in objects, lists, or records.
type Dot interface {
	Resolve(rt.Runtime) (Dotted, error)
}

// Dotted is an intermediate stage in picking values
type Dotted interface {
	Peek(rt.Runtime, g.Value) (g.Value, error)
	Poke(run rt.Runtime, target, value g.Value) error
	writeTo(b *strings.Builder)
}

type DottedPath []Dotted

type DotField string
type DotIndex int // zero-based index

func (dl DottedPath) String() string {
	var b strings.Builder
	for _, el := range dl {
		el.writeTo(&b)
	}
	return b.String()
}

// Resolving is somewhat superfluous:
// while it could happen during peek and poke
// having all of the elements resolved ahead of time is better for debugging.
func ResolvePath(run rt.Runtime, dots []Dot) (ret DottedPath, err error) {
	if cnt := len(dots); cnt > 0 {
		path := make(DottedPath, 0, cnt)
		for _, el := range dots {
			if p, e := el.Resolve(run); e != nil {
				if str := path.String(); len(str) == 0 {
					err = e
				} else {
					err = errutil.New(e, "with partial path", str)
				}
				break
			} else {
				path = append(path, p)
			}
		}
		if err == nil {
			ret = path
		}
	}
	return
}

func (op *AtField) Resolve(run rt.Runtime) (ret Dotted, err error) {
	if field, e := safe.GetText(run, op.Field); e != nil {
		err = cmdError(op, e)
	} else {
		ret = DotField(field.String())
	}
	return
}

func (op *AtIndex) Resolve(run rt.Runtime) (ret Dotted, err error) {
	if idx, e := safe.GetNumber(run, op.Index); e != nil {
		err = cmdError(op, e)
	} else {
		ret = DotIndex(idx.Int() - 1)
	}
	return
}

// raw string
func (dot DotField) Field() string {
	return string(dot)
}

// print friendly string
func (dot DotField) writeTo(b *strings.Builder) {
	b.WriteRune('.')
	b.WriteString(dot.Field())
}

func (dot DotField) Peek(run rt.Runtime, val g.Value) (ret g.Value, err error) {
	if aff := val.Affinity(); !affine.HasFields(aff) {
		err = errutil.New(aff, "doesn't have fields")
	} else if el, e := val.FieldByName(string(dot)); e != nil {
		err = e
	} else {
		// backwards compat: records used to produce default values internally
		// now we have to do that manually
		if aff != affine.Record || el != nil {
			ret = el
		} else if k, e := run.GetKindByName(el.Type()); e != nil {
			err = e
		} else if rec := val.Record(); rec == nil {
			err = errutil.New("can't peek a nil value; how'd we get here?")
		} else {
			// use the record interface to avoid a copy.
			newVal := g.RecordOf(k.NewRecord())
			if e := rec.SetNamedField(string(dot), newVal); e != nil {
				err = e
			} else {
				ret = newVal
			}
		}
	}
	return
}

func (dot DotField) Poke(run rt.Runtime, target, newValue g.Value) (err error) {
	if aff := target.Affinity(); !affine.HasFields(aff) {
		err = errutil.New(aff, "doesn't have fields")
	} else if e := target.SetFieldByName(dot.Field(), newValue); e != nil {
		err = e
	}
	return
}

// raw int
func (dot DotIndex) Index() int {
	return int(dot)
}

// print friendly string
func (dot DotIndex) writeTo(b *strings.Builder) {
	b.WriteRune('[')
	b.WriteString(strconv.Itoa(dot.Index()))
	b.WriteRune(']')
}

func (dot DotIndex) Peek(run rt.Runtime, val g.Value) (ret g.Value, err error) {
	if aff := val.Affinity(); !affine.IsList(aff) {
		err = errutil.New(aff, "isn't a list")
	} else if i, e := safe.Range(dot.Index(), 0, val.Len()); e != nil {
		err = e
	} else {
		ret = val.Index(i)
	}
	return
}

func (dot DotIndex) Poke(run rt.Runtime, target, newValue g.Value) (err error) {
	if aff := target.Affinity(); !affine.IsList(aff) {
		err = errutil.New(aff, "isn't a list")
	} else if i, e := safe.Range(dot.Index(), 0, target.Len()); e != nil {
		err = e
	} else if e := target.SetIndex(i, newValue); e != nil {
		err = e
	}
	return
}
