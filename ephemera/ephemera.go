package ephemera

import (
	"bytes"
	"database/sql"
	"strings"

	"git.sr.ht/~ionous/iffy/ephemera/eph"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/tables"
)

type Recorder struct {
	srcId int64
	cache *tables.Cache
}

func NewRecorder(db *sql.DB) *Recorder {
	return &Recorder{cache: tables.NewCache(db)}
}

func (k *Recorder) SetSource(srcURI string) *Recorder {
	k.srcId = k.cache.MustGetId(eph_source, srcURI)
	return k
}

// NewName records a user-specified string, including a category and location,
// and returns a unique identifier for it.
// Category is likely one of kind, noun, aspect, attribute, property, relation.
// The format of the location ofs depends on the data source.
func (k *Recorder) NewName(name, category, ofs string) (ret eph.Named) {
	return k.NewDomainName(eph.Named{}, name, category, ofs)
}

func (k *Recorder) NewDomainName(domain eph.Named, name, category, ofs string) (ret eph.Named) {
	// normalize names in an attempt to simplify lookup of some names
	// many tests would have to be adjusted to be able to handle normalization wholesale
	// so for now make this opt-in.
	norm := strings.TrimSpace(name)
	namedId := k.cache.MustGetId(eph_named, norm, name, category, domain, k.srcId, ofs)
	return eph.MakeName(namedId, norm)
}

type Prog struct{ eph.Named }

func (k *Recorder) NewProg(typeName string, cmd jsn.Marshalee) (ret Prog, err error) {
	var buf bytes.Buffer
	if e := cout.Marshal(&buf, cmd); e != nil {
		err = e
	} else {
		id := k.cache.MustGetId(eph_prog, k.srcId, typeName, buf.Bytes())
		ret = Prog{eph.MakeName(id, typeName)}
	}
	return
}

var None eph.Named

// NewAlias provides a new name for another name.
func (k *Recorder) NewAlias(alias, actual eph.Named) {
	k.cache.Must(eph_alias, alias, actual)
}

// NewAspect explicitly declares the existence of an aspect.
func (k *Recorder) NewAspect(aspect eph.Named) {
	k.cache.Must(eph_aspect, aspect)
}

// NewCertainty supplies a kind with a default trait.
func (k *Recorder) NewCertainty(certainty string, trait, kind eph.Named) {
	// usually fast horses.
	k.cache.Must(eph_certainty, certainty, trait, kind)
}

// NewDefault supplies a kind with a default value;
// see also NewValue
func (k *Recorder) NewDefault(kind, field eph.Named, value interface{}) {
	// horses height 5.
	k.cache.Must(eph_default, kind, field, value)
}

// NewKind connects a kind (plural) to its parent kind (singular).
// ex. cats are a kind of animal.
func (k *Recorder) NewKind(kind, parent eph.Named) {
	k.cache.Must(eph_kind, kind, parent)
}

// NewNoun connects a noun to its kind (singular).
func (k *Recorder) NewNoun(noun, kind eph.Named) {
	k.cache.Must(eph_noun, noun, kind)
}

//  a pattern or pattern parameter
func (k *Recorder) NewPatternDecl(pattern, param, patternType eph.Named, affinity string) {
	k.cache.Must(eph_pattern, pattern, param, patternType, affinity, Prog{})
}

//  a pattern initializer
func (k *Recorder) NewPatternInit(pattern, param, patternType eph.Named, affinity string, prog Prog) {
	k.cache.Must(eph_pattern, pattern, param, patternType, affinity, prog)
}

// a reference to a pattern declared elsewhere
func (k *Recorder) NewPatternRef(pattern, param, patternType eph.Named, affinity string) {
	k.cache.Must(eph_pattern, pattern, param, patternType, affinity, -1)
}

func (k *Recorder) NewPatternRule(pattern, target, domain eph.Named, handler Prog) {
	k.cache.Must(eph_rule, pattern, target, domain, handler)
}

// NewPlural maps the plural form of a name to its singular form.
func (k *Recorder) NewPlural(plural, singular eph.Named) {
	k.cache.Must(eph_plural, plural, singular)
}

// NewField property in the named kind.
func (k *Recorder) NewField(kind, prop eph.Named, primType, primAff string) {
	k.cache.Must(eph_field, kind, prop, primType, primAff)
}

// NewRelation defines a connection between a primary and secondary kind.
func (k *Recorder) NewRelation(relation, primaryKind, secondaryKind eph.Named, cardinality string) {
	k.cache.Must(eph_relation, relation, primaryKind, secondaryKind, cardinality)
}

// NewRelative connects two named nouns using a verb stem.
func (k *Recorder) NewRelative(primary, stem, secondary, domain eph.Named) {
	k.cache.Must(eph_relative, primary, stem, secondary, domain)
}

func (k *Recorder) NewTestProgram(test eph.Named, prog Prog) {
	k.cache.Must(eph_check, test, prog)
}

func (k *Recorder) NewTestExpectation(test eph.Named, testType string, expect string) {
	k.cache.Must(eph_expect, test, testType, expect)
}

// NewTrait records a member of an aspect and its order ( rank. )
func (k *Recorder) NewTrait(trait, aspect eph.Named, rank int) {
	k.cache.Must(eph_trait, trait, aspect, rank)
}

// NewValue assigns the property of a noun a value;
// traits can be assigned by naming the individual trait and setting a true ( or false ) value.
func (k *Recorder) NewValue(noun, prop eph.Named, value interface{}) {
	k.cache.Must(eph_value, noun, prop, value)
}

// NewRelative connects two specific nouns using a verb.
func (k *Recorder) NewVerb(stem, relation eph.Named, verb string) {
	k.cache.Must(eph_verb, stem, relation, verb)
}

var eph_alias = tables.Insert("eph_alias", "idNamedAlias", "idNamedActual")
var eph_aspect = tables.Insert("eph_aspect", "idNamedAspect")
var eph_certainty = tables.Insert("eph_certainty", "certainty", "idNamedTrait", "idNamedKind")
var eph_check = tables.Insert("eph_check", "idNamedTest", "idProg")
var eph_default = tables.Insert("eph_default", "idNamedKind", "idNamedProp", "value")
var eph_expect = tables.Insert("eph_expect", "idNamedTest", "testType", "expect")
var eph_field = tables.Insert("eph_field", "idNamedKind", "idNamedField", "primType", "primAff")
var eph_rule = tables.Insert("eph_rule", "idNamedPattern", "idNamedTarget", "idNamedDomain", "idProg")
var eph_kind = tables.Insert("eph_kind", "idNamedKind", "idNamedParent")
var eph_named = tables.Insert("eph_named", "name", "og", "category", "domain", "idSource", "offset")
var eph_noun = tables.Insert("eph_noun", "idNamedNoun", "idNamedKind")
var eph_pattern = tables.Insert("eph_pattern", "idNamedPattern", "idNamedParam", "idNamedType", "affinity", "idProg")
var eph_plural = tables.Insert("eph_plural", "idNamedPlural", "idNamedSingluar")
var eph_prog = tables.Insert("eph_prog", "idSource", "progType", "prog")
var eph_relation = tables.Insert("eph_relation", "idNamedRelation", "idNamedKind", "idNamedOtherKind", "cardinality")
var eph_relative = tables.Insert("eph_relative", "idNamedHead", "idNamedStem", "idNamedDependent", "idNamedDomain")
var eph_source = tables.Insert("eph_source", "src")
var eph_trait = tables.Insert("eph_trait", "idNamedTrait", "idNamedAspect", "rank")
var eph_value = tables.Insert("eph_value", "idNamedNoun", "idNamedProp", "value")
var eph_verb = tables.Insert("eph_verb", "idNamedStem", "idNamedRelation", "verb")
