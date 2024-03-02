package jess

type localResults struct {
	Primary   []DesiredNoun
	Secondary []DesiredNoun // usually just one, except for nxm relations
	Macro     Macro
}

func makeResult(macro Macro, reverse bool, a, b []DesiredNoun) localResults {
	if reverse {
		a, b = b, a
	}
	return localResults{
		Primary:   a,
		Secondary: b,
		Macro:     macro,
	}
}

// fix: maybe move this into mdl as a "noun builder" in order to reduce duplicate queries?
// ( or switch mdl to a resource interface, where Noun is a pre-validated endpoint )
type DesiredNoun struct {
	Article     string // the predefined, lowercase article
	Flags       ArticleFlags
	Count       int    // for counted nouns: "seven (apples)"
	Noun        string // name of the noun ( if it existed )
	DesiredName string // text as specified by author ( if the noun didnt exist )
	Traits      []string
	Kinds       []string // it's possible, if rare, to apply multiple kinds
	// ex. The container called the coffin is a closed openable thing.
}
