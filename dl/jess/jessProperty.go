package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type PropertyValue interface {
	Assignment() rt.Assignment
}

type FieldName = string

// match the name of a property
func TryPropertyName(q JessContext, in InputState, kind string,
	accept func(Property, InputState), reject func(error),
) {
	// the article is optional, even at the start of a sentence where grammar often demands it.
	TryArticle(q, in, func(article *Article, next InputState) {
		q.Try(After(weaver.PropertyPhase), func(w weaver.Weaves, run rt.Runtime) {
			if field, width := q.FindField(kind, next.words); width > 0 {
				accept(Property{
					Article:   article,
					fieldName: field,
					Matched:   next.Cut(width),
				}, next.Skip(width))
			}
		}, reject)
	}, reject)
}

// match and write a property value
func generatePropertyValue(q JessContext, in InputState,
	noun, field string, isPlural bool,
	accept func(PropertyValue, InputState),
	reject func(error),
) {
	flags := AllowSingular
	if isPlural {
		flags = AllowPlural
	}
	if pv, ok := matchPropertyValue(q, &in, flags); !ok {
		reject(FailedMatch{"didn't understand the value", in})
	} else {
		q.Try(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) {
			if e := w.AddNounValue(noun, field, pv.Assignment()); e != nil {
				reject(e)
			} else {
				accept(pv, in)
			}
		}, reject)
	}
}
