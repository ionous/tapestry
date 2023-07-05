package groktest

import "git.sr.ht/~ionous/tapestry/support/grok"

type MacroList struct {
	SpanList
	types []grok.MacroType
	names []string
}

func (n MacroList) FindMacro(ws []grok.Word) (ret grok.MacroInfo, okay bool) {
	if i, skip := n.FindPrefix(ws); skip > 0 {
		var match grok.Span = n.SpanList[i]
		ret = grok.MacroInfo{
			Name:  n.names[i],
			Match: match,
			Type:  n.types[i],
		}
		okay = true
	}
	return
}

func PanicMacros(typePhraseMacro ...any) (out MacroList) {
	const width = 3
	cnt := len(typePhraseMacro) / width
	out.SpanList = make(SpanList, cnt)
	out.types = make([]grok.MacroType, cnt)
	out.names = make([]string, cnt)
	for i := 0; i < cnt; i++ {
		x := i * width
		out.types[i] = typePhraseMacro[x+0].(grok.MacroType)
		out.SpanList[i] = PanicSpan(typePhraseMacro[x+1].(string))
		out.names[i] = typePhraseMacro[x+2].(string)
	}
	return
}
