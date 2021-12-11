package tables

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
