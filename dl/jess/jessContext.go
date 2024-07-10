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

// adds flags to the query ( via or )
func AddContext(q Query, flags int) (ret Query) {
	if ctx, ok := q.(queryContext); ok {
		// since the query context isnt a reference
		// this creates new flags for the scope
		ctx.flags |= flags
		ret = ctx
	} else {
		ret = queryContext{Query: q, flags: flags}
	}
	return
}

// remove flags from the query
func ClearContext(q Query, flags int) (ret Query) {
	if ctx, ok := q.(queryContext); !ok {
		ret = q // unchanged, b/c only a context can have flags.
	} else {
		ctx.flags &= ^flags
		ret = ctx
	}
	return
}

type queryContext struct {
	Query
	flags int
}

func (q queryContext) GetContext() int {
	return q.flags
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
