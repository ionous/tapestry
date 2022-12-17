package idl

import (
	"git.sr.ht/~ionous/tapestry/tables"
)

// see sql table definitions and additional notes in "tapestry/tables/idl.sql"

var Op = tables.Insert("idl_op", "name", "package", "uses", "closed")
var Sig = tables.Insert("idl_sig", "op", "slot", "hash", "signature")
var Str = tables.Insert("idl_str", "op", "label", "value")
var Enum = tables.Insert("idl_enum", "op", "label", "value")
var Swap = tables.Insert("idl_swap", "op", "label", "type")
var Term = tables.Insert("idl_term", "op", "term", "label", "type", "private", "optional", "repeats")
