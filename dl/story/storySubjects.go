package story

type Nouns struct {
	Subjects, Objects []string
	Objectifying      bool // phrases discuss noun subjects by default
}

// add to the known recent nouns over the course of the passed function.
// "subjects" are the main focus of the sentence, often the ones mentioned first (lhs).
func (n *Nouns) CollectSubjects(fn func() error) error {
	n.Subjects = nil
	n.Objectifying = false
	return fn()
}

// add to the known recent nouns over the course of the passed function.
// "objects" are the support nouns in a sentence, often mentioned last (rhs).
func (n *Nouns) CollectObjects(fn func() error) error {
	n.Objects = nil
	n.Objectifying = true
	err := fn()
	n.Objectifying = false
	return err
}

func (n *Nouns) pList() (ret *[]string) {
	if n.Objectifying {
		ret = &n.Objects
	} else {
		ret = &n.Subjects
	}
	return
}

func (n *Nouns) Add(name string) {
	pn := n.pList()
	(*pn) = append((*pn), name)
}
