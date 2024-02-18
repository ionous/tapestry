package jessdb

import (
	"database/sql"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

// implements jess.Match; returned by dbSource.
type dbMatch struct {
	Id int
	match.Span
}

// implements jess.Query; returned by dbSource.
type dbSource struct {
	db                      *tables.Cache
	domain                  string
	aspectPath, patternPath string
}

func (d *dbSource) FindKind(ws match.Span) (ret jess.Matched, retWidth int) {
	if res, e := d.findKind(ws); e != nil {
		log.Println("FindKind", e)
	} else if res != nil {
		ret, retWidth = res, res.NumWords()
	}
	return
}

func (d *dbSource) FindTrait(ws match.Span) (ret jess.Matched, retWidth int) {
	if res, e := d.findTrait(ws); e != nil {
		log.Println("FindTrait", e)
	} else if res != nil {
		ret, retWidth = res, res.NumWords()
	}
	return
}

func (d *dbSource) FindMacro(ws match.Span) (ret jess.Macro, retWidth int) {
	if res, e := d.findMacro(ws); e == nil {
		ret, retWidth = res, res.NumWords()
	} else if e != sql.ErrNoRows {
		log.Println("FindMacro", e)
	}
	return
}

// if the passed words starts with a kind,
// return the number of words in  that match.
func (d *dbSource) findKind(ws match.Span) (ret jess.Matched, err error) {
	// to ensure a whole word match, during query the names of the kinds are appended with blanks
	// and so we also give the phrase a final blank in case the phrase is a single word.
	if len(ws) > 0 {
		words := strings.ToLower(ws.String()) + blank
		var found struct {
			id   int
			name string
		}
		// excludes patterns (and macros) from matching these physical kinds
		if mp, e := d.getPatternPath(); e != nil {
			log.Println(e)
		} else {
			e := d.db.QueryRow(`
	with kinds(id, name, alt) as (
		select mk.rowid, mk.kind, mk.singular
		from mdl_kind mk
 		join domain_tree
 		on (uses = domain)
 		where base = ?1
 		and not instr(',' || mk.path, ?3 )
	)
	select id, name  from (
		select id, name, substr(?2 ,0, length(name)+2) as words
		from kinds
		where words = (name || ' ')
		union all 
		select id, alt, substr(?2, 0, length(alt)+2) as words
		from kinds
		where alt is not null and words = (alt || ' ')
	)
	order by length(name) desc
	limit 1`,
				d.domain, words, mp).Scan(&found.id, &found.name)
			switch e {
			case nil:
				width := strings.Count(found.name, blank) + 1
				ret = dbMatch{
					Id:   found.id,
					Span: ws[:width],
				}
			case sql.ErrNoRows:
				// return nothing.
			default:
				log.Println(e)
			}
		}
	}
	return
}

const blank = " "
const space = ' '

func (d *dbSource) getPatternPath() (ret string, err error) {
	return getPath(d.db, kindsOf.Pattern, &d.patternPath)
}

func (d *dbSource) getAspectPath() (ret string, err error) {
	return getPath(d.db, kindsOf.Aspect, &d.aspectPath)
}

func getPath(db *tables.Cache, kind kindsOf.Kinds, out *string) (ret string, err error) {
	if len(*out) > 0 {
		ret = *out
	} else {
		var path string
		e := db.QueryRow(`
		select (',' || rowid || ',') 
		from mdl_kind where kind = ?1
		limit 1
		`, kind.String()).Scan(&path)
		switch e {
		case nil:
			ret, *out = path, path
		case sql.ErrNoRows:
			err = errutil.New("couldn't determine", kind)
		default:
			err = e
		}
	}
	return
}

// if the passed words starts with a trait,
// return the number of words in that match.
// tbd: technically there's some possibility that there might be three traits:
// "wood", "veneer", and "wood veneer" -- subset names
// with the first two applying to one kind, and the third applying to a different kind;
// all in scope.  this would always match the second -- even if its not applicable.
// ( i guess that's where commas can be used by the user to separate things )
func (d *dbSource) findTrait(ws match.Span) (ret jess.Matched, err error) {
	if ap, e := d.getAspectPath(); e != nil {
		err = e
	} else if len(ws) > 0 {
		words := strings.ToLower(ws.String()) + blank
		var found struct {
			name string
		}
		e := d.db.QueryRow(`
		with traits(name) as (
		select mf.field
		from mdl_kind mk
		join domain_tree dt
			on (dt.uses = mk.domain)
		join mdl_field mf 
			on(mf.kind = mk.rowid)	
		where dt.base = ?1
		and instr(',' || mk.path, ?2 )
	)
	select name from (
		select name, substr(?3 ,0, length(name)+2) as words
		from traits
		where words = (name || ' ')
	)
	order by length(name) desc
	limit 1`,
			d.domain, ap, words).Scan(&found.name)
		switch e {
		case nil:
			width := strings.Count(found.name, blank) + 1
			ret = match.Span(ws[:width])
		case sql.ErrNoRows:
			// return nothing.
		default:
			err = e
		}
	}
	return
}

func (d *dbSource) findMacro(ws match.Span) (ret jess.Macro, err error) {
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
		if e := d.db.QueryRow(`
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
	limit 1`, d.domain, words).Scan(
			&found.kid, &found.name, &found.phrase, &found.reversed, &found.result); e != nil {
			err = e
		} else if parts, e := tables.QueryStrings(d.db,
			`select affinity 
		from mdl_field 
		where kind=?1
		order by rowid`, found.kid); e != nil {
			err = e
		} else if numFields := len(parts) - int(found.result.Int32); numFields <= 0 {
			err = errutil.Fmt("most macros should have two fields and one result; has %d fields and %d returns",
				numFields, found.result.Int32)
		} else {
			var flag jess.MacroType
			if numFields == 1 {
				flag = jess.Macro_PrimaryOnly
			} else {
				a, b := affine.Affinity(parts[0]), affine.Affinity(parts[1])
				if a == affine.Text {
					if b == affine.Text {
						err = errutil.New("one one not supported?")
					} else if b == affine.TextList {
						flag = jess.Macro_ManySecondary
					} else {
						err = errutil.New("unexpected aff", b)
					}
				} else if a == affine.TextList {
					if b == affine.Text {
						flag = jess.Macro_ManyPrimary
					} else if b == affine.TextList {
						flag = jess.Macro_ManyMany
					} else {
						err = errutil.New("unexpected aff", b)
					}
				} else {
					err = errutil.New("unexpected aff", a)
				}
			}
			if err == nil {
				width := strings.Count(found.phrase, blank) + 1
				ret = jess.Macro{
					Name:     found.name,
					Matched:  match.Span(ws[:width]),
					Type:     flag,
					Reversed: found.reversed,
				}
			}
		}
	}
	return
}
