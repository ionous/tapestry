package jess

const (
	// only allow simple names when matching.
	PlainNameMatching = iota
	// when matching names, some phrases imply the creation of new nouns
	// this flag prevents those from matching
	ExcludeNounCreation
	// when matching names, try to match nouns
	IncludeExistingNouns
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

func matchNouns(q Query) bool {
	flags := q.GetContext()
	return (flags & IncludeExistingNouns) != 0
}

func allowNounCreation(q Query) bool {
	flags := q.GetContext()
	return (flags&PlainNameMatching) == 0 &&
		(flags&ExcludeNounCreation) == 0 &&
		(flags&IncludeExistingNouns) == 0
}
