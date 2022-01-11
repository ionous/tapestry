package eph

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"github.com/ionous/errutil"
)

// a record literal is defined by a kind, and a sparse set of named values
type localRecord struct {
	k   *ScopedKind          // our kind, containing all the field definitions
	rec *literal.FieldValues // our record, containing a sparse list of values
	at  string               // fix? the at here is not very useful b/c it doesnt tell us where the individual values came from
	// but where would we store that in the db, unless we either allowed a table of multiple at values ( maybe with path )
	// or, we recorded individual values for records and assembled them at runtime.
}

func (rp *localRecord) isValid() bool {
	return rp.k != nil
}

// store the passed value at the passed path ( which points to some field nested within this record pair. )
// except for the leaf, each element of the path is expected to be the field of a sub-record.
// the leaf can be any type of field, so long as it can store the passed value.
// the noun name is passed for logging.
func (rp *localRecord) writeValue(noun, at, field string, path []string, val literal.LiteralValue) (err error) {
	if nestedRec, e := rp.ensureRecords(at, path); e != nil {
		err = e
	} else {
		err = nestedRec.nestedWrite(noun, field, at, val)
	}
	return
}

func (rp *localRecord) nestedWrite(noun, field, at string, val literal.LiteralValue) (err error) {
	if fd, e := rp.k.findCompatibleField(field, val.Affinity()); e != nil {
		err = e
	} else {
		// redo the field and value if setting a trait
		if aspect := fd.name; aspect != field {
			field, val = aspect, &literal.TextValue{Text: field}
		}
		// if an old field exists, compare.
		if oldVal, ok := rp.findField(fd.name); ok {
			err = compareValue(noun, field, at, oldVal, val)
		} else {
			// otherwise add the new field
			rp.appendField(fd.name, val)
		}
	}
	return
}

// find or create each record, and return the innermost one.
func (rp *localRecord) ensureRecords(at string, path []string) (ret localRecord, err error) {
	it := *rp
	for _, field := range path {
		if fd, e := it.k.findCompatibleField(field, affine.Record); e != nil {
			err = e
			break
		} else if nextKind, ok := it.k.domain.GetKind(fd.class); !ok {
			err = errutil.Fmt("couldnt find record of %q", fd.class)
			break
		} else if oldVal, ok := it.findField(fd.name); !ok {
			nextRec := new(literal.FieldValues)
			it.appendField(fd.name, nextRec)
			it = localRecord{nextKind, nextRec, at}
		} else if nextRec, ok := oldVal.(*literal.FieldValues); !ok {
			err = errutil.New("field value isnt a record")
			break
		} else {
			it = localRecord{nextKind, nextRec, at}
		}
	}
	if err == nil {
		ret = it
	}
	return
}

func (rp *localRecord) appendField(name string, newVal literal.LiteralValue) {
	rp.rec.Contains = append(rp.rec.Contains, literal.FieldValue{
		Field: name,
		Value: newVal,
	})
}

// find the value of the named field within the passed (sparse) record.
func (rp *localRecord) findField(field string) (ret literal.LiteralValue, okay bool) {
	for _, ft := range rp.rec.Contains {
		if ft.Field == field {
			ret, okay = ft.Value, true
			break
		}
	}
	return
}

// uses stringer ( of all things :/ ) to compare potentially conflicting values.
// FIX: if you have to do it this way... why not serialize it to compact format?
// you could even store the whole literal that way -- saving some duplication of effort during write.
func compareValue(noun, field, at string, oldValue, newValue literal.LiteralValue) error {
	why, was, wants := Redefined, field, field
	type stringer interface{ String() string }
	if try, ok := newValue.(stringer); ok {
		if curr, ok := oldValue.(stringer); ok {
			if try, curr := try.String(), curr.String(); try == curr {
				was, wants, why = curr, try, Duplicated
			}
		}
	}
	key := MakeKey(noun, field)
	return newConflict(
		key,
		why,
		Definition{key, at, was},
		wants,
	)
}
