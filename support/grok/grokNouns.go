package grok

import "github.com/ionous/errutil"

type genFlag int

const (
	// fix: anonymous kinds should be permitted to target the most recently named noun
	AllowMany genFlag = 1 << iota
	AllowAnonymous
	OnlyOne
	OnlyNamed
)

// these are three options for a noun:
//  1. the [traits] (kind) called the (name of the noun.)
//  2. the [traits] (kind) [ more nouns... ]
//  3. the (name of a noun) [ more nouns... ]
//
// only 1 & 2 support traits; only 2 & 3 support additional nouns.
//
// note: inform doesn't support leading anonymous nouns ( ex. "the car is in the garage" )
// they point out: it's not clear whether that indicates the most recent noun, or some new generic noun.
// however, trailing anonymous nouns are allowed. ( ex. "in the garage is a car" )
func genNouns(out *[]Noun, ws []Word, flag genFlag) (err error) {
	for nextName := ws; len(nextName) > 0; {
		if _, skip := known.determiners.findPrefix(nextName); skip >= len(nextName) {
			err = makeWordError(nextName[0], "expected some sort of name")
		} else {
			det, name := nextName[:skip], nextName[skip:]
			nextName = nil // by default nothing else after this.
			if ts, e := parseTraitSet(name); e != nil {
				err = e
				break
			} else {
				postTraits := name[ts.wordCount:]
				if len(ts.kind) == 0 {
					// case 3: no kindness detected; throw the traits out
					ts, postTraits = traitSet{}, name
					if flag&AllowMany != 0 {
						if before, after, e := anyAnd(name); e != nil {
							err = e
							break
						} else if after > 0 {
							name, nextName = name[:before], name[after:]
						}
					}
				} else if called := len(postTraits) > 0 && postTraits[0].equals(keywords.called); !called {
					// case 2: an anonymous kind.
					name, det = nil, nil
				} else if d, n, e := chopName(postTraits[1:]); e != nil {
					err = e
					break
				} else {
					// case 1: any bits after "called" become the determiner and name
					det, name = d, n
					flag = flag & ^AllowMany // tbd: why couldn't "called" couldn't be smarter to split on "and"?
				}

				// more nouns may be allowed after "and"
				if flag&AllowMany != 0 {
					if _, after, e := anyAnd(postTraits); e != nil {
						err = e
						break
					} else if after > 0 {
						nextName = postTraits[after:]
					}
				}

				if len(name) == 0 && flag&AllowAnonymous == 0 {
					err = errutil.New("anonymous nouns not allowed.")
					break
				} else {
					// turn the name into a noun:
					*out = append(*out, Noun{
						Det:    det,
						Name:   name,
						Traits: ts.traits,
						Kinds:  ts.kinds(),
					})
				}
			}
		}
	}
	return
}

// the entire passed text is a name ( possibly with a prefix to start )
// errors only if the name is completely empty.
func chopName(ws []Word) (retDet, retName []Word, err error) {
	if cnt := len(ws); cnt == 0 {
		err = errutil.New("empty name")
	} else if _, skip := known.determiners.findPrefix(ws); skip >= len(ws) {
		err = makeWordError(ws[0], "no name found")
	} else {
		retDet, retName = ws[:skip], ws[skip:]
	}
	return
}
