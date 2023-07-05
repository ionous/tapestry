package grokdb

import (
	"database/sql"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/support/groktest"
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

// implements grok.Match; returned by dbSource.
type dbMatch struct {
	Id int
	grok.Span
}

// implements grok.Grokker; returned by dbSource.
type dbSource struct {
	db                      *tables.Cache
	domain                  string
	aspectPath, patternPath string
}

func (d *dbSource) FindDeterminer(ws []grok.Word) grok.Match {
	// FIX: should come from the db
	return det.FindMatch(ws)
}

var det = groktest.PanicSpans("the", "a", "an", "some", "our")

// if the passed words starts with a kind,
// return the number of words in  that match.
func (d *dbSource) FindKind(ws []grok.Word) (ret grok.Match) {
	// to ensure a whole word match, during query the names of the kinds are appended with blanks
	// and so we also give the phrase a final blank in case the phrase is a single word.
	words := grok.WordsWithSep(ws, '_') + "_"
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
		where words = (name || '_')
		union all 
		select id, alt, substr(?2, 0, length(alt)+2) as words
		from kinds
		where alt is not null and words = (alt || '_')
	)
	order by length(name) desc
	limit 1`,
			d.domain, words, mp).Scan(&found.id, &found.name)
		switch e {
		case nil:
			width := strings.Count(found.name, "_") + 1
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
	return
}

func (d *dbSource) getPatternPath() (ret string, err error) {
	return getPath(d.db, kindsOf.Macro, &d.patternPath)
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
		select (','||rowid||',') 
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
func (d *dbSource) FindTrait(ws []grok.Word) (ret grok.Match) {
	if ap, e := d.getAspectPath(); e != nil {
		panic(e) // maybe should be returning error
	} else {
		words := grok.WordsWithSep(ws, '_') + "_"
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
		where words = (name || '_')
	)
	order by length(name) desc
	limit 1`,
			d.domain, ap, words).Scan(&found.name)
		switch e {
		case nil:
			width := strings.Count(found.name, "_") + 1
			ret = grok.Span(ws[:width])
		case sql.ErrNoRows:
			// return nothing.
		default:
			panic(e)
		}
	}
	return
}

// if the passed words starts with a macro,
// return information about that match
func (d *dbSource) FindMacro(ws []grok.Word) (ret grok.MacroInfo, okay bool) {
	if m, e := d.findMacro(ws); e != nil && e != sql.ErrNoRows {
		panic(e)
	} else if e == nil {
		ret, okay = m, true
	}
	return
}

func (d *dbSource) findMacro(ws []grok.Word) (ret grok.MacroInfo, err error) {
	// uses spaces instead of underscores...
	words := grok.WordsWithSep(ws, ' ') + " "
	var found struct {
		kid    int64  // id of the kind
		name   string // name of the kind/macro
		phrase string // string of the macro phrase
		rev    bool
	}
	if e := d.db.QueryRow(`
	select mk.rowid, mk.kind, phrase, reversed 
	from mdl_phrase mg
	join mdl_kind mk 
		on (mk.rowid = mg.macro)
	join domain_tree dt
		on (dt.uses = mg.domain)
	where base = ?1
	and (phrase || ' ') = substr(?2 ,0, length(phrase)+2)
	order by length(phrase) desc
	limit 1`, d.domain, words).Scan(
		&found.kid, &found.name, &found.phrase, &found.rev); e != nil {
		err = e
	} else if parts, e := tables.QueryStrings(d.db,
		`select affinity 
		from mdl_field 
		where kind=?1
		order by rowid`, found.kid); e != nil {
		err = e
	} else if numParts := len(parts); numParts != 0 && numParts != 3 {
		err = errutil.Fmt("invalid macro; should have two fields and one return; has %d", numParts)
	} else {
		var flag grok.MacroType
		if numParts > 0 {
			a, b := affine.Affinity(parts[0]), affine.Affinity(parts[1])
			if found.rev {
				a, b = b, a
			}
			if a == affine.Text {
				if b == affine.Text {
					err = errutil.New("one one not supported?")
				} else if b == affine.TextList {
					flag = grok.OneToMany
				} else {
					err = errutil.New("unexpected aff", b)
				}
			} else if a == affine.TextList {
				if b == affine.Text {
					flag = grok.ManyToOne
				} else if b == affine.TextList {
					flag = grok.ManyToMany
				} else {
					err = errutil.New("unexpected aff", b)
				}
			} else {
				err = errutil.New("unexpected aff", a)
			}
		}
		if err == nil {
			width := strings.Count(found.phrase, " ") + 1
			ret = grok.MacroInfo{
				Name:  found.name,
				Match: grok.Span(ws[:width]),
				Type:  flag,
			}
		}
	}
	return
}
