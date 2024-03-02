package mdl

import (
	"database/sql"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

type MatchedMacro struct {
	Macro  Macro
	Phrase string // phrase that was used to match
}

type MatchedKind struct {
	Name  string        // the name of the kind in the db
	Base  kindsOf.Kinds // which of the built-in kinds is the returned kind most like?
	Match string        // the word(s) (singular or plural) used to match
}

type MatchedNoun struct {
	Name  string // the name of the noun in the db
	Match string // the words (singular or plural) used to match
}

// searches for the kind which uses the most number of words from the front of the passed string.
// for example: the words "container of the apocalypse" would match the kind "container."
// assumes the words are lowercased and whitespace is normalized.
func (pen *Pen) GetPartialKind(str string) (ret MatchedKind, err error) {
	// to ensure a whole word match, during query the names of the kinds are appended with blanks
	// and so we also give the phrase a final blank in case the phrase is a single word.
	if len(str) > 0 {
		var k kindInfo
		var match string
		words := str + blank
		switch e := pen.db.QueryRow(`
	with kinds(id, name, alt, path) as (
		select mk.rowid, mk.kind, mk.singular, ',' || mk.path
		from mdl_kind mk
 		join domain_tree
 		on (uses = domain)
 		where base = ?1
	)
	select id, name, match, path from (
		select id, name, name as match, substr(?2 ,0, length(name)+2) as words, path
			from kinds
		where words = (name || ' ')
		union all 
		select id, name, alt as match, substr(?2, 0, length(alt)+2) as words, path
			from kinds
		where alt is not null and words = (alt || ' ')
	)
	order by length(name) desc
	limit 1`,
			pen.domain, words).Scan(&k.id, &k.name, &match, &k.path); e {
		case nil:
			var base kindsOf.Kinds
			fullPath := k.fullpath() // walk backwards to grab "actions" before "patterns"
			for i := len(kindsOf.DefaultKinds) - 1; i >= 0; i-- {
				k := kindsOf.DefaultKinds[i]
				if strings.HasSuffix(fullPath, pen.getPath(k)) {
					base = k
					break
				}
			}
			ret = MatchedKind{
				Name: k.name,
				Base: base,
				// 99.99% of the time the singular and plural kind will be the same number of words
				// so match is superfluous. potentially could assert that expectation on creation.
				// although having arbitrary aliases for kinds (in a separate table) might be nicest.
				Match: match,
			}
		case sql.ErrNoRows:
			// return nothing when unmatched
		default:
			err = e
		}
	}
	return
}

const blank = " "
const space = ' '

// match the passed words with the known fields of all in-scope kinds.
func (pen *Pen) GetPartialField(str string) (ret string, err error) {
	if len(str) > 0 {
		words := str + blank
		if e := pen.db.QueryRow(`
	with fields(name) as (
		select distinct mf.field
		from mdl_kind mk
		join domain_tree dt
			on (dt.uses = mk.domain)
		join mdl_field mf 
			on(mf.kind = mk.rowid)	
		where dt.base = ?1
		and mf.type is null 
	)
	select name from (
		select name, substr(?2 ,0, length(name)+2) as words
		from fields
		where words = (name || ' ')
	)
	order by length(name) desc
	limit 1`,
			pen.domain, words).Scan(&ret); e != sql.ErrNoRows {
			err = e // could be nil or error
		}
	}
	return
}

func (pen *Pen) GetPartialNoun(name, kind string) (ret MatchedNoun, err error) {
	if len(name) > 0 {
		words := name + blank
		if e := pen.db.QueryRow(`
		with nouns(noun, name) as (
			select mn.noun, my.name 
			from mdl_name my
			join mdl_noun mn
				on (mn.rowid = my.noun)
			join mdl_kind mk 
				on (mn.kind = mk.rowid)
			join domain_tree dt
				on (dt.uses = my.domain)
			where base = ?1
			and my.rank >= 0
			and (?3 = "" or mk.kind = ?3)
			order by my.rank, my.rowid asc
		)
		-- for each possible pair chop a chunk of words from our input string
		-- that's the length of the noun name, to see if it matches the noun name.
		select noun, name from (
			select noun, name, substr(?2 ,0, length(name)+2) as words
			from nouns
			where words = (name || ' ')
		)
		order by length(name) desc
		limit 1`,
			pen.domain, words, kind).
			Scan(&ret.Name, &ret.Match); e != sql.ErrNoRows {
			err = e // could be nil or error
		}
	}
	return
}

// tbd: technically there's some possibility that there might be three traits:
// "wood", "veneer", and "wood veneer" -- subset names
// with the first two applying to one kind, and the third applying to a different kind;
// all in scope.  this would always match the second -- even if its not applicable.
// ( i guess that's where commas can be used by the user to separate things )
func (pen *Pen) GetPartialTrait(str string) (ret string, err error) {
	if ap, e := pen.getAspectPath(); e != nil {
		err = e
	} else if len(str) > 0 {
		words := str + blank
		if e := pen.db.QueryRow(`
	with traits(name) as (
		select distinct mf.field
		from mdl_kind mk
		join domain_tree dt
			on (dt.uses = mk.domain)
		join mdl_field mf 
			on(mf.kind = mk.rowid)	
		where dt.base = ?1
		-- ?2 is the aspect path  to filter for traits
		and instr(',' || mk.path, ?2 )
	)
	select name from (
		select name, substr(?3 ,0, length(name)+2) as words
		from traits
		where words = (name || ' ')
	)
	order by length(name) desc
	limit 1`,
			pen.domain, ap, words).Scan(&ret); e != sql.ErrNoRows {
			err = e // could be nil or error
		}
	}
	return
}

func (pen *Pen) GetPartialMacro(str string) (ret MatchedMacro, err error) {
	if len(str) == 0 {
		err = sql.ErrNoRows
	} else {
		words := str + blank
		var found struct {
			kid      int64  // id of the kind
			name     string // name of the kind/macro
			phrase   string // string of the macro phrase
			reversed bool
			result   sql.NullInt32 // from mdl_pat, number of result fields ( 0 or 1 )
		}
		if e := pen.db.QueryRow(`
	select mk.rowid, mk.kind, mg.phrase, mg.reversed, length(mp.result)>0
	from mdl_phrase mg
	join mdl_kind mk 
		on (mk.rowid = mg.macro)
	join mdl_pat mp 
		on (mp.kind = mg.macro)
	join domain_tree dt
		on (dt.uses = mg.domain)
	where base = ?1
	and (phrase || ' ') = substr(?2 ,0, length(phrase)+2)
	order by length(phrase) desc
	limit 1`, pen.domain, words).Scan(
			&found.kid, &found.name, &found.phrase, &found.reversed, &found.result); e != nil {
			err = e
		} else if parts, e := tables.QueryStrings(pen.db,
			`select affinity 
		from mdl_field 
		where kind=?1
		order by rowid`, found.kid); e != nil {
			err = e
		} else if numFields := len(parts) - int(found.result.Int32); numFields <= 0 {
			err = errutil.Fmt("most macros should have two fields and one result; has %d fields and %d returns",
				numFields, found.result.Int32)
		} else {
			var flag MacroType
			if numFields == 1 {
				flag = Macro_PrimaryOnly
			} else {
				a, b := affine.Affinity(parts[0]), affine.Affinity(parts[1])
				if a == affine.Text {
					if b == affine.Text {
						err = errutil.New("one one not supported?")
					} else if b == affine.TextList {
						flag = Macro_ManySecondary
					} else {
						err = errutil.New("unexpected aff", b)
					}
				} else if a == affine.TextList {
					if b == affine.Text {
						flag = Macro_ManyPrimary
					} else if b == affine.TextList {
						flag = Macro_ManyMany
					} else {
						err = errutil.New("unexpected aff", b)
					}
				} else {
					err = errutil.New("unexpected aff", a)
				}
			}
			if err == nil {
				ret = MatchedMacro{
					Phrase: found.phrase,
					Macro: Macro{
						Name:     found.name,
						Type:     flag,
						Reversed: found.reversed,
					},
				}
			}
		}
	}
	return
}
