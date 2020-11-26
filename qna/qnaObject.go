package qna

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type qnaObject struct {
	g.Nothing
	n  *Runner // for pointing back to the field cache
	id string
}

func newObjectValue(run *Runner, v interface{}) (ret snapper, err error) {
	if id, ok := v.(string); !ok {
		err = errutil.Fmt("expected id value, got %v(%T)", v, v)
	} else {
		ret = &qnaObject{n: run, id: id}
	}
	return
}

func (q *qnaObject) Snapshot(rt.Runtime) (ret g.Value, err error) {
	ret = q
	return
}

func (q *qnaObject) Affinity() affine.Affinity {
	return affine.Object //
}

func (q *qnaObject) String() (ret string) {
	return q.id
}

func (q *qnaObject) Type() (ret string) {
	if v, e := q.FieldByName(object.Kind); e != nil {
		ret = "object{}"
	} else {
		ret = v.String()
	}
	return
}

func (q *qnaObject) FieldByName(field string) (ret g.Value, err error) {
	// fix temp:
	var key keyType
	switch field {
	case object.Name, object.Kind, object.Kinds:
		// sigh
		key = makeKey(field, q.id)
	default:
		key = makeKey(q.id, field)
	}
	return q.n.getField(key)
}

func (q *qnaObject) SetFieldByName(field string, val g.Value) (err error) {
	if len(field) == 0 {
		err = errutil.Fmt("no field specified")
	} else if writable := field[0] != object.Prefix; !writable {
		err = errutil.Fmt("can't change reserved field %q", field)
	} else {
		key := makeKey(q.id, field)
		err = q.n.setField(key, val)
	}
	return
}
