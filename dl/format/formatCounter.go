package format

// A slot used internally for generating unique names during weave.
type Counter interface {
	GetInternalCounter() *string
}
