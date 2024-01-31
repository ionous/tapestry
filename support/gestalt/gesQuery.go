package gestalt

import (
	"log"

	"git.sr.ht/~ionous/tapestry/support/grok"
)

type Query struct {
	g           grok.Grokker
	Quiet       bool
	LastArticle grok.Article
	LastMacro   grok.Macro
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
	if res, e := q.g.FindArticle(ws); e != nil {
		cnt = -1
		q.log("CountArticle", e)
	} else if res.Count > 0 {
		cnt = -1
	} else {
		q.LastArticle = res
		cnt = res.Len()
	}
	return
}

func (q Query) FindArticle(ws grok.Span) (ret grok.Article, cnt int) {
	if res, e := q.g.FindArticle(ws); e != nil {
		cnt = -1
		q.log("FindArticle", e)
	} else if m := res.Match; m != nil {
		ret, cnt = res, m.NumWords()
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

func (q Query) FindMacro(ws grok.Span) (cnt int) {
	if res, e := q.g.FindMacro(ws); e != nil {
		cnt = -1
		q.log("FindMacro", e)
	} else {
		q.LastMacro = res
		cnt = res.Len()
	}
	return
}

func (q Query) log(fn string, e error) {
	if !q.Quiet {
		log.Println(fn, e)
	}
}
