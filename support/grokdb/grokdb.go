package grokdb

import (
	"database/sql"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/support/groktest"
)

// implements grok.Match; returned by dbg.
type dbMatch struct {
	Id int
	grok.Span
}

// implements grok.Grokker; returned by dbg.
type dbg struct {
	db     *sql.DB
	domain string
}

func (d *dbg) FindDeterminer(ws []grok.Word) grok.Match {
	// FIX: should come from the db
	return det.FindMatch(ws)
}

var det = groktest.PanicSpans("the", "a", "an", "some", "our")

// if the passed words starts with a kind,
// return the number of words in  that match.
func (d *dbg) FindKind(ws []grok.Word) (ret grok.Match) {
	// to ensure a whole word match, during query the names of the kinds are appended with blanks
	// and so we also give the phrase a final blank in case the phrase is a single word.z1
	words := grok.WordsWithSep(ws, '_') + "_"
	var found struct {
		id   int
		name string
	}
	e := d.db.QueryRow(`
	with kinds(id, name, alt) as (
		select mk.rowid, mk.kind, mk.singular
		from mdl_kind mk
 		join domain_tree
 		on (uses = domain)
 		where base = ?1
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
		d.domain, words).Scan(&found.id, &found.name)
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
	return
}

// if the passed words starts with a trait,
// return the number of words in  that match.
func (d *dbg) FindTrait(ws []grok.Word) (ret grok.Match) {
	// FIX: should come from the db
	return traits.FindMatch(ws)
}

var traits = groktest.PanicSpans("closed",
	"open",
	"openable",
	"transparent",
	"fixed in place")

// if the passed words starts with a macro,
// return information about that match
func (d *dbg) FindMacro(ws []grok.Word) (ret grok.MacroInfo, okay bool) {
	return macros.FindMacro(ws)
}

var macros = groktest.PanicMacros(
	// tbd: flags need more thought.
	grok.ManyToOne, "kind of", // for "a closed kind of container"
	grok.ManyToOne, "kinds of", // for "are closed containers"
	grok.ManyToOne, "a kind of", // for "a kind of container"
	// other macros
	grok.OneToMany, "on", // on the x are the w,y,z
	grok.OneToMany, "in",
	//
	grok.ManyToMany, "suspicious of",
)
