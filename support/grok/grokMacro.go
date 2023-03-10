package grok

type macroType int

const (
	ManyToOne macroType = iota
	OneToMany
	ManyToMany
)

type macroList struct {
	spanList
	types []macroType
}

func (ml macroList) get(i int) ([]Word, macroType) {
	return ml.spanList[i], ml.types[i]
}

func panicMacros(pairs ...any) (out macroList) {
	cnt := len(pairs) / 2
	out.spanList = make(spanList, cnt)
	out.types = make([]macroType, cnt)
	for i := 0; i < cnt; i++ {
		out.types[i] = pairs[i*2+0].(macroType)
		out.spanList[i] = panicHash(pairs[i*2+1].(string))
	}
	return
}
