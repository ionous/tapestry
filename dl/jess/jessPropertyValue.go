package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type PvFlags int

const (
	AllowPlural PvFlags = (1 << iota)
	AllowSingular
)

// match and apply a value
func TryPropertyValue(q JessContext, in InputState, flags PvFlags,
	accept func(PropertyValue, InputState),
	reject func(error),
) {
	q.Try(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) {
		if a, ok := matchPropertyValue(q, &in, flags); !ok {
			reject(FailedMatch{"didn't understand the value in this context", in})
		} else {
			accept(a, in)
		}
	}, reject)
}

func matchPropertyValue(q JessContext, in *InputState, flags PvFlags) (ret PropertyValue, okay bool) {
	// fix? quoted texts will eat singular quoted texts
	// to allow both singular and plural, you'd really want the affinity
	// of the property ( and then base this parsing on that )
	plural, singular := (flags&AllowPlural) != 0, (flags&AllowSingular) != 0
	if next, qs := *in, (QuotedTexts{}); plural && qs.Match(q, &next) {
		ret = &qs
		*in, okay = next, true
	} else {
		if singular {
			if text := (QuotedText{}); text.Match(q, &next) {
				ret = &text
				*in, okay = next, true
			} else if num := (MatchingNum{}); num.Match(q, &next) {
				ret = &num
				*in, okay = next, true
			} else if noun := (ExistingNoun{}); noun.Match(q, &next) {
				ret = &noun
				*in, okay = next, true
			} else if kind := (Kind{}); kind.Match(q, &next) {
				ret = &kind
				*in, okay = next, true
			}
		}
	}
	return
}
