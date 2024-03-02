package jess

const (
	// only allow simple names when matching.
	PlainNameMatching = iota
	ExcludeNounMatching
	//
	CheckIndefiniteArticles
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

func matchNouns(q Query) bool {
	flags := q.GetContext()
	return (flags&PlainNameMatching) == 0 &&
		(flags&ExcludeNounMatching) == 0
}

func useIndefinite(q Query) bool {
	flags := q.GetContext()
	return (flags & CheckIndefiniteArticles) != 0
}
