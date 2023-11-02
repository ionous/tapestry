package rift_test

import (
	"reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func match(t *testing.T, name string, have, want any) (err error) {
	if strings.HasPrefix(name, `x `) {
		// commenting out tests causes go fmt to replace spaces with tabs. *sigh*
		t.Log("skipping", name)
	} else if e, ok := have.(error); ok {
		err = errutil.Fmt("ng failed %q %v", name, e)
	} else if d := pretty.Diff(want, have); len(d) != 0 {
		err = errutil.Fmt("ng mismatched %q want: %v have: %v diff: %v",
			name, want, have, d)
	} else {
		t.Logf("ok success: %q %T %v", name, have, have)
	}
	return
}

// replace statename with reflection lookup
// could be put in a charm helper package
func init() {
	charm.StateName = func(n charm.State) (ret string) {
		if s, ok := n.(interface{ String() string }); ok {
			ret = s.String()
		} else if n == nil {
			ret = "null"
		} else {
			ret = reflect.TypeOf(n).Elem().Name()
		}
		return
	}
}
