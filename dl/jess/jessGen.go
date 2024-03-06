package jess

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func startsUpper(str string) bool {
	first, _ := utf8.DecodeRuneInString(str)
	return unicode.IsUpper(first) // this works okay even if the string was empty
}

// even one name can generate several nouns ( ex. "two things" )
// after gets called for each one.
func genNoun(rar Registrar, ns []DesiredNoun, after postGenOne) error {
	return genNouns(rar, ns, nil, func(ns, _ []DesiredNoun) (err error) {
		for _, n := range ns {
			if e := after(n.Noun); e != nil {
				err = e
				break
			}
		}
		return
	})
}

type postGenOne func(a string) error
type postGenMany func(a, b []DesiredNoun) error

func genNouns(rar Registrar, a, b []DesiredNoun, after postGenMany) (err error) {
	return rar.PostProcess(GenerateValues, func(Query) (err error) {
		if e := generateValues(rar, a); e != nil {
			err = e
		} else if e := generateValues(rar, b); e != nil {
			err = e
		} else {
			err = after(a, b)
		}
		return
	})
}

func registerKinds(rar Registrar, noun string, kinds []string) (err error) {
	for _, k := range kinds {
		if e := rar.AddNounKind(noun, k); e != nil && !errors.Is(e, mdl.Duplicate) {
			err = e
			break
		}
	}
	return
}

func generateValues(rar Registrar, ns []DesiredNoun) (err error) {
	for _, n := range ns {
		if e := n.generateValues(rar); e != nil {
			err = e
			break
		}
	}
	return
}

// creates a noun as a placeholder
// later, a pass ensures that all placeholder nouns have been given kinds;
// or it upgrades them to things.
// to simplify the code, this happens even if the kind might possibly be known.
func ensureNoun(q Query, rar Registrar, name match.Span) (ret string, created bool, err error) {
	if noun, w := q.FindNoun(name, ""); w > 0 {
		ret = noun
	} else {
		name := name.String()
		noun := inflect.Normalize(name)
		if e := rar.AddNounKind(noun, ""); e != nil {
			err = e // if duplicate, FindNoun should have triggered; so return all errors
		} else if e := registerNames(rar, noun, name); e != nil {
			err = e
		} else {
			ret = noun
			created = true
		}
	}
	return
}

func registerNames(rar Registrar, noun, name string) (err error) {
	names := mdl.MakeNames(name)
	for i, n := range names {
		if e := rar.AddNounName(noun, n, i); e != nil {
			err = e
			break
		}
	}
	return
}
