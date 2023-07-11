package grok

import "github.com/ionous/errutil"

// `In](ws: the coffin> are <some coins, a notebook, and a gripping hand.)`,
func macroPhrase(known Grokker, macro Macro, ws []Word) (ret Results, err error) {
	out := Results{Macro: macro}
	sources, targets := &out.Sources, &out.Targets
	at, cnt := 0, len(ws)
	for ; at < cnt; at++ {
		if w := ws[at]; w.equals(keywords.is) || w.equals(keywords.are) {
			lhs, rhs := ws[:at], ws[at+1:]
			if e := grokNouns(known, sources, rhs, AllowMany|AllowAnonymous); e != nil {
				err = errutil.New("parsing right side nouns", e)
			} else if e := grokNouns(known, targets, lhs, OnlyOne|OnlyNamed); e != nil {
				err = errutil.New("parsing left side nouns", e)
			}
			break // either way, done.
		}
	}
	if err == nil {
		if nothingLeft := at == len(ws); nothingLeft {
			err = makeWordError(ws[0], "no is statement found")
		} else {
			ret = out
		}
	}
	return
}
