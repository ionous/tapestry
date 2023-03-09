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

// these are two options for a Noun:
//  1. the closed ( and openable ) container  ( anonymous | called the box )
//  2. the name of a Noun [ and 1|2 ]
//
// only #1 supports traits; only #2 supports multiple nouns.
//
// note: inform doesn't support leading anonymous nouns ( ex. "a car is in the garage" )
// because it's not clear whether that means the last car, or some new generic car;
// however, "a car called Genevieve is in the garage" and "in the garage is a car" are allowed.
func genNouns(out *[]Noun, ws []Word, flag genFlag) (err error) {
	for nextName := ws; len(nextName) > 0; {
		if _, skip := known.determiners.findPrefix(nextName); skip >= len(nextName) {
			err = makeWordError(ws[0], "expected some sort of name")
		} else {
			det, name := nextName[:skip], nextName[skip:]
			nextName = nil // by default nothing else after this.
			if ts, e := parseTraitSet(name); e != nil {
				err = e
				break
			} else {
				// case 1: have a kind.
				if len(ts.kind) != 0 {
					// fix? there really is no reason "called" couldn't be smarter to split on "and"
					// being able to parse quoted text really helps us here.
					if d, n, e := chopCalled(name[ts.wordCount:]); e != nil {
						err = e
						break
					} else if len(n) == 0 && flag&AllowAnonymous == 0 {
						err = makeWordError(name[0], "expected a name")
						break
					} else {
						// the bits after called become the determiner and name
						// if they are empty, they're anonymous
						det, name = d, n
					}
				} else {
					// case 2: no kindness detected.
					// kind of odd, but we throw the traits out
					// and instead scan for "and" to allow a new name afterwards.
					ts = traitSet{}
					if before, after, e := anyAnd(name); e != nil {
						err = e
						break
					} else if after > 0 {
						name, nextName = name[:before], name[after:]
					}
				}
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
	return
}

// anonymous Noun: ex. "[the] container"
// vs. "[the] container called (the box)."
func chopCalled(ws []Word) (retDet, retName []Word, err error) {
	if len(ws) > 0 {
		if called := ws[0].equals(keywords.called); !called {
			err = makeWordError(ws[0], "unknown trailing text")
		} else {
			// the (det) and name directly follow the Word "called"
			retDet, retName, err = chopName(ws[1:])
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
