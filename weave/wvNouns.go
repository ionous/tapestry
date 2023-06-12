package weave

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

func (cat *Catalog) AssertNounKind(opNoun, opKind string) error {
	return cat.Schedule(assert.NounPhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		_, name := d.StripDeterminer(opNoun)
		_, kind := d.StripDeterminer(opKind)

		if noun, ok := UniformString(name); !ok {
			err = InvalidString(opNoun)
		} else if kn, ok := UniformString(kind); !ok {
			err = InvalidString(opKind)
		} else if e := cat.writer.Noun(d.name, noun, d.singularize(kn), at); e != nil {
			err = e
		} else {
			cat.domainNouns[domainNoun{d.name, noun}] = &ScopedNoun{domain: d, name: noun}
			err = d.makeNames(noun, name, at)
		}
		return
	})
}

func (d *Domain) makeNames(noun, name, at string) (err error) {
	q := d.catalog.writer
	// if the original got transformed into underscores
	// write the original name (ex. "toy boat" vs "toy_boat" )
	var out []string
	if clip := strings.TrimSpace(name); clip != noun {
		out = append(out, clip)
	}
	out = append(out, noun)

	// generate additional names by splitting the lowercase uniform name on the underscores:
	split := strings.FieldsFunc(noun, lang.IsBreak)
	if cnt := len(split); cnt > 1 {
		// in case the name was reduced due to multiple separators
		if breaks := strings.Join(split, "_"); breaks != noun {
			out = append(out, breaks)
		}
		// write individual words in increasing rank ( ex. "boat", then "toy" )
		// note: trailing words are considered "stronger"
		// because adjectives in noun names tend to be first ( ie. "toy boat" )
		for i := len(split) - 1; i >= 0; i-- {
			word := split[i]
			out = append(out, word)
		}
	}
	for i, name := range out {
		// ignore duplicate errors here.
		// since these are generated, there's probably very little the user could do about them.
		if e := q.Name(d.name, noun, name, i, at); e != nil && !errors.Is(e, mdl.Duplicate) {
			err = e
			break
		}
	}
	return
}
