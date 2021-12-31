package mdl

import (
  "git.sr.ht/~ionous/tapestry/tables"
)

// see sql table definitions and additional notes in "model.sql"

// the pattern half of Start; "domain, kind, field" are a pointer into Field
var Assign = tables.Insert("mdl_assign", "domain", "kind", "field", "value")

// author tests of stories
var Check = tables.Insert("mdl_check", "domain", "name", "value", "affinity", "prog", "at")

// domain name and materialized parents separated by commas
var Domain = tables.Insert("mdl_domain", "domain", "path", "at")

// fix? the domain exists to uniquely identify the kind,
// it's not the (sub)domain in which the field was declared.
var Field = tables.Insert("mdl_field", "domain", "kind", "field", "affinity", "type", "at")

var Grammar = tables.Insert("mdl_grammar", "domain", "name", "prog", "at")

// singular name of kind and materialized hierarchy of ancestors separated by commas
var Kind = tables.Insert("mdl_kind", "domain", "kind", "path", "at")

// words for authors and game players refer to nouns
// follows the domain rules of Noun.
var Name = tables.Insert("mdl_name", "domain", "noun", "name", "rank", "at")

// the domain tells the scope in which the noun was defined
// ( the same as - or a child of - the domain of the kind )
var Noun = tables.Insert("mdl_noun", "domain", "noun", "kind", "at")

// domain captures the scope in which the pairing was defined.
// within that scope: the noun, relation, and otherNoun are all unique names --
// even if they are not unique globally, and even if they a broader/different scope than the pair's domain.
var Pair = tables.Insert("mdl_pair", "domain", "relKind", "oneNoun", "otherNoun", "at")

// doesn't store "at" because its kind already defines that
var Pat = tables.Insert("mdl_pat", "domain", "kind", "labels", "result")

// a plural word ("many") can have at most one singular definition per domain
// ie. "people" and "persons" are valid plurals of "person",
// but "people" as a singular can only be defined as "person" ( not also "human" )
var Plural = tables.Insert("mdl_plural", "domain", "many", "one", "at")

// relation and constraint between two kinds of nouns
//  fix? the data is duplicated in kinds and fields... should this be removed?
// might also consider adding a "cardinality" field to the relation kind, and then use init for individual relations
var Rel = tables.Insert("mdl_rel", "domain", "relKind", "oneKind", "otherKind", "cardinality", "at")

//
var Rule = tables.Insert("mdl_rule", "domain", "pattern", "target", "phase", "filter", "prog", "at")

// the noun half of what was Start.
// "domain, noun, field" reference a join of Noun and Kind to get a filtered Field.
var Value = tables.Insert("mdl_value", "domain", "noun", "field", "value", "at")
