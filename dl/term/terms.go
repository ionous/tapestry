package term

import (
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

// Terms implements a Scope mapping names to specified parameters.
type Terms struct {
	fields []g.Field
	values []g.Value
}

// fix: creating the kind should really happen during assembly
func (ps *Terms) NewRecord(kinds g.Kinds) (ret *g.Record, err error) {
	ret = g.NewInternalRecord(kinds, ps.fields, ps.values)
	return
}

func (ps *Terms) AddValue(name string, v g.Value) int {
	field := g.Field{name, v.Affinity(), v.Type()}
	ps.fields = append(ps.fields, field)
	ps.values = append(ps.values, v)
	return len(ps.fields) - 1
}

// converts an object value to an object id
// a nil kind is okay -- it allows any type
func ConvertObject(run rt.Runtime, obj g.Value, kind string) (ret g.Value, err error) {
	if !safe.Compatible(obj, kind, false) {
		err = errutil.New("object", obj, "not compatible with", kind)
	} else {
		ret = g.ObjectAsText(obj)
	}
	return
}

// converts a text value to a valid object id
func ConvertName(run rt.Runtime, n string, kind string) (ret g.Value, err error) {
	// look up the named object...
	if len(n) == 0 {
		ret = g.StringFrom("", "object="+kind)
	} else if obj, e := safe.ObjectFromString(run, n); e != nil {
		err = e
	} else {
		ret, err = ConvertObject(run, obj, kind)
	}
	return
}
