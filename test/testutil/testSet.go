package testutil

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
)

func SetRecord(d *rt.Record, pairs ...any) (err error) {
	for i, cnt := 0, len(pairs); i < cnt; i = i + 2 {
		if n, ok := pairs[i].(string); !ok {
			err = errors.New("couldnt convert field")
		} else {
			if v, e := ValueOf(pairs[i+1]); e != nil {
				err = fmt.Errorf("%w couldnt convert value", e)
				break
			} else if e := d.SetNamedField(n, v); e != nil {
				err = e
				break
			}
		}
	}
	return
}

// ValueOf returns a new generic value wrapper based on analyzing the passed value.
func ValueOf(i any) (ret rt.Value, err error) {
	switch i := i.(type) {
	case bool:
		ret = rt.BoolOf(i)
	case int:
		ret = rt.IntOf(i)
	case float64:
		ret = rt.FloatOf(i)
	case string:
		ret = rt.StringOf(i)
	case []float64:
		ret = rt.FloatsOf(i)
	case []string:
		ret = rt.StringsOf(i)
	case *rt.Record:
		ret = rt.RecordOf(i)
	default:
		err = fmt.Errorf("unhandled value %v(%T)", i, i)
	}
	return
}
