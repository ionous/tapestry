package groktest

import "git.sr.ht/~ionous/tapestry/support/grok"

type MacroList struct {
	grok.SpanList
	info []macroInfo
}

type macroInfo struct {
	name      string
	macroType grok.MacroType
	reversed  bool
}

func (n MacroList) FindMacro(ws []grok.Word) (ret grok.Macro, err error) {
	if i, skip := n.FindPrefix(ws); skip > 0 {
		var match grok.Span = n.SpanList[i]
		info := n.info[i]
		ret = grok.Macro{
			Matched:  match,
			Name:     info.name,
			Type:     info.macroType,
			Reversed: info.reversed,
		}
	}
	return
}

func PanicMacros(phraseMacroTypeRev ...any) (out MacroList) {
	const width = 4
	cnt := len(phraseMacroTypeRev) / width
	out.SpanList = make(grok.SpanList, cnt)
	out.info = make([]macroInfo, cnt)
	for i := 0; i < cnt; i++ {
		x := i * width
		out.SpanList[i] = grok.PanicSpan(phraseMacroTypeRev[x+0].(string))
		out.info[i] = macroInfo{
			name:      phraseMacroTypeRev[x+1].(string),
			macroType: phraseMacroTypeRev[x+2].(grok.MacroType),
			reversed:  phraseMacroTypeRev[x+3].(bool),
		}
	}
	return
}
