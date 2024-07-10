package match

import (
	"fmt"
	"strings"
)

// return the name after removing leading articles
// eats any errors it encounters and returns the original name
func StripArticle(str string) (ret string) {
	var okay bool
	if parts, e := TokenizeString(str); e == nil && len(parts) > 1 {
		if _, width := FindCommonArticles(parts); width > 0 {
			rest := parts[width:]
			if s, n := Stringify(rest); n == len(rest) {
				ret, okay = s, true
			}
		}
	}
	if !okay {
		ret = str
	}
	return
}

// turn all of the passed tokens into a helpful string representation
func DebugStringify(ts []TokenValue) (ret string) {
	var out strings.Builder
	for _, tv := range ts {
		if out.Len() > 0 && tv.Token != Stop && tv.Token != Comma {
			out.WriteRune(' ')
		}
		out.WriteString(tv.String())
	}
	return out.String()
}

// turn a series of string tokens into a space padded string
// returns the number of string tokens consumed.
func Stringify(ts []TokenValue) (ret string, width int) {
	return toString(ts, false)
}

// turn a series of string tokens into a normalized string
// returns the number of string tokens consumed.
// somewhat dubious because it skips inflect.Normalize
func Normalize(ts []TokenValue) (ret string, width int) {
	return toString(ts, true)
}

func toString(ts []TokenValue, lower bool) (ret string, width int) {
	var out strings.Builder
	for i, tv := range ts {
		if i == 0 && tv.Token == Quoted {
			addString(&out, tv, lower)
			width++
			break
		} else if tv.Token != String {
			break
		} else {
			if out.Len() > 0 {
				out.WriteRune(' ')
			}
			addString(&out, tv, lower)
			width++
		}
	}
	return out.String(), width
}

func addString(out *strings.Builder, tv TokenValue, lower bool) {
	str := tv.String()
	if lower {
		str = strings.ToLower(str)
	}
	out.WriteString(str)
}

// same as Normalize but errors if all of the tokens weren't consumed.
func NormalizeAll(ts []TokenValue) (ret string, err error) {
	if str, n := Normalize(ts); n == len(ts) {
		ret = str
	} else {
		out := DebugStringify(ts)
		err = fmt.Errorf("couldn't normalize %q", out)
	}
	return
}

// for now, the common articles are a fixed set.
// when the author specifies some particular indefinite article for a noun
// that article only gets used for printing the noun;
// it doesn't enhance the parsing of the story.
// ( it would take some work to lightly hold the relation between a name and an article
// then parse a sentence matching names to nouns in the
// fwiw: the articles in inform also seems to be predetermined in this way.  )
func FindCommonArticles(ts []TokenValue) (ret Span, width int) {
	if m, skip := determiners.FindPrefix(ts); skip > 0 {
		ret, width = m, skip
	}
	return
}

// fix? i feel like this should be part of package inflect instead
var determiners = PanicSpans("the", "a", "an", "some", "our")
