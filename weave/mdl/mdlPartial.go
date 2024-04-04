package mdl

import (
	"database/sql"
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

type MatchedKind struct {
	Name  string        // the name of the kind in the db
	Base  kindsOf.Kinds // which of the built-in kinds is the returned kind most like?
	Match string        // the word(s) (singular or plural) used to match
}

func (m MatchedKind) WordCount() int {
	return countWords(m.Match)
}

type MatchedField struct {
	Name string // name of the trait in the db
}

// same logic as MatchedTrait
func (m MatchedField) WordCount() int {
	return countWords(m.Name)
}

type MatchedNoun struct {
	Name  string // the id of the noun in the db
	Kind  string // the noun's kind
	Match string // the words used to match
}

func (m MatchedNoun) WordCount() int {
	return countWords(m.Match)
}

type MatchedTrait struct {
	Name string // name of the trait in the db
}

// the returned name is the name of the trait from the db
// it was used to match the front of the passed string
// so the words in the trait are the words in the string.
func (m MatchedTrait) WordCount() int {
	return countWords(m.Name)
}

func countWords(str string) (ret int) {
	if len(str) > 0 {
		ret = 1 + strings.Count(str, " ")
	}
	return
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
// return the full name of the field that matched.
// an unmatched noun returns the empty string and no error.
func (pen *Pen) GetPartialField(kind, field string) (ret MatchedField, err error) {
	if kind, e := pen.findRequiredKind(kind); e != nil {
		err = e
	} else if len(field) == 0 {
		err = errors.New("get partial field requires a non empty string")
	} else {
		words := field + blank
		if e := pen.db.QueryRow(`
	with fields(name) as (
		select field as name 		 
		from mdl_field mf 
		join mdl_kind mk 
			-- does our ancestry (X) contain any of these kinds (Y)
			on ((mf.kind = mk.rowid) and instr(@ancestry, ',' || mk.rowid || ',' ))
		left join mdl_kind mt 
			on (mt.rowid = mf.type)
	)
	select name from (
		select name, substr(@words ,0, length(name)+2) as words
		from fields
		where words = (name || ' ')
	)
	order by length(name) desc
	limit 1`,
			sql.Named("ancestry", kind.fullpath()),
			sql.Named("words", words)).Scan(&ret.Name); e != sql.ErrNoRows {
			err = e // could be nil or error
		}
	}
	return
}

// if specified, kind must match exactly
// an unmatched noun returns the empty string and no error.
func (pen *Pen) GetPartialNoun(name, kind string) (ret MatchedNoun, err error) {
	if len(name) > 0 {
		words := name + blank
		if e := pen.db.QueryRow(`
		with nouns(noun, name, kind) as (
			select mn.noun, my.name, mk.kind
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
		select noun, name, kind from (
			select noun, name, kind, substr(?2 ,0, length(name)+2) as words
			from nouns
			where words = (name || ' ')
		)
		order by length(name) desc
		limit 1`,
			pen.domain, words, kind).
			Scan(&ret.Name, &ret.Match, &ret.Kind); e != sql.ErrNoRows {
			err = e // could be nil or error
		}
	}
	return
}

// an unmatched trait returns the empty string and no error
// tbd: technically there's some possibility that there might be three traits:
// "wood", "veneer", and "wood veneer" -- subset names
// with the first two applying to one kind, and the third applying to a different kind;
// all in scope.  this would always match the second -- even if its not applicable.
// ( i guess that's where commas can be used by the user to separate things )
func (pen *Pen) GetPartialTrait(str string) (ret MatchedTrait, err error) {
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
			pen.domain, ap, words).Scan(&ret.Name); e != sql.ErrNoRows {
			err = e // could be nil or error
		}
	}
	return
}
