package regen

import "strings"

type Type struct {
	m map[string]interface{}
	// cache: which is kind of questionable in a one-off tool, but whatever.
	params  []Param
	exstr   *bool
	cmt     *string
	english bool // returned by tokens
	tokens  []string
	groups  []string
	slots   []string
}

func NewType(m map[string]interface{}) *Type {
	return &Type{m: m}
}

func (t *Type) Name() string {
	return StringOf("name", t.m)
}

func (t *Type) Uses() string {
	return StringOf("uses", t.m)
}

func (t *Type) Slots() []string {
	if t.slots == nil {
		w := MapOf("with", t.m)
		t.slots = ListOf("slots", w)
	}
	return t.slots
}

// all but the primary group
func (t *Type) OtherGroups() []string {
	return t.AllGroups()[1:]
}

func (t *Type) AllGroups() []string {
	if t.groups == nil {
		t.groups = ListOf("group", t.m)
	}
	return t.groups
}

// return a json literal: either an array or a quoted string
func (t *Type) LiteralComment() string {
	if t.cmt == nil {
		var cmt string
		if desc := ListOf("desc", t.m); desc != nil {
			cmt = MarshalIndent(desc)
		} else if desc := StringOf("desc", t.m); len(desc) > 0 {
			cmt = MarshalIndent(desc)
		} else if desc := MapOf("desc", t.m); desc != nil {
			// desc is a description with the form: "label: short description. long description...
			label := sentenceOf("label", desc)
			short := sentenceOf("short", desc)
			long := sentenceOf("long", desc)
			var ar []string
			if len(label) > 0 {
				x := strings.ToLower(strings.ReplaceAll(label, " ", "_"))
				if n := t.Name(); x != n+"." {
					ar = append(ar, label)
				}
			}
			if short != "" {
				ar = append(ar, short)
			}
			if long != "" {
				ar = append(ar, long)
			}
			if len(ar) > 0 {
				cmt = MarshalIndent(ar)
			}
		}
		t.cmt = &cmt
	}
	return *t.cmt
}

// for str: is the str type limited to just the options provided?
func (t *Type) Exclusive() bool {
	if t.exstr == nil {
		var exstr bool
		w := MapOf("with", t.m)
		if ps := MapOf("params", w); len(ps) > 0 {
			// if  nothing was excluded from options;
			// the name of the str wasnt in the params;
			// so the choices arent flexible.
			opts := t.Options()
			exstr = len(opts) == len(ps)
		}
		t.exstr = &exstr
	}
	return *t.exstr
}

// for str: list of possible choices
// filters the name of the str itself ( which is used by typespec to indicate any choice is okay )
func (t *Type) Options() []Param {
	if t.params == nil {
		var opt []Param
		n := t.Name()
		w := MapOf("with", t.m)
		ps := MapOf("params", w)
		self := Tokenize(n)
		_, tokens := t.Tokens()
		for _, k := range tokens {
			if len(k) > 0 && k[0] == '$' {
				v := MapOf(k, ps)
				if self != k {
					opt = append(opt, Param{
						Name:  StringOf("value", v),
						Label: StringOf("label", v),
					})
				}
			}
		}
		t.params = opt
	}
	return t.params
}

// for swap
func (t *Type) Picks() []Param {
	if t.params == nil {
		picks := make([]Param, 0)
		w := MapOf("with", t.m)
		ps := MapOf("params", w)
		_, tokens := t.Tokens()
		for _, k := range tokens {
			if len(k) > 0 && k[0] == '$' {
				v := MapOf(k, ps)
				picks = append(picks, Param{
					Name:  Detokenize(k),
					Label: StringOf("label", v),
					Type:  StringOf("type", v),
				})
			}
		}
		t.params = picks
	}
	return t.params
}

// for flow, the name used in story file commands
func (t *Type) Flow() (ret string) {
	if english, ts := t.Tokens(); !english && len(ts) > 0 {
		if s := ts[0]; len(s) > 0 && s[0] != '$' {
			// if s != t.Name() {
			ret = s
			// }
		}
	}
	return
}

// for flow: is the first label anonymous
// (note: all english phrases are this way automatically )
func (t *Type) Trim() (ret bool) {
	if english, _ := t.Tokens(); english {
		ret = true
	} else if ts := t.Terms(); len(ts) > 0 {
		ret = ts[0].Label == "_"
	}
	return
}

// for flow: the possible options
func (t *Type) Terms() []Param {
	if t.params == nil {
		terms := make([]Param, 0)
		w := MapOf("with", t.m)
		ps := MapOf("params", w)
		_, tokens := t.Tokens()
		for _, k := range tokens {
			if len(k) > 0 && k[0] == '$' {
				v := MapOf(k, ps)
				terms = append(terms, Param{
					Name:     Detokenize(k),
					Label:    StringOf("label", v),
					Type:     StringOf("type", v),
					Repeats:  BoolOf("repeats", v),
					Optional: BoolOf("optional", v),
				})
			}
		}
		t.params = terms
	}
	return t.params
}

// return all token keys and labels
func (t *Type) Tokens() (bool, []string) {
	if t.tokens == nil {
		var story, modeling bool
		for _, k := range t.AllGroups() {
			if k == "story" {
				story = true
			} else if k == "modeling" {
				modeling = true
			}
		}
		w := MapOf("with", t.m)
		t.tokens = ListOf("tokens", w)
		t.english = story && !modeling
	}
	return t.english, t.tokens
}

// for flow: the "english" phrasing
// only exists currently for some of the story specs
// so this limits itself particularly to those
func (t *Type) Phrase() (ret string) {
	english, tokens := t.Tokens()
	if english {
		var b strings.Builder
		for _, k := range tokens {
			if len(k) > 0 {
				if k[0] != '$' {
					b.WriteString(k)
				} else {
					b.WriteRune('{')
					b.WriteString(strings.ToLower(k))
					b.WriteRune('}')
				}
			}
		}
		ret = b.String()
	}
	return
}
