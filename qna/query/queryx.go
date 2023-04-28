package query

// Queryx adds builder checks to the runtime interface
// tbd if this should be in a separate struct or package
// for now, its nice to have all the db bits together.

type Queryx interface {
	Query
	// check the active domains for rival definitions
	// fix? ideally would be able to check *before* activation.
	FindActiveConflicts(func(domain, key, value, at string) error) error

	FindPluralDefinitions(plural string, cb func(domain, one, at string) error) error
}
