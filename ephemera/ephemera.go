package ephemera

import (
	"database/sql"
	"strings"

	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/tables"
)

type Recorder struct {
	srcId          int64
	cache          *tables.Cache
	normalizeNames bool
}

func NewNormalizingRecorder(srcURI string, db *sql.DB) *Recorder {
	rec := NewRecorder(srcURI, db)
	rec.normalizeNames = true
	return rec
}

// backwards compatibility method;
// tests are written without expecting name normalization so...
func NewRecorder(srcURI string, db *sql.DB) (ret *Recorder) {
	cache := tables.NewCache(db)
	srcId := cache.Must(eph_source, srcURI)
	return &Recorder{srcId, cache, false}
}

// NewName records a user-specified string, including a category and location,
// and returns a unique identifier for it.
// Category is likely one of kind, noun, aspect, attribute, property, relation.
// The format of the location ofs depends on the data source.
func (r *Recorder) NewName(name, category, ofs string) (ret Named) {
	norm := strings.TrimSpace(name)
	// many tests would have to be adjusted to be able to handle normalization wholesale
	// so for now make this opt-in.
	if r.normalizeNames {
		switch category {
		case tables.NAMED_TEST, tables.NAMED_PATTERN:
			norm = lang.Camelize(norm)
		}
	}
	namedId := r.cache.Must(eph_named, norm, category, r.srcId, ofs, name)
	return Named{namedId, norm}
}

type Prog struct{ Named }

// fix: this should probably take "ofs" just like NewName does.
func (r *Recorder) NewProg(rootType string, blob []byte) (ret Prog) {
	id := r.cache.Must(eph_prog, r.srcId, rootType, blob)
	ret = Prog{Named{id, rootType}}
	return
}

var None Named

// NewAlias provides a new name for another name.
func (r *Recorder) NewAlias(alias, actual Named) {
	r.cache.Must(eph_alias, alias, actual)
}

// NewAspect explicitly declares the existence of an aspect.
func (r *Recorder) NewAspect(aspect Named) {
	r.cache.Must(eph_aspect, aspect)
}

// NewCertainty supplies a kind with a default trait.
// usually fast horses.
func (r *Recorder) NewCertainty(certainty string, trait, kind Named) {
	// usually fast horses.
	r.cache.Must(eph_certainty, certainty, trait, kind)
}

// NewDefault supplies a kind with a default value.
// height horses 5.
func (r *Recorder) NewDefault(kind, field Named, value interface{}) {
	// height horses 5.
	r.cache.Must(eph_default, kind, field, value)
}

// NewKind connects a kind (plural) to its parent kind (singular).
// cats are a kind of animal.
func (r *Recorder) NewKind(kind, parent Named) {
	r.cache.Must(eph_kind, kind, parent)
}

// NewNoun connects a noun to its kind (singular).
func (r *Recorder) NewNoun(noun, kind Named) {
	r.cache.Must(eph_noun, noun, kind)
}

// declare a pattern or pattern parameter
func (r *Recorder) NewPatternDecl(pattern, param, patternType Named) {
	r.cache.Must(eph_pattern, pattern, param, patternType, true)
}

//
func (r *Recorder) NewPatternRef(pattern, param, patternType Named) {
	r.cache.Must(eph_pattern, pattern, param, patternType, false)
}

func (r *Recorder) NewPatternRule(pattern Named, handler Prog) {
	r.cache.Must(eph_rule, pattern, handler)
}

// NewPlural maps the plural form of a name to its singular form.
func (r *Recorder) NewPlural(plural, singluar Named) {
	r.cache.Must(eph_plural, plural, singluar)
}

// NewField property in the named kind.
func (r *Recorder) NewField(kind, prop Named, primType string) {
	r.cache.Must(eph_field, primType, kind, prop)
}

// NewRelation defines a connection between a primary and secondary kind.
func (r *Recorder) NewRelation(relation, primaryKind, secondaryKind Named, cardinality string) {
	r.cache.Must(eph_relation, relation, primaryKind, secondaryKind, cardinality)
}

// NewRelative connects two specific nouns using a verb stem.
func (r *Recorder) NewRelative(primary, stem, secondary Named) {
	r.cache.Must(eph_relative, primary, stem, secondary)
}

func (r *Recorder) NewTest(test Named, prog Prog, expect string) {
	r.cache.Must(eph_check, test, prog, expect)
}

// NewTrait records a member of an aspect and its order ( rank. )
func (r *Recorder) NewTrait(trait, aspect Named, rank int) {
	r.cache.Must(eph_trait, trait, aspect, rank)
}

// NewValue assigns the property of a noun a value;
// traits can be assigned by naming the individual trait and setting a true ( or false ) value.
func (r *Recorder) NewValue(noun, prop Named, value interface{}) {
	r.cache.Must(eph_value, noun, prop, value)
}

// NewRelative connects two specific nouns using a verb.
func (r *Recorder) NewVerb(stem, relation Named, verb string) {
	r.cache.Must(eph_verb, stem, relation, verb)
}

var eph_alias = tables.Insert("eph_alias", "idNamedAlias", "idNamedActual")
var eph_aspect = tables.Insert("eph_aspect", "idNamedAspect")
var eph_certainty = tables.Insert("eph_certainty", "certainty", "idNamedTrait", "idNamedKind")
var eph_check = tables.Insert("eph_check", "idNamedTest", "idProg", "expect")
var eph_default = tables.Insert("eph_default", "idNamedKind", "idNamedProp", "value")
var eph_field = tables.Insert("eph_field", "primType", "idNamedKind", "idNamedField")
var eph_rule = tables.Insert("eph_rule", "idNamedPattern", "idProg")
var eph_kind = tables.Insert("eph_kind", "idNamedKind", "idNamedParent")
var eph_named = tables.Insert("eph_named", "name", "category", "idSource", "offset", "og")
var eph_noun = tables.Insert("eph_noun", "idNamedNoun", "idNamedKind")
var eph_pattern = tables.Insert("eph_pattern", "idNamedPattern", "idNamedParam", "idNamedType", "decl")
var eph_plural = tables.Insert("eph_plural", "idNamedPlural", "idNamedSingluar")
var eph_prog = tables.Insert("eph_prog", "idSource", "type", "prog")
var eph_relation = tables.Insert("eph_relation", "idNamedRelation", "idNamedKind", "idNamedOtherKind", "cardinality")
var eph_relative = tables.Insert("eph_relative", "idNamedHead", "idNamedStem", "idNamedDependent")
var eph_source = tables.Insert("eph_source", "src")
var eph_trait = tables.Insert("eph_trait", "idNamedTrait", "idNamedAspect", "rank")
var eph_value = tables.Insert("eph_value", "idNamedNoun", "idNamedProp", "value")
var eph_verb = tables.Insert("eph_verb", "idNamedStem", "idNamedRelation", "verb")
