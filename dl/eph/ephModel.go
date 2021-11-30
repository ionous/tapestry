package eph

import (
	"git.sr.ht/~ionous/iffy/tables"
)

// database/sql like interface
type Writer interface {
	Write(q string, args ...interface{}) error
}

// domain name and materialized parents separated by commas
var mdl_domain = tables.Insert("mdl_domain", "domain", "path")

// a plural word ("many") can have at most one singular definition per domain
// ie. "people" and "persons" are valid plurals of "person",
// but "people" as a singular can only be defined as "person" ( not "cat" )
var mdl_plural = tables.Insert("mdl_plural", "domain", "many", "one")

// singular name of kind and materialized hierarchy of ancestors separated by commas
var mdl_kind = tables.Insert("mdl_kind", "domain", "kind", "path")

// while the domain does tell us where the field was defined, it primarily exists to differentiate rival kinds
var mdl_field = tables.Insert("mdl_field", "domain", "kind", "field", "affinity", "type", "at")
