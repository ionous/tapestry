package jess

import (
	"log"

	"git.sr.ht/~ionous/tapestry/support/grok"
)

type Query struct {
	g     grok.Grokker
	Quiet bool
}

func MakeQuery(g grok.Grokker) Query {
	return Query{g: g}
}

// returns -1 on error
func (q Query) SkipSeparators(ws InputState) (retWidth int) {
	if sep, e := grok.CommaAnd(ws); e != nil {
		retWidth = -1
		q.error("SkipSeparators", e)
	} else {
		retWidth = sep.Len()
	}
	return
}

// ignores counted nouns
func (q Query) SkipArticle(ws InputState) (retWidth int) {
	_, retWidth = q.FindArticle(ws)
	return
}

func (q Query) FindArticle(ws InputState) (ret grok.Article, retWidth int) {
	if m, e := grok.FindCommonArticles(grok.Span(ws)); e != nil {
		retWidth = -1
		q.error("FindArticle", e)
	} else if m != nil {
		ret, retWidth = grok.Article{Matched: m}, m.NumWords()
	}
	return
}

func (q Query) FindKind(ws InputState) (ret grok.Matched, retWidth int) {
	if res, e := q.g.FindKind(grok.Span(ws)); e != nil {
		retWidth = -1
		q.error("FindKind", e)
	} else if res != nil {
		ret, retWidth = res, res.NumWords()
	}
	return
}

func (q Query) FindTrait(ws InputState) (ret grok.Matched, retWidth int) {
	if res, e := q.g.FindTrait(grok.Span(ws)); e != nil {
		retWidth = -1
		q.error("FindTrait", e)
	} else if res != nil {
		ret, retWidth = res, res.NumWords()
	}
	return
}

func (q Query) FindMacro(ws InputState) (ret grok.Macro, retWidth int) {
	if res, e := q.g.FindMacro(grok.Span(ws)); e != nil {
		retWidth = -1
		q.error("FindMacro", e)
	} else {
		ret, retWidth = res, res.Len()
	}
	return
}

func (q Query) error(fn string, e error) {
	if !q.Quiet {
		log.Println(fn, e)
	}
}

func (q Query) log(msg string) {
	if !q.Quiet {
		log.Println(msg)
	}
}
