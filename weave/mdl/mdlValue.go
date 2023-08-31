package mdl

import (
	"database/sql"
	"strings"

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

// the owner of a value
// noun can be zero for default values of kind
type ownerInfo struct {
	kind, noun int64
	debugName  string
}

func (pen *Pen) addDefaultValue(kind kindInfo, name string, value assign.Assignment) (err error) {
	interim := isProvisional(value)
	if field, e := pen.findField(kind.class(), name); e != nil {
		err = e
	} else if value, e := field.rewriteTrait(name, value); e != nil {
		err = errutil.New("can't assign trait to noun")
	} else if aff := assign.GetAffinity(value); aff != field.aff {
		err = errutil.Fmt("mismatched affinity, cant assign %s to %s", aff, field.aff)
	} else if out, e := marshalAssignment(value); e != nil {
		err = e
	} else {
		err = pen.addValue(kind.ownerInfo(), interim, field, field.name, "", out)
	}
	return
}

func (pen *Pen) addFieldValue(noun, name string, value assign.Assignment) (err error) {
	interim := isProvisional(value)
	if noun, e := pen.findRequiredNoun(noun, nounWithKind); e != nil {
		err = e
	} else if field, e := pen.findField(noun.class(), name); e != nil {
		err = e
	} else if value, e := field.rewriteTrait(name, value); e != nil {
		err = errutil.New("can't assign trait to noun")
	} else if aff := assign.GetAffinity(value); aff != field.aff {
		err = errutil.Fmt("mismatched affinity, cant assign %s to %s", aff, field.aff)
	} else if out, e := marshalAssignment(value); e != nil {
		err = e
	} else {
		err = pen.addValue(noun.ownerInfo(), interim, field, field.name, "", out)
	}
	return
}

func (pen *Pen) addPathValue(noun string, parts []string, value literal.LiteralValue) (err error) {
	interim := isProvisional(value)
	if noun, e := pen.findRequiredNoun(noun, nounWithKind); e != nil {
		err = e
	} else if outer, inner, e := pen.digField(noun, parts); e != nil {
		err = e
	} else if end := len(parts) - 1; parts[end] != inner.name {
		err = errutil.New("can't add traits to records of nouns")
	} else if aff := literal.GetAffinity(value); aff != inner.aff {
		err = errutil.Fmt("affinity %s is incompatible with %s field %q in kind %q",
			aff, inner.aff, inner.name, noun.kind)
	} else if out, e := marshalout(value); e != nil {
		err = e
	} else {
		root, dot := parts[0], strings.Join(parts[1:], ".")
		err = pen.addValue(noun.ownerInfo(), interim, outer, root, dot, out)
	}
	return
}

func (pen *Pen) addValue(owner ownerInfo, interim bool, outer fieldInfo, root, dot, value string) (err error) {
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
		owner.noun, outer.id, opt,
	); e != nil {
		err = errutil.New("database error", e)
	} else if e := tables.ScanAll(rows, func() (err error) {
		if prev.dot.String != dot {
			err = errutil.Fmt(`%w writing value for %s, had value for %s.`,
				Conflict, debugJoin(owner.debugName, root, dot), debugJoin(owner.debugName, root, prev.dot.String))
		} else if prev.value != value {
			err = errutil.Fmt(`%w mismatched value for %s.`,
				Conflict, debugJoin(owner.debugName, root, dot))
		} else {
			err = errutil.Fmt(`%w value for %s.`,
				Duplicate, debugJoin(owner.debugName, root, dot))
		}
		return
	}, &prev.dot, &prev.value); e != nil {
		err = eatDuplicates(pen.warn, e)
	} else {
		if _, e := pen.db.Exec(mdl_value, owner.noun, outer.id, opt, value, interim, pen.at); e != nil {
			err = e
		}
	}
	return
}

// note: values are written per noun, not per domain
// fix? some values are references to objects in the form "#domain::noun" -- should the be changed to ids?
var mdl_value = tables.Insert("mdl_value", "noun", "field", "dot", "value", "provisional", "at")

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
func marshalAssignment(val assign.Assignment) (ret string, err error) {
	// questionable: since we know the type of the field
	// storing the assignment wrapper is redundant.
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
