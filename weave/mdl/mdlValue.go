package mdl

import (
	"database/sql"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

func MakePath(path ...string) string {
	return strings.Join(path, ".")
}

func (pen *Pen) addDefaultValue(kind kindInfo, name string, value assign.Assignment) (err error) {
	if field, e := pen.findField(kind.class(), name); e != nil {
		err = e
	} else if value, e := field.rewriteTrait(name, value); e != nil {
		err = errutil.Fmt("can't assign trait %q to kind %q", name, kind.name)
	} else if out, provisional, e := marshalAssignment(value, field.aff); e != nil {
		err = e
	} else {
		err = pen.addKindValue(kind, !provisional, field, out)
	}
	return
}

func (pen *Pen) addFieldValue(noun, name string, value assign.Assignment) (err error) {
	if noun, e := pen.findRequiredNoun(noun, nounWithKind); e != nil {
		err = e
	} else if field, e := pen.findField(noun.class(), name); e != nil {
		err = e
	} else if value, e := field.rewriteTrait(name, value); e != nil {
		err = errutil.New("can't assign trait to noun")
	} else if out, provisional, e := marshalAssignment(value, field.aff); e != nil {
		err = e
	} else if noun.domain != pen.domain {
		err = DomainValueError{noun.name, field.name, value}
	} else {
		err = pen.addNounValue(noun, !provisional, field, field.name, "", out)
	}
	return
}

type DomainValueError struct {
	Noun, Field string
	Value       assign.Assignment
}

func (e DomainValueError) Error() string {
	return errutil.Sprint("initial values for noun %q (%q) must be in the same domain as its declaration.",
		e.Noun, e.Field)
}

// dot values are required to be literals.
func (pen *Pen) addPathValue(noun string, parts []string, value literal.LiteralValue) (err error) {
	if noun, e := pen.findRequiredNoun(noun, nounWithKind); e != nil {
		err = e
	} else if outer, inner, e := pen.digField(noun, parts); e != nil {
		err = e
	} else if end := len(parts) - 1; parts[end] != inner.name {
		err = errutil.New("can't add traits to records of nouns")
	} else if out, provisional, e := marshalLiteral(value, inner.aff); e != nil {
		err = e
	} else {
		field, dot := parts[0], strings.Join(parts[1:], ".")
		err = pen.addNounValue(noun, !provisional, outer, field, dot, out)
	}
	return
}

// writing to fields inside a record is permitted so long as the record itself has not been written to.
// overwriting a field with a record is allowed from another domain.
func (pen *Pen) addNounValue(noun nounInfo, final bool, outer fieldInfo, field, dot, value string) (err error) {
	opt := sql.NullString{
		String: dot,
		Valid:  len(dot) > 0,
	}
	// search for existing paths which conflict:
	// could be the same path, or could be a record written as a whole
	// now being written as a part; or vice versa.
	// OR the exact match ( ex. a duplicate )
	var prev struct {
		dot   sql.NullString
		value string
	}
	// find cases where the new path starts a previous path,
	// or a previous path starts the new path.
	// instr(X,Y) - searches X for Y.
	if rows, e := pen.db.Query(`
		select mv.dot, mv.value 
		from mdl_value mv 
		where mv.noun = @1
		and mv.field = @2
		and (
			 (1 == instr(case when mv.dot is null then "." else '.' || mv.dot || '.' end, 
			             case when   @3   is null then "." else '.' ||   @3   || '.' end)) or 
		   (1 == instr(case when   @3   is null then "." else '.' ||   @3   || '.' end, 
		               case when mv.dot is null then "." else '.' || mv.dot || '.' end))
		)`,
		noun.id, outer.id, opt,
	); e != nil {
		err = errutil.New("database error", e)
	} else if e := tables.ScanAll(rows, func() (err error) {
		if prev.dot.String != dot {
			err = errutil.Fmt(`%w writing noun value for %s, had value for %s.`,
				Conflict, debugJoin(noun.name, field, dot), debugJoin(noun.name, field, prev.dot.String))
		} else if prev.value != value {
			err = errutil.Fmt(`%w mismatched noun value for %s.`,
				Conflict, debugJoin(noun.name, field, dot))
		} else {
			err = errutil.Fmt(`%w noun value for %s.`,
				Duplicate, debugJoin(noun.name, field, dot))
		}
		return
	}, &prev.dot, &prev.value); e != nil {
		err = eatDuplicates(pen.warn, e)
	} else if noun.domain != pen.domain {
		// this to simplify domain management (ex. would have to check rival values)
		// and avoids questions about what happens to values at the *end* of domain
		// (ex. do the values revert back to their previous dynamic value?
		//  or, are they forced to the values at the start of the parent scene, etc. )
		err = errutil.Fmt("assignments to noun %q (at %q) must be in the domain %q, was %q",
			noun.name, field, noun.domain, pen.domain)
	} else {
		if _, e := pen.db.Exec(mdl_value, noun.id, outer.id, opt, value, final, pen.at); e != nil {
			err = e
		}
	}
	return
}

// note: values are written per noun, not per domain
// fix? some values are references to objects in the form "#domain::noun" -- should the be changed to ids?
var mdl_value = tables.Insert("mdl_value", "noun", "field", "dot", "value", "final", "at")

// writing to fields inside a record is permitted so long as the record itself has not been written to.
// overwriting a field with a record is allowed from another domain.
func (pen *Pen) addKindValue(kind kindInfo, final bool, field fieldInfo, value string) (err error) {
	// search for existing paths which conflict:
	// could be the same path, or could be a record written as a whole
	// now being written as a part; or vice versa.
	// OR the exact match ( ex. a duplicate )
	var prev struct {
		value string
		final bool
	}
	// find cases where the new path starts a previous path,
	// or a previous path starts the new path.
	// instr(X,Y) - searches X for Y.
	if rows, e := pen.db.Query(`
		select mv.value, mv.final
		from mdl_value_kind mv 
		where mv.kind = @1
		and mv.field = @2`,
		kind.id, field.id,
	); e != nil {
		err = errutil.New("database error", e)
	} else if e := tables.ScanAll(rows, func() (err error) {
		if prev.final {
			if prev.value != value {
				err = errutil.Fmt(`%w mismatched kind value for %s.`,
					Conflict, debugJoin(kind.name, field.name, ""))
			} else {
				err = errutil.Fmt(`%w kind value for %s.`,
					Duplicate, debugJoin(kind.name, field.name, ""))
			}
		}
		return
	}, &prev.value, &prev.final); e != nil {
		err = eatDuplicates(pen.warn, e)
	} else if field.domain != pen.domain {
		// this to simplify domain management (ex. would have to check rival values)
		err = errutil.Fmt("the domain of the assignment (%s) must match the field %q domain (%s)",
			pen.domain, field.name, field.domain)
	} else if _, e := pen.db.Exec(mdl_value_kind, kind.id, field.id, value, final, pen.at); e != nil {
		err = e
	}
	return
}

var mdl_value_kind = tables.Insert("mdl_value_kind", "kind", "field", "value", "final", "at")

func debugJoin(noun, field, path string) string {
	var b strings.Builder
	b.WriteRune('\'')
	b.WriteString(noun)
	b.WriteRune('.')
	b.WriteString(field)
	if len(path) > 0 {
		b.WriteRune('.')
		b.WriteString(path)
	}
	b.WriteRune('\'')
	return b.String()
}

// matches with decode.parseEval
func marshalAssignment(val assign.Assignment, wantAff affine.Affinity) (ret string, provisional bool, err error) {
	// questionable: since we know the type of the field
	// storing the assignment wrapper is redundant.
	if a, ok := val.(ProvisionalAssignment); ok {
		provisional = true
		val = a.Assignment
	}
	if aff := assign.GetAffinity(val); aff != wantAff {
		err = errutil.Fmt("assignment is %s not %s", aff, wantAff)
	} else {
		// strip off the From section to avoid serializing redundant info
		switch v := val.(type) {
		case *assign.FromBool:
			ret, err = marshalout(v.Value)
		case *assign.FromNumber:
			ret, err = marshalout(v.Value)
		case *assign.FromText:
			ret, err = marshalout(v.Value)
		case *assign.FromRecord:
			ret, err = marshalout(v.Value)
		case *assign.FromNumList:
			ret, err = marshalout(v.Value)
		case *assign.FromTextList:
			ret, err = marshalout(v.Value)
		case *assign.FromRecordList:
			ret, err = marshalout(v.Value)
		default:
			err = errutil.New("unknown type")
		}
	}
	return
}

func marshalLiteral(val literal.LiteralValue, wantAff affine.Affinity) (ret string, provisional bool, err error) {
	// questionable: since we know the type of the field
	// storing the assignment wrapper is redundant.
	if a, ok := val.(ProvisionalLiteral); ok {
		provisional = true
		val = a.LiteralValue
	}
	if aff := literal.GetAffinity(val); aff != wantAff {
		err = errutil.Fmt("literal is %s not %s", aff, wantAff)
	} else {
		ret, err = marshalout(val)
	}
	return
}

// shared generic marshal prog to text
func marshalout(cmd any) (ret string, err error) {
	if cmd != nil {
		if op, ok := cmd.(jsn.Marshalee); !ok {
			err = errutil.Fmt("can only marshal autogenerated types (%T)", cmd)
		} else {
			ret, err = marshalop(op)
		}
	}
	return
}

func marshalprog(prog []rt.Execute) (ret string, err error) {
	if len(prog) > 0 {
		slice := rt.Execute_Slice(prog)
		if out, e := marshalop(&slice); e != nil {
			err = e
		} else {
			ret = out
		}
	}
	return
}

func marshalop(op jsn.Marshalee) (string, error) {
	// fix:shouldn't this be core?
	return cout.Marshal(op, literal.CompactEncoder)
}
