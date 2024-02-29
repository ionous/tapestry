package jessdb

import (
	"database/sql"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type Source struct {
	inner *mdl.Modeler
}

func NewSource(db *mdl.Modeler) Source {
	return Source{db}
}

func countWords(str string) (ret int) {
	if len(str) > 0 {
		ret = 1 + strings.Count(str, " ")
	}
	return
}

func (x *Source) MatchSpan(domain string, span match.Span) (jess.Generator, error) {
	w := dbWrapper{x.inner.Pin(domain, "jess")}
	return jess.Match(w, span)
}

// implements jess.Query; returned by dbWrapper.
type dbWrapper struct {
	*mdl.Pen
}

func (d dbWrapper) GetContext() int {
	return 0
}

func (d dbWrapper) FindKind(ws match.Span, out *kindsOf.Kinds) (ret string, width int) {
	str := strings.ToLower(ws.String()) // fix: so not a fan of these string lowers
	if m, e := d.GetPartialKind(str); e != nil {
		log.Println("FindKind", e)
	} else {
		if out != nil {
			*out = m.Base
		}
		ret, width = m.Name, countWords(m.Match)
	}
	return
}

func (d dbWrapper) FindTrait(ws match.Span) (ret string, width int) {
	str := strings.ToLower(ws.String())
	if trait, e := d.GetPartialTrait(str); e != nil {
		log.Println("FindTrait", e)
	} else {
		// the returned name is the name of the trait from the db
		// it was used to match the front of the passed string
		// so the words in the trait are the words in the string.
		ret, width = trait, countWords(trait)
	}
	return
}

func (d dbWrapper) FindField(ws match.Span) (ret string, width int) {
	str := strings.ToLower(ws.String())
	if field, e := d.GetPartialField(str); e != nil {
		log.Println("FindField", e)
	} else {
		// re: countWords, same logic as find trait.
		ret, width = field, countWords(field)
	}
	return
}

func (d dbWrapper) FindMacro(ws match.Span) (ret jess.Macro, width int) {
	str := strings.ToLower(ws.String())
	if m, e := d.GetPartialMacro(str); e == nil {
		ret, width = m.Macro, countWords(m.Phrase)
	} else if e != sql.ErrNoRows {
		log.Println("FindMacro", e)
	}
	return
}

func (d dbWrapper) FindNoun(ws match.Span) (ret string, width int) {
	str := strings.ToLower(ws.String())
	if m, e := d.GetPartialNoun(str); e != nil {
		log.Println("FindNoun", e)
	} else {
		ret, width = m.Name, countWords(m.Match)
	}
	return
}
