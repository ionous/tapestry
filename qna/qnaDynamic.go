package qna

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/pack"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// indicates a value stored by the user
// as opposed to a default read from the db.
// ( so save knows what to save )
type UserValue struct {
	rt.Value
}

// implements TextMarshaler
func (u UserValue) MarshalText() ([]byte, error) {
	res := pack.PackValue(u.Value)
	return []byte(res), nil // fix this cast?
}

var errMissing = errors.New("not in cache")

// hrm... needs runner to eval assignments
func (run *Runner) unpackDynamicValue(key query.Key, aff affine.Affinity, cls string) (ret rt.Value, err error) {
	if cached, ok := run.dynamicVals.Get(key); !ok {
		err = errMissing
	} else {
		switch c := cached.(type) {
		case error:
			err = c

		case UserValue:
			ret = c.Value

		case []byte:
			// a packed user value.
			// read from the save db, only unpack once and then save back:
			if v, e := pack.UnpackValue(run, c, aff, cls); e != nil {
				err = e
			} else {
				// store it as a user value so save will write it back again
				run.dynamicVals.Store(key, UserValue{v})
				ret = v
			}

		case rt.Value:
			// database literal
			ret = c

		case rt.Assignment:
			// woven values, re-evaluated on each call.
			// evaluate the assignment to get the current value
			// tbd: should there be a "this" pushed into scope?
			if v, e := safe.GetAssignment(run, c); e != nil {
				err = e
			} else {
				ret, err = safe.RectifyText(run, v, aff, cls)
			}

		default:
			err = fmt.Errorf("unexpected type in object cache %T for %v", c, key)
		}
	}
	return
}
