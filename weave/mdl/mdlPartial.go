package mdl

import (
	"database/sql"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

type MacroMatch struct {
	Macro Macro
	Width int
}

func (pen *Pen) GetPartialKind(ws match.Span, out *kindsOf.Kinds) (ret string, err error) {
	// to ensure a whole word match, during query the names of the kinds are appended with blanks
	// and so we also give the phrase a final blank in case the phrase is a single word.
	if len(ws) > 0 {
		var k kindInfo
		str := ws.String()
		words := strings.ToLower(str) + blank
		switch e := pen.db.QueryRow(`
	with kinds(id, name, alt, path) as (
		select mk.rowid, mk.kind, mk.singular, ',' || mk.path
		from mdl_kind mk
 		join domain_tree
 		on (uses = domain)
 		where base = ?1
	)
	select id, name, path from (
		select id, name, substr(?2 ,0, length(name)+2) as words, path
		from kinds
		where words = (name || ' ')
		union all 
		select id, name, substr(?2, 0, length(alt)+2) as words, path
		from kinds
		where alt is not null and words = (alt || ' ')
	)
	order by length(name) desc
	limit 1`,
			pen.domain, words).Scan(&k.id, &k.name, &k.path); e {
		case nil:
			ret = k.name
			if out != nil {
				fullPath := k.fullpath()
				for _, k := range kindsOf.DefaultKinds {
					if strings.HasSuffix(fullPath, pen.getPath(k)) {
						*out = k
						break
					}
				}
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
// return the number of words in that match.
func (pen *Pen) GetPartialField(ws match.Span) (ret string, err error) {
	if len(ws) > 0 {
		words := strings.ToLower(ws.String()) + blank
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

// return the number of words in that match.
// tbd: technically there's some possibility that there might be three traits:
// "wood", "veneer", and "wood veneer" -- subset names
// with the first two applying to one kind, and the third applying to a different kind;
// all in scope.  this would always match the second -- even if its not applicable.
// ( i guess that's where commas can be used by the user to separate things )
func (pen *Pen) GetPartialTrait(ws match.Span) (ret string, err error) {
	if ap, e := pen.getAspectPath(); e != nil {
		err = e
	} else if len(ws) > 0 {
		words := strings.ToLower(ws.String()) + blank
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

func (pen *Pen) GetPartialMacro(ws match.Span) (ret MacroMatch, err error) {
	// uses spaces instead of underscores...
	if len(ws) == 0 {
		err = sql.ErrNoRows
	} else {
		words := strings.ToLower(ws.String()) + blank
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
				width := strings.Count(found.phrase, blank) + 1
				ret = MacroMatch{
					Width: width,
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
