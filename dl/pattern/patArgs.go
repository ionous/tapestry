package pattern

import (
	"strconv"

	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

func DetermineArgs(run rt.Runtime, rec *g.Record, labels []string, args []rt.Arg) (err error) {
	// note: set indexed field assigns without copying
	var labelIndex int
	var fieldIndex int
	//
	for i, a := range args {
		n := a.Name
		// search for a matching label.
		if len(n) == 0 {
			err = errutil.New("unnamed arg at", i)
		} else if n[0] == '$' {
			// validate positional arguments make sense
			if argIndex(labelIndex) != n {
				break
			}
			fieldIndex, labelIndex = labelIndex, labelIndex+1
		} else {
			// search in increasing order for the next label that matches the specified argument
			// this is our soft way of allowing patterns to participate in fluid like specs.
			if i := findLabel(labels, n, labelIndex); i < 0 {
				err = errutil.New("has mismatched arg.", i, n)
				break
			} else {
				fieldIndex, labelIndex = i, i+1
			}
		}
		//
		field := rec.Kind().Field(fieldIndex)
		if val, e := safe.GetAssignedValue(run, a.From); e != nil {
			err = errutil.New("error determining arg", i, n, e)
			break
		} else if v, e := filterText(run, field, val); e != nil {
			err = errutil.New("error narrowing arg", i, n, e)
			break
		} else if e := rec.SetIndexedField(fieldIndex, v); e != nil {
			err = errutil.New("error setting arg", i, n, e)
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
