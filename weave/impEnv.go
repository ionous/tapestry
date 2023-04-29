package weave

type Environ struct {
	ActivityDepth int
}

// tbd? could we look upwards to ask whether we are in a given block?
func (env *Environ) InProgram() bool {
	return env.ActivityDepth > 0
}
