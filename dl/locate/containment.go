package locate

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

//go:generate stringer -type=Containment
type Containment int

const (
	Supports Containment = iota
	Contains
	Wears
	Carries
	Has
)

type LocationOf struct {
	Target rt.ObjectEval
}

func (l *LocationOf) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if rel, ok := run.GetRelation("locale"); !ok {
		err = errutil.New("parent-child relation not found")
	} else if obj, e := l.Target.GetObject(run); e != nil {
		err = errutil.New("couldnt find target", l.Target)
	} else {
		table := rel.GetTable()
		ret, err = nextParent(run, table, obj)
	}
	return
}

func GetAncestors(run rt.Runtime, child rt.Object) (ret rt.ObjectStream, err error) {
	if rel, ok := run.GetRelation("locale"); !ok {
		err = errutil.New("parent-child relation not found")
	} else {
		// find the parent for child:
		// so, find first child in the secondary index.
		// the easiest way, would just be to make get relation return table.
		// we have some minor filtering on "Object" to id.
		table := rel.GetTable()
		if next, e := nextParent(run, table, child); e != nil {
			err = e
		} else {
			ret = stream.NewObjectStream(func() (ret interface{}, okay bool) {
				if next != nil {
					if n, e := nextParent(run, table, next); e != nil {
						ret, okay = stream.Error(e)
						next = nil
					} else {
						ret, okay = stream.Value(next)
						next = n
					}
				}
				return
			})
		}
	}
	return
}

// note: can return nil object.
func nextParent(run rt.Runtime, table *index.Table, child rt.Object) (ret rt.Object, err error) {
	if i, ok := table.Secondary.FindFirst(0, child.Id().Name); ok {
		if pid := table.Secondary.Rows[i].Minor; len(pid) > 0 {
			if parent, ok := run.GetObject(pid); !ok {
				err = errutil.New("couldnt find parent", pid)
			} else {
				ret = parent
			}
		}
	}
	return
}
