package qna

import (
	"database/sql"
	"strconv"

	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/tables"
)

// Fields implements rt.Fields: key,field,value storage for nouns, kinds, and patterns.
// It reads its data from the play database and caches the results in memory.
type Fields struct {
	pairs mapType
	valueOf,
	patternAt,
	progBytes,
	countOf,
	ancestorsOf,
	kindOf,
	aspectOf,
	nameOf,
	idOf *sql.Stmt
}

type keyType struct {
	owner, member string
}

func (k *keyType) dot() string {
	return k.owner + "." + k.member
}

type mapType map[keyType]interface{}

type mapTarget struct {
	key   keyType
	pairs mapType
	value interface{}
}

func (k *mapTarget) Scan(v interface{}) (err error) {
	// bytes will need special processing ( copies )
	k.pairs[k.key], k.value = v, v
	return
}

func NewFields(db *sql.DB) (ret *Fields, err error) {
	var ps tables.Prep
	f := &Fields{
		pairs: make(mapType),
		valueOf: ps.Prep(db,
			`select value 
				from run_value 
				where noun=? and field=? 
				order by tier asc nulls last limit 1`),
		patternAt: ps.Prep(db,
			`select param from mdl_pat where pattern=? and idx=?`),
		progBytes: ps.Prep(db,
			`select bytes 
				from mdl_prog
				where name = ?
				and type = ?
				limit 1`),
		countOf: ps.Prep(db,
			`select count() from run_noun where noun=?`),
		ancestorsOf: ps.Prep(db,
			`select kind || ( case path when '' then ('') else (',' || path) end ) as path
				from mdl_noun mn 
				join mdl_kind mk 
					using (kind)
				where noun=?`),
		kindOf: ps.Prep(db,
			`select kind from mdl_noun where noun=?`),
		// return the name of the aspect of the specified trait, or the empty string.
		aspectOf: ps.Prep(db,
			`select ifnull(max(aspect),"") from mdl_noun_traits 
				where (noun||'.'||trait)=?`),
		// given an id, find the name
		nameOf: ps.Prep(db,
			`select name 
				from mdl_name
				join run_noun
					using (noun)
				where noun=?
				order by rank
				limit 1`),
		// given a name, find the id
		idOf: ps.Prep(db,
			`select noun 
				from mdl_name
				join run_noun
					using (noun)
				where name=?
				order by rank
				limit 1`),
	}
	if e := ps.Err(); e != nil {
		err = e
	} else {
		ret = f
	}
	return
}

func (n *Fields) SetField(obj, field string, v interface{}) (err error) {
	if len(field) == 0 || field[0] == object.Prefix || field == object.Name {
		err = errutil.Fmt("can't change reserved field %q", field)
	} else {
		// fix, future: verify type and existence?
		key := newKey(obj, field)
		// check if the specified field is a trait
		if a, e := n.GetField(key.dot(), object.Aspect); e != nil {
			err = e
		} else {
			// no, just set the field normally.
			if aspect := a.(string); len(aspect) == 0 {
				n.pairs[key] = v
			} else {
				// yes, then we want to change the aspect not the trait
				if val, ok := v.(bool); !ok || !val {
					err = errutil.Fmt("%q can only be set to true; have %T(%v)", key, v, v)
				} else {
					// set
					err = n.SetField(obj, aspect, field)
				}
			}
		}
	}
	return
}

func newKey(obj, field string) keyType {
	// FIX FIX FIX --
	// operations generating get field should be registering the field as a name
	// and, as best as possible, relating obj to field for property verification
	// name translation should be done there.
	if len(field) > 0 && field[0] != object.Prefix {
		field = lang.Camelize(field)
	}
	return keyType{obj, field}
}

func newKeyForEval(obj, typeName string) keyType {
	return keyType{obj, typeName}
}

func newKeyWithIndex(obj string, idx int) keyType {
	return keyType{obj, "$" + strconv.Itoa(idx)}
}

// pv is a pointer to a pattern instance, and we copy its contents in.
func (n *Fields) GetEvalByName(name string, pv interface{}) (err error) {
	outVal := r.ValueOf(pv).Elem() // outVal is a pattern instance who's members get overwritten
	rtype := outVal.Type()
	// note: newKey camelCases, while go types are PascalCase
	// this automatically keeps them from conflicting.
	key := newKeyForEval(name, rtype.Name())
	if val, ok := n.pairs[key]; ok {
		store := r.ValueOf(val)
		outVal.Set(store)
	} else {
		var store interface{}
		switch e := n.progBytes.QueryRow(key.owner, key.member).Scan(&tables.GobScanner{outVal}); e {
		case nil:
			store = outVal.Interface()
		case sql.ErrNoRows:
			err = fieldNotFound{key.owner, key.member}
		default:
			err = e
		}
		n.pairs[key] = store
	}
	return
}

func (n *Fields) GetField(obj, field string) (ret interface{}, err error) {
	key := newKey(obj, field)
	if val, ok := n.pairs[key]; ok {
		ret = val
	} else {
		// note: uses the normalized member name, not the raw parameter name
		switch field := key.member; field {
		case object.Name:
			// search for the object name using the object's id
			ret, err = n.getCachingQuery(key, n.nameOf, obj)

		case object.Id:
			// search for the object id by a partial object name
			ret, err = n.getCachingQuery(key, n.idOf, obj)

		case object.Aspect:
			ret, err = n.getCachingQuery(key, n.aspectOf, obj)

		case object.Kind:
			ret, err = n.getCachingQuery(key, n.kindOf, obj)

		case object.Kinds:
			ret, err = n.getCachingQuery(key, n.ancestorsOf, obj)

		case object.Exists:
			// searches for an exact name match
			ret, err = n.getCachingQuery(key, n.countOf, obj)

		default:
			// see if the user is asking for the status of a trait
			if a, e := n.GetField(key.dot(), object.Aspect); e != nil {
				err = e
			} else if aspect := a.(string); len(aspect) > 0 {
				ret, err = n.getCachingStatus(obj, aspect, field)
			} else {
				ret, err = n.getCachingField(key)
			}
		}
	}

	return
}

// returns the name of a field based on an index
// ex. especially for resolving positional pattern parameters into names.
func (n *Fields) GetFieldByIndex(obj string, idx int) (ret string, err error) {
	if idx <= 0 {
		err = errutil.New("GetFieldByIndex out of range", idx)
	} else {
		// first, lookup the parameter name
		key := newKeyWithIndex(obj, idx)
		// we use the cache to keep $(idx) -> param name.
		val, ok := n.pairs[key]
		if !ok {
			val, err = n.getCachingQuery(key, n.patternAt, obj, idx)
		}
		if field, ok := val.(string); !ok {
			err = fieldNotFound{key.owner, key.member}
		} else {
			ret = field
		}
	}
	return
}

// return true if the object's aspect equals the specified trait.
func (n *Fields) getCachingStatus(obj, aspect, trait string) (ret bool, err error) {
	if val, e := n.GetField(obj, aspect); e != nil {
		err = e
	} else {
		ret = val == trait
	}
	return
}

func (n *Fields) GetCachingField(obj, field string) (ret interface{}, err error) {
	key := newKey(obj, field)
	return n.getCachingField(key)
}

func (n *Fields) getCachingField(key keyType) (ret interface{}, err error) {
	// FIX? needs more work to determine if the field really exists
	// ex. possibly a union query of class field with a nil value
	if v, e := n.getCachingQuery(key, n.valueOf, key.owner, key.member); e == nil {
		ret = v
	} else if _, ok := e.(fieldNotFound); !ok {
		err = e
	} else {
		n.pairs[key] = nil
		ret = nil
	}
	return
}

// getCachingQuery uses the rowscanner to write the results of a query into the cache
func (n *Fields) getCachingQuery(key keyType, q *sql.Stmt, args ...interface{}) (ret interface{}, err error) {
	tgt := mapTarget{key: key, pairs: n.pairs}
	switch e := q.QueryRow(args...).Scan(&tgt); e {
	case nil:
		ret = tgt.value
	case sql.ErrNoRows:
		err = fieldNotFound{key.owner, key.member}
	default:
		err = e
	}
	return
}

var notImplemented = errutil.New("not implemented")
