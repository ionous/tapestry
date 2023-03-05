package story

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

type NamedNoun interface {
	ImportNoun(*imp.Importer) error
}

type SingularNoun interface {
	NamedNoun
	NounName() string
	UniformString() (string, error)
}

func CollectSubjectNouns(k *imp.Importer, els []NamedNoun) error {
	return k.Env().Recent.Nouns.CollectSubjects(func() error {
		return ImportNamedNouns(k, els)
	})
}

func CollectObjectNouns(k *imp.Importer, els []NamedNoun) error {
	return k.Env().Recent.Nouns.CollectObjects(func() error {
		return ImportNamedNouns(k, els)
	})
}

func ImportNamedNouns(k *imp.Importer, els []NamedNoun) (err error) {
	for _, el := range els {
		if e := el.ImportNoun(k); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func ImportNouns(k *imp.Importer, nouns []string) (ret []string, err error) {
	for _, noun := range nouns {
		// FIX: this should be happening during the weave, not during import.
		article, word := lang.SliceArticle(noun)
		name := NounNamed{Name: NounName{word}}

		var legacy NamedNoun
		if len(article) == 0 {
			legacy = &name
		} else {
			legacy = &CommonNoun{
				Determiner: Determiner{article},
				Noun:       name,
			}
		}
		if e := legacy.ImportNoun(k); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = append(ret, word)
		}
	}
	return
}

// declare a noun class that has several default fields
func declareNounClass(k *imp.Importer) {
	if once := "noun"; k.Once(once) {
		k.WriteOnce(&eph.EphKinds{Kind: "objects", Ancestor: kindsOf.Kind.String()})
		// common or proper nouns ( rabbit, vs. Roger )
		k.AddImplicitAspect("noun_types", "objects", "common_named", "proper_named", "counted")
		// whether a player can refer to an object by its name.
		k.AddImplicitAspect("private_names", "objects", "publicly_named", "privately_named")
	}
}

// ex. "two triangles" -> triangle is a kind of thing
// fix? consider a specific counted noun phrase; the noun phrase needs more work.
// also, we probably want noun stacks not individually duplicated names
func (op *CountedNouns) ImportNoun(k *imp.Importer) (err error) {
	if once := "printed_name"; k.Once(once) {
		k.WriteOnce(&eph.EphKinds{Kind: "objects", Contain: []eph.EphParams{{Name: "printed_name", Affinity: eph.Affinity{eph.Affinity_Text}}}})
	}
	if cnt, ok := lang.WordsToNum(op.Count); !ok {
		err = errutil.New("couldnt turn", op.Count, "into a number")
	} else if cnt > 0 {

		// note: kind is phrased in the singular here when count is 1, plural otherwise.
		// but, because of "Recent.Nouns" processing we have to generate some sort of noun name *immediately*
		// ( itd be nice to have a more start and stop importer, where we could delay processing of branches of the tree. )
		kindOrKinds := op.Kinds.String()
		names := make([]string, cnt)
		for i := 0; i < cnt; i++ {
			noun := k.NewCounter(kindOrKinds, nil)
			k.Env().Recent.Nouns.Add(noun) // for relations, etc.
			k.WriteEphemera(&eph.EphValues{Noun: noun, Field: "counted", Value: B(true)})
			names[i] = noun
		}
		k.WriteEphemera(eph.PhaseFunction{eph.AncestryPhase,
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
					if e := d.AddEphemera(at, &eph.EphKinds{Kind: kinds, Ancestor: "thing"}); e != nil {
						err = e
					} else {
						for _, n := range names {
							if e := d.AddEphemera(at, &eph.EphNouns{Noun: n, Kind: kindOrKinds}); e != nil {
								err = e
							} else if e := d.AddEphemera(at, &eph.EphAliases{
								// so that typing "triangle" means "triangles_1"
								ShortName: n, Aliases: []string{kind},
							}); e != nil {
								err = e
								break
							} else if e := d.AddEphemera(at, &eph.EphValues{
								// so that printing "triangles_1" yields "triangle"
								// FIX: itd make a lot more sense to have a default value for the kind
								Noun: n, Field: "printed_name", Value: T(kind),
							}); e != nil {
								err = e
								break
							}
						}
					}
				}
				return
			}},
		)
	}
	return
}

func (op *CommonNoun) NounName() string {
	return op.Noun.NounName()
}

func (op *CommonNoun) UniformString() (string, error) {
	return op.Noun.UniformString()
}

func (op *CommonNoun) ImportNoun(k *imp.Importer) (err error) {
	declareNounClass(k)
	detStr, detFound := composer.FindChoice(&op.Determiner, op.Determiner.Str)
	// setup the indefinite article
	if !detFound {
		// create a "indefinite article" field for all objects
		if k.Once("named_noun") {
			k.WriteOnce(&eph.EphKinds{Kind: "objects", Contain: []eph.EphParams{{Name: "indefinite_article", Affinity: eph.Affinity{eph.Affinity_Text}}}})
		}
		// set the indefinite article field
		k.WriteEphemera(&eph.EphValues{Noun: op.Noun.NounName(), Field: "indefinite_article", Value: T(detStr)})
	}
	op.Noun.addNoun(k, detStr)
	return
}

func (op *NounNamed) NounName() string {
	return op.Name.NounName()
}

func (op *NounNamed) UniformString() (string, error) {
	return op.Name.UniformString()
}

func (op *NounNamed) ImportNoun(k *imp.Importer) (err error) {
	declareNounClass(k)
	op.addNoun(k, "our")
	return
}

func (op *NounName) NounName() string {
	return strings.TrimSpace(op.Str)
}

func (op *NounName) UniformString() (ret string, err error) {
	if u, ok := eph.UniformString(op.Str); !ok {
		err = eph.InvalidString(op.Str)
	} else {
		ret = u
	}
	return
}

func (op *NounNamed) addNoun(k *imp.Importer, detStr string) {
	// strip extraneous spaces that exist for obscure mainline reasons;
	// testing ToUpper against space ( below ) for capitals was making nouns starting with spaces proper named.
	noun := op.NounName()
	k.Env().Recent.Nouns.Add(noun)

	// pick common or proper based on noun capitalization.
	// fix: implicitly generated facts should be considered preliminary so that authors can override them.
	if detStr == "our" {
		if first, _ := utf8.DecodeRuneInString(noun); unicode.ToUpper(first) == first {
			k.WriteEphemera(&eph.EphValues{Noun: noun, Field: "proper_named", Value: B(true)})
		}
	}
}
