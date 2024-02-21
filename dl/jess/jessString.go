package jess

// a string that satisfies the Matched interface
type matchedString struct {
	str       string
	wordCount int
}

func (m matchedString) String() string {
	return m.str
}

func (m matchedString) NumWords() int {
	return m.wordCount
}
