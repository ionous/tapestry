package jess

const (
	// exclude kinds when matching names
	PlainNameMatching = iota
	// kinds|aspects|records
	// by default, only matches kindsOf.Kinds ( concrete and abstract nouns )
	PropertyKinds
)

// set the query flags to the passed flags
// tbd: if this should or the existing flags; currently it doesnt
func AddContext(q Query, flags int) (ret Query) {
	if ctx, ok := q.(queryContext); ok {
		// since the query context is an object not a pointer
		// this creates new flags for the scope
		ctx.flags = flags
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
