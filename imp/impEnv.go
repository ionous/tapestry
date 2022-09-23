package imp

type Environ struct {
	Recent struct {
		// Scene, Aspect, Test string
		// Nouns[]? Relation, Trait
		// string or string
		Nouns Nouns
	}
	ActivityDepth int
}

// tbd? could we look upwards to ask whether we are in a given block?
func (env *Environ) InProgram() bool {
	return env.ActivityDepth > 0
}
