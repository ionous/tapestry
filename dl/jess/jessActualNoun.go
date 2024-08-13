package jess

// the noun that matched ( as opposed to the name that matched )
type ActualNoun struct {
	Name string
	Kind string
}

func (an ActualNoun) IsValid() bool {
	return len(an.Name) > 0
}
