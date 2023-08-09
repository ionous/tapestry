package assign

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// args run in the scope of their parent context
// they write to the record that will become the new context
func MakeRecord(run rt.Runtime, recordName string, args ...Arg) (ret *g.Record, err error) {
	// create a container to hold results of args, locals, and the pending return value
	if kind, e := run.GetKindByName(recordName); e != nil {
		err = e
	} else if rec := kind.NewRecord(); len(args) == 0 {
		// no args specified? oh this is the easy way.
		ret = rec
	} else if lf, e := safe.NewLabelFinder(run, kind); e != nil {
		err = e
	} else {
		for i, a := range args {
			if i >= kind.NumField() {
				err = errutil.New("too many args", i, "making record", kind)
			} else if at, e := lf.FindNext(a.Name); e != nil {
				err = errutil.Fmt("%w while reading arg %d(%s)", e, i, a.Name)
				break
			} else if at < 0 {
				break
			} else if src, e := safe.GetAssignment(run, a.Value); e != nil {
				err = errutil.Fmt("%w while reading arg %d(%s)", e, i, a.Name)
				break
			} else if val, e := safe.AutoConvert(run, kind.Field(at), src); e != nil {
				err = e
				break
			} else if e := rec.SetIndexedField(at, val); e != nil {
				// note: set indexed field assigns without copying
				// but get value copies out, so this should be okay.
				err = errutil.Fmt("%w while setting arg %d(%s)", e, i, a.Name)
				break
			}
		}
		if err == nil {
			ret = rec
		} else {
			err = errutil.Fmt("%w for %q", err, recordName)
		}
	}
	return
}
