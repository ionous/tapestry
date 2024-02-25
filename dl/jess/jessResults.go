package jess

type localResults struct {
	Primary   []resultName
	Secondary []resultName // usually just one, except for nxm relations
	Macro     Macro
}

func makeResult(macro Macro, reverse bool, a, b []resultName) localResults {
	if reverse {
		a, b = b, a
	}
	return localResults{
		Primary:   a,
		Secondary: b,
		Macro:     macro,
	}
}

type resultName struct {
	Article articleResult
	Match   Matched
	Exact   bool // when the phrase contains "called", we shouldn't fold the noun into other similarly named nouns.
	Traits  []string
	Kinds   []string // it's possible, if rare, to apply multiple kinds
	// ex. The container called the coffin is a closed openable thing.
}

func (n resultName) String() (ret string) {
	if m := n.Match; m != nil {
		ret = m.String()
	}
	return
}

type articleResult struct {
	Matched
	Count int // for counted nouns: "seven (apples)"
}

func (a articleResult) NumWords() (ret int) {
	if a.Matched != nil {
		ret = a.Matched.NumWords()
	}
	return
}

func (a articleResult) String() (ret string) {
	if a.Matched != nil {
		ret = a.Matched.String()
	}
	return
}
