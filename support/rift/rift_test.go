package rift_test

import (
	"reflect"
	"strings"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func compare(have any, want any) (err error) {
	if haveErr, ok := have.(error); !ok {
		if d := pretty.Diff(want, have); len(d) != 0 {
			err = errutil.Fmt("mismatched want: %v have: %v diff: %v",
				want, have, d)
		}
	} else {
		if expectErr, ok := want.(error); !ok {
			err = errutil.Fmt("failed %v", haveErr)
		} else if !strings.HasPrefix(haveErr.Error(), expectErr.Error()) {
			err = errutil.Fmt("failed %v, expected %v", haveErr, expectErr)
		}
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
