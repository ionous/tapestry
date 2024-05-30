package safe

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

type LabelFinder struct {
	kind         *rt.Kind
	labels       []string
	next         int
	noMoreBlanks bool // indexed fields are allowed until the first named key
}

func NewLabelFinder(run rt.Runtime, kind *rt.Kind) (ret *LabelFinder, err error) {
	// could all this be determined at assembly time?
	if labels, e := run.GetField(meta.PatternLabels, kind.Name()); e != nil {
		err = e
	} else {
		ret = &LabelFinder{kind: kind, labels: labels.Strings()}
	}
	return
}

// returns nil on success; updates internals
func (lf *LabelFinder) FindNext(k string) (ret int, err error) {
	// blank names are positional arguments
	if key := inflect.Normalize(k); len(key) == 0 {
		if lf.noMoreBlanks {
			err = fmt.Errorf("unexpected blank label %q", lf.next)
		} else if now := lf.next; now >= lf.kind.FieldCount() {
			err = fmt.Errorf("too many args %d making record %s", now, lf.kind.Name())
		} else {
			ret, lf.next = now, lf.next+1
		}
	} else if at := findLabel(lf.labels, key, lf.next); at < 0 {
		err = fmt.Errorf("no matching arg for %q in labels %q", lf.kind.Name(), strings.Join(lf.labels, ","))
	} else {
		var fn string
		if at < lf.kind.FieldCount() {
			fn = lf.kind.Field(at).Name
		}
		if fn == key {
			ret, lf.next = at, at+1
			lf.noMoreBlanks = true
		} else {
			err = fmt.Errorf("mismatched field %q", fn)
		}
	}
	return
}

// search in increasing order for the next label that matches the specified argument
// this is our soft way of allowing patterns to participate in fluid like specs with optional values.
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
