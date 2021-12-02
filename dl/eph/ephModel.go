package eph

import (
	"git.sr.ht/~ionous/iffy/tables"
)

// database/sql like interface
type Writer interface {
	Write(q string, args ...interface{}) error
}

// domain name and materialized parents separated by commas
var mdl_domain = tables.Insert("mdl_domain", "domain", "path", "at")

// a plural word ("many") can have at most one singular definition per domain
// ie. "people" and "persons" are valid plurals of "person",
// but "people" as a singular can only be defined as "person" ( not "cat" )
var mdl_plural = tables.Insert("mdl_plural", "domain", "many", "one")

// singular name of kind and materialized hierarchy of ancestors separated by commas
var mdl_kind = tables.Insert("mdl_kind", "domain", "kind", "path", "at")

// fix? the domain exists to uniquely identify the kind,
// it's not the (sub)domain in which the field was declared.
var mdl_field = tables.Insert("mdl_field", "domain", "kind", "field", "affinity", "type", "at")

//
var mdl_rel = tables.Insert("mdl_rel", "domain", "relation", "kind", "cardinality", "otherKind", "at")

// the domain tells the domain in which the noun was defined
// ( which is implicitly the same as or the child of the domain for kind )
var mdl_noun = tables.Insert("mdl_noun", "domain", "noun", "kind", "at")

// words for authors and game players refer to nouns
// follows the domain rules of mdl_noun.
var mdl_name = tables.Insert("mdl_name", "domain", "noun", "name", "rank", "at")
