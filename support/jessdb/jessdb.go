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
	if res, e := d.GetPartialKind(ws, out); e != nil {
		log.Println("FindKind", e)
	} else {
		// fix: this is probably wrong --
		// it shouldnt be the width of the final kind;
		// it should be the number of words used to find that kind
		ret, width = res, countWords(res)
	}
	return
}

func (d dbWrapper) FindTrait(ws match.Span) (ret string, width int) {
	if res, e := d.GetPartialTrait(ws); e != nil {
		log.Println("FindTrait", e)
	} else {
		ret, width = res, countWords(res)
	}
	return
}

func (d dbWrapper) FindField(ws match.Span) (ret string, width int) {
	if res, e := d.GetPartialField(ws); e != nil {
		log.Println("FindField", e)
	} else {
		ret, width = res, countWords(res)
	}
	return
}

func (d dbWrapper) FindMacro(ws match.Span) (ret jess.Macro, width int) {
	if res, e := d.GetPartialMacro(ws); e == nil {
		ret, width = res.Macro, res.Width
	} else if e != sql.ErrNoRows {
		log.Println("FindMacro", e)
	}
	return
}

func (d dbWrapper) FindClosestNoun(ws match.Span) (ret string, width int) {
	if n, e := d.GetClosestNoun(ws.String()); e != nil {
		log.Println("FindClosestNoun", e)
	} else {
		ret, width = n, ws.NumWords()
	}
	return
}

func (d dbWrapper) FindExactNoun(ws match.Span) (ret string, width int) {
	if n, e := d.GetExactNoun(ws.String()); e != nil {
		log.Println("FindExactNoun", e)
	} else {
		ret, width = n, ws.NumWords()
	}
	return
}
