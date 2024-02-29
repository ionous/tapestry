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
	Matched string
	Exact   bool // when the phrase contains "called", we shouldn't fold the noun into other similarly named nouns.
	Traits  []string
	Kinds   []string // it's possible, if rare, to apply multiple kinds
	// ex. The container called the coffin is a closed openable thing.
}

func (n resultName) String() (ret string) {
	return n.Matched
}

func (n resultName) GetKind(defaultKind string) (ret string) {
	if len(n.Kinds) > 0 {
		ret = n.Kinds[0]
	} else {
		ret = defaultKind
	}
	return
}

type articleResult struct {
	Matched string
	Count   int // for counted nouns: "seven (apples)"
}

func (a articleResult) String() string {
	return a.Matched
}
