package gestalt

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
func (q Query) SkipSeparators(ws grok.Span) (cnt int) {
	if sep, e := grok.CommaAnd(ws); e != nil {
		cnt = -1
		q.log("SkipSeparators", e)
	} else {
		cnt = sep.Len()
	}
	return
}

// ignores counted nouns
func (q Query) SkipArticle(ws grok.Span) (cnt int) {
	_, cnt = q.FindArticle(ws)
	return
}

func (q Query) FindArticle(ws grok.Span) (ret grok.Article, cnt int) {
	if m, e := grok.FindCommonArticles(ws); e != nil {
		cnt = -1
		q.log("FindArticle", e)
	} else if m != nil {
		ret, cnt = grok.Article{Match: m}, m.NumWords()
	}
	return
}

func (q Query) FindKind(ws grok.Span) (ret grok.Match, cnt int) {
	if res, e := q.g.FindKind(ws); e != nil {
		cnt = -1
		q.log("FindKind", e)
	} else if res != nil {
		ret, cnt = res, res.NumWords()
	}
	return
}

func (q Query) FindTrait(ws grok.Span) (ret grok.Match, cnt int) {
	if res, e := q.g.FindTrait(ws); e != nil {
		cnt = -1
		q.log("FindTrait", e)
	} else if res != nil {
		ret, cnt = res, res.NumWords()
	}
	return
}

func (q Query) FindMacro(ws grok.Span) (ret grok.Macro, cnt int) {
	if res, e := q.g.FindMacro(ws); e != nil {
		cnt = -1
		q.log("FindMacro", e)
	} else {
		ret = res
		cnt = res.Len()
	}
	return
}

func (q Query) log(fn string, e error) {
	if !q.Quiet {
		log.Println(fn, e)
	}
}
