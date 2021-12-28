package pattern

import (
	"strconv"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/meta"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

// args run in the scope of their parent context
// they write to the record that will become the new context
func newPattern(run rt.Runtime, name string, parts []string, args []rt.Arg) (rec *g.Record, err error) {
	// create a container to hold results of args, locals, and the pending return value
	if k, e := run.GetKindByName(name); e != nil {
		err = e
	} else {
		rec = k.NewRecord()
		var labelIndex, fieldIndex int
		//
		for i, a := range args {
			n := lang.SpecialUnderscore(a.Name) // because args can be $1, $2, etc. [ fix: but why arent they correct format already?? ]
			// search for a matching label.
			if len(n) == 0 {
				err = errutil.New("unnamed arg at", i)
			} else if n[0] == '$' {
				// validate positional arguments make sense
				// ( see also: EphPatterns.Assemble )
				if argIndex(labelIndex) != n {
					break
				}
				fieldIndex, labelIndex = labelIndex, labelIndex+1
			} else {
				// search in increasing order for the next label that matches the specified argument
				// this is our soft way of allowing patterns to participate in fluid like specs with optional values.
				if at := findLabel(parts, n, labelIndex); at < 0 {
					err = errutil.New("no matching label for arg", i, n, "in", parts)
					break
				} else {
					var fn string
					if at < k.NumField() {
						fn = k.Field(at).Name
					}
					if fn == n {
						fieldIndex, labelIndex = at, at+1
					} else {
						err = errutil.Fmt("mismatched field(%s) for arg(%s) in %q", fn, n, name)
						break
					}
				}
			}
			// note: set indexed field assigns without copying
			// so these two values can become shared. (ex. lists)
			if src, e := safe.GetAssignedValue(run, a.From); e != nil {
				err = errutil.New(e, "while reading arg", i, n)
				break
			} else if val, e := autoConvert(run, k.Field(fieldIndex), src); e != nil {
				err = e
				break
			} else if e := rec.SetIndexedField(fieldIndex, val); e != nil {
				err = errutil.New(e, "while setting arg", i, n)
				break
			}
		}
	}
	return
}

// fix? allows callers to use positional arguments
// for lists could have a special RunWithVarArgs that uses a custom determineArgs
// or, allow blank names to match any arg --
// note: templates currently use positional args too.
func argIndex(i int) string {
	return "$" + strconv.Itoa(i+1)
}

// returns -1 if not found
func findLabel(labels []string, name string, startingAt int) (ret int) {
	ret = -1 // provisionally
	for i, cnt := startingAt, len(labels); i < cnt; i++ {
		if l := labels[i]; l == name {
			ret = i
			break
		}
	}
	return
}

// fix: handle differences b/t kinds, aspects, etc.?
func autoConvert(run rt.Runtime, ft g.Field, val g.Value) (ret g.Value, err error) {
	if needsConversion := ft.Affinity == affine.Text && len(ft.Type) > 0 &&
		val.Affinity() == affine.Text && len(val.Type()) == 0; !needsConversion {
		ret = val
	} else {
		// set indexed field validates the ft.Type and the val.Type match
		// we just have to give it the proper value in the first place.
		ret, err = run.GetField(meta.ObjectId, val.String())
	}
	return
}
