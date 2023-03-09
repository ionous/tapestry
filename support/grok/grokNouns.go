package grok

import "github.com/ionous/errutil"

type genFlag int

const (
	// fix: anonymous kinds should be permitted to target the most recently named Noun
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
	for nextName, visits := ws, 0; len(nextName) > 0; visits++ {
		if allowMany := flag&AllowMany != 0; !allowMany && visits > 0 {
			err = makeWordError(nextName[0], "only expected one noun")
		} else if _, skip := known.determiners.findPrefix(nextName); skip >= len(nextName) {
			err = makeWordError(ws[0], "expected some sort of name")
		} else {
			det, name := nextName[:skip], nextName[skip:]
			nextName = nil // by default nothing else after this.
			if ts, e := parseTraitSet(name); e != nil {
				err = e
				break
			} else {
				if len(ts.kind) == 0 {
					// case 3: no kindness detected; throw the traits out
					// instead, if multiple nouns are allowed, separate by "and".
					// ( alt: if this were a done via a wrapper, skip this for "called" )
					if ts = (traitSet{}); flag&AllowMany != 0 {
						if before, after, e := anyAnd(name); e != nil {
							err = e
							break
						} else if after > 0 {
							name, nextName = name[:before], name[after:]
						}
					}
				} else {
					postTraits := name[ts.wordCount:]
					// case 1: any bits after "called" become the determiner and name
					if len(postTraits) > 0 && postTraits[0].equals(keywords.called) {
						if d, n, e := chopName(postTraits[1:]); e != nil {
							err = e
						} else {
							flag = flag & ^AllowMany // fix? why couldn't "called" couldn't be smarter to split on "and"
							det, name = d, n
						}
					} else {
						// case 2: there is no name, just kind; more nouns may be allowed after "and"
						if name, det = nil, nil; flag&AllowMany != 0 {
							if _, after, e := anyAnd(postTraits); e != nil {
								err = e
								break
							} else if after > 0 {
								nextName = postTraits[after:]
							}
						}
					}
				}
				if len(name) == 0 && flag&AllowAnonymous == 0 {
					err = makeWordError(name[0], "expected a name")
					break
				} else {
					// turn the name into a Noun:
					*out = append(*out, Noun{
						det:    det,
						name:   name,
						traits: ts.traits,
						kinds:  ts.kinds(),
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
