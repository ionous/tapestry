package groktest

import "git.sr.ht/~ionous/tapestry/support/grok"

type MacroList struct {
	SpanList
	types []grok.MacroType
}

func (n MacroList) FindMacro(ws []grok.Word) (ret grok.MacroInfo, okay bool) {
	if i, skip := n.FindPrefix(ws); skip > 0 {
		var match grok.Span = n.SpanList[i]
		ret = grok.MacroInfo{
			Match: match,
			Type:  n.types[i],
		}
		okay = true
	}

	return
}

func PanicMacros(pairs ...any) (out MacroList) {
	cnt := len(pairs) / 2
	out.SpanList = make(SpanList, cnt)
	out.types = make([]grok.MacroType, cnt)
	for i := 0; i < cnt; i++ {
		out.types[i] = pairs[i*2+0].(grok.MacroType)
		out.SpanList[i] = PanicSpan(pairs[i*2+1].(string))
	}
	return
}
