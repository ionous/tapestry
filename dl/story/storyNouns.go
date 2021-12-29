package story

import (
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/eph"
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
	} else if cnt > 0 {
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
	// note: kind is phrased in the singular here when count is 1, plural otherwise.
	// but, because of "Recent.Nouns" processing we have to generate some sort of noun name *immediately*
	// ( itd be nice to have a more start and stop importer, where we could delay processing of branches of the tree. )
	kindOrKinds := op.Name.String()
	names := make([]string, cnt)
	for i := 0; i < cnt; i++ {
		noun := k.newCounter(kindOrKinds, reader.Position{})
		k.Env().Recent.Nouns.Add(noun) // for relations, etc.
		k.Write(&eph.EphValues{Noun: noun, Field: "counted", Value: B(true)})
		names[i] = noun
	}
	k.Write(eph.PhaseFunction{eph.AncestryPhase,
		func(c *eph.Catalog, d *eph.Domain, at string) (err error) {
			// by now, plurals will be determined, so we can determine which is which.
			var kind, kinds string
			if cnt == 1 {
				kind = kindOrKinds
				kinds, err = d.Pluralize(kindOrKinds)
			} else {
				kinds = kindOrKinds
				kind, err = d.Singularize(kindOrKinds)
			}
			if err == nil {
				if e := d.AddEphemera(eph.EphAt{at, &eph.EphKinds{Kinds: kinds, From: "thing"}}); e != nil {
					err = e
				} else {
					for _, n := range names {
						if e := d.AddEphemera(eph.EphAt{at, &eph.EphNouns{Noun: n, Kind: kindOrKinds}}); e != nil {
							err = e
						} else if e := d.AddEphemera(eph.EphAt{at, &eph.EphAliases{
							// so that typing "triangle" means "triangles_1"
							ShortName: n, Aliases: []string{kind},
						}}); e != nil {
							err = e
							break
						} else if e := d.AddEphemera(eph.EphAt{at, &eph.EphValues{
							// so that printing "triangles_1" yields "triangle"
							// FIX: itd make a lot more sense to have a default value for the kind
							Noun: n, Field: "printed_name", Value: T(kind),
						}}); e != nil {
							err = e
							break
						}
					}
				}
			}
			return
		}},
	)
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
		k.Write(&eph.EphValues{Noun: noun, Field: "indefinite_article", Value: T(detStr)})
	}
	// pick common or proper based on noun capitalization.
	// fix: implicitly generated facts should be considered preliminary so that authors can override them.
	if detStr == "our" {
		if first, _ := utf8.DecodeRuneInString(noun); unicode.ToUpper(first) == first {
			k.Write(&eph.EphValues{Noun: noun, Field: "proper_named", Value: B(true)})
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
			k.Write(&eph.EphValues{Noun: noun, Field: trait.String(), Value: B(true)})
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
			k.Write(&eph.EphValues{Noun: noun, Field: trait.String(), Value: B(true)})
		}
	}
	return
}
