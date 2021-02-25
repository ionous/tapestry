package tables

// primType
const (
	PRIM_TEXT   = "text"   // string
	PRIM_DIGI   = "number" // number
	PRIM_BOOL   = "bool"   // boolean (rare, more usually aspect)
	PRIM_ASPECT = "aspect" // string
	PRIM_TRAIT  = "trait"  // string
)

// eph_named category
// fix: there shouldnt be a a central name table
const (
	NAMED_ARGUMENT     = "argument"
	NAMED_ASPECT       = "aspect"
	NAMED_CERTAINTY    = "certainty"
	NAMED_EVENT        = "event"
	NAMED_FIELD        = "field"
	NAMED_KIND         = "singular_kind"
	NAMED_KINDS        = "kind" // FIX: why are only the auto-generated types using this?
	NAMED_PLURAL_KINDS = "plural_kinds"
	NAMED_PROPERTY     = "prop" // field, trait, or aspect
	NAMED_NOUN         = "noun"
	NAMED_PATTERN      = "pattern"
	NAMED_LOCAL        = "local"
	NAMED_RETURN       = "return"
	NAMED_PARAMETER    = "parameter"
	NAMED_RELATION     = "relation"
	NAMED_VERB         = "verb"
	NAMED_SCENE        = "scene"
	NAMED_TEST         = "test"
	NAMED_TRAIT        = "trait"
	NAMED_TYPE         = "type" // autogenerated types: bool, expr, comp, prog
)

// cardinality
const (
	ONE_TO_ONE   = "one_one"
	ONE_TO_MANY  = "one_any"
	MANY_TO_ONE  = "any_one"
	MANY_TO_MANY = "any_any"
)

// certainty
const (
	USUALLY = "usually"
	ALWAYS  = "always"
	SELDOM  = "seldom"
	NEVER   = "never"
)
