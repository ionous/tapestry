package parser

import (
	"bytes"
	"strconv"
)

// Results used by the parser include, a list of results, a resolved object, a resolved action, etc.
// On success, the parser generally returns a ResultList as its primary result.
type Result interface {
	// the number of words used to match this result.
	WordsMatched() int
}

type ResolvedAction struct {
	Name string
}

// ResolvedActor
// ResolvedNumber
// ResolvedWords

type ResolvedMulti struct {
	Nouns     []NounInstance
	WordCount int
}

type ResolvedNoun struct {
	NounInstance NounInstance
	Words        []string // what the user said to identify the object
}

type ResolvedWord struct {
	Word string
}

func (f ResolvedAction) String() string {
	return "Action: " + f.Name
}
func (f ResolvedAction) WordsMatched() int {
	return 0
}

//
func (f ResolvedMulti) String() string {
	var b bytes.Buffer
	b.WriteString("Nouns(")
	b.WriteString(strconv.Itoa(len(f.Nouns)))
	b.WriteString("): ")
	for i, res := range f.Nouns {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(res.Id().String())
	}
	return b.String()
}
func (f ResolvedMulti) WordsMatched() int {
	return f.WordCount
}

//
func (f ResolvedNoun) String() string {
	return "Noun: " + f.NounInstance.Id().String()
}
func (f ResolvedNoun) WordsMatched() int {
	return len(f.Words)
}

//
func (f ResolvedWord) String() string {
	return "Word: " + f.Word
}
func (f ResolvedWord) WordsMatched() int {
	return 1
}
