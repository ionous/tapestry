package maps

// write to a map
type Builder interface {
	// add the passed key, value to the in progress map
	Add(key string, val any) Builder
	// return the completed map
	Map() any
}

// a function which returns a new builder
type BuilderFactory func(reserve int) Builder
