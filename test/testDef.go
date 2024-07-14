package test

import "strconv"

// a simple noun
type Things struct {
}

// local variables for grouping tests
type Locals struct {
	Objects   []string
	Settings  []GroupSettings
	Collation GroupCollation
}

// the pattern assign grouping builds a list of group settings from a list of objects
type GroupSettings struct {
	// the name of the object this record describes
	// ex. tile X, or tile W ( from a scrabble set )
	Name string
	// objects with the same label are considered to be in the same group.
	// ex. "scrabble tiles"
	Label string
	// whether and how to print objects.
	GroupOptions
}

// GroupOptions defines how objects in groups should display.
type GroupOptions int

//go:generate stringer -type=GroupOptions
const (
	// indicates we dont want the individual objects in the group
	// ex. the scrabble tiles, the usual utensils, several things.
	WithoutObjects GroupOptions = iota
	// indicates individual objects should have articles.
	// ex. the X and the W tiles from a Scrabble set.
	ObjectsWithArticles
	// indicates the individual objects shouldnt use articles
	// ex. tiles X and W from a Scrabble set.
	ObjectsWithoutArticles
)

// the pattern collate groups builds a group collation from a list of group settings
type GroupCollation struct {
	Groups []GroupedObjects
}

type GroupedObjects struct {
	Settings GroupSettings // the settings of the first object in the group
	Objects  []string      // the list of objects with the same settings
}

// a noun with an aspect
type Messages struct {
	Neatness Neatness
}

// make a fake aspect for testing
type Neatness int

const (
	Neat Neatness = iota
	Scuffed
	Trampled
	NumNeatness = iota
)

func (i Neatness) String() string {
	switch i {
	case Neat:
		return "neat"
	case Scuffed:
		return "scuffed"
	case Trampled:
		return "trampled"
	default:
		return "Neatness(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
