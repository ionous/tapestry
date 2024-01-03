package en

import (
	"strings"
	"unicode"

	"github.com/ionous/inflect"
)

// Singularize attempts to return the singular form of the passed assumed plural string.
func Singularize(s string) (ret string) {
	if len(s) > 0 {
		ret = inflect.Singularize(s)
	}
	return
}

// Pluralize attempts to return the plural form of the passed assumed singular string.
func Pluralize(s string) (ret string) {
	if len(s) > 0 {
		ret = inflect.Pluralize(s)
	}
	return
}

// Capitalize returns a new string, starting the first word with a capital.
func Capitalize(s string) (ret string) {
	if len(s) > 0 {
		// fix? capitalize doesnt handle leading spaces well.
		// what should it do?
		var lead string
		if i := strings.IndexFunc(s, func(u rune) bool {
			return !unicode.IsSpace(u)
		}); i >= 0 {
			lead, s = s[:i], s[i:]
		}
		ret = lead + inflect.Capitalize(strings.ToLower(s))
	}
	return
}

// SentenceCase returns the passed string in lowercase, starting new sentences with capital letters.
func SentenceCase(s string) string {
	sentences := strings.Split(s, ". ")
	for i, s := range sentences {
		sentences[i] = Capitalize(s)
	}
	return strings.Join(sentences, ". ")
}

// IsCapitalized returns true if the passed string starts with an upper case letter.
func IsCapitalized(n string) (ret bool) {
	for _, r := range n {
		ret = unicode.IsUpper(r)
		break
	}
	return
}

// Titlecase returns a new string, starting every word with a capital.
func Titlecase(s string) (ret string) {
	if len(s) > 0 {
		ret = inflect.Titleize(s)
	}
	return
}

// StartsWith returns true if the passed string starts with any one of the passed strings in set
func xStartsWith(s string, set ...string) (ok bool) {
	for _, x := range set {
		if strings.HasPrefix(s, x) {
			ok = true
			break
		}
	}
	return ok
}

// Elide cuts strings after a certain length replacing the rest with ...
// todo: probably better to elide after words.
func xElide(s string, cutAfter int) (ret string) {
	const ellipse = "..."
	if cnt := len(s); cnt > 0 && cnt <= cutAfter {
		ret = s
	} else if min := len(ellipse); cnt < min {
		ret = ellipse
	} else {
		ret = s[:cutAfter] + ellipse
	}
	return
}

// StartsWithVowel returns true if the passed strings starts with a vowel or vowel sound.
// http://www.mudconnect.com/SMF/index.php?topic=74725.0
func xStartsWithVowel(str string) (vowelSound bool) {
	s := strings.ToUpper(str)
	if xStartsWith(s, "A", "E", "I", "O", "U") {
		if !xStartsWith(s, "EU", "EW", "ONCE", "ONE", "OUI", "UBI", "UGAND", "UKRAIN", "UKULELE", "ULYSS", "UNA", "UNESCO", "UNI", "UNUM", "URA", "URE", "URI", "URO", "URU", "USA", "USE", "USI", "USU", "UTA", "UTE", "UTI", "UTO") {
			vowelSound = true
		}
	} else if xStartsWith(s, "HEIR", "HERB", "HOMAGE", "HONEST", "HONOR", "HONOUR", "HORS", "HOUR") {
		vowelSound = true
	}
	return vowelSound
}

// IsSpace reports whether the rune is a space character as defined by lower ascii.
// '\t', '\n', '\v', '\f', '\r', ' '
// it specifically excludes non-breaking spaces.
func IsSpace(r rune) (ret bool) {
	switch r {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		ret = true
	}
	return
}

// Fields is similar to strings.Fields except it follows the rules of en.IsSpace
func Fields(s string) []string {
	return strings.FieldsFunc(s, IsSpace)
}

func Join(els []string) string {
	return strings.Join(els, " ")
}

// given a normalized phrase, return the first word
// ( exists to try to break up dependence on space as the separator of words )
func FirstWord(s string) string {
	return strings.Split(s, " ")[0]
}

// Normalize lowercases the passed string, trims spaces, and eats some kinds of punctuation.
// Ascii underscores (_) are treated as whitespace, ascii dashes (-) are kept; all other unicode punctuation gets removed.
// Whitespace gets removed at the front and end of the strings;
// any remaining groups of one or more ws characters get replaced by a single space.
func Normalize(s string) string {
	var out strings.Builder
	var spaced bool
	for _, r := range s {
		if r == '_' || IsSpace(r) {
			spaced = true
		} else if r == '-' || !unicode.IsPunct(r) {
			if out.Len() > 0 && spaced {
				out.WriteRune(' ')
			}
			out.WriteRune(unicode.ToLower(r))
			spaced = false
		}
	}
	return out.String()
}

// MixedCaseToSpaces takes a MixedCaseWord and splits it into lowercase words ( ex. mixed case word )
func MixedCaseToSpaces(s string) string {
	var out strings.Builder
	var prev bool
	for _, r := range s {
		l := unicode.ToLower(r)
		cap := l != r
		if !prev && cap && out.Len() > 0 {
			out.WriteRune(' ')
		}
		out.WriteRune(l)
		prev = cap
	}
	return out.String()
}
