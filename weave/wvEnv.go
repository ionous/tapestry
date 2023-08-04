package weave

// custom storage for importers ( ie. package story )
// alt: could use context and pass through all functions to be more go like.
type Env map[string]any

// can inc with zero to retrieve the current value
func (m Env) Inc(name string, inc int) int {
	var prev int
	if a, ok := m[name]; ok {
		prev = a.(int)
	}
	next := prev + inc
	m[name] = next
	return next
}
