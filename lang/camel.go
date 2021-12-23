package lang

import (
	"strings"
	"unicode"
)

type BreakcaseOptions uint8

const (
	OPT_SKIPIDS BreakcaseOptions = 1 << iota
	OPT_LOWER
	OPT_UPPER
)

// Breakcase turns runs of whitespace into single underscores. It does not change casing.
// "some   BIG Example" => "some_BIG_Example".
func Breakcase(name string) string {
	return OptionCase(name, 0)
}

func SpecialBreakcase(name string) string {
	return OptionCase(name, OPT_SKIPIDS)
}

func LowerBreakcase(name string) string {
	return OptionCase(name, OPT_LOWER)
}

func UpperBreakcase(name string) string {
	return OptionCase(name, OPT_UPPER)
}

// eventually, these transforms will happen at assembly time
func OptionCase(name string, opts BreakcaseOptions) (ret string) {
	if len(name) == 0 || ((opts&OPT_SKIPIDS) != 0 && (name[0] == '#' || name[0] == '$')) {
		ret = name
	} else {
		var b strings.Builder
		var needBreak, canBreak bool
		for _, r := range name {
			brakes := IsBreak(r)
			if ignorable := !brakes && IsIgnorable(r); ignorable {
				// consolidate all ignorable characters...
				// eventually writing a break if we've written something of note.
				needBreak = canBreak
			} else {
				// dont write a consolidated break if we are writing an explicit break
				// ex. "  _" -> just write a single underscore, not two.
				if needBreak && !brakes {
					b.WriteRune(breaker)
				}
				if (opts & OPT_LOWER) != 0 {
					r = unicode.ToLower(r)
				} else if (opts & OPT_UPPER) != 0 {
					r = unicode.ToUpper(r)
				}
				b.WriteRune(r)
				canBreak = !brakes
				needBreak = false
			}
		}
		ret = b.String()
	}
	return
}

// IsBreak returns true for the set of characters which breaks words in breakcase
func IsBreak(r rune) bool {
	return r == breaker
}

// IsIgnorable returns true for the set of characters which will be consolidated by breakcase
func IsIgnorable(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSpace(r)
}

const breaker = '_'

// HasBadPunct returns true for non-breakcase punctuation and spaces.
func HasBadPunct(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		badMatch := !IsBreak(r) && IsIgnorable(r)
		return badMatch
	}) >= 0
}

// Fields is similar to strings.Fields except this splits on dashes, case changes, spaces, and the like:
// the same rules as Camelize
func Fields(name string) []string {
	p := combineCase(name, false, false)
	return p.flush(true)
}

// Turns names from mixed case to lowercase breakcase
func Underscore(name string) string {
	var b strings.Builder
	for i, el := range Fields(name) {
		if i > 0 {
			b.WriteRune('_')
		}
		b.WriteString(strings.ToLower(el))
	}
	return b.String()
}

// fix? backwards compat: skip strings starting with # or $
func SpecialUnderscore(name string) (ret string) {
	if len(name) == 0 || name[0] == '#' || name[0] == '$' {
		ret = name
	} else {
		ret = Underscore(name)
	}
	return
}

// fix. this horrible algorithm sure needs to change
// itd be fine i think to split first and then combine with word rules even
// possibly worth considering ditching camelCasing anyway:
// only support de-camelization into lower fields for template matching.
func combineCase(name string, changeFirst, changeAny bool) *parts {
	type word int
	const (
		noword word = iota
		letter
		number
	)
	var parts parts
	inword, wasUpper := noword, 0
	changeCase := changeFirst
	for _, r := range name {
		if r == '_' || unicode.In(r, unicode.Hyphen) || unicode.IsSpace(r) {
			inword = noword
			wasUpper = 0
			continue
		}
		if unicode.IsDigit(r) {
			if sameWord := inword == number; !sameWord {
				parts.flush(true)
			}
			parts.WriteRune(r)
			wasUpper = 0
			inword = number
		} else if unicode.IsLetter(r) {
			// classify some common word changes
			lower := unicode.ToLower(r)
			currUpper := lower != r
			// what should we do after a string of uppercase characters?
			var afterUpper bool
			if changeAny { // switch to make old behavior happy.
				afterUpper = (wasUpper > 0 && !currUpper)
			} else {
				afterUpper = (wasUpper == 1 && !currUpper)
			}
			sameWord := (inword == letter) && ((wasUpper > 0 == currUpper) || afterUpper)
			// everything gets lowered
			if currUpper && changeCase {
				r = lower
			}
			changeCase = true
			if !sameWord {
				parts.flush(wasUpper <= 1)
				wasUpper = 0

				// hack for camelCasing.
				if len(parts.arr) > 0 && changeAny {
					r = unicode.ToUpper(r)
				}
			}
			parts.WriteRune(r) // docs say err is always nil
			if currUpper {
				wasUpper++
			} else {
				wasUpper = 0
			}
			inword = letter
		}
	}
	return &parts
}

type parts struct {
	queued rune
	str    strings.Builder
	arr    []string
}

func (p *parts) WriteRune(r rune) {
	if prev := p.queued; prev > 0 {
		p.str.WriteRune(prev)
	}
	p.queued = r
}

func (p *parts) flush(all bool) []string {
	if prev := p.queued; all && prev > 0 {
		p.str.WriteRune(prev)
		p.queued = 0
	}
	if p.str.Len() > 0 {
		p.arr = append(p.arr, p.str.String())
		p.str.Reset()
	}
	return p.arr
}

func (p *parts) join() string {
	return strings.Join(p.flush(true), "")
}
