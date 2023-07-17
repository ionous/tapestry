package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

type kindCat struct {
	domain *Domain
	kind   string
}

// a record literal is defined by a kind, and a sparse set of named values
// stores literal values because they are serializable
// ( as opposed to generic values which aren't. )
type localRecord struct {
	k kindCat

	rec *literal.RecordValue // our record, containing a sparse list of values
	at  string               // fix? the at here is not very useful b/c it doesnt tell us where the individual values came from
	// but where would we store that in the db, unless we either allowed a table of multiple at values ( maybe with path )
	// or, we recorded individual values for records and assembled them at runtime.
}

// inside the record we store slightly narrower info
type innerRecord struct {
	k kindCat

	fieldValues *[]literal.FieldValue // a sparse list of values
}

func (in *innerRecord) findCompatibleField(field string, affinity affine.Affinity) (retName, retCls string, err error) {
	k := in.k
	cat := k.domain.cat
	return cat.Pin(k.domain.name, "").FindCompatibleField(k.kind, field, affinity)
}

func (rp *localRecord) isValid() bool {
	return rp.k.domain != nil
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

func (in *innerRecord) nestedWrite(noun, field, at string, val literal.LiteralValue) (err error) {
	if name, _, e := in.findCompatibleField(field, val.Affinity()); e != nil {
		err = e
	} else {
		// redo the value if setting a trait
		if name != field {
			val = &literal.TextValue{Value: field}
		}
		// if an old field exists, compare.
		if oldVal, ok := in.findField(name); ok {
			err = compareValue(noun, name, at, oldVal, val)
		} else {
			// otherwise add the new field
			in.appendField(name, val)
		}
	}
	return
}

// find or create each record, and return the innermost one.
func (rp *localRecord) ensureRecords(at string, path []string) (ret innerRecord, err error) {
	it := innerRecord{rp.k, &rp.rec.Fields}
	// drill inward, creating sub records if needed.
	for _, field := range path {
		if name, cls, e := it.findCompatibleField(field, affine.Record); e != nil {
			err = e
			break
		} else {
			nextKind := kindCat{it.k.domain, cls}
			if oldVal, ok := it.findField(name); !ok {
				nextRec := new(literal.FieldList)
				it.appendField(name, nextRec)
				it = innerRecord{nextKind, &nextRec.Fields}
			} else if nextRec, ok := oldVal.(*literal.FieldList); !ok {
				err = errutil.New("field value isnt a record")
				break
			} else {
				it = innerRecord{nextKind, &nextRec.Fields}
			}
		}
	}
	if err == nil {
		ret = it
	}
	return
}

func (in *innerRecord) appendField(name string, newVal literal.LiteralValue) {
	*in.fieldValues = append(*in.fieldValues, literal.FieldValue{
		Field: name,
		Value: newVal,
	})
}

// find the value of the named field within the passed (sparse) record.
func (in *innerRecord) findField(field string) (ret literal.LiteralValue, okay bool) {
	for _, ft := range *in.fieldValues {
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
func compareValue(noun, field, at string, oldValue, newValue literal.LiteralValue) (err error) {
	var why error = mdl.Conflict
	was, wants := field, field
	type stringer interface{ String() string }
	if try, ok := newValue.(stringer); ok {
		if curr, ok := oldValue.(stringer); ok {
			if try, curr := try.String(), curr.String(); try == curr {
				was, wants, why = curr, try, mdl.Duplicate
			}
		}
	}
	return errutil.Fmt("%w in field %q of noun %q was %v wants %v", why, noun, field, was, wants)
}
