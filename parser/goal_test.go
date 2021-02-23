package parser_test

// Goal - tests the results of a parsed statement.
type Goal interface {
	Goal() Goal // marker: returns self
}

// ActionGoal - expects the named action and the specified nouns.
type ActionGoal struct {
	Action string
	Nouns  []string
}

// ClarifyGoal - expects that the parser ended ambiguously.
// supplies a word to keep it going.
type ClarifyGoal struct {
	// do we print the text here or not?
	// it might be nice for testing sake --
	// What do you want to examine
	// What do you want to look at?
	// and note, yu eed the matched "verb"?
	NounInstance string
}

// ErrorGoal - expects a specific error string.
type ErrorGoal struct {
	Error string
}

func (a *ActionGoal) Goal() Goal  { return a }
func (a *ClarifyGoal) Goal() Goal { return a }
func (a *ErrorGoal) Goal() Goal   { return a }
