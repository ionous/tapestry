package story

import (
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/ephemera/eph"
	"git.sr.ht/~ionous/iffy/ephemera/reader"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/tables"
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
		// common or proper nouns ( rabbit, vs. Roger )
		k.NewImplicitAspect("noun_types", "objects", "common_named", "proper_named", "counted")
		// whether a player can refer to an object by its name.
		k.NewImplicitAspect("private_names", "objects", "publicly_named", "privately_named")
	}
	//
	if cnt, ok := lang.WordsToNum(op.Determiner.Str); !ok {
		err = op.ReadNamedNoun(k)
	} else {
		err = op.ReadCountedNoun(k, cnt)
	}
	return
}

// fix? consider a specific counted noun phrase;
// the noun phrase needs more work.
// also, we probably want noun stacks not individually duplicated names
func (op *NamedNoun) ReadCountedNoun(k *Importer, cnt int) (err error) {
	// declare the existence of the field "printed name"
	if once := "printed_name"; k.Once(once) {
		domain := k.Env().Game.Domain
		things := k.NewDomainName(domain, "objects", tables.NAMED_KINDS, once)
		field := k.NewDomainName(domain, "printed_name", tables.NAMED_FIELD, once)
		k.NewField(things, field, tables.PRIM_TEXT, "")
	}

	// generate the singular noun name and singular kind name
	// ex. "two triangles" -> triangle is a kind of thing
	baseName := op.Name.String()
	if cnt > 1 {
		baseName = lang.Singularize(baseName)
	}
	at := op.Name.At.String()
	namedSingularKind := k.NewName(baseName, tables.NAMED_KIND, at)

	// ensure that there's a kind of this name. pluralKinds, singularParent
	pluralKinds := lang.Pluralize(baseName)
	k.NewKind(k.NewName(lang.Breakcase(pluralKinds), tables.NAMED_PLURAL_KINDS, at),
		k.NewName("thing", tables.NAMED_KIND, at))
	//
	countedTypeTrait := k.NewName("counted", tables.NAMED_TRAIT, at)
	printedNameProp := k.NewName("printed_name", tables.NAMED_FIELD, at)

	for i := 0; i < cnt; i++ {
		countedNoun := k.newCounter(baseName, reader.Position{})
		noun := k.NewName(countedNoun.Offset, "noun", at)
		k.NewNoun(noun, namedSingularKind)
		k.NewValue(noun, countedTypeTrait, true)
		k.NewValue(noun, printedNameProp, baseName)
		// needed for relations, etc.
		k.Env().Recent.Nouns.Add(noun)
	}
	return
}

func (op *NamedNoun) ReadNamedNoun(k *Importer) (err error) {
	if noun, e := NewNounName(k, op.Name); e != nil {
		err = e
	} else {
		k.Env().Recent.Nouns.Add(noun)
		// pick common or proper based on noun capitalization.
		// fix: implicitly generated facts should be considered preliminary
		// so that authors can override them.
		var traitStr string
		detStr, detFound := composer.FindChoice(&op.Determiner, op.Determiner.Str)
		if detStr == "our" {
			if first, _ := utf8.DecodeRuneInString(noun.String()); unicode.ToUpper(first) == first {
				traitStr = "proper_named"
			}
		}
		if len(traitStr) > 0 {
			typeTrait := k.NewName(traitStr, tables.NAMED_TRAIT, op.Name.At.String())
			k.NewValue(noun, typeTrait, true)
		}

		// record any custom determiner
		if !detFound {
			// set the indefinite article field
			article := k.NewName("indefinite_article", tables.NAMED_FIELD, op.Name.At.String())
			k.NewValue(noun, article, detStr)

			// create a "indefinite article" field for all objects
			if once := "named_noun"; k.Once(once) {
				domain := k.Env().Game.Domain
				objects := k.NewDomainName(domain, "objects", tables.NAMED_KINDS, once)
				indefinite := k.NewDomainName(domain, "indefinite_article", tables.NAMED_FIELD, once)
				k.NewField(objects, indefinite, tables.PRIM_TEXT, "")
			}
		}
	}
	return
}

// ex. "[the box] (is a) (closed) (container) ((on) (the beach))"
func (op *KindOfNoun) ImportNouns(k *Importer) (err error) {
	if kind, e := NewSingularKind(k, op.Kind); e != nil {
		err = e
	} else {
		//
		var traits []eph.Named
		if ts := op.Trait; ts != nil {
			for _, t := range ts {
				if t, e := NewTrait(k, t); e != nil {
					err = errutil.Append(err, e)
				} else {
					traits = append(traits, t)
				}
			}
		}
		if err == nil {
			// we collected the nouns and delayed processing them till now.
			for _, noun := range k.Env().Recent.Nouns.Subjects {
				k.NewNoun(noun, kind)
				for _, trait := range traits {
					k.NewValue(noun, trait, true) // the value of the trait for the noun is true
				}
			}

			//
			if op.NounRelation != nil {
				err = op.NounRelation.ImportNouns(k)
			}
		}
	}
	return
}

// ex. [the cat and the hat] (are) (in) (the book)
// ex. [Hector and Maria] (are) (suspicious of) (Santa and Santana).
func (op *NounRelation) ImportNouns(k *Importer) (err error) {
	if rel, e := NewRelation(k, op.Relation); e != nil {
		err = e
	} else if e := k.Env().Recent.Nouns.CollectObjects(func() (err error) {
		return ImportNamedNouns(k, op.Nouns)
	}); e != nil {
		err = e
	} else {
		domain := k.Env().Current.Domain
		for _, subject := range k.Env().Recent.Nouns.Subjects {
			for _, object := range k.Env().Recent.Nouns.Objects {
				k.NewRelative(object, rel, subject, domain)
			}
		}
	}
	return
}

//
func (op *NounTraits) ImportNouns(k *Importer) (err error) {
	for _, t := range op.Trait {
		if trait, e := NewTrait(k, t); e != nil {
			err = e
			break
		} else {
			for _, noun := range k.Env().Recent.Nouns.Subjects {
				k.NewValue(noun, trait, true) // the value of the trait for the noun is true
			}
		}
	}
	return
}
