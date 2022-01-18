package pattern

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
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
		err = InitRecordFromArgs(run, rec, parts, args)
	}
	return
}

// fix: ideally this would become a much simpler "MakeRecord" --
// possible in package safe ( package generic cant rely on package rt. )
// first -- parts[] would have to go, which means we'd need to know which pieces of a pattern are "public"
// or which pieces are args and locals; possibly by separating pattern into three records.
// it'd also be nice to move "positional args" out into the caller maybe
func InitRecordFromArgs(run rt.Runtime, rec *g.Record, parts []string, args []rt.Arg) (err error) {
	var labelIndex, fieldIndex int
	k := rec.Kind()
	//
	for i, a := range args {
		// because args can be $1, $2, etc.
		// ( mainly for list operations [ fix: but why arent they in the correct format already?? ] )
		// search for a matching label.
		if n := a.Name; len(n) == 0 {
			err = errutil.New("unnamed arg at", i)
		} else if a.Name[0] == '$' {
			// validate positional arguments make sense
			// ( see also: EphPatterns.Assemble )
			if argIndex(labelIndex) != n {
				break
			}
			fieldIndex, labelIndex = labelIndex, labelIndex+1
		} else {
			n := lang.Underscore(n)
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
					err = errutil.Fmt("mismatched field(%s) for arg(%s) in %q", fn, n, k.Name())
					break
				}
			}
		}
		// note: set indexed field assigns without copying
		// so these two values can become shared. (ex. lists)
		if src, e := safe.GetAssignedValue(run, a.From); e != nil {
			err = errutil.New(e, "while reading arg", i, a.Name)
			break
		} else if val, e := safe.AutoConvert(run, k.Field(fieldIndex), src); e != nil {
			err = e
			break
		} else if e := rec.SetIndexedField(fieldIndex, val); e != nil {
			err = errutil.New(e, "while setting arg", i, a.Name)
			break
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
