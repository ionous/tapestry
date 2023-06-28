package grok

import "strings"

// fix: should be customizable

type info struct {
	determiners, kinds, traits spanList
	macros                     macroList
}

func (n *info) FindDeterminer(words []Word) (skip int) {
	_, skip = n.determiners.findPrefix(words)
	return
}

func (n *info) FindKind(words []Word) (skip int) {
	_, skip = n.kinds.findPrefix(words)
	return
}

func (n *info) FindTrait(words []Word) (skip int) {
	_, skip = n.traits.findPrefix(words)
	return
}

type MacroInfo struct {
	macroType MacroType
	str       string
	width     int
}

func (m *MacroInfo) Type() MacroType {
	return m.macroType
}

func (m *MacroInfo) String() string {
	return m.str
}

// number of words
func (m *MacroInfo) Width() int {
	return m.width
}

func (n *info) FindMacro(words []Word) (ret MacroInfo, okay bool) {
	if at, skipped := n.macros.findPrefix(words); skipped > 0 {
		w, t := n.macros.get(at)
		ret = MacroInfo{
			width:     skipped,
			macroType: t,
			str:       wordsToString(w),
		}
		okay = true
	}
	return
}

func wordsToString(w []Word) (ret string) {
	var b strings.Builder
	for i, w := range w {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(w.String())
	}
	return b.String()
}

var known = info{
	determiners: panicSpans([]string{
		"the", "a", "an", "some", "our",
		// ex. kettle of fish
		"a kettle of",
	}),
	macros: panicMacros(
		// tbd: flags need more thought.
		ManyToOne, "kind of", // for "a closed kind of container"
		ManyToOne, "kinds of", // for "are closed containers"
		ManyToOne, "a kind of", // for "a kind of container"
		// other macros
		OneToMany, "on", // on the x are the w,y,z
		OneToMany, "in",
		//
		ManyToMany, "suspicious of",
	),
	kinds: panicSpans([]string{
		"thing", "things",
		"container", "containers",
		"supporter", "supporters",
	}),
	traits: panicSpans([]string{
		"closed",
		"open",
		"openable",
		"transparent",
		"fixed in place",
	}),
}

var keywords = struct {
	and, are, called, comma, has, is uint64
}{
	and:    plainHash("and"),
	are:    plainHash("are"),
	called: plainHash("called"),
	comma:  plainHash(","),
	has:    plainHash("has"),
	is:     plainHash("is"),
}

// fix: would be better to have a push interface so we can just add things as we go
// this is easier for development though
type Results struct {
	Sources []Noun
	Targets []Noun // usually just one, except for nxm relations
	Macro   MacroInfo
}

type Noun struct {
	Det    []Word
	Name   []Word
	Traits [][]Word
	Kinds  [][]Word // it's possible, if rare, to apply multiple kinds
	// ex. The container called the coffin is a closed openable thing.
}

func Grok(p string) (ret Results, err error) {
	out := &Results{}
	if words, e := hashWords(p); e != nil {
		err = e
	} else {
		// scan for "is/are" or a macro verb, which ever comes first;
		// the order can reverse subjects and objects.
		for i, w := range words {
			if w.equals(keywords.is) || w.equals(keywords.are) {
				err = beingPhrase(out, words[:i], words[i+1:])
				break
			} else {
				if macro, ok := known.FindMacro(words[i:]); ok {
					out.Macro = macro
					err = macroPhrase(out, words[i+macro.Width():])
					break
				}
			}
		}
	}
	if err == nil {
		ret = *out
	}
	return
}
