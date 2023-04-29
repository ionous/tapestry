package weave

type Nouns struct {
	Subjects []string
}

// add to the known recent nouns over the course of the passed function.
// "subjects" are the main focus of the sentence, often the ones mentioned first (lhs).
func (n *Nouns) CollectSubjects(fn func() error) error {
	n.Subjects = nil // reset
	return fn()
}

func (n *Nouns) Add(name string) {
	n.Subjects = append(n.Subjects, name)
}
