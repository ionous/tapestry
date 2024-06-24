package query

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
)

type CheckData struct {
	Name   string
	Domain string
	Aff    affine.Affinity
	Prog   Bytes `json:",omitempty"` // a serialized rt.Execute_Slice
	Value  Bytes `json:",omitempty"` // a serialized literal.Value
}

type FieldData struct {
	Name     string
	Affinity affine.Affinity
	Class    string
	Init     Bytes `json:",omitempty"`
}

type NounInfo struct {
	Domain, Noun, Kind string // noun is unique identifier within the domain.
}

type RuleData struct {
	Name    string
	Stop    bool
	Jump    int
	Updates bool
	Prog    Bytes `json:",omitempty"` // a serialized rt.Execute_Slice
}

type ValueData struct {
	Field string
	Path  string
	Value Bytes `json:",omitempty"` // a serialized assignment or literal
}

type Query interface {
	IsDomainActive(name string) (bool, error)
	ActivateDomains(name string) (prev, next []string, err error)
	ReadChecks(actuallyJustThisOne string) ([]CheckData, error)
	// every field used by the passed kind
	FieldsOf(exactKind string) ([]FieldData, error)
	// given a plural or singular kind
	// return all ancestors starting with the kind itself.
	KindOfAncestors(kindOrKinds string) ([]string, error)
	// search using a short name
	NounInfo(shortname string) (NounInfo, error)
	// return the friendly name of the exact named noun
	NounName(fullname string) (string, error)
	// find the parser aliases for this noun
	NounNames(fullname string) ([]string, error)
	// a single field can contain a set of recursive spare values;
	// so this returns pairs of path, value.
	NounValues(fullname, field string) ([]ValueData, error)
	// all nouns of the indicated kind
	NounsByKind(exactKind string) ([]string, error)
	// the empty string if not found
	PluralToSingular(plural string) (string, error)
	// the empty string if not found
	PluralFromSingular(singular string) (string, error)
	// includes the parameters, followed by the result
	// the result can be a blank string for execute statements
	PatternLabels(pat string) ([]string, error)
	RulesFor(pat string) ([]RuleData, error)
	ReciprocalsOf(rel, id string) ([]string, error)
	RelativesOf(rel, id string) ([]string, error)
	// errors if nothing changed
	// doesnt check to see if the relation is valid;
	// the caller should do that.
	Relate(rel, noun, otherNoun string) error
	//
	LoadGame(path string) (CacheMap, error)
	SaveGame(path string, dynamicValues CacheMap) error
	Random(inclusiveMin, exclusiveMax int) int
}

func (n *NounInfo) IsValid() bool {
	return len(n.Noun) != 0
}

func (n *NounInfo) String() (ret string) {
	if !n.IsValid() {
		ret = "<unknown object>"
	} else {
		ret = strings.Join([]string{n.Domain, n.Noun}, "::")
	}
	return
}
