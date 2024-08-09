package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// see: TryNounPropertyValue
func (op *NounPropertyValue) PromiseMatcher() PromiseMatcher {
	return op
}

// `The pen has (the) description (of) "mightier than the sword."`
// `The bottle has age 42 and the description "A plain glass bottle."`
func TryNounPropertyValue(q JessContext, in InputState,
	accept func(PromiseMatcher), reject func(error),
) {
	// the word "has/have" splits the noun from the property
	if tgtProp, ok := keywordSplit(in, keywords.Has, keywords.Have); !ok {
		reject(FailedMatch{"noun property phrases require the word 'have' or 'has'", in})
	} else {
		// match a name to a noun ( or generate one )
		TryNamedNoun(q, tgtProp.lhs, func(nn NamedNoun, an ActualNoun, postName InputState) {
			if postName.Len() != 0 {
				reject(FailedMatch{"noun property had unexpected words after the name", postName})
			} else {
				// given the kind, match the property names and values.
				// ( required to separate it from unquoted nouns or kinds as values )
				TryPropertyPossessions(q, tgtProp.rhs, an, func(ps PropertyPossessions) {
					accept(&NounPropertyValue{
						NamedNoun:           nn,
						Has:                 Words{Matched: tgtProp.matched},
						PropertyPossessions: ps,
					})
				}, reject)
			}
		}, reject)
	}
}

// recursively generate a list of property values.
// expects to eat the entire input, and therefore doesn't return next InputState.
// ex. `age of 42, color red, and the description "A plain glass bottle."`
func TryPropertyPossessions(q JessContext, in InputState,
	an ActualNoun,
	accept func(PropertyPossessions),
	reject func(error),
) {
	// matches an optional article, and field name.
	TryPropertyName(q, in, an.Kind, func(prop Property, valueOf InputState) {
		if valueOf.Len() == 0 {
			reject(FailedMatch{"it seems you were trying to define a property without a value", in})
		} else {
			// optionally, the word "of" can separate the property name and value
			of := matchOf(&valueOf)
			// find the value
			TryPropertyValue(q, valueOf, an, AllowSingular, func(pv PropertyValue, rest InputState) {
				// try to apply that value to the noun
				// FIX: this should be delayed until the phrase is completely matched.
				q.Try(weaver.AnyPhase, func(w weaver.Weaves, run rt.Runtime) {
					if e := w.AddNounValue(an.Name, prop.fieldName, pv.Assignment()); e != nil {
						reject(e)
					}
				}, reject)
				TryAdditionalPossessions(q, rest, an,
					func(more *AdditionalPossessions) {
						accept(PropertyPossessions{
							Property:              prop,
							Of:                    of,
							PropertyValue:         pv,
							AdditionalPossessions: more,
						})
					}, reject)
			}, reject)
		}
	}, reject)
}

// recursively a series of property values.
// expects to eat the entire input, and therefore doesn't return next InputState.
// because the values are optional, the input can be empty,
// in which case the accepted additional values are nil.
// "ex. "and the description ...."
func TryAdditionalPossessions(q JessContext, in InputState,
	an ActualNoun,
	accept func(*AdditionalPossessions),
	reject func(error),
) {
	if in.Len() == 0 {
		accept(nil)
	} else {
		var ca CommaAnd
		if !ca.InputMatch(&in) {
			reject(FailedMatch{"unknown words following values", in})
		} else {
			TryPropertyPossessions(q, in, an, func(pv PropertyPossessions) {
				accept(&AdditionalPossessions{
					CommaAnd:            ca,
					PropertyPossessions: pv,
				})
			}, reject)
		}
	}
}

func matchOf(in *InputState) (ret *Words) {
	if width := in.MatchWord(keywords.Of); width > 0 {
		*in, ret = in.Skip(width), &Words{Matched: in.Cut(width)}
	}
	return
}
