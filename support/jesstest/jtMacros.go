package jesstest

import (
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/support/match"
)

type MacroList struct {
	match.SpanList
	info []macroInfo
}

type macroInfo struct {
	name      string
	macroType jess.MacroType
	reversed  bool
}

func (n MacroList) FindMacro(ws []match.Word) (ret jess.Macro, width int) {
	if i, skip := n.FindPrefix(ws); skip > 0 {
		var match match.Span = n.SpanList[i]
		info := n.info[i]
		width, ret = len(match), jess.Macro{
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
	out.SpanList = make(match.SpanList, cnt)
	out.info = make([]macroInfo, cnt)
	for i := 0; i < cnt; i++ {
		x := i * width
		out.SpanList[i] = match.PanicSpan(phraseMacroTypeRev[x+0].(string))
		out.info[i] = macroInfo{
			name:      phraseMacroTypeRev[x+1].(string),
			macroType: phraseMacroTypeRev[x+2].(jess.MacroType),
			reversed:  phraseMacroTypeRev[x+3].(bool),
		}
	}
	return
}
