package assign

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// args run in the scope of their parent context
// they write to the record that will become the new context
func MakeRecord(run rt.Runtime, kind string, args ...Arg) (ret *g.Record, err error) {
	// create a container to hold results of args, locals, and the pending return value
	if kind, e := run.GetKindByName(kind); e != nil {
		err = e
	} else if rec := kind.NewRecord(); len(args) == 0 {
		// no args specified? oh this is the easy way.
		ret = rec
	} else {
		lf := labelFinder{kind: kind}
		for i, a := range args {
			if at, e := lf.findNext(run, i, a); e != nil {
				err = e
				break
			} else if at < 0 {
				break
			} else if src, e := GetSafeAssignment(run, a.Value); e != nil {
				err = errutil.New(e, "while reading arg", i, a.Name)
				break
			} else if val, e := safe.AutoConvert(run, kind.Field(at), src); e != nil {
				err = e
				break
			} else if e := rec.SetIndexedField(at, val); e != nil {
				// note: set indexed field assigns without copying
				// but get value copies out, so this should be okay.
				err = errutil.Fmt("%e while setting arg %v(%d) at %d", e, a.Name, i, at)
				break
			}
		}
		if err == nil {
			ret = rec
		}
	}
	return
}

type labelFinder struct {
	kind         *g.Kind
	labels       []string
	next         int
	noMoreBlanks bool // error checking
}

// returns nil on success; updates internals
func (lf *labelFinder) findNext(run rt.Runtime, i int, a Arg) (ret int, err error) {
	// blank names are positional arguments
	if n := a.Name; len(n) == 0 {
		if lf.noMoreBlanks {
			err = errutil.New("unexpected blank label", i)
		} else {
			ret, lf.next = lf.next, lf.next+1
		}
	} else {
		// otherwise, find the named argument
		if labels, e := lf.getLabels(run); e != nil {
			err = e
		} else {
			n := lang.Underscore(n)
			// search in increasing order for the next label that matches the specified argument
			// this is our soft way of allowing patterns to participate in fluid like specs with optional values.
			if at := findLabel(labels, n, lf.next); at < 0 {
				err = errutil.New("no matching label for arg", i, n, "in", lf.labels)
			} else {
				var fn string
				if at < lf.kind.NumField() {
					fn = lf.kind.Field(at).Name
				}
				if fn == n {
					ret, lf.next = at, at+1
					lf.noMoreBlanks = true
				} else {
					err = errutil.Fmt("mismatched field(%s) for arg(%s) in %q", fn, n, lf.kind.Name())
				}
			}
		}
	}
	return
}

// could all this be determined at assembly time?s
func (lf *labelFinder) getLabels(run rt.Runtime) (ret []string, err error) {
	if lf.labels != nil {
		ret = lf.labels
	} else if labels, e := run.GetField(meta.PatternLabels, lf.kind.Name()); e != nil {
		err = e
	} else {
		lf.labels = labels.Strings()
		ret = lf.labels
	}
	return
}

// returns -1 if not found, but startingAt if there are no labels at all
// ( no labels indicates a CallPattern is being used for record initialization )
func findLabel(labels []string, name string, startingAt int) (ret int) {
	if cnt := len(labels); cnt == 0 {
		ret = startingAt
	} else {
		ret = -1 // provisionally
		for i := startingAt; i < cnt; i++ {
			if l := labels[i]; l == name {
				ret = i
				break
			}
		}
	}
	return
}
