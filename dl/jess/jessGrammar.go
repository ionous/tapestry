package jess

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
)

// generate an user input parser from a string written in a special format.
// words separated by forward slashes indicate a choice between words;
// double dashes as a choice means defer trying to match a word;
// the equal sign is used as an escape for forward slashes, open brackets, and dashes;
// brackets are used to indicate noun names of the specified type.
// ex.  `hang [objects] on/onto/-- [objects]`,
// public for testing
func BuildPhrase(phrase string) (ret grammar.ScannerMaker, err error) {
	if parts := strings.Fields(phrase); len(parts) > 0 {
		series := make([]grammar.ScannerMaker, 0, len(parts))
		for _, str := range parts {
			if cnt := len(str); str[0] == openBracket && str[cnt-1] == closeBracket {
				kind := str[1 : cnt-1]
				series = append(series, &grammar.Noun{Kind: kind})
			} else {
				if words, e := buildWords(str); e != nil {
					err = e
				} else {
					series = append(series, &grammar.Words{Words: words})
				}
			}
		}
		if err == nil {
			ret = &grammar.Sequence{Series: series}
		}
	}
	return
}

const (
	escapeRune   = '='
	chooseRune   = '/'
	dash         = '-'
	openBracket  = '['
	closeBracket = ']'
)

// tdd: should there be a "match any one thing", "match anything or nothing", etc?
// or even just regex with groups used as the nouns
func buildWords(w string) (ret []string, err error) {
	var dashCount int
	var str strings.Builder
	var escape bool
Loop:
	for _, n := range w {
		if escape {
			switch n {
			case escapeRune, chooseRune, dash, openBracket:
				str.WriteRune(n)
				escape = false
			default:
				err = errors.New("only the equal sign, forward slash, or dash can follow an equal sign.")
				break Loop
			}
		} else {
			if n == dash {
				if dashCount < 2 && str.Len() == 0 {
					dashCount++
				} else {
					err = errors.New(doubleDashMsg)
					break Loop
				}
			} else if dashCount > 0 {
				err = errors.New("to avoid ambiguity, double dashes can only appear as the last option.")
			} else {
				switch n {
				default:
					str.WriteRune(n)
				case openBracket:
					// fix: we should allow substitutions as choices.
					err = errors.New("bracketed text indicates substitutions; substitutions can only be used as standalone words.")
					break Loop
				case escapeRune:
					escape = true
				case chooseRune:
					if str.Len() == 0 {
						err = errors.New(`word seems to be empty. 
to match a forward slash, precede it with an equal sign.
to match no word at all, use the trailing double dash syntax.`)
						break Loop
					}
					ret = append(ret, str.String())
					str.Reset()
				}
			}
		}
	}
	if useBlank := dashCount == 2; dashCount > 0 && !useBlank {
		err = errors.New(doubleDashMsg)
	} else if useBlank {
		ret = append(ret, "") // blank word to indicate a skip
	} else {
		ret = append(ret, str.String())
	}
	return
}

const doubleDashMsg = "to avoid ambiguity with the trailing double dash syntax, precede actual dashes by an equal sign."
