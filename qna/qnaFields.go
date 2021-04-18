package qna

import (
	"database/sql"
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/tables"
	"github.com/ionous/errutil"
)

// Fields implements rt.Fields: key,field,value storage for nouns, kinds, and patterns.
// It reads its data from the play database and caches the results in memory.
type Fields struct {
	activeDomains,
	activeNouns,
	activeNounList,
	ancestorsOf,
	aspectOf,
	//countOf,
	fieldsOf,
	isLike,
	kindOf,
	nameOf,
	objOf,
	patternOf,
	reciprocalOf,
	relateTo,
	relativeKinds,
	relativesOf,
	rulesFor,
	startOf,
	traitsFor,
	typeOf,
	updatePairs,
	valueOf *sql.Stmt
}

func NewFields(db *sql.DB) (ret *Fields, err error) {
	var ps tables.Prep
	f := &Fields{
		activeDomains: ps.Prep(db,
			`select 1 from run_domain where active and domain=?`),
		activeNouns: ps.Prep(db,
			// instr(X,Y) finds the first occurrence of string Y in string X
			`select 1 from 
			mdl_noun mn 
			join run_domain rd 
			where rd.active and instr(mn.noun, '#' || rd.domain || '::') = 1
			and mn.noun=?`),
		activeNounList: ps.Prep(db,
			`select mn.noun
			from mdl_noun mn 
			join mdl_kind mk 
					using (kind) 
			join run_domain rd 
			where instr(mk.kind || ',' || mk.path || ',', ?|| ',')  
			and	rd.active and instr(mn.noun, '#' || rd.domain || '::') = 1`),
		// return kind path for the named noun
		ancestorsOf: ps.Prep(db,
			`select kind || ( case path when '' then ('') else (',' || path) end ) as path
				from mdl_noun mn 
				join mdl_kind mk 
					using (kind)
				where noun=?`),
		// return the name of the aspect of the specified trait, or the empty string.
		aspectOf: ps.Prep(db,
			`select aspect
				from mdl_noun_traits 
				where (noun||'.'||trait)=?`),
		// the starting value for every field of a kind.
		// ( this is not the thing to use for objects.
		//   objects might have instance initializers or runtime )
		fieldsOf: ps.Prep(db,
			`select field, type, affinity,
					( select value 
							from mdl_start mv 
							where mv.owner=mf.kind 
							and mv.field=mf.field ) 
						as value
			from mdl_field mf
			where kind=?
			order by rowid`),
		isLike: ps.Prep(db,
			`select ? like ?`),
		kindOf: ps.Prep(db,
			`select kind
				from mdl_noun 
				where noun=?`),
		// given an id, find the name
		nameOf: ps.Prep(db,
			`select name
				from mdl_name
				join mdl_noun
					using (noun)
				where (rank>=0) and (noun=?)
				order by rank
				limit 1`),
		// given a name, find the id.
		// we filter out parser understandings (which have ranks < 0)
		// FIX: shouldnt this be limited to activeNouns?
		objOf: ps.Prep(db,
			`select noun
				from mdl_name
				join mdl_noun
					using (noun)
				where (rank>=0) and (UPPER(name)=UPPER(?))
				order by rank
				limit 1`),
		patternOf: ps.Prep(db,
			`select name, labels, result
			from mdl_pat
			where UPPER(name) = UPPER(?1)
			order by name != ?1`),
		reciprocalOf: ps.Prep(db,
			`select noun from run_pair where active and otherNoun=?1 and relation=?2`),
		// use the sqlite like function to match
		relateTo: ps.Prep(db,
			`with next as (
				select ?1 as noun, ?2 as otherNoun, ?3 as relation, ?4 as cardinality
			)
			insert or replace into run_pair
			select prev.noun, relation, prev.otherNoun, 0
			from next
			join run_pair prev 
				using (relation)
			where  ((prev.noun = next.noun and next.cardinality glob '*_one') or
					(prev.otherNoun = next.otherNoun and next.cardinality glob 'one_*')) 
			union all 
			select next.noun, relation, next.otherNoun, 1
			from next`),
		relativeKinds: ps.Prep(db,
			`select mr.kind, mr.otherKind, mr.cardinality
				from mdl_rel mr 
				where relation=?`),
		relativesOf: ps.Prep(db,
			`select otherNoun from run_pair where active and noun=?1 and relation=?2`),
		rulesFor: ps.Prep(db,
			`select 
					coalesce(mr.name, 'rule' || mr.rowid) as name,
					mr.phase,
					mr.prog
				from mdl_rule mr
				where (mr.owner=?1) and (mr.target=?2) and
					(ifnull(mr.domain,'') is '' or 
						(select 1 from run_domain rd 
						 where (rd.active and (rd.domain = mr.domain))))
					order by abs(phase), mr.rowid`),
		traitsFor: ps.Prep(db,
			`select trait
				from mdl_aspect 
				where aspect=?
				order by rank`),
		// returns either $aspect, or a kind's ancestry.
		// it is case-aware for the sake of patterns
		typeOf: ps.Prep(db,
			`select * from (
				select aspect as name,'$aspect'  as path
				from mdl_aspect 
				where rank = 0 and UPPER(aspect) = UPPER(?1)
				union all 
				select kind, kind || ( case path when '' then ('') else (',' || path) end ) 
				from mdl_kind 
				where UPPER(kind) = UPPER(?1)
				)
				order by name != ?1`),
		// instead of separately deleting old values and inserting new ones;
		// we insert and replace active ones.
		updatePairs: ps.Prep(db,
			`with next as (
			select noun, otherNoun, relation, cardinality 
			from mdl_pair 
			join mdl_rel mr 
				using (relation)
			where ?=ifnull(domain, 'entire_game')
			)
			insert or replace into run_pair
			select prev.noun, relation, prev.otherNoun, 0
				from next
				join run_pair prev 
					using (relation)
				where  ((prev.noun = next.noun and next.cardinality glob '*_one') or
						(prev.otherNoun = next.otherNoun and next.cardinality glob 'one_*')) 
			union all
			select next.noun, relation, next.otherNoun, 1 
			from next`),
		valueOf: ps.Prep(db,
			`select value, type
				from run_value 
				where noun=? and field=? 
				order by tier asc nulls last limit 1`),
	}
	if e := ps.Err(); e != nil {
		err = e
	} else {
		ret = f
	}
	return
}

func (n *Fields) UpdatePairs(domain string) (ret int, err error) {
	if res, e := n.updatePairs.Exec(domain); e != nil {
		err = e
	} else {
		ret = tables.RowsAffected(res)
	}
	return
}

func (n *Runner) SetField(target, rawField string, val g.Value) (err error) {
	field := lang.SpecialBreakcase(rawField)
	//
	if len(target) == 0 || len(field) == 0 {
		err = errutil.Fmt("invalid targeted field '%s.%s'", target, rawField)
	} else {
		switch target {
		case object.Variables:
			err = n.Stack.SetFieldByName(field, val)

		case object.Option:
			err = n.options.SetOption(field, val)

		case object.Counter:
			if aff := val.Affinity(); aff != affine.Number {
				err = errutil.Fmt("counter expected a number '%s.%s', got %s", target, rawField, aff)
			} else {
				counter := lang.Breakcase(rawField)
				n.counters[counter] = val.Int()
			}

		default:
			if target[0] == object.Prefix {
				err = errutil.Fmt("can't change reserved field '%s.%s'", target, rawField)
			} else {
				key := makeKey(target, field)
				err = n.setField(key, val)
			}
		}
	}
	return
}

func (n *Runner) setField(key keyType, val g.Value) (err error) {
	// first, check if the specified field refers to a trait
	// ( by asking for the name of an aspect. ex. "player.is_in_darkness" )
	switch aspect, e := n.GetField(object.Aspect, key.dot()); e.(type) {
	default:
		err = e // there was an unknown error
	case nil:
		if aff := val.Affinity(); aff != affine.Bool {
			err = errutil.New("can only set a trait with booleans, have", aff)
		} else if trait, e := oppositeDay(n, aspect.String(), key.field, val.Bool()); e != nil {
			err = e
		} else {
			// recurse to the g.Unknown path.
			key := keyType{key.target, aspect.String()} // ex. "player.in_darkness"
			err = n.setField(key, g.StringOf(trait))
		}

	case g.Unknown:
		// didnt refer to a trait, so just set the field normally.
		// ( to set the field, we get the field to verify it exists, and to check affinity )
		if q, e := n.getOrCache(key.target, key.field, n.queryFieldValue); e != nil {
			err = e
		} else if a := q.Affinity(); a != val.Affinity() {
			err = errutil.New("value is not", a)
		} else if v, e := g.CopyValue(val); e != nil {
			err = e
		} else {
			n.values[key] = staticValue{a, v}
		}
	}
	return
}

func oppositeDay(ks g.Kinds, aspect, trait string, b bool) (ret string, err error) {
	if b {
		ret = trait
	} else if k, e := ks.GetKindByName(aspect); e != nil {
		err = e
	} else if cnt := k.NumField(); cnt != 2 {
		err = errutil.Fmt("couldn't determine the opposite of %s.%s", aspect, trait)
	} else if i := k.FieldIndex(trait); i < 0 {
		err = errutil.Fmt("couldn't find the trait %s.%s", aspect, trait)
	} else {
		field := k.Field((i + 1) & 1)
		ret = field.Name
	}
	return
}

func (n *Runner) GetField(target, rawField string) (ret g.Value, err error) {
	switch target {
	case object.Aspect:
		// used internally: return the name of an aspect for a noun's trait
		// rawField looks like: #test::apple.w
		nounDotTrait := rawField
		ret, err = n.getOrCache(object.Aspect, nounDotTrait, func(key keyType) (ret rt.Assignment, err error) {
			var val string
			if e := n.fields.aspectOf.QueryRow(nounDotTrait).Scan(&val); e != nil {
				err = e
			} else {
				ret = staticValue{affine.Text, val}
			}
			return
		})

	case object.Counter:
		// fix: i think at some point we should have a global $counters object
		// with integer fields for everything, and we read/write serialize normally.
		counter := lang.Breakcase(rawField)
		i := n.counters[counter]
		ret = g.IntOf(i)

	case object.Domain:
		// fix,once there's a domain hierarchy:
		// store the active path and test using find in path.
		var b bool
		domain := lang.Breakcase(rawField)
		switch e := n.fields.activeDomains.QueryRow(domain).Scan(&b); e {
		case nil, sql.ErrNoRows:
			ret = g.BoolOf(b)
		default:
			err = errutil.New(target, e)
		}

	case object.Id:
		// fix: object.Value should go away...
		if tmp, e := n.GetField(object.Value, rawField); e != nil {
			err = e
		} else {
			ret = g.StringOf(tmp.String())
		}

	case object.Kind:
		objId := rawField
		ret, err = n.getOrCache(object.Kind, objId, func(key keyType) (ret rt.Assignment, err error) {
			var val string
			if e := n.fields.kindOf.QueryRow(objId).Scan(&val); e != nil {
				err = e
			} else {
				ret = staticValue{affine.Text, val}
			}
			return
		})

	case object.Kinds:
		objId := rawField
		ret, err = n.getOrCache(object.Kinds, objId, func(key keyType) (ret rt.Assignment, err error) {
			var val string
			if e := n.fields.ancestorsOf.QueryRow(objId).Scan(&val); e != nil {
				err = e
			} else {
				ret = staticValue{affine.Text, val}
			}
			return
		})

	case object.Name:
		// given an id, make sure the object should be available,
		// then return its author given name.
		objId := rawField
		if !n.activeNouns.isActive(objId) {
			err = g.UnknownObject(objId)
		} else {
			ret, err = n.getOrCache(object.Name, objId, func(key keyType) (ret rt.Assignment, err error) {
				var val string
				if e := n.fields.nameOf.QueryRow(objId).Scan(&val); e != nil {
					err = e
				} else {
					ret = staticValue{affine.Text, val}
				}
				return
			})
		}

	case object.Nouns:
		kind := lang.Breakcase(rawField)
		if rows, e := n.fields.activeNounList.Query(kind); e != nil {
			err = errutil.New(target, e)
		} else {
			var nouns []string
			var noun string
			if tables.ScanAll(rows, func() (err error) {
				nouns = append(nouns, noun)
				return
			}, &noun); e != nil {
				err = e
			} else {
				ret = g.StringsOf(nouns)
			}
		}

	case object.Option:
		ret, err = n.options.Option(rawField)

	case object.Running:
		b := n.currentPatterns.isRunning(rawField)
		ret = g.IntOf(b)

	case object.Value:
		// fix: internal object handling needs some love; i dont much like the # test.
		if strings.HasPrefix(rawField, "#") {
			objId := rawField
			if !n.activeNouns.isActive(objId) {
				// fix: differentiate b/t unknown and unavailable?
				err = g.UnknownObject(objId)
			} else {
				ret, err = n.getOrCache(object.Value, objId, func(key keyType) (ret rt.Assignment, err error) {
					ret = &qnaObject{n: n, id: objId}
					return
				})
			}
		} else {
			// given a name, find an object (id) and make sure it should be available
			// note: currently we're able to get names with spaces here " apple", so we breakcase it.
			objName := lang.Breakcase(rawField)
			ret, err = n.getOrCache(object.Value, objName, func(key keyType) (ret rt.Assignment, err error) {
				var id string
				if e := n.fields.objOf.QueryRow(objName).Scan(&id); e != nil {
					err = e
				} else {
					if !n.activeNouns.isActive(id) {
						err = g.UnknownObject(id)
					} else {
						ret = &qnaObject{n: n, id: id}
					}
				}
				return
			})
		}

	case object.Variables:
		varName := lang.SpecialBreakcase(rawField)
		ret, err = n.Stack.FieldByName(varName)

	default:
		varName := lang.SpecialBreakcase(rawField)
		key := makeKey(target, varName)
		if q, ok := n.values[key]; ok {
			ret, err = q.GetAssignedValue(n)
		} else {
			// first: loop. ask if we are trying to find the value of a trait. ( noun.trait )
			switch aspectOfTrait, e := n.GetField(object.Aspect, key.dot()); e.(type) {
			default:
				err = e
			case nil:
				// we found the aspect name from the trait
				// now we need to ask for the current value of the aspect
				aspectName := aspectOfTrait.String()
				if q, e := n.getOrCache(key.target, aspectName, n.queryFieldValue); e != nil {
					err = e
				} else {
					// return whether the object's aspect equals the specified trait.
					// ( we dont cache this value because multiple things can change it )
					ret = g.BoolOf(key.field == q.String())
				}
			case g.Unknown:
				// it wasnt a trait, so query the field value
				// fix: b/c its more common, should we do this first?
				ret, err = n.getOrCache(key.target, key.field, n.queryFieldValue)
			}
			return
		}
	}
	return
}

// FIX: see about extracting this into a small helper like qnaStart, etc.
// its possible that each of the key values ( object.Kind, etc. ) all should have their own micro class
// they could share the same map object maybe if need be.
// check the cache before asking the database for info
func (n *Runner) getOrCache(target, field string, queryFn func(key keyType) (ret rt.Assignment, err error)) (ret g.Value, err error) {
	key := makeKey(target, field)
	// first try to get the cached value
	if v, ok := n.values[key]; ok {
		ret, err = v.GetAssignedValue(n)
	} else {
		// no? call the query.
		switch val, e := queryFn(key); e {
		case nil: // success!
			n.values[key] = val
			ret, err = val.GetAssignedValue(n)

		case sql.ErrNoRows: // no data.
			v := n.values.storeError(key, key.unknown())
			ret, err = v.GetAssignedValue(n)

		default: // some other error.
			err = errutil.New("runtime error:", e)
		}
	}
	return
}

// query the db for the value of an noun's field
func (n *Runner) queryFieldValue(key keyType) (ret rt.Assignment, err error) {
	var i interface{}
	var a affine.Affinity
	if e := n.fields.valueOf.QueryRow(key.target, key.field).Scan(&i, &a); e != nil {
		err = e
	} else {
		ret, err = decodeValue(a, i)
	}
	return
}
