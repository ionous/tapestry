package story

import (
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/eph"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/reader"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt/kindsOf"
	"github.com/ionous/errutil"
)

type NounImporter interface {
	ImportNouns(k *Importer) (err error)
}

func (*Pronoun) ImportNouns(*Importer) (err error) {
	// FIX: pronoun(s) can indicate plurality
	return
}

func (op *NounPhrase) ImportNouns(k *Importer) (err error) {
	if imp, ok := op.Value.(NounImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		err = imp.ImportNouns(k)
	}
	return
}

func ImportNamedNouns(k *Importer, els []NamedNoun) (err error) {
	for _, el := range els {
		if e := el.ImportNouns(k); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (op *NamedNoun) ImportNouns(k *Importer) (err error) {
	// declare a noun class that has several default fields
	if once := "noun"; k.Once(once) {
		k.WriteOnce(&eph.EphKinds{Kinds: "objects", From: kindsOf.Kind.String()})
		// common or proper nouns ( rabbit, vs. Roger )
		k.AddImplicitAspect("noun_types", "objects", "common_named", "proper_named", "counted")
		// whether a player can refer to an object by its name.
		k.AddImplicitAspect("private_names", "objects", "publicly_named", "privately_named")
	}
	//
	if cnt, ok := lang.WordsToNum(op.Determiner.Str); !ok {
		err = op.ReadNamedNoun(k)
	} else {
		err = op.ReadCountedNoun(k, cnt)
	}
	return
}

// ex. "two triangles" -> triangle is a kind of thing
// fix? consider a specific counted noun phrase; the noun phrase needs more work.
// also, we probably want noun stacks not individually duplicated names
func (op *NamedNoun) ReadCountedNoun(k *Importer, cnt int) (err error) {
	if once := "printed_name"; k.Once(once) {
		k.WriteOnce(&eph.EphKinds{Kinds: "objects", Contain: []eph.EphParams{{Name: "printed_name", Affinity: eph.Affinity{eph.Affinity_Text}}}})
	}

	kind := op.Name.String()
	k.Write(&eph.EphKinds{Kinds: kind, From: "thing"})
	for i := 0; i < cnt; i++ {
		noun := k.newCounter(kind, reader.Position{})
		k.Write(&eph.EphNouns{Noun: noun, Kind: kind})
		k.Write(&eph.EphValues{Noun: noun, Field: "counted", Value: &literal.BoolValue{true}})
		// fix: we probably want the kind to be singular... can we make a "dependent" command like that?
		// a program in ephemera space... or a eph.literal that we manually look for when processing values?
		k.Write(&eph.EphValues{Noun: noun, Field: "printed_name", Value: &literal.TextValue{kind}})
		// needed for relations, etc.
		k.Env().Recent.Nouns.Add(noun)
	}
	return
}

func (op *NamedNoun) ReadNamedNoun(k *Importer) (err error) {
	noun := op.Name.String()
	k.Env().Recent.Nouns.Add(noun)
	detStr, detFound := composer.FindChoice(&op.Determiner, op.Determiner.Str)
	// setup the indefinite article
	if !detFound {
		// create a "indefinite article" field for all objects
		if k.Once("named_noun") {
			k.WriteOnce(&eph.EphKinds{Kinds: "objects", Contain: []eph.EphParams{{Name: "indefinite_article", Affinity: eph.Affinity{eph.Affinity_Text}}}})
		}
		// set the indefinite article field
		k.Write(&eph.EphValues{Noun: noun, Field: "indefinite_article", Value: &literal.TextValue{detStr}})
	}
	// pick common or proper based on noun capitalization.
	// fix: implicitly generated facts should be considered preliminary so that authors can override them.
	if detStr == "our" {
		if first, _ := utf8.DecodeRuneInString(noun); unicode.ToUpper(first) == first {
			k.Write(&eph.EphValues{Noun: noun, Field: "proper_named", Value: &literal.BoolValue{true}})
		}
	}
	return
}

// ex. "[the box] (is a) (closed) (container) ((on) (the beach))"
func (op *KindOfNoun) ImportNouns(k *Importer) (err error) {
	// we collected the nouns and delayed processing them till now.
	kind := op.Kind.String()
	for _, noun := range k.Env().Recent.Nouns.Subjects {
		k.Write(&eph.EphNouns{Noun: noun, Kind: kind})
		for _, trait := range op.Trait {
			k.Write(&eph.EphValues{Noun: noun, Field: trait.String(), Value: &literal.BoolValue{true}})
		}
	}
	if op.NounRelation != nil {
		err = op.NounRelation.ImportNouns(k)
	}
	return
}

// ex. [the cat and the hat] (are) (in) (the book)
// ex. [Hector and Maria] (are) (suspicious of) (Santa and Santana).
func (op *NounRelation) ImportNouns(k *Importer) (err error) {
	if e := k.Env().Recent.Nouns.CollectObjects(func() (err error) {
		return ImportNamedNouns(k, op.Nouns)
	}); e != nil {
		err = e
	} else {
		rel := op.Relation.String()
		for _, subject := range k.Env().Recent.Nouns.Subjects {
			for _, object := range k.Env().Recent.Nouns.Objects {
				k.Write(&eph.EphRelatives{Rel: rel, Noun: object, OtherNoun: subject})
			}
		}
	}
	return
}

//
func (op *NounTraits) ImportNouns(k *Importer) (err error) {
	for _, trait := range op.Trait {
		for _, noun := range k.Env().Recent.Nouns.Subjects {
			k.Write(&eph.EphValues{Noun: noun, Field: trait.String(), Value: &literal.BoolValue{true}})
		}
	}
	return
}
