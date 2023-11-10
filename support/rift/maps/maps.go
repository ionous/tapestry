package maps

// write to a map
type Builder interface {
	// add the passed pair to the in-progress map
	// returns a new builder ( not guaranteed to be the original one )
	// future: add uniqueness check and error
	Add(key string, val any) Builder
	// return the completed map
	Map() any
}

// a function which returns a new builder
// reserve indicates whether to keep space for an blank key
// (ie. comments)
type BuilderFactory func(reserve bool) Builder
