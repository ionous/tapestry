package jess

// see: TryNounPropertyValue
func (op *NounPropertyValue) GetBuilder() Builder {
	return &op.builder
}

// `The pen has (the) description (of) "mightier than the sword."`
// `The bottle has age 42 and the description "A plain glass bottle."`
func TryNounPropertyValue(q JessContext, in InputState,
	accept func(PromisedMatcher), reject func(error),
) {
	// the word "has/have" splits the noun from the property
	if tgtProp, ok := keywordSplit(in, keywords.Has, keywords.Have); !ok {
		reject(FailedMatch{"noun property phrases require the word 'have' or 'has'", in})
	} else {
		// match a name to a noun ( or generate one )
		TryPropertyNoun(q, tgtProp.lhs, func(pn PropertyNoun, postName InputState) {
			pb := PropertyBuilder{Context: q, noun: pn}
			if postName.Len() != 0 {
				reject(FailedMatch{"noun property had unexpected words after the name", postName})
			} else {
				// given the kind, match the property names and values.
				// ( required to separate it from unquoted nouns or kinds as values )
				TryPropertyPossessions(&pb, tgtProp.rhs, func(ps PropertyPossessions) {
					accept(&NounPropertyValue{
						PropertyNoun:        pn,
						Has:                 Words{Matched: tgtProp.matched},
						PropertyPossessions: ps,
						builder:             pb,
					})
				}, reject)
			}
		}, reject)
	}
}

// recursively generate a list of property values.
// expects to eat the entire input, and therefore doesn't return next InputState.
// ex. `age of 42, color red, and the description "A plain glass bottle."`
func TryPropertyPossessions(pb *PropertyBuilder, in InputState,
	accept func(PropertyPossessions),
	reject func(error),
) {
	// matches an optional article, and field name.
	TryPropertyName(pb.Context, in, pb.GetKind(), func(prop Property, valueOf InputState) {
		if valueOf.Len() == 0 {
			reject(FailedMatch{"it seems you were trying to define a property without a value", in})
		} else {
			// optionally, the word "of" can separate the property name and value
			of := matchOf(&valueOf)
			// find the value
			TryPropertyValue(pb.Context, valueOf, AllowSingular, func(pv PropertyValue, rest InputState) {
				pb.addProperty(prop.fieldName, pv.Assignment())
				TryAdditionalPossessions(pb, rest,
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
func TryAdditionalPossessions(pb *PropertyBuilder, in InputState,
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
			TryPropertyPossessions(pb, in, func(pv PropertyPossessions) {
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
