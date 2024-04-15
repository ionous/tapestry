package jess

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func startsUpper(rs []match.TokenValue) (okay bool) {
	if len(rs) > 0 && rs[0].Token == match.String {
		str := rs[0].Value.(string)
		first, _ := utf8.DecodeRuneInString(str)
		okay = unicode.IsUpper(first) // this works okay even if the string was empty
	}
	return
}

// even one name can generate several nouns ( ex. "two things" )
// after gets called for each one.
func genNounValues(u Scheduler, ns []DesiredNoun, after postGenOne) (err error) {
	return u.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if e := writeNounValues(w, ns); e != nil {
			err = e
		} else if after != nil {
			for _, n := range ns {
				if e := after(n.Noun); e != nil {
					err = e
					break
				}
			}
		}
		return
	})
}

type postGenOne func(a string) error
type postGenMany func(a, b []DesiredNoun) error

func writeKinds(w weaver.Weaves, noun string, kinds []string) (err error) {
	for _, k := range kinds {
		if e := w.AddNounKind(noun, k); e != nil && !errors.Is(e, weaver.Duplicate) {
			err = e
			break
		}
	}
	return
}

func writeNounValues(w weaver.Weaves, ns []DesiredNoun) (err error) {
	for _, n := range ns {
		if e := n.writeNounValues(w); e != nil {
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
func ensureNoun(q Query, w weaver.Weaves, ts []match.TokenValue, props *NounProperties) (ret string, created bool, err error) {
	if noun, width := q.FindNoun(ts, nil); width > 0 {
		ret = noun
	} else if name, count := match.Stringify(ts); count != len(ts) {
		out := match.DebugStringify(ts)
		err = fmt.Errorf("not all of name consumed? %q", out)
	} else {
		noun := inflect.Normalize(name)
		defaultKind := Objects
		if ks := props.Kinds; len(ks) > 0 {
			defaultKind = ks[0]
			props.Kinds = ks[1:]
		}
		if e := w.AddNounKind(noun, defaultKind); e != nil {
			err = e // if duplicate, FindNoun should have triggered; so return all errors
		} else if e := registerNames(w, noun, name); e != nil {
			err = e
		} else {
			ret = noun
			created = true
		}
	}
	return
}

// fix: make names could use tokens directly
func registerNames(w weaver.Weaves, noun, name string) (err error) {
	names := mdl.MakeNames(name)
	for i, n := range names {
		if e := w.AddNounName(noun, n, i); e != nil {
			err = e
			break
		}
	}
	return
}
