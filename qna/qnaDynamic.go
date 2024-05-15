package qna

import (
	"database/sql"
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/pack"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/tables"
)

type dynamicVals struct {
	cache
}

// indicates a value stored by the user
// as opposed to a default read from the db.
// ( so save knows what to save )
type UserValue struct {
	rt.Value
}

var errMissing = errors.New("not in cache")

// read noun data from the db ( rt.run_value )
// and store them as bytes; run.readNounValue will unpack them
func (d *dynamicVals) readValues(db *sql.DB) (err error) {
	var key qkey
	var b []byte
	out := make(cacheMap)
	if rows, e := db.Query(
		`select domain, noun, field, value from run_value
		where domain != ''`); e != nil {
		err = e
	} else if e := tables.ScanAll(rows, func() (_ error) {
		out[key] = b
		return
	}, &key.domain, &key.target, &key.field, &b); e != nil {
		err = e
	} else {
		d.cache.store = out
	}
	return
}

// write all dynamic values to the database using the prepared 'runValue' statement.
func (d dynamicVals) writeValues(w writeCb) (err error) {
	for key, cached := range d.store {
		// all stored values are variant values ( unless they are errors )
		// doesnt save the errors... we likely wouldn't get here if there were any
		if userVal, ok := cached.(UserValue); ok {
			if str, e := pack.PackValue(userVal); e != nil {
				err = e
			} else if e := w(key.domain, key.target, key.field, str); e != nil {
				err = e
				break
			}
		}
	}
	return
}

// hrm... needs runner to eval assignments
func (run *Runner) unpackDynamicValue(key qkey, aff affine.Affinity, cls string) (ret rt.Value, err error) {
	if cached, ok := run.dynamicVals.store[key]; !ok {
		err = errMissing
	} else {
		switch c := cached.(type) {
		case error:
			err = c

		case UserValue:
			ret = c.Value

		case []byte:
			// read from the save db, only unpack once and then save back:
			if v, e := pack.UnpackValue(run, c, aff, cls); e != nil {
				err = e
			} else {
				// store it as a user value so save will write it back again
				run.dynamicVals.store[key] = UserValue{v}
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
			err = fmt.Errorf("unexpected type in object cache %T for noun %q field %q", c, key.target, key.field)
		}
	}
	return
}
