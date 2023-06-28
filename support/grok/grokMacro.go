package grok

type MacroType int

const (
	ManyToOne MacroType = iota
	OneToMany
	ManyToMany
)

type macroList struct {
	spanList
	types []MacroType
}

func (ml macroList) get(i int) ([]Word, MacroType) {
	return ml.spanList[i], ml.types[i]
}

func panicMacros(pairs ...any) (out macroList) {
	cnt := len(pairs) / 2
	out.spanList = make(spanList, cnt)
	out.types = make([]MacroType, cnt)
	for i := 0; i < cnt; i++ {
		out.types[i] = pairs[i*2+0].(MacroType)
		out.spanList[i] = panicHash(pairs[i*2+1].(string))
	}
	return
}
