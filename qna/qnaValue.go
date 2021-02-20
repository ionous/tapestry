package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

// the database stores raw values we normalize them to assignments
type staticValue struct {
	affinity affine.Affinity
	val      interface{}
}

func (f staticValue) Affinity() affine.Affinity {
	return f.affinity
}

func (f staticValue) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	switch i, a := f.val, f.affinity; a {
	case affine.Bool:
		switch v := i.(type) {
		case nil:
			ret = g.False // zero value for unhandled defaults in sqlite
		case bool:
			ret = g.BoolOf(v)
		case int64:
			// sqlite, boolean values can be represented as 1/0
			ret = g.BoolOf(v == 0)
		default:
			err = errutil.Fmt("unknown %s %T", a, v)
		}

	case affine.Number:
		switch v := i.(type) {
		case nil:
			ret = g.Zero // zero value for unhandled defaults in sqlite
		case int64:
			ret = g.IntOf(int(v))
		case float64:
			ret = g.FloatOf(v)
		default:
			err = errutil.Fmt("unknown %s %T", a, v)
		}

	case affine.NumList:
		switch vs := i.(type) {
		case []float64:
			ret = g.FloatsOf(vs)
		default:
			err = errutil.Fmt("unknown %s %T", a, vs)
		}

	case affine.Text:
		switch v := i.(type) {
		case nil:
			ret = g.Empty // zero value for unhandled defaults in sqlite
		case string:
			ret = g.StringOf(v)
		default:
			err = errutil.Fmt("unknown %s %T", a, v)
		}

	case affine.TextList:
		switch vs := i.(type) {
		case []string:
			ret = g.StringsOf(vs)
		default:
			err = errutil.Fmt("unknown %s %T", a, vs)
		}

	case affine.Record:
		if v, ok := i.(*g.Record); !ok {
			err = errutil.Fmt("unknown %s %T", a, i)
		} else {
			ret = g.RecordOf(v)
		}

	// we could either disallow direct record list storage, or:
	// store the requested kind for storage.
	// case affine.RecordList:
	// 	switch vs := i.(type) {
	// 	case []*g.Record:
	// 		ret = g.RecordsOf(vs)
	// 	default:
	// 		err = errutil.Fmt("unknown %s %T", a, vs)
	// 	}

	default:
		err = errutil.New("unhandled affinity", a)
	}
	return
}
