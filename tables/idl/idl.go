package idl

import (
	"git.sr.ht/~ionous/tapestry/tables"
)

// see sql table definitions and additional notes in "tapestry/tables/idl.sql"

var Op = tables.Insert("idl_op", "name", "package", "spec")
var Sig = tables.Insert("idl_sig", "op", "slot", "hash", "body")
var Enum = tables.Insert("idl_enum", "op", "value")
var Term = tables.Insert("idl_term", "op", "name", "label", "type", "private", "optional", "repeats")
var Markup = tables.Insert("idl_markup", "op", "key", "value")
