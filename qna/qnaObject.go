package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

// qnaObject implements the generic Value interface
// FIX? should it ( and meta.Value, affine.Object ) go away?
//
// currently, it feels like we have two ways to access object values
// 		1. run.Get/SetField(object, field)
// 		2. run.Get/SetField(meta.Value, objectName) -> returns this qnaObject which acts like a g.Record
// right now the code is only using #2 --- but why? why have both?
//
// related issues:
// 	  1. object references ( both in properties and local variables )
//    2. object local variables -- vs. object text variables storing ( #names )
//    3. affine.Object
//
// to me if feels like either nouns should become named records,
// or: objects should only be exposed via GetField
//
// TBD: where does tapestry treat object as a value? what are the returns from meta.Value passed to.
//
type qnaObject struct {
	g.PanicValue
	run    *Runner // for pointing back to the field cache
	domain string
	name   string
	kind   *g.Kind
}

func (obj *qnaObject) Affinity() affine.Affinity {
	return affine.Object
	//  FIX: is there any reason not to return record???
	/// ( record could implement String(), and regular records could be "unnamed" empty string )
	// well... Record isnt an interface --
	// so we'd have to "snapshot" the whole noun at once --
	// the only reason that's bad at all is "speed"
	// ( but picking db values one by one is not necessarily fast anyway )
	// and save ( but we could keep a list of all values )
}

func (obj *qnaObject) String() (ret string) {
	// hrmm...
	return "#" + obj.domain + "::" + obj.name
}

func (obj *qnaObject) Kind() (ret *g.Kind) {
	return obj.kind

}
func (obj *qnaObject) Type() (ret string) {
	return obj.kind.Name()
}

func (obj *qnaObject) FieldByName(rawField string) (ret g.Value, err error) {
	field := lang.Underscore(rawField) // fix: why here?
	//
	if i := obj.kind.FieldIndex(field); i < 0 {
		err = g.UnknownField(obj.name, rawField)
	} else {
		// just a regular field?
		if ft := obj.kind.Field(i); ft.Name == field {
			ret, err = getObjectField(obj.run, obj.domain, obj.name, field)
		} else {
			// asking for a trait; so ft.Name is now the aspect field
			if v, e := getObjectField(obj.run, obj.domain, obj.name, ft.Name); e != nil {
				err = e
			} else {
				// return true if the aspect field holds the particular requested field
				trait := v.String()
				ret = g.BoolFrom(trait == field, "" /*"trait"*/)
			}
		}
	}
	return
}

func (obj *qnaObject) SetFieldByName(rawField string, val g.Value) (err error) {
	field := lang.Underscore(rawField)
	if i := obj.kind.FieldIndex(field); i < 0 {
		err = g.UnknownField(obj.name, rawField)
	} else {
		// just a regular field?
		if ft := obj.kind.Field(i); ft.Name == field {
			setObjectField(obj.run, obj.domain, obj.name, field, val)
		} else {
			// asking for a trait
			// FIX: records dont have opposite day so this seems ... unfair.
			// FIX: im also curious about aspects that only have one trait, and ... blank ( nothing ).
			if aff := val.Affinity(); aff != affine.Bool {
				err = errutil.New("can only set a trait with booleans, have", aff)
			} else if trait, e := oppositeDay(obj.run, ft.Name, field, val.Bool()); e != nil {
				err = e
			} else {
				// set the aspect to the value of the requested trait
				setObjectField(obj.run, obj.domain, obj.name, ft.Name, g.StringFrom(trait /*trait*/, ""))
			}
		}
	}
	return
}

// FIX? object properties can hold templates ( ex. a room description ) so this generate, not hold, values.
// but mdl_value is storing literals right now ( and note: affinity there is redundant -- its the affinity of the field. )
// [ ideas: you could switch on the first character being { or something -- escaping { by putting it in a command
//   or, you could make everything constant a literal value -- should be fine just for the second ]
func getObjectField(run *Runner, domain, noun, field string) (ret g.Value, err error) {
	if c, e := run.values.cache(func() (ret interface{}, err error) {
		if _, v, e := run.qdb.NounValue(noun, field); e != nil {
			err = e
		} else if str, ok := v.(string); !ok {
			err = errutil.New("todo: see notes")
		} else {
			ret = g.StringOf(str)
		}
		return
	}, domain, noun, field); e != nil {
		err = e
	} else {
		ret = c.(g.Value)
	}
	return
}

// both obj and field are normalized, and field is not a trait
func setObjectField(run *Runner, domain, noun, field string, val g.Value) (err error) {
	if v, e := g.CopyValue(val); e != nil {
		err = e
	} else {
		// FIX: unless the object also includes domain, and we push to
		key := makeKey(domain, noun, field)
		run.nounValues[key] = cachedValue{v: v}
	}
	return
}
