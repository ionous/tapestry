package jess

const (
	// only allow simple names when matching.
	PlainNameMatching = (1 << iota)
	// this limits matching to kinds which can be instanced
	// so that names which match other kinds can still become nouns
	// for instance, there can be a pattern called "on" and a verb called "on".
	MatchKindsOfKinds
	MatchKindsOfAspects
	//
	MatchPronouns
	// "called" checks for when the indefinite article
	// is *not* an indefinite article (a/an), and records it.
	// printing references to the noun will the specified article.
	CheckIndefiniteArticles
	// log each match automatically; used for testing
	LogMatches
)

// adds flags to the query
func AddContext(q JessContext, flags int) JessContext {
	q.flags |= flags
	return q
}

// remove flags from the query
func ClearContext(q JessContext, flags int) (ret JessContext) {
	q.flags &= ^flags
	return q
}

func matchKinds(q Query) bool {
	flags := q.GetContext()
	return (flags & PlainNameMatching) == 0
}

func matchPronouns(q Query) bool {
	flags := q.GetContext()
	return (flags & MatchPronouns) != 0
}

func matchKindsOfKinds(q Query) bool {
	flags := q.GetContext()
	return (flags & MatchKindsOfKinds) != 0
}

func matchKindsOfAspects(q Query) bool {
	flags := q.GetContext()
	return (flags & MatchKindsOfAspects) != 0
}

func useIndefinite(q Query) bool {
	flags := q.GetContext()
	return (flags & CheckIndefiniteArticles) != 0
}

func useLogging(q Query) bool {
	flags := q.GetContext()
	return (flags & LogMatches) != 0
}
