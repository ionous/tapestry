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
// note: inform doesn't support leading anonymous nouns ( ex. "the thing is in the garage" )
// they point out: it's not clear whether that indicates the most recent noun, or some new generic noun.
// however, trailing anonymous nouns are allowed. ( ex. "in the garage is a thing" )
func grokNouns(known Grokker, out *[]Name, ws []Word, flag genFlag) (err error) {
	for nextName := ws; len(nextName) > 0; {
		if det, e := known.FindArticle(nextName); e != nil {
			err = e
		} else if skip := MatchedLen(det.Matched); skip >= len(nextName) {
			err = makeWordError(nextName[0], "expected some sort of name")
		} else {
			name := nextName[skip:]
			nextName = nil // by default nothing else after this.
			if ts, e := ParseTraitSet(known, name); e != nil {
				err = e
				break
			} else {
				var exact bool
				postTraits := name[ts.WordCount:]
				if !ts.hasKind() {
					// case 3: no kindness detected; throw the traits out
					ts, postTraits = TraitSet{}, name
					if flag&AllowMany != 0 {
						if before, after, e := anyAnd(name); e != nil {
							err = e
							break
						} else if after > 0 {
							name, nextName = name[:before], name[after:]
						}
					}
				} else {
					// does it have a "called ..." some name trailing phrase?
					called := len(postTraits) > 0 && postTraits[0].equals(Keyword.Called)
					if !called {
						// case 2a: a counted kind: "two cats are on the bed."
						// case 2b: an anonymous kind: "a container is in the lobby."
						det.Matched = nil // erases the article, leaves the count if any ( ex. 2 )
						name = nil
					} else {
						if det.Count > 0 {
							err = errutil.New("can't name counted nouns")
							break
						} else {
							if d, n, e := chopArticle(known, postTraits[1:]); e != nil {
								err = e
								break
							} else {
								// case 1: any bits after "called" become the determiner and name
								det, name = d, n
								flag = flag & ^AllowMany // tbd: why couldn't "called" couldn't be smarter to split on "and"?
								exact = true
							}
						}
					}
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

				if len(name) == 0 && det.Count == 0 && (flag&AllowAnonymous == 0) {
					err = errutil.New("anonymous nouns not allowed.")
					break
				} else {
					// turn the name into a noun:
					*out = append(*out, Name{
						Article: det,
						Span:    name,
						Traits:  ts.Traits,
						Kinds:   ts.kinds(),
						Exact:   exact,
					})
				}
			}
		}
	}
	return
}

// the entire passed text is a name ( possibly with a prefix to start )
// errors only if the name is completely empty.
func chopArticle(known Grokker, ws []Word) (retDet Article, retName []Word, err error) {
	if cnt := len(ws); cnt == 0 {
		err = errutil.New("empty name")
	} else if det, e := known.FindArticle(ws); e != nil {
		err = e
	} else if skip := MatchedLen(det.Matched); skip >= len(ws) {
		err = makeWordError(ws[0], "no name found")
	} else {
		retDet, retName = det, ws[skip:]
	}
	return
}
