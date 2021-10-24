package importer

import (
	"errors"

	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/ionous/errutil"
)

// expects that the very next element is a key
func Key(want string, fn func(string) error) *chart.StateMix {
	return &chart.StateMix{
		OnKey: func(lede, have string) (err error) {
			if have == want {
				err = fn(lede)
			} else {
				err = chart.Unhandled(errutil.Sprint("importer key: wanted", want, "have", have))
			}
			return
		},
	}
}

// expects that the very next value is a block
func Block(want string, fn func(jsn.Block) error) *chart.StateMix {
	return &chart.StateMix{
		OnBlock: func(b jsn.Block) (err error) {
			if have := b.GetType(); have == want {
				err = fn(b)
			} else {
				err = chart.Unhandled(errutil.Sprint("importer block: wanted", want, "have", have))
			}
			return
		},
	}
}

// expects that the very next value is a simple value
func Value(want string, fn func(interface{}) error) *chart.StateMix {
	return &chart.StateMix{
		OnValue: func(have string, v interface{}) (err error) {
			if have == want {
				err = fn(v)
			} else {
				err = chart.Unhandled(errutil.Sprint("importer value: wanted", want, "have", have))
			}
			return
		},
	}
}

func Eventually(want jsn.State) *chart.StateMix {
	var unhandled chart.Unhandled
	return &chart.StateMix{
		OnBlock: func(b jsn.Block) (err error) {
			if e := want.MarshalBlock(b); !errors.As(e, &unhandled) {
				err = e
			}
			return
		},
		OnKey: func(a, b string) (err error) {
			if e := want.MarshalKey(a, b); !errors.As(e, &unhandled) {
				err = e
			}
			return
		},
		OnValue: func(a string, b interface{}) (err error) {
			if e := want.MarshalValue(a, b); !errors.As(e, &unhandled) {
				err = e
			}
			return
		},
	}
}
